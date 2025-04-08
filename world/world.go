package world

import (
	"PowerClassic/entity"
	"PowerClassic/events"
	"PowerClassic/packets"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/rs/zerolog/log"
	"sort"
)

type WorldRunnable struct {
	Run func(ctx *actor.Context, w *World)
}
type ChunkRunnable struct {
	Run func(ctx *actor.Context, c *Chunk)
}

//type CreateFlatWorld struct{}
//type CreateFlatWorldResponse struct {
//	Error error
//}

type World struct {
	pid    *actor.PID
	Width  int
	Height int
	Depth  int
	Chunks map[[2]int]*Chunk

	entityMap map[entity.Entity]*Chunk

	entities []entity.Entity

	idManager *entity.EntityIdManager

	eng *actor.Engine
}

func NewWorld(eng *actor.Engine, width, height, depth int) *World {
	chunks := make(map[[2]int]*Chunk)
	chunksX := (width + ChunkWidth - 1) / ChunkWidth
	chunksZ := (depth + ChunkDepth - 1) / ChunkDepth

	for cx := 0; cx < chunksX; cx++ {
		for cz := 0; cz < chunksZ; cz++ {
			c := NewChunk(cx, cz, height)

			c.pid = eng.Spawn(func() actor.Receiver {
				return c
			}, "chunk")

			chunks[[2]int{cx, cz}] = c
		}
	}
	return &World{
		Width:     width,
		Height:    height,
		Depth:     depth,
		Chunks:    chunks,
		idManager: entity.NewEntityIdManager(),
		eng:       eng,
	}
}

func (w *World) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *WorldRunnable:
		msg.Run(ctx, w)
	}
}

func (w *World) BroadcastEntityRunnable(ctx *actor.Context, runnable *entity.EntityRunnable) {
	for _, e := range w.entities {
		ctx.Send(e.GetPid(), runnable)
	}
}

func (w *World) GetPID() *actor.PID {
	return w.pid
}

func (w *World) SetPID(pid *actor.PID) {
	w.pid = pid
}

func (w *World) spawnEntityActor(ctx *actor.Context, e entity.Entity) {
	eid := ctx.SpawnChild(func() actor.Receiver {
		return e
	}, "entity")
	e.SetPid(eid)

}

// TODO: i don't like passing event here, come back to this
func (w *World) AddEntity(ctx *actor.Context, unsafe_E entity.Entity, evt *events.PlayerIdentificationEvent) {
	const SelfID = 255

	w.entities = append(w.entities, unsafe_E)
	w.spawnEntityActor(ctx, unsafe_E)

	ctx.Send(unsafe_E.GetPid(), &entity.EntityRunnable{Run: func(ctx *actor.Context, e entity.Entity) {
		w.checkChunkChange(ctx, e)

		se, isSe := e.(entity.SessionedEntity)
		if isSe {
			pks, err := w.getLevelDataChunkPackets()
			if err != nil {
				log.Err(fmt.Errorf("error getting level data chunk packets: %v", err))
				return
			}

			se.SendPacket(&packets.LevelInitializePacket{})
			for _, pk := range pks {
				se.SendPacket(pk)
			}

			se.SendPacket(&packets.LevelFinalizePacket{
				X: 1024,
				Y: 64,
				Z: 1024,
			})
		}

		e.Teleport(ctx, *evt.SpawnX(), *evt.SpawnY(), *evt.SpawnZ())

		if isSe {
			for _, other := range w.entities {
				if sessionedOther, ok := other.(entity.SessionedEntity); ok {
					id := e.Id()
					if other.Id() == e.Id() {
						id = SelfID
					}

					sessionedOther.SendPacket(&packets.SpawnPlayerPacket{
						PlayerId:   id,
						PlayerName: e.GetName(),
						X:          e.X(),
						Y:          e.Y(),
						Z:          e.Z(),
						Yaw:        0,
						Pitch:      0,
					})
				}

				if other.Id() != e.Id() {
					otherCopy := other

					ctx.Send(otherCopy.GetPid(), &entity.EntityRunnable{Run: func(ctx *actor.Context, existingE entity.Entity) {
						se.SendPacket(&packets.SpawnPlayerPacket{
							PlayerId:   existingE.Id(),
							PlayerName: existingE.GetName(),
							X:          existingE.X(),
							Y:          existingE.Y(),
							Z:          existingE.Z(),
							Yaw:        0,
							Pitch:      0,
						})
					}})
				}
			}
		}
	}})
}

