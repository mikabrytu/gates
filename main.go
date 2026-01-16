package main

import (
	game_events "gates/events"
	"gates/scenes"
	"gates/values"

	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
)

func main() {
	game()
}

func game() {
	gomesengine.Init("RPG", int32(values.SCREEN_SIZE.X), int32(values.SCREEN_SIZE.Y))
	game_events.Init()

	settings()
	scenes.RunMap()
	//scenes.RunCombat()

	gomesengine.Run()
}

func settings() {
	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_ESCAPE, func(data any) {
		lifecycle.Kill()
	})
}
