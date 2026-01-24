package events

import "github.com/Papiermond/eventbus"

type GameOverEvent struct {
	Message string
}

func (e GameOverEvent) GetType() eventbus.EventType {
	return GAME_OVER_EVENT
}
