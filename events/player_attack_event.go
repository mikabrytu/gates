package game_events

import "github.com/Papiermond/eventbus"

type PlayerAttackEvent struct {
	Damage  int
	Message string
}

func (e PlayerAttackEvent) GetType() eventbus.EventType {
	return PLAYER_ATTACK_EVENT
}
