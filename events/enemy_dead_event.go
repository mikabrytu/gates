package game_events

import "github.com/Papiermond/eventbus"

type EnemyDeadEvent struct {
	Message string
}

func (e EnemyDeadEvent) GetType() eventbus.EventType {
	return ENEMY_DEAD_EVENT
}
