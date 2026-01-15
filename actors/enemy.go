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
	"github.com/mikabrytu/gomes-engine/math"
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
var enemy_go *lifecycle.GameObject
var enemy_sprite *render.Sprite
var enemy_health *systems.Health
var enemy_damage_ui_text *render.Font
var enemy_hp_rect utils.RectSpecs
var enemy_hp_max_width int
var enemy_anim_position utils.RectSpecs
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
				enemy_take_damage(e.(game_events.PlayerAttackEvent).Damage)
			})
		},
		Update: func() {
			enemy_hp_rect.Width = (enemy_hp_max_width * enemy_health.GetCurrent()) / enemy_specs.HP
		},
		Render: func() {
			render.DrawRect(enemy_hp_rect, render.Red)
			enemy_sprite.UpdateRect(enemy_anim_position)
		},
	})
}

func LoadEnemy(specs EnemySpecs) {
	message := fmt.Sprintf("Changing enemy specs. Loading %v", specs.Name)
	println(message)

	if enemy_sprite != nil && enemy_specs.Name != specs.Name {
		enemy_sprite.UpdateImage(specs.Image_Path, render.Transparent)
		enemy_sprite.Disable()
	}

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
	enemy_anim_position = rect

	enemy_hp_max_width = rect.Width
	enemy_hp_rect = rect
	enemy_hp_rect.PosY -= 24
	enemy_hp_rect.Height = 16

	enemy_damage_ui_text = render.NewFont(values.FONT_SPECS, values.SCREEN_SIZE)
	enemy_damage_ui_text.Init("10", render.Transparent, math.Vector2{X: 0, Y: 0})
	enemy_damage_ui_text.AlignText(render.TopCenter, math.Vector2{X: 0, Y: 32})

	if enemy_sprite == nil {
		enemy_sprite = render.NewSprite(
			enemy_specs.Name,
			enemy_specs.Image_Path,
			rect,
			render.White,
		)
	} else {
		enemy_sprite.Enable()
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
	enemy_sprite.Disable()
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

func enemy_take_damage(player_damage int) {
	if !enemy_is_alive {
		return
	}

	raw := player_damage / enemy_specs.Defense
	damage := utils1.CalcDamange(raw, raw/2)

	print(fmt.Sprintf("%sPlayer is attacking with %d damage %s\n", values.Green, damage, values.Reset))

	enemy_health.TakeDamage(int(damage))
	enemy_scale(-1)

	enemy_damage_ui_text.UpdateText(fmt.Sprint(damage))
	enemy_damage_ui_text.UpdateColor(render.White)
	time.AfterFunc(time.Millisecond*1200, func() {
		enemy_damage_ui_text.UpdateColor(render.Transparent)
	})
	time.AfterFunc(time.Millisecond*400, func() {
		enemy_scale(1)
	})

	if enemy_health.GetCurrent() <= 0 {
		enemy_stop()
	}
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
			damage := enemy_specs.Attack_Damage
			enemy_scale(1)
			time.AfterFunc(time.Millisecond*400, func() {
				enemy_scale(-1)
			})

			game_events.Bus.Publish(game_events.EnemyAttackEvent{
				Damage: damage,
			})
		}
	}
}

func enemy_scale(direction int) {
	if !enemy_is_alive {
		return
	}

	sprite_rect := enemy_sprite.GetRect()
	enemy_anim_position.PosX = sprite_rect.PosX - 64*direction
	enemy_anim_position.PosY = sprite_rect.PosY - 64*direction
	enemy_anim_position.Width = sprite_rect.Width + 128*direction
	enemy_anim_position.Height = sprite_rect.Height + 128*direction
}
