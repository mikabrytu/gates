package events

import "github.com/Papiermond/eventbus"

type PlayerCanLevelUpEvent struct{}

func (e PlayerCanLevelUpEvent) GetType() eventbus.EventType {
	return PLAYER_CAN_LEVEL_UP_EVENT
}
