package world

import (
	"PowerClassic/entity"
	"PowerClassic/messages"
	"PowerClassic/packets"
	"PowerClassic/player"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/rs/zerolog/log"
	"sort"
)

type CreateFlatWorld struct{}
type CreateFlatWorldResponse struct {
	Error error
}

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
	case *CreateFlatWorld:
		err := w.createFlatWorld()
		if err != nil {
			log.Error().Err(err).Msg("failed to create flat world")
		}
		ctx.Respond(&CreateFlatWorldResponse{
			Error: err,
		})
	case *messages.MovePlayer:
		log.Info().Msg("move player on world")
		tpk := &packets.PlayerTeleportPacket{
			PlayerId: msg.ID,
			X:        msg.X,
			Y:        msg.Y,
			Z:        msg.Z,
			Yaw:      0,
			Pitch:    0,
		}
		for _, child := range ctx.Children() {
			ctx.Send(child, tpk)
		}
	case *messages.AddEntity:
		w.addEntity(msg, ctx)
	}
}

func (w *World) GetPID() *actor.PID {
	return w.pid
}

func (w *World) SetPID(pid *actor.PID) {
	w.pid = pid
}

func (w *World) addEntity(msg *messages.AddEntity, ctx *actor.Context) {

	w.entities = append(w.entities, msg.E)

	eid := ctx.SpawnChild(func() actor.Receiver {
		return msg.E
	}, "entity")

	msg.E.SetPid(eid)
	pos, err := msg.E.GetPosition(ctx)

	if err != nil {
		panic(err)
	}
	err = w.moveEntity(ctx, msg.E, pos)
	if err != nil {
		log.Err(fmt.Errorf("error moving entity to chunk: %v", err))
		return
	}

	_, ok := msg.E.(entity.SessionedEntity)
	if ok {
		pks, err := w.getLevelDataChunkPackets()
		if err != nil {
			log.Err(fmt.Errorf("error getting level data chunk packets: %v", err))
			return
		}
		ctx.Send(eid, &messages.LevelData{Pks: pks})
	}
	msg.E.Teleport(*msg.Evt.SpawnX(), *msg.Evt.SpawnY(), *msg.Evt.SpawnZ())

	/*TODO: investigate if a server could support > 255 players seperated by multiple chunks of distance, despawning clients, and reissuing network ids*/
	if p, ok := msg.E.(*player.Player); ok {
		for _, e := range w.entities {
			ctx.Send(e.GetPid(), &packets.SpawnPlayerPacket{
				PlayerId:   msg.E.Id(),
				PlayerName: p.GetName(),
				X:          *msg.Evt.SpawnX(),
				Y:          *msg.Evt.SpawnY(),
				Z:          *msg.Evt.SpawnZ(),
				Yaw:        0,
				Pitch:      0,
			})

			ctx.Send(p.GetPid(), &packets.SpawnPlayerPacket{
				PlayerId:   e.Id(),
				PlayerName: p.GetName(),
				X:          e.Unsafe_X(), //TODO:
				Y:          e.Unsafe_Y(),
				Z:          e.Unsafe_Z(),
				Yaw:        0,
				Pitch:      0,
			})
		}
	}

}

func (w *World) GetNextEntityId() byte {
	return w.idManager.NextEntityId()
}

func (w *World) createFlatWorld() error {
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

func (w *World) moveEntity(ctx *actor.Context, e entity.Entity, pos *entity.GetPositionResponse) error {
	chunkX := int(pos.X) / ChunkWidth
	chunkZ := int(pos.Z) / ChunkDepth
	newChunkKey := [2]int{chunkX, chunkZ}

	newChunk, ok := w.Chunks[newChunkKey]
	if !ok {
		return fmt.Errorf("target chunk (%d, %d) not found", chunkX, chunkZ)
	}

	currentChunk, exists := w.entityMap[e]
	if exists && currentChunk == newChunk {
		return nil
	}

	ctx.Send(newChunk.pid, &messages.JoinChunk{E: e, Pos: pos})

	if exists {
		ctx.Send(currentChunk.pid, &messages.LeaveChunk{E: e})
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
