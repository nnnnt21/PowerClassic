package entity

import (
	"github.com/anthdm/hollywood/actor"
)

type ChunkProvider interface {
}

type Entity interface {
	Id() byte
	Unsafe_X() float32
	Unsafe_Y() float32
	Unsafe_Z() float32
	GetPosition(ctx *actor.Context) (*GetPositionResponse, error)
	GetPositionEng(eng *actor.Engine) (*GetPositionResponse, error)
	Teleport(x, y, z float32)
	SetPid(pid *actor.PID)
	GetPid() *actor.PID
	actor.Receiver
}

type SessionedEntity interface {
	Disconnect(reason string)
}

// TODO: move this, temp to avoid circular dependency
type GetPosition struct{}
type GetPositionResponse struct {
	X, Y, Z, Pitch, Yaw float32
}
