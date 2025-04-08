package messages

import (
	"PowerClassic/entity"
	"PowerClassic/events"
	"PowerClassic/packets"
)

type Disconnect struct{ Reason string }
type AddEntity struct {
	E   entity.Entity
	Evt *events.PlayerIdentificationEvent
}
type Teleport struct {
	X, Y, Z, Pitch, Yaw float32
}
type LevelData struct {
	Pks []*packets.LevelDataChunkPacket
}

type MovePlayer struct {
	ID                  uint8
	X, Y, Z, Pitch, Yaw float32
}
