package events

import "github.com/Papiermond/eventbus"

type GameRestartEvent struct {
	Message string
}

func (e GameRestartEvent) GetType() eventbus.EventType {
	return GAME_RESTART_EVENT
}
