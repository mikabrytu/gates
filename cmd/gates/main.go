package main

import (
	"gates/config"
	game_events "gates/internal/events"
	"gates/internal/game"

	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
)

func main() {
	gomesengine.Init("Gates", int32(config.SCREEN_SIZE.X), int32(config.SCREEN_SIZE.Y))
	game_events.Init()

	settings()
	game.Init()

	gomesengine.Run()
}

func settings() {
	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_ESCAPE, func(data any) {
		lifecycle.Kill()
	})
}
