package main

import (
	"time"

	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type generic func()

const PLAYER_ATTACK_EVENT string = "PLAYER_ATTACK_EVENT"

var SCREEN_SIZE math.Vector2 = math.Vector2{X: 800, Y: 600}

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

func player() {
	size := 32
	rect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - (size / 2),
		PosY:   SCREEN_SIZE.Y - size - 16,
		Width:  size,
		Height: size,
	}
	can_attack := true

	register(rect, render.Green, func() {
		events.Subscribe(events.INPUT_MOUSE_CLICK_DOWN, func(params ...any) error {
			if !can_attack {
				return nil
			}

			println("Player attack")
			can_attack = false
			events.Emit(PLAYER_ATTACK_EVENT)

			// Recovery delay
			time.AfterFunc(time.Millisecond*2000, func() {
				can_attack = true
			})

			return nil
		})
	})
}

func enemy() {
	size := 256
	rect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - (size / 2),
		PosY:   32,
		Width:  size,
		Height: size,
	}

	register(rect, render.Blue, func() {
		events.Subscribe(PLAYER_ATTACK_EVENT, func(params ...any) error {
			println("Enemy damaged")

			return nil
		})
	})
}

func register(rect utils.RectSpecs, color render.Color, start generic) *lifecycle.GameObject {
	hpRect := rect
	hpRect.PosY -= 24
	hpRect.Height = 16

	return lifecycle.Register(&lifecycle.GameObject{
		Start: start,
		Render: func() {
			// Draw HP
			render.DrawRect(hpRect, render.Red)
			// Draw Sprite
			render.DrawRect(rect, color)
		},
	})
}
