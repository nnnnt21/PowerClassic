package entity

import (
	"PowerClassic/network"
	"github.com/anthdm/hollywood/actor"
)

type Pos struct {
	X, Y, Z float32
}
type PositionSnapshot *Pos

type ChunkProvider interface {
}

type EntityRunnable struct {
	Run func(ctx *actor.Context, e Entity)
}

type Entity interface {
	Id() byte
	GetName() string
	GetPosition() PositionSnapshot
	Teleport(ctx *actor.Context, x, y, z float32)
	SetPid(pid *actor.PID)
	GetPid() *actor.PID
	actor.Receiver
}

type SessionedEntity interface {
	SendPacket(pkt network.SerializablePacket)
	Disconnect(reason string)
}
