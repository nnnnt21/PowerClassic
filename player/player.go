package player

import (
	"PowerClassic/entity"
	"PowerClassic/event"
	"PowerClassic/network"
	"PowerClassic/packets"
	"PowerClassic/session"
	"PowerClassic/world"
	"github.com/anthdm/hollywood/actor"
	"github.com/rs/zerolog/log"
)

type PlayerRunnable struct {
	Run func(ctx *actor.Context, p *Player)
}

type Player struct {
	pid *actor.PID

	id uint8
	x  float32
	y  float32
	z  float32

	yaw   float32
	pitch float32

	s *session.Session

	eng *actor.Engine

	initialSpawnX float32
	initialSpawnY float32
	initialSpawnZ float32
}

func NewPlayer(s *session.Session, id byte, x, y, z float32, eng *actor.Engine) *Player {
	p := &Player{
		id: id,
		x:  x,
		y:  y,
		z:  z,

		eng: eng,

		s: s,

		initialSpawnX: x,
		initialSpawnY: y,
		initialSpawnZ: z,
	}

	s.SetHandler(NewDefaultPlayerHandler(p))

	return p
}

func (p *Player) SendPacket(pkt network.SerializablePacket) {
	p.s.AddOutgoingPacket(pkt)
}

func (p *Player) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *PlayerRunnable:
		msg.Run(ctx, p)
	case *entity.EntityRunnable:
		msg.Run(ctx, p)
	default:
		log.Debug().Msgf("Received message of type %T\n", msg)
	}
}

func (p *Player) GetName() string {
	return p.s.Data().Username
}

func (p *Player) SetPid(pid *actor.PID) {
	p.pid = pid
}

func (p *Player) GetPid() *actor.PID {
	return p.pid
}

func (p *Player) GetEventBus() *event.EventBus {
	return p.s.GetEventBus()
}

func (p *Player) sendLevelData(pks []*packets.LevelDataChunkPacket) {
	p.s.AddOutgoingPacket(&packets.LevelInitializePacket{})
	for _, pk := range pks {
		p.s.AddOutgoingPacket(pk)
	}

	p.s.AddOutgoingPacket(&packets.LevelFinalizePacket{
		X: 1024,
		Y: 64,
		Z: 1024,
	})
}

func (p *Player) SetPosition(x, y, z float32) {
	p.x = x
	p.y = y
	p.z = z
}
func (p *Player) SendPosition(e entity.Entity) {
	p.s.AddOutgoingPacket(&packets.PlayerTeleportPacket{
		PlayerId: e.Id(),
		X:        e.X(),
		Y:        e.Y(),
		Z:        e.Z(),
		Yaw:      0,
		Pitch:    0,
	})
}
func (p *Player) Teleport(ctx *actor.Context, x, y, z float32) {
	p.x = x
	p.y = y
	p.z = z

	p.s.AddOutgoingPacket(&packets.PlayerTeleportPacket{
		PlayerId: 255,
		X:        x,
		Y:        y,
		Z:        z,
		Yaw:      0,
		Pitch:    0,
	})

	ctx.Send(ctx.Parent(), &world.WorldRunnable{func(ctx *actor.Context, w *world.World) {
		w.BroadcastEntityRunnable(ctx, &entity.EntityRunnable{func(ctx *actor.Context, e entity.Entity) {
			if ep, ok := e.(*Player); ok {
				ep.SendPosition(p)
			}
		}})
	}})
}

func (p *Player) Disconnect(reason string) {
	p.s.AddOutgoingPacket(&packets.DisconnectPacket{
		Reason: reason,
	})
}

func (p *Player) Id() byte {
	return p.id
}

func (p *Player) X() float32 {
	return p.x
}

func (p *Player) Y() float32 {
	return p.y
}

func (p *Player) Z() float32 {
	return p.z
}

var _ entity.Entity = (*Player)(nil)
var _ entity.SessionedEntity = (*Player)(nil)