func (w *World) GetNextEntityId() byte {
	return w.idManager.NextEntityId()
}

func (w *World) CreateFlatWorld() error {
	for _, chunk := range w.Chunks {
		for x := 0; x < ChunkWidth; x++ {
			for z := 0; z < ChunkDepth; z++ {
				if err := chunk.SetBlock(x, 2, z, 1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// TODO: move to actor message
func (w *World) Unsafe_getBlock(x, y, z int) (byte, error) {
	chunkX := x / ChunkWidth
	chunkZ := z / ChunkDepth
	localX := x % ChunkWidth
	localZ := z % ChunkDepth
	chunk, ok := w.Chunks[[2]int{chunkX, chunkZ}]
	if !ok {
		return 0, fmt.Errorf("chunk not found at (%d,%d)", chunkX, chunkZ)
	}
	return chunk.GetBlock(localX, y, localZ)
}

// TODO: move to actor message
func (w *World) Unsafe_setBlock(x, y, z int, block byte) error {
	chunkX := x / ChunkWidth
	chunkZ := z / ChunkDepth
	localX := x % ChunkWidth
	localZ := z % ChunkDepth
	chunk, ok := w.Chunks[[2]int{chunkX, chunkZ}]
	if !ok {
		return fmt.Errorf("chunk not found at (%d,%d)", chunkX, chunkZ)
	}
	return chunk.SetBlock(localX, y, localZ, block)
}

func (w *World) checkChunkChange(ctx *actor.Context, e entity.Entity) error {
	chunkX := int(e.X()) / ChunkWidth
	chunkZ := int(e.Z()) / ChunkDepth
	newChunkKey := [2]int{chunkX, chunkZ}

	newChunk, ok := w.Chunks[newChunkKey]
	if !ok {
		return fmt.Errorf("target chunk (%d, %d) not found", chunkX, chunkZ)
	}

	currentChunk, exists := w.entityMap[e]
	if exists && currentChunk == newChunk {
		return nil
	}
	ctx.Send(newChunk.pid, ChunkRunnable{Run: func(ctx *actor.Context, c *Chunk) {
		c.JoinChunk(e, ctx)
	}})

	if exists {
		ctx.Send(currentChunk.pid, ChunkRunnable{Run: func(ctx *actor.Context, c *Chunk) {
			c.LeaveChunk(e, ctx)
		}})
	}
	return nil
}

func (w *World) getLevelDataChunkPackets() ([]*packets.LevelDataChunkPacket, error) {
	totalBlocks := w.Width * w.Height * w.Depth
	levelData := make([]byte, totalBlocks)
	index := 0

	for y := 0; y < w.Height; y++ {
		for z := 0; z < w.Depth; z++ {
			for x := 0; x < w.Width; x++ {
				block, err := w.Unsafe_getBlock(x, y, z)
				if err != nil {
					return nil, err
				}
				levelData[index] = block
				index++
			}
		}
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(totalBlocks))

	rawData := append(header, levelData...)

	var compBuf bytes.Buffer
	gz := gzip.NewWriter(&compBuf)
	if _, err := gz.Write(rawData); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	compressedData := compBuf.Bytes()

	const PacketSize = 1024
	var pks []*packets.LevelDataChunkPacket
	totalLen := len(compressedData)
	for i := 0; i < totalLen; i += PacketSize {
		end := i + PacketSize
		if end > totalLen {
			end = totalLen
		}
		chunk := compressedData[i:end]
		if len(chunk) < PacketSize {
			padded := make([]byte, PacketSize)
			copy(padded, chunk)
			chunk = padded
		}

		var percent byte
		if end == totalLen {
			percent = 255
		} else {
			percent = byte((end * 255) / totalLen)
		}
		pkt := &packets.LevelDataChunkPacket{
			ChunkLength:     int16(end - i),
			ChunkData:       chunk,
			PercentComplete: percent,
		}
		pks = append(pks, pkt)
	}

	sort.Slice(pks, func(i, j int) bool {
		return pks[i].PercentComplete < pks[j].PercentComplete
	})

	return pks, nil
}

var _ actor.Receiver = (*World)(nil)
