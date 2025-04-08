package player

import (
	"PowerClassic/entity"
	"PowerClassic/events"
	"PowerClassic/packets"
	"PowerClassic/session"
	"PowerClassic/world"
	"github.com/anthdm/hollywood/actor"
	"github.com/rs/zerolog/log"
)

type DefaultPlayerHandler struct {
	session.DefaultPacketHandler

	player *Player
}

func NewDefaultPlayerHandler(player *Player) *DefaultPlayerHandler {
	return &DefaultPlayerHandler{player: player}
}

func (h *DefaultPlayerHandler) HandleMovement(eng *actor.Engine, pkt *packets.PlayerTeleportPacket) {

	eng.Send(h.player.GetPid(), &PlayerRunnable{func(ctx *actor.Context, p *Player) {
		evt := events.NewPlayerMoveEvent(p.X(), p.Y(), p.Z(), &pkt.X, &pkt.Y, &pkt.Z)

		log.Debug().Msgf("player move received")

		h.player.GetEventBus().Fire("PlayerMove", evt)
		if evt.IsCancelled() {
			h.player.Teleport(ctx, evt.FromX(), evt.FromY(), evt.FromZ())
			return
		}
		p.SetPosition(pkt.X, pkt.Y, pkt.Z)
		ctx.Send(ctx.Parent(), &world.WorldRunnable{func(ctx *actor.Context, w *world.World) {
			w.BroadcastEntityRunnable(ctx, &entity.EntityRunnable{func(ctx *actor.Context, e entity.Entity) {
				if e == h.player {
					return
				}
				if ep, ok := e.(*Player); ok {
					ep.SendPosition(h.player)
				}
			}})
		}})
	}})

}

var _ session.PacketHandler = (*DefaultPlayerHandler)(nil)
