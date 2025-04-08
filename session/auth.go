package session

import (
	"PowerClassic/events"
	"PowerClassic/packets"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	DefaultPacketHandler

	session *Session
}

func NewAuthHandler(s *Session) *AuthHandler {
	return &AuthHandler{session: s}
}

func (h *AuthHandler) HandlePlayerIdentificationPacket(pk *packets.PlayerIdentificationPacket) {
	h.session.data.Username = pk.Username
	h.session.data.ProtocolVersion = pk.ProtocolVersion
	h.session.data.VerificationKey = pk.VerificationKey

	evt := &events.PlayerIdentificationEvent{
		Packet: pk,
	}

	log.Debug().Msgf("Player connecting with the username %s protocol version: %d", h.session.data.Username, h.session.data.ProtocolVersion)

	h.session.eventBus.Fire("PlayerIdentification", evt)
	if evt.IsCancelled() {
		h.session.AddOutgoingPacket(&packets.DisconnectPacket{Reason: evt.DisconnectReason()})
		return
	}

	if evt.SpawnX() != nil && evt.SpawnY() != nil && evt.SpawnZ() != nil && evt.SpawnWorld() != nil {
		h.session.PromoteSession(evt)
		return
	}
	h.session.AddOutgoingPacket(&packets.DisconnectPacket{Reason: "Server misconfigured, must handle spawn event"})

}

var _ PacketHandler = (*AuthHandler)(nil)
