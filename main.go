package main

import (
	"gates/actors"
	"gates/actors/enemies"
	"gates/actors/weapons"
	game_events "gates/events"
	"gates/values"

	"github.com/Papiermond/eventbus"
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
	game_events.Init()

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
	// Engine events
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

	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_SPACE, func(params ...any) error {
		sequence()
		return nil
	})

	// Game events
	game_events.Bus.Subscribe(game_events.GAME_OVER_EVENT, func(e eventbus.Event) {
		game_state = Stopped
	})

	game_events.Bus.Subscribe(game_events.ENEMY_DEAD_EVENT, func(e eventbus.Event) {
		if game_state != Running {
			println(values.Yellow + "Trying to kill an enemy while game is not running. Current state:" + string(game_state) + values.Reset)
			return
		}

		println("MAIN::" + e.(game_events.EnemyDeadEvent).Message)

		game_state = Waiting
		rounds += 1

		// Player Level Up
		if rounds == 2 || rounds == 4 || rounds == 7 {
			actors.PlayerLevelUp()
		}
	})
}

func sequence() {
	println("Sequence called")

	if game_state == Preparing {
		actors.LoadEnemy(enemies.Rat)
		actors.Enemy()
		actors.Player()

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

		game_state = Running
		game_events.Bus.Publish(game_events.GameRestartEvent{
			Message: "Restarting game!",
		})
	}
}
