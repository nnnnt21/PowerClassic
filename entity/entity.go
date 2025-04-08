package entity

import (
	"PowerClassic/network"
	"github.com/anthdm/hollywood/actor"
)

type ChunkProvider interface {
}

type EntityRunnable struct {
	Run func(ctx *actor.Context, e Entity)
}

type Entity interface {
	Id() byte
	GetName() string
	X() float32
	Y() float32
	Z() float32
	Teleport(ctx *actor.Context, x, y, z float32)
	SetPid(pid *actor.PID)
	GetPid() *actor.PID
	actor.Receiver
}

type SessionedEntity interface {
	SendPacket(pkt network.SerializablePacket)
	Disconnect(reason string)
}
