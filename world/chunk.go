package world

import (
	"PowerClassic/entity"
	"fmt"
	"github.com/anthdm/hollywood/actor"
)

const (
	ChunkWidth = 16
	ChunkDepth = 16
)

// virtual chunk of the world, classic client doesn't seem to know about a chunk but we can still use this to seperate processing, or maybe there is an extension for chunks? We will see later
type Chunk struct {
	pid            *actor.PID
	ChunkX, ChunkZ int
	Height         int
	Blocks         []byte
	entities       map[entity.Entity]struct{}
}

func NewChunk(chunkX, chunkZ, height int) *Chunk {
	total := ChunkWidth * height * ChunkDepth
	return &Chunk{
		ChunkX:   chunkX,
		ChunkZ:   chunkZ,
		Height:   height,
		Blocks:   make([]byte, total),
		entities: make(map[entity.Entity]struct{}),
	}
}

func (c *Chunk) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *ChunkRunnable:
		msg.Run(ctx, c)
	}
}

func (c *Chunk) JoinChunk(e entity.Entity, ctx *actor.Context) {
	c.entities[e] = struct{}{}
}

func (c *Chunk) LeaveChunk(e entity.Entity, ctx *actor.Context) {
	delete(c.entities, e)
}

func (c *Chunk) index(localX, y, localZ int) (int, error) {
	if localX < 0 || localX >= ChunkWidth || y < 0 || y >= c.Height || localZ < 0 || localZ >= ChunkDepth {
		return 0, fmt.Errorf("coordinates out of range in chunk (%d, %d): (%d, %d, %d)", c.ChunkX, c.ChunkZ, localX, y, localZ)
	}
	return localX + localZ*ChunkWidth + y*ChunkWidth*ChunkDepth, nil
}

func (c *Chunk) GetBlock(localX, y, localZ int) (byte, error) {
	idx, err := c.index(localX, y, localZ)
	if err != nil {
		return 0, err
	}
	return c.Blocks[idx], nil
}

func (c *Chunk) SetBlock(localX, y, localZ int, block byte) error {
	idx, err := c.index(localX, y, localZ)
	if err != nil {
		return err
	}
	c.Blocks[idx] = block
	return nil
}

var _ actor.Receiver = (*Chunk)(nil)
