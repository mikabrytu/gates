package main

import (
	"gates/actors"
	"gates/values"

	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
)

func main() {
	gomesengine.Init("RPG", int32(values.SCREEN_SIZE.X), int32(values.SCREEN_SIZE.Y))

	settings()
	scene()

	gomesengine.Run()
}

func settings() {
	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_ESCAPE, func(params ...any) error {
		lifecycle.Kill()
		return nil
	})
}

func scene() {
	actors.Player()
	actors.Enemy()
}
