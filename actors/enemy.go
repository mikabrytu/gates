package actors

import (
	game_events "gates/events"
	"gates/values"

	"github.com/Papiermond/eventbus"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type EnemySpecs struct {
	Name            string
	Image_Path      string
	Size            int
	HP              int
	Attack_Interval int
	Attack_Damage   int
	Defense         int
}

var enemy_specs EnemySpecs
var enemy_sprite *render.Sprite
var enemy_sprite_rect utils.RectSpecs
var enemy_hp_rect utils.RectSpecs

func Enemy() {
	enemy_init()

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			enemy_sprite = render.NewSprite(
				enemy_specs.Name,
				enemy_specs.Image_Path,
				enemy_sprite_rect,
				render.White,
			)
			enemy_sprite.Init()

			game_events.Bus.Subscribe(game_events.PLAYER_ATTACK_EVENT, func(e eventbus.Event) {
				attack := e.(game_events.PlayerAttackEvent)
				println("Enemy take damage", attack.Damage)
			})
		},
		Render: func() {
			render.DrawRect(enemy_hp_rect, render.Red)
		},
	})
}

func LoadEnemy(specs EnemySpecs) {
	enemy_specs = specs
}

func enemy_init() {
	if enemy_specs.Name == "" {
		panic("Enemy specs not loaded!")
	}

	println("Initializing", enemy_specs.Name)

	enemy_sprite_rect = utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - (enemy_specs.Size / 2),
		PosY:   32,
		Width:  enemy_specs.Size,
		Height: enemy_specs.Size,
	}

	enemy_hp_rect = enemy_sprite_rect
	enemy_hp_rect.PosY -= 24
	enemy_hp_rect.Height = 16
}
