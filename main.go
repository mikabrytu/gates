package main

import (
	"gates/actors"
	"gates/actors/enemies"
	"gates/actors/weapons"
	"gates/values"

	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
)

type GameState int

const (
	Running GameState = iota
	Preparing
	Waiting
	Stopped
)

var game_state GameState
var rounds int

func main() {
	gomesengine.Init("RPG", int32(values.SCREEN_SIZE.X), int32(values.SCREEN_SIZE.Y))

	settings()
	listeners()

	println(values.Blue + "Choose your weapon: 1 - Sword | 2 - Fire Spell | 3 - Bow" + values.Reset)

	game_state = Preparing
	gomesengine.Run()
}

func settings() {
	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_ESCAPE, func(params ...any) error {
		lifecycle.Kill()
		return nil
	})
}

func listeners() {
	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_1, func(params ...any) error {
		if game_state != Preparing {
			return nil
		}

		actors.PlayerLoadWeapon(weapons.Sword)
		sequence()
		return nil
	})

	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_2, func(params ...any) error {
		if game_state != Preparing {
			return nil
		}

		actors.PlayerLoadWeapon(weapons.SpellFire)
		sequence()
		return nil
	})

	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_3, func(params ...any) error {
		if game_state != Preparing {
			return nil
		}

		actors.PlayerLoadWeapon(weapons.Bow)
		sequence()
		return nil
	})

	events.Subscribe(values.GAME_OVER_EVENT, func(params ...any) error {
		game_state = Stopped
		return nil
	})

	events.Subscribe(values.ENEMY_DEAD_EVENT, func(params ...any) error {
		game_state = Waiting
		rounds += 1

		// Player Level Up
		if rounds == 2 || rounds == 4 || rounds == 7 {
			actors.PlayerLevelUp()
		}

		return nil
	})

	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_SPACE, func(params ...any) error {
		sequence()
		return nil
	})
}

func sequence() {
	println("Sequence called")

	if game_state == Preparing {
		actors.Player()
		actors.LoadEnemy(enemies.Rat)
		actors.Enemy()

		game_state = Running
		println("Starting game...")
	}

	if game_state == Waiting {
		if rounds >= 3 && rounds < 7 {
			actors.LoadEnemy(enemies.Skeleton)
		}

		if rounds >= 7 {
			actors.LoadEnemy(enemies.Dragon)
		}

		println("Restarting game!")
		game_state = Running
		events.Emit(values.GAME_RESTART_EVENT)
	}
}
