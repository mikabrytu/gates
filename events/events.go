package game_events

import "github.com/Papiermond/eventbus"

var Bus eventbus.EventBus

func Init() {
	Bus = eventbus.New()
}
