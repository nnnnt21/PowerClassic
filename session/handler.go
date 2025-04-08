package session

import (
	"PowerClassic/packets"
	"github.com/anthdm/hollywood/actor"
)

type PacketHandler interface {
	HandlePlayerIdentificationPacket(pkt *packets.PlayerIdentificationPacket)
	HandleMovement(eng *actor.Engine, pkt *packets.PlayerTeleportPacket)
}

type DefaultPacketHandler struct{}

func (d DefaultPacketHandler) HandlePlayerIdentificationPacket(pkt *packets.PlayerIdentificationPacket) {
}

func (d DefaultPacketHandler) HandleMovement(eng *actor.Engine, pkt *packets.PlayerTeleportPacket) {
}
