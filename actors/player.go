package actors

import (
	"gates/systems"
	"gates/values"
	"time"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var player_sprite *render.Sprite
var player_health *systems.Health
var player_weapon_rect utils.RectSpecs
var player_hp_rect utils.RectSpecs
var player_attack_damage int = 10
var player_max_hp int
var player_max_hp_width int
var player_can_attack bool = true

func Player() {
	player_init()

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			player_sprite.Init()

			events.Subscribe(events.INPUT_MOUSE_CLICK_DOWN, func(params ...any) error {
				player_click_listener()
				return nil
			})

			events.Subscribe(values.ENEMY_ATTACK_EVENT, func(params ...any) error {
				damage := params[0].([]any)[0].([]any)[0].(int32)
				player_take_damage_listener(int(damage))

				return nil
			})
		},
		Update: func() {
			player_hp_rect.Width = (player_max_hp_width * player_health.GetCurrent()) / player_max_hp
		},
		Render: func() {
			render.DrawRect(player_hp_rect, render.Green)
		},
	})
}

func player_init() {
	player_max_hp = 1000
	player_max_hp_width = 512

	player_weapon_rect = utils.RectSpecs{
		PosX:   values.SCREEN_SIZE.X - 256 - 64,
		PosY:   values.SCREEN_SIZE.Y - 512,
		Width:  256,
		Height: 512,
	}

	player_hp_rect = utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - (player_max_hp_width / 2),
		PosY:   values.SCREEN_SIZE.Y - 40,
		Width:  player_max_hp_width,
		Height: 24,
	}

	player_sprite = render.NewSprite(
		"Player Weapon",
		"assets/images/dagger.jpg",
		player_weapon_rect,
		render.Transparent,
	)

	player_health = systems.InitHealth(player_max_hp)
}

func player_click_listener() {
	if !player_can_attack {
		return
	}

	player_can_attack = false
	events.Emit(values.PLAYER_ATTACK_EVENT, int32(player_attack_damage))

	temp_rect := player_weapon_rect
	temp_rect.PosX = (values.SCREEN_SIZE.X / 2) - 128
	player_sprite.UpdateRect(temp_rect)

	// Recovery delay
	time.AfterFunc(time.Millisecond*500, func() {
		player_can_attack = true
		player_sprite.UpdateRect(player_weapon_rect)
	})
}

func player_take_damage_listener(damage int) {
	player_health.TakeDamage(damage)

	if player_health.GetCurrent() <= 0 {
		events.Emit(values.GAME_OVER_EVENT)
	}
}
