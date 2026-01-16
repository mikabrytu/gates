package game_events

import "github.com/Papiermond/eventbus"

type EnemyBreakParalysisEvent struct{}

func (e EnemyBreakParalysisEvent) GetType() eventbus.EventType {
	return ENEMY_BREAK_PARALYSIS_EVENT
}
