package actors

import (
	"fmt"
	"gates/systems"
	"gates/values"
	"math/rand/v2"
	"time"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var player_sprite *render.Sprite
var player_health *systems.Health
var player_skills *systems.Skill
var player_weapon_rect utils.RectSpecs
var player_hp_rect utils.RectSpecs
var player_max_hp int
var player_max_hp_width int
var player_can_attack bool = true
var player_can_level_up = false

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

			events.Subscribe(events.INPUT_KEYBOARD_PRESSED_1, func(params ...any) error {
				player_level_up_listener(1)
				return nil
			})

			events.Subscribe(events.INPUT_KEYBOARD_PRESSED_2, func(params ...any) error {
				player_level_up_listener(2)
				return nil
			})

			events.Subscribe(events.INPUT_KEYBOARD_PRESSED_3, func(params ...any) error {
				player_level_up_listener(3)
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

func PlayerLevelUp() {
	player_can_level_up = true
	println("LEVEL UP. Select 1 to increase STR, 2 to increase INT and 3 to increase SPD")
}

func player_init() {
	player_max_hp = 100
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
	player_skills = systems.NewSkill()
}

func player_damage() int {
	sword := 6
	//bow := 3
	//spell_fire := 10

	weapon := sword
	mod := player_skills.STR

	dice_roll := rand.IntN(weapon-(weapon/2)) + (weapon / 2)
	var damage int = dice_roll + (player_skills.STR * mod)

	println(values.Green + fmt.Sprintf("Player is dealing %v damage to enemy", damage) + values.Reset)
	return damage
}

func player_click_listener() {
	if !player_can_attack {
		return
	}

	player_can_attack = false
	events.Emit(values.PLAYER_ATTACK_EVENT, int32(player_damage()))

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

func player_level_up_listener(skill int) {
	if !player_can_level_up {
		return
	}

	level_up := systems.Skill{}
	if skill == 1 {
		level_up.STR = 1
	}
	if skill == 2 {
		level_up.INT = 1
	}
	if skill == 3 {
		level_up.SPD = 1
	}

	player_skills.LevelUp(level_up)
	player_can_level_up = false
	player_health.Reset()

	println(fmt.Sprintf("Player Current Level: %d {STR: %d, INT: %d, SPD: %v} \n",
		player_skills.GetLevel(),
		player_skills.STR,
		player_skills.INT,
		player_skills.SPD))
}
