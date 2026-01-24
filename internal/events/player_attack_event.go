package events

import (
	"gates/pkg/spell"

	"github.com/Papiermond/eventbus"
)

type PlayerAttackEvent struct {
	Damage  int
	Effect  spell.Effect
	Message string
}

func (e PlayerAttackEvent) GetType() eventbus.EventType {
	return PLAYER_ATTACK_EVENT
}
