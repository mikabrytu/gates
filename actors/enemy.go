package actors

import (
	"fmt"
	game_events "gates/events"
	"gates/systems"
	utils1 "gates/utils"
	"gates/values"
	"time"

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
	Attack_Damage   Dice
}

var enemy_specs EnemySpecs
var enemy_go *lifecycle.GameObject
var enemy_sprite *render.Sprite
var enemy_health *systems.Health
var enemy_hp_rect utils.RectSpecs
var enemt_attack_circle utils.CircleSpecs
var enemy_hp_max_width int
var enemy_render_attack bool = false
var enemy_is_alive bool = false
var enemy_attack_done = make(chan bool)

func Enemy() {
	// CRASH: There's a bug with the event system that blocks the main thread when subscribing to more events on this file (maybe package?)
	// I need to debug the event system to see what's going on

	enemy_init()
	enemy_respawn()

	game_events.Bus.Subscribe(game_events.GAME_OVER_EVENT, func(e eventbus.Event) {
		if !enemy_is_alive {
			return
		}

		println(e.(game_events.GameOverEvent).Message)
		enemy_stop()
	})

	enemy_go = lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			enemy_sprite.Init()

			game_events.Bus.Subscribe(game_events.PLAYER_ATTACK_EVENT, func(e eventbus.Event) {
				if !enemy_is_alive {
					return
				}

				damage := e.(game_events.PlayerAttackEvent).Damage
				enemy_health.TakeDamage(int(damage))

				if enemy_health.GetCurrent() <= 0 {
					enemy_stop()
				}
			})
		},
		Update: func() {
			enemy_hp_rect.Width = (enemy_hp_max_width * enemy_health.GetCurrent()) / enemy_specs.HP
		},
		Render: func() {
			render.DrawRect(enemy_hp_rect, render.Red)

			if enemy_render_attack {
				render.DrawCircle(enemt_attack_circle, render.Red)
			}
		},
	})
}

func LoadEnemy(specs EnemySpecs) {
	message := fmt.Sprintf("Changing enemy specs. Loading %v", specs.Name)
	println(message)

	enemy_specs = specs
}

func enemy_init() {
	message := fmt.Sprintf("Initializing %v", enemy_specs.Name)
	println(message)

	enemy_is_alive = true

	rect := utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - (enemy_specs.Size / 2),
		PosY:   32,
		Width:  enemy_specs.Size,
		Height: enemy_specs.Size,
	}

	enemy_hp_max_width = rect.Width
	enemy_hp_rect = rect
	enemy_hp_rect.PosY -= 24
	enemy_hp_rect.Height = 16

	enemt_attack_circle = utils.CircleSpecs{
		PosX:   values.SCREEN_SIZE.X / 2,
		PosY:   rect.PosY + rect.Height + 64,
		Radius: 64,
	}

	if enemy_sprite == nil {
		enemy_sprite = render.NewSprite(
			enemy_specs.Name,
			enemy_specs.Image_Path,
			rect,
			render.White,
		)
	} else {
		enemy_sprite.UpdateRect(rect)
		enemy_sprite.UpdateImage(enemy_specs.Image_Path)
		enemy_sprite.Init()
	}

	if enemy_health == nil {
		enemy_health = systems.InitHealth(enemy_specs.HP)
	} else {
		enemy_health.ChangeMax(enemy_specs.HP)
		enemy_health.Reset()
	}

	go enemy_attack_task(enemy_specs.Attack_Interval)
}

func enemy_stop() {
	go func() {
		enemy_attack_done <- true
	}()

	enemy_is_alive = false
	enemy_sprite.ClearSprite()
	lifecycle.Disable(enemy_go)

	go func() {
		message := values.Yellow + "Enemy " + enemy_specs.Name + " is dead" + values.Reset
		game_events.Bus.Publish(game_events.EnemyDeadEvent{
			Message: message,
		})
	}()
}

func enemy_respawn() {
	game_events.Bus.Subscribe(game_events.GAME_RESTART_EVENT, func(e eventbus.Event) {
		if enemy_is_alive {
			println("Enemy is already alive. Skipping respawn")
			return
		}

		if enemy_go != nil {
			go func() {
				enemy_attack_done = make(chan bool)
			}()

			enemy_init()
			lifecycle.Enable(enemy_go)
		}
	})
}

func enemy_attack_task(interval int) {
	println("Starting enemy attack task...")
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))

	for {
		select {
		case <-enemy_attack_done:
			println("Stopping enemy attack")
			ticker.Stop()
			return
		case <-ticker.C:
			dice := enemy_specs.Attack_Damage.GetValue()
			damage := utils1.CalcDamange(dice, dice/2)

			message := values.Red + fmt.Sprintf("Enemy attacks with %d damage", damage) + values.Reset
			println(message)

			enemy_render_attack = true
			time.AfterFunc(time.Millisecond*800, func() {
				enemy_render_attack = false
			})

			game_events.Bus.Publish(game_events.EnemyAttackEvent{
				Damage:  damage,
				Message: message,
			})
		}
	}
}
