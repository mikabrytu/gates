package game_events

import (
	"gates/spells"

	"github.com/Papiermond/eventbus"
)

type PlayerAttackEvent struct {
	Damage  int
	Effect  spells.Effect
	Message string
}

func (e PlayerAttackEvent) GetType() eventbus.EventType {
	return PLAYER_ATTACK_EVENT
}
