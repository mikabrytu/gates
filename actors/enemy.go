package actors

import (
	"fmt"
	game_events "gates/events"
	"gates/systems"
	"gates/values"

	"github.com/Papiermond/eventbus"
	"github.com/mikabrytu/gomes-engine/events"
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

var enemy_go *lifecycle.GameObject
var enemy_specs EnemySpecs
var enemy_health *systems.Health
var enemy_sprite *render.Sprite
var enemy_sprite_rect utils.RectSpecs
var enemy_hp_rect utils.RectSpecs
var enemy_is_alive bool

func Enemy() {
	enemy_init()
	enemy_respawn()

	events.Subscribe(events.Game, game_events.GAME_RESTART_EVENT, func(data any) {
		restart := data.(game_events.GameRestartEvent)
		println(values.Red + restart.Message + values.Reset)

		enemy_respawn()
	})

	enemy_go = lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			enemy_sprite.Init()

			game_events.Bus.Subscribe(game_events.PLAYER_ATTACK_EVENT, func(e eventbus.Event) {
				if !enemy_is_alive {
					return
				}

				enemy_takes_damage(e.(game_events.PlayerAttackEvent).Damage)
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

	// Set values
	enemy_sprite_rect = utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - (enemy_specs.Size / 2),
		PosY:   32,
		Width:  enemy_specs.Size,
		Height: enemy_specs.Size,
	}

	enemy_sprite = render.NewSprite(
		enemy_specs.Name,
		enemy_specs.Image_Path,
		enemy_sprite_rect,
		render.White,
	)

	enemy_hp_rect = enemy_sprite_rect
	enemy_hp_rect.PosY -= 24
	enemy_hp_rect.Height = 16

	enemy_is_alive = true

	// Init systems
	enemy_health = systems.InitHealth(enemy_specs.HP)
}

func enemy_respawn() {
	if enemy_is_alive {
		return
	}

	println(values.Red + enemy_specs.Name + " is respawning" + values.Reset)

	enemy_is_alive = true
	enemy_sprite.Enable()
	lifecycle.Enable(enemy_go)
}

func enemy_takes_damage(damage int) {
	println(values.Green + "Player attacks with " + fmt.Sprint(damage) + " damage" + values.Reset)

	enemy_health.TakeDamage(damage)

	if enemy_health.GetCurrent() <= 0 {
		enemy_dead()
	}
}

func enemy_dead() {
	println(values.Red + enemy_specs.Name + " is dead." + values.Reset)

	enemy_is_alive = false
	enemy_sprite.Disable()
	lifecycle.Disable(enemy_go)

	events.Emit(events.Game, game_events.EnemyDeadEvent{Message: enemy_specs.Name + " is dead"})
}
