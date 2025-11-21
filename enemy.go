package main

import (
	"gates/systems"
	"time"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var enemy_go *lifecycle.GameObject
var enemy_sprite *render.Sprite
var enemy_health *systems.Health
var enemy_hp_rect utils.RectSpecs
var enemt_attack_circle utils.CircleSpecs
var enemy_hp_max int
var enemy_hp_max_width int
var enemy_render_attack bool = false
var enemy_attack_done = make(chan bool)

func enemy() {
	enemy_init()
	enemy_respawn()

	enemy_go = lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			enemy_sprite.Init()

			events.Subscribe(PLAYER_ATTACK_EVENT, func(params ...any) error {
				damage := params[0].([]any)[0].([]any)[0].(int32)
				enemy_health.TakeDamage(int(damage))

				if enemy_health.GetCurrent() <= 0 {
					go func() {
						enemy_attack_done <- true
					}()

					enemy_sprite.ClearSprite()
					lifecycle.Disable(enemy_go)
				}

				return nil
			})
		},
		Update: func() {
			enemy_hp_rect.Width = (enemy_hp_max_width * enemy_health.GetCurrent()) / enemy_hp_max
		},
		Render: func() {
			render.DrawRect(enemy_hp_rect, render.Red)

			if enemy_render_attack {
				render.DrawCircle(enemt_attack_circle, render.Red)
			}
		},
	})
}

func enemy_init() {
	size := 230
	rect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - size,
		PosY:   32,
		Width:  (size + 5) * 2,
		Height: size * 2,
	}

	enemy_hp_max = 50
	enemy_hp_max_width = rect.Width
	enemy_hp_rect = rect
	enemy_hp_rect.PosY -= 24
	enemy_hp_rect.Height = 16

	enemt_attack_circle = utils.CircleSpecs{
		PosX:   SCREEN_SIZE.X / 2,
		PosY:   rect.PosY + rect.Height + 64,
		Radius: 64,
	}

	if enemy_sprite == nil {
		enemy_sprite = render.NewSprite(
			"Dragon",
			"assets/images/dragon.png",
			rect,
			render.Transparent,
		)
	} else {
		enemy_sprite.Init()
	}

	if enemy_health == nil {
		enemy_health = systems.InitHealth(enemy_hp_max)
	} else {
		enemy_health.Reset()
	}

	go enemy_attack_task()
}

func enemy_respawn() {
	events.Subscribe(events.INPUT_KEYBOARD_PRESSED_SPACE, func(params ...any) error {
		if enemy_go != nil {
			go func() {
				println("channel value should be false by now")
				enemy_attack_done = make(chan bool)
			}()

			enemy_init()
			lifecycle.Enable(enemy_go)
		}

		return nil
	})
}

func enemy_attack_task() {
	ticker := time.NewTicker(time.Millisecond * 3000)

	for {
		select {
		case <-enemy_attack_done:
			println("Stopping enemy attack")
			ticker.Stop()
			return
		case <-ticker.C:
			println("Enemy attack on a set interval")

			enemy_render_attack = true
			time.AfterFunc(time.Millisecond*800, func() {
				enemy_render_attack = false
			})

			events.Emit(ENEMY_ATTACK_EVENT, int32(10))
		}
	}
}
