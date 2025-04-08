package player

import (
	"PowerClassic/events"
	"PowerClassic/messages"
	"PowerClassic/packets"
	"PowerClassic/session"
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

	pos, err := h.player.GetPositionEng(eng)

	if err != nil {
		panic(err)
	}

	evt := events.NewPlayerMoveEvent(pos.X, pos.Y, pos.Z, &pkt.X, &pkt.Y, &pkt.Z)

	log.Debug().Msgf("player move received")

	h.player.GetEventBus().Fire("PlayerMove", evt)
	if evt.IsCancelled() {
		h.player.Teleport(evt.FromX(), evt.FromY(), evt.FromZ())
		return
	}
	eng.Send(h.player.pid, &messages.MovePlayer{
		ID:    h.player.id,
		X:     pkt.X,
		Y:     pkt.Y,
		Z:     pkt.Z,
		Pitch: 0,
		Yaw:   0,
	})
}

var _ session.PacketHandler = (*DefaultPlayerHandler)(nil)
