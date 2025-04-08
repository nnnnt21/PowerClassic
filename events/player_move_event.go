package events

import "PowerClassic/event"

type PlayerMoveEvent struct {
	event.BaseEvent

	fromX float32
	toX   *float32

	fromY float32
	toY   *float32

	fromZ float32
	toZ   *float32
}

func (p PlayerMoveEvent) FromX() float32 {
	return p.fromX
}

func (p PlayerMoveEvent) FromY() float32 {
	return p.fromY
}

func (p PlayerMoveEvent) FromZ() float32 {
	return p.fromZ
}

func NewPlayerMoveEvent(fromx, fromy, fromz float32, tox, toy, toz *float32) *PlayerMoveEvent {
	return &PlayerMoveEvent{
		fromX: fromx,
		toX:   tox,
		fromY: fromy,
		toY:   toy,
		fromZ: fromz,
		toZ:   toz,
	}
}
