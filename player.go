package main

import (
	"time"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

func player() {
	hpRect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - 256,
		PosY:   SCREEN_SIZE.Y - 40,
		Width:  512,
		Height: 24,
	}
	weaponRect := utils.RectSpecs{
		PosX:   SCREEN_SIZE.X - 256 - 64,
		PosY:   SCREEN_SIZE.Y - 512,
		Width:  256,
		Height: 512,
	}

	sprite := render.NewSprite(
		"Player Weapon",
		"assets/images/dagger.jpg",
		weaponRect,
		render.Transparent,
	)

	can_attack := true

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			sprite.Init()

			events.Subscribe(events.INPUT_MOUSE_CLICK_DOWN, func(params ...any) error {
				if !can_attack {
					return nil
				}

				println("Player attack")
				can_attack = false
				events.Emit(PLAYER_ATTACK_EVENT)

				tempRect := weaponRect
				tempRect.PosX = (SCREEN_SIZE.X / 2) - 128
				sprite.UpdateRect(tempRect)

				// Recovery delay
				time.AfterFunc(time.Millisecond*500, func() {
					can_attack = true
					sprite.UpdateRect(weaponRect)
				})

				return nil
			})
		},
		Render: func() {
			render.DrawRect(hpRect, render.Green)
		},
	})
}
