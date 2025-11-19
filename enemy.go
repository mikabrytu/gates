package main

import (
	"gates/systems"
	"time"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var enemy_render_attack bool = false

func enemy() {
	size := 230
	rect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - size,
		PosY:   32,
		Width:  (size + 5) * 2,
		Height: size * 2,
	}

	hp_max := 50
	hp_max_width := rect.Width
	hp_rect := rect
	hp_rect.PosY -= 24
	hp_rect.Height = 16

	attack_circle := utils.CircleSpecs{
		PosX:   SCREEN_SIZE.X / 2,
		PosY:   rect.PosY + rect.Height + 64,
		Radius: 64,
	}

	sprite := render.NewSprite(
		"Dragon",
		"assets/images/dragon.png",
		rect,
		render.Transparent,
	)

	health := systems.InitHealth(hp_max)

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			sprite.Init()

			go enemy_attack_task()

			events.Subscribe(PLAYER_ATTACK_EVENT, func(params ...any) error {
				damage := params[0].([]any)[0].([]any)[0].(int32)
				health.TakeDamage(int(damage))

				return nil
			})
		},
		Update: func() {
			hp_rect.Width = (hp_max_width * health.GetCurrent()) / hp_max
		},
		Render: func() {
			render.DrawRect(hp_rect, render.Red)

			if enemy_render_attack {
				render.DrawCircle(attack_circle, render.Red)
			}
		},
	})
}

func enemy_attack_task() {
	for range time.Tick(time.Millisecond * 5000) {
		println("Enemy attack on a set interval")

		enemy_render_attack = true
		time.AfterFunc(time.Millisecond*800, func() {
			enemy_render_attack = false
		})

		events.Emit(ENEMY_ATTACK_EVENT, int32(10))
	}
}
