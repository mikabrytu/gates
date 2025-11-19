package main

import (
	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
)

const PLAYER_ATTACK_EVENT string = "PLAYER_ATTACK_EVENT"
const ENEMY_ATTACK_EVENT string = "ENEMY_ATTACK_EVENT"

var SCREEN_SIZE math.Vector2 = math.Vector2{X: 1600, Y: 960}

func main() {
	gomesengine.Init("RPG", int32(SCREEN_SIZE.X), int32(SCREEN_SIZE.Y))

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
	player()
	enemy()
}
