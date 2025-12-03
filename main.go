package main

import (
	"fmt"
	"gates/actors"
	"gates/actors/enemies"
	"gates/actors/weapons"
	game_events "gates/events"
	"gates/values"

	"github.com/Papiermond/eventbus"
	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/ui"
)

type GameState int

const (
	Running GameState = iota
	Preparing
	Waiting
	Stopped
)

var game_state GameState
var message_font *ui.Font
var rounds int

func main() {
	gomesengine.Init("RPG", int32(values.SCREEN_SIZE.X), int32(values.SCREEN_SIZE.Y))
	game_events.Init()

	settings()
	listeners()
	show_weapon_text()

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
		switch game_state {
		case Preparing:
			actors.PlayerLoadWeapon(weapons.Sword)
			sequence()
		case Waiting:
			message_font.UpdateColor(render.Transparent)
		}

		return nil
	})

	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_2, func(params ...any) error {
		switch game_state {
		case Preparing:
			actors.PlayerLoadWeapon(weapons.SpellFire)
			sequence()
		case Waiting:
			message_font.UpdateColor(render.Transparent)
		}

		return nil
	})

	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_3, func(params ...any) error {
		switch game_state {
		case Preparing:
			actors.PlayerLoadWeapon(weapons.Bow)
			sequence()
		case Waiting:
			message_font.UpdateColor(render.Transparent)
		}

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
			println(values.Yellow + "Trying to kill an enemy while game is not running. Current state:" + fmt.Sprint(game_state) + values.Reset)
			return
		}

		println("MAIN::" + e.(game_events.EnemyDeadEvent).Message)

		game_state = Waiting
		rounds += 1

		// Player Level Up
		actors.PlayerLevelUp()
		show_level_up_text()
	})
}

func sequence() {
	println("Sequence called")

	if game_state == Preparing {
		message_font.UpdateColor(render.Transparent)

		actors.LoadEnemy(enemies.Rat)
		actors.Enemy()
		actors.Player()

		game_state = Running
		println("Starting game...")
	}

	if game_state == Waiting {
		if rounds == 1 {
			actors.LoadEnemy(enemies.Skeleton)
		}

		if rounds == 2 {
			actors.LoadEnemy(enemies.Dragon)
		}

		game_state = Running
		game_events.Bus.Publish(game_events.GameRestartEvent{
			Message: "Restarting game!",
		})
	}
}

func show_weapon_text() {
	message := "Choose your weapon: 1 - Sword | 2 - Fire Spell | 3 - Bow"
	// println(message)

	specs := ui.FontSpecs{
		Name: "Pixelboy",
		Path: "assets/fonts/pixeboy-font/Pixeboy-z8XGD.ttf",
		Size: 32,
	}
	message_font = ui.NewFont(specs, values.SCREEN_SIZE)
	message_font.Init(message, render.Blue, math.Vector2{X: 0, Y: 0})
	message_font.AlignText(ui.MiddleCenter, math.Vector2{X: 0, Y: 0})
}

func show_level_up_text() {
	message := "LEVEL UP. Choose a skill to increase: 1 - STR | 2 - INT | 3 - SPD"
	message_font.UpdatePosition(math.Vector2{
		X: 0,
		Y: (values.SCREEN_SIZE.Y / 2) - 16,
	})
	message_font.UpdateText(message)
	message_font.UpdateColor(render.Blue)
	message_font.AlignText(ui.MiddleCenter, math.Vector2{X: 0, Y: 0})
}
