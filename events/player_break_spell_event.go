package game_events

import "github.com/Papiermond/eventbus"

type PlayerBreakSpellEvent struct {
	Message string
}

func (e PlayerBreakSpellEvent) GetType() eventbus.EventType {
	return PLAYER_BREAK_SPELL_EVENT
}
