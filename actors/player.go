package actors

import (
	"fmt"
	game_events "gates/events"
	"gates/systems"
	"gates/values"
	"math/rand/v2"
	"time"

	"github.com/Papiermond/eventbus"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type Weapon struct {
	Name     string
	Damage   int
	Recovery int
	Modifier systems.Attribute
}

var player_sprite *render.Sprite
var player_health *systems.Health
var player_skills *systems.Skill
var player_current_weapon Weapon
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

			game_events.Bus.Subscribe(game_events.ENEMY_ATTACK_EVENT, func(e eventbus.Event) {
				damage := e.(game_events.EnemyAttackEvent).Damage
				player_take_damage_listener(int(damage))
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
	println(values.Blue + "LEVEL UP. Select 1 to increase STR, 2 to increase INT and 3 to increase SPD" + values.Reset)
}

func PlayerLoadWeapon(weapon Weapon) {
	player_current_weapon = weapon
	println(fmt.Sprintf("Player weapon loaded: %s", weapon.Name))
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

	println("Player initialized")
}

func player_damage() int {
	weapon := player_current_weapon.Damage

	mod := 1
	switch player_current_weapon.Modifier {
	case systems.STR:
		mod = player_skills.STR
	case systems.INT:
		mod = player_skills.INT
	case systems.SPD:
		mod = player_skills.SPD
	default:
		println("Couldn't find match between weapon attribute and player skills. Defaulting mod to 1...")
	}

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
	game_events.Bus.Publish(game_events.PlayerAttackEvent{
		Damage: player_damage(),
	})

	temp_rect := player_weapon_rect
	temp_rect.PosX = (values.SCREEN_SIZE.X / 2) - 128
	player_sprite.UpdateRect(temp_rect)

	// Recovery delay
	time.AfterFunc(time.Millisecond*time.Duration(player_current_weapon.Recovery), func() {
		player_can_attack = true
	})

	// Weapon animation
	time.AfterFunc(time.Millisecond*350, func() {
		player_sprite.UpdateRect(player_weapon_rect)
	})
}

func player_take_damage_listener(damage int) {
	player_health.TakeDamage(damage)

	if player_health.GetCurrent() <= 0 {
		game_events.Bus.Publish(game_events.GameOverEvent{
			Message: "Player is dead",
		})
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
