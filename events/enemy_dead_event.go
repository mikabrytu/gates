package game_events

import "github.com/Papiermond/eventbus"

type EnemyDeadEvent struct {
	Message string
}

func (e EnemyDeadEvent) GetType() eventbus.EventType {
	return eventbus.EventType(ENEMY_DEAD_EVENT)
}
