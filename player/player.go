package player

import (
	"PowerClassic/entity"
	"PowerClassic/event"
	"PowerClassic/messages"
	"PowerClassic/packets"
	"PowerClassic/session"
	"errors"
	"github.com/anthdm/hollywood/actor"
	"github.com/rs/zerolog/log"
	"time"
)

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

// TODO: need a better way to do this
func (p *Player) GetPosition(ctx *actor.Context) (*entity.GetPositionResponse, error) {
	resp := ctx.Request(p.pid, &entity.GetPosition{}, time.Second*10)

	res, err := resp.Result()
	if err != nil {
		return nil, err
	}
	if res, ok := res.(*entity.GetPositionResponse); ok {
		return res, nil
	}
	return nil, errors.New("invalid response type")
}

// TODO: need a better way to do this
func (p *Player) GetPositionEng(eng *actor.Engine) (*entity.GetPositionResponse, error) {
	resp := eng.Request(p.pid, &entity.GetPosition{}, time.Second*10)

	res, err := resp.Result()
	if err != nil {
		return nil, err
	}
	if res, ok := res.(*entity.GetPositionResponse); ok {
		return res, nil
	}
	return nil, errors.New("invalid response type")
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

func (p *Player) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.LevelData:
		p.sendLevelData(msg.Pks)
	case *messages.Teleport:
		p.teleport(msg)
	case *messages.Disconnect:
		p.disconnect(msg)
	case *messages.MovePlayer:
		p.move(msg)
		ctx.Forward(ctx.Parent())
	case *entity.GetPosition:
		log.Info().Msgf("Received GetPosition %+v", msg)
		ctx.Respond(&entity.GetPositionResponse{
			X:     p.x,
			Y:     p.y,
			Z:     p.z,
			Pitch: 0,
			Yaw:   0,
		})
	case *packets.PlayerTeleportPacket:
		if msg.PlayerId == p.id {
			log.Info().Msgf("Received self PlayerTeleportPacket %+v", msg)
			return
		}
		p.s.AddOutgoingPacket(msg)
	case *packets.SpawnPlayerPacket:
		if msg.PlayerId == p.id {
			cop := *msg
			cop.PlayerId = 255
			p.s.AddOutgoingPacket(&cop)
			return
		}
		p.s.AddOutgoingPacket(msg)
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

func (p *Player) move(msg *messages.MovePlayer) {
	p.x = msg.X
	p.y = msg.Y
	p.z = msg.Z
}
func (p *Player) Teleport(x, y, z float32) {
	p.eng.Send(p.pid, &messages.Teleport{
		X:     x,
		Y:     y,
		Z:     z,
		Pitch: 0,
		Yaw:   0,
	})
}
func (p *Player) teleport(msg *messages.Teleport) {
	p.x = msg.X
	p.y = msg.Y
	p.z = msg.Z

	p.s.AddOutgoingPacket(&packets.PlayerTeleportPacket{
		PlayerId: 255,
		X:        msg.X,
		Y:        msg.Y,
		Z:        msg.Z,
		Yaw:      0,
		Pitch:    0,
	})
}

func (p *Player) Disconnect(reason string) {
	p.eng.Send(p.pid, &messages.Disconnect{Reason: reason})
}
func (p *Player) disconnect(msg *messages.Disconnect) {
	p.s.AddOutgoingPacket(&packets.DisconnectPacket{
		Reason: msg.Reason,
	})
}

func (p *Player) Id() byte {
	return p.id
}

func (p *Player) Unsafe_X() float32 {
	return p.x
}

func (p *Player) Unsafe_Y() float32 {
	return p.y
}

func (p *Player) Unsafe_Z() float32 {
	return p.z
}

var _ entity.Entity = (*Player)(nil)
var _ entity.SessionedEntity = (*Player)(nil)
