package game_events

import "github.com/Papiermond/eventbus"

type EnemyAttackEvent struct {
	Damage  int
	Message string
}

func (e EnemyAttackEvent) GetType() eventbus.EventType {
	return ENEMY_ATTACK_EVENT
}
