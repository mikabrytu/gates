package main

import (
	game_events "gates/events"
	"gates/scenes"
	"gates/values"

	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
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
	render.SetBackgroundColor(render.Color{R: 25, G: 20, B: 43, A: 255})
	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_ESCAPE, func(data any) {
		lifecycle.Kill()
	})
}
