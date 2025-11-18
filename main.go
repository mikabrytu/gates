package main

import (
	gomesengine "github.com/mikabrytu/gomes-engine"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var SCREEN_SIZE math.Vector2 = math.Vector2{X: 800, Y: 600}

func main() {
	gomesengine.Init("RPG", int32(SCREEN_SIZE.X), int32(SCREEN_SIZE.Y))

	settings()
	drawScene()

	gomesengine.Run()
}

func settings() {
	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_ESCAPE, func(params ...any) error {
		lifecycle.Kill()
		return nil
	})
}

func drawScene() {
	playerSize := 32
	playerRect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - (playerSize / 2),
		PosY:   SCREEN_SIZE.Y - playerSize - 16,
		Width:  playerSize,
		Height: playerSize,
	}

	enemySize := 256
	enemyRect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - (enemySize / 2),
		PosY:   32,
		Width:  enemySize,
		Height: enemySize,
	}

	register(enemyRect, render.Blue)
	register(playerRect, render.Green)
}

func register(rect utils.RectSpecs, color render.Color) {
	hpRect := rect
	hpRect.PosY -= 24
	hpRect.Height = 16

	lifecycle.Register(&lifecycle.GameObject{
		Render: func() {
			// Draw HP
			render.DrawRect(hpRect, render.Red)
			// Draw Sprite
			render.DrawRect(rect, color)
		},
	})
}
