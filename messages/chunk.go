package messages

import "PowerClassic/entity"

type JoinChunk struct {
	E   entity.Entity
	Pos *entity.GetPositionResponse
}
type LeaveChunk struct {
	E entity.Entity
}
