package main

import (
	"gates/systems"
	"time"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var max_hp int = 100
var player_can_attack bool = true
var player_weapon_rect utils.RectSpecs = utils.RectSpecs{
	PosX:   SCREEN_SIZE.X - 256 - 64,
	PosY:   SCREEN_SIZE.Y - 512,
	Width:  256,
	Height: 512,
}
var player_sprite render.Sprite = *render.NewSprite(
	"Player Weapon",
	"assets/images/dagger.jpg",
	player_weapon_rect,
	render.Transparent,
)
var player_health *systems.Health = systems.InitHealth(max_hp)

func player() {

	max_hp_width := 512
	hp_rect := utils.RectSpecs{
		PosX:   (SCREEN_SIZE.X / 2) - (max_hp_width / 2),
		PosY:   SCREEN_SIZE.Y - 40,
		Width:  max_hp_width,
		Height: 24,
	}

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			player_sprite.Init()

			events.Subscribe(events.INPUT_MOUSE_CLICK_DOWN, func(params ...any) error {
				player_click_listener()
				return nil
			})

			events.Subscribe(ENEMY_ATTACK_EVENT, func(params ...any) error {
				damage := params[0].([]any)[0].([]any)[0].(int32)
				player_take_damage_listener(int(damage))

				return nil
			})
		},
		Update: func() {
			hp_rect.Width = (max_hp_width * player_health.GetCurrent()) / max_hp
		},
		Render: func() {
			render.DrawRect(hp_rect, render.Green)
		},
	})
}

func player_click_listener() {
	if !player_can_attack {
		return
	}

	player_can_attack = false
	events.Emit(PLAYER_ATTACK_EVENT, int32(5))

	temp_rect := player_weapon_rect
	temp_rect.PosX = (SCREEN_SIZE.X / 2) - 128
	player_sprite.UpdateRect(temp_rect)

	// Recovery delay
	time.AfterFunc(time.Millisecond*500, func() {
		player_can_attack = true
		player_sprite.UpdateRect(player_weapon_rect)
	})
}

func player_take_damage_listener(damage int) {
	player_health.TakeDamage(damage)
}
