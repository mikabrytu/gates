package game_events

import "github.com/Papiermond/eventbus"

type GameRestartEvent struct {
	Message string
}

func (e GameRestartEvent) GetType() eventbus.EventType {
	return eventbus.EventType(GAME_RESTART_EVENT)
}
