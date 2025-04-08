package events

import (
	"PowerClassic/entity"
	"PowerClassic/event"
	"PowerClassic/packets"
)

type PlayerIdentificationEvent struct {
	event.BaseEvent
	Packet *packets.PlayerIdentificationPacket

	spawnX *float32
	spawnY *float32
	spawnZ *float32

	spawnWorld entity.ChunkProvider

	disconnectReason string
}

func (evt *PlayerIdentificationEvent) SpawnWorld() entity.ChunkProvider {
	return evt.spawnWorld
}
func (evt *PlayerIdentificationEvent) SetSpawnWorld(w entity.ChunkProvider) {
	evt.spawnWorld = w
}
func (evt *PlayerIdentificationEvent) SpawnX() *float32 {
	return evt.spawnX
}
func (evt *PlayerIdentificationEvent) SpawnY() *float32 {
	return evt.spawnY
}
func (evt *PlayerIdentificationEvent) SpawnZ() *float32 {
	return evt.spawnZ
}
func (evt *PlayerIdentificationEvent) DisconnectReason() string {
	return evt.disconnectReason
}
func (evt *PlayerIdentificationEvent) SetSpawnX(spawnX *float32) {
	evt.spawnX = spawnX
}
func (evt *PlayerIdentificationEvent) SetSpawnY(spawnY *float32) {
	evt.spawnY = spawnY
}
func (evt *PlayerIdentificationEvent) SetSpawnZ(spawnZ *float32) {
	evt.spawnZ = spawnZ
}
func (evt *PlayerIdentificationEvent) Disconnect(reason string) {
	evt.Cancel()
	evt.disconnectReason = reason
}
