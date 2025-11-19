package main

import (
	"time"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

func enemy() {
	size := 230
	rect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - size,
		PosY:   32,
		Width:  (size + 5) * 2,
		Height: size * 2,
	}

	hpRect := rect
	hpRect.PosY -= 24
	hpRect.Height = 16

	sprite := render.NewSprite(
		"Dragon",
		"assets/images/dragon.png",
		rect,
		render.Transparent,
	)

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			sprite.Init()

			go enemy_attack_task()

			events.Subscribe(PLAYER_ATTACK_EVENT, func(params ...any) error {
				println("Enemy damaged")
				return nil
			})
		},
		Render: func() {
			render.DrawRect(hpRect, render.Red)
		},
	})
}

func enemy_attack_task() {
	for range time.Tick(time.Millisecond * 5000) {
		println("Enemy attack on a set interval")
	}
}
