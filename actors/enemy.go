package actors

import (
	"fmt"
	"gates/actors/enemies"
	game_events "gates/events"
	"gates/spells"
	"gates/systems"
	utils1 "gates/utils"
	"gates/values"
	"time"

	"github.com/Papiermond/eventbus"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var enemy_specs enemies.EnemySpecs
var enemy_go *lifecycle.GameObject
var enemy_sprite *render.Sprite
var enemy_health *systems.Health
var enemy_damage_ui_text *render.Font
var enemy_hp_rect utils.RectSpecs
var enemy_anim_position utils.RectSpecs
var enemy_effect_icons []spells.EffectIcon
var enemy_hp_max_width int
var enemy_stack_count int
var enemy_burn_damage int
var enemy_is_alive bool = false
var enemy_is_burn = false
var enemy_is_cold = false
var enemy_is_paralized = false
var enemy_attack_done = make(chan bool)
var enemy_burn_done = make(chan bool)

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
				damage := e.(game_events.PlayerAttackEvent).Damage
				effect := e.(game_events.PlayerAttackEvent).Effect

				enemy_take_damage(damage, effect)
			})

			events.Subscribe(events.Game, game_events.PLAYER_BREAK_SPELL_EVENT, func(data any) {
				println(values.Magenta + "Enemy knows the player is not concentrating anymore. All effects should stop now" + values.Reset)
				enemy_reset_effects()
			})
		},
		Update: func() {
			enemy_hp_rect.Width = (enemy_hp_max_width * enemy_health.GetCurrent()) / enemy_specs.HP
		},
		Render: func() {
			render.DrawRect(enemy_hp_rect, render.Red)
			enemy_sprite.UpdateRect(enemy_anim_position)

			if enemy_is_burn || enemy_is_cold || enemy_is_paralized {
				for i, icon := range enemy_effect_icons {
					icon.Rect.PosX = enemy_hp_rect.PosX + (i * (32 + 8))
					icon.Rect.PosY = enemy_hp_rect.PosY + enemy_hp_rect.Height + 8
					icon.Rect.Width = 32
					icon.Rect.Height = 32

					render.DrawRect(icon.Rect, icon.Color)
				}
			}
		},
	})
}

func LoadEnemy(specs enemies.EnemySpecs) {
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

	enemy_effect_icons = make([]spells.EffectIcon, 0)

	enemy_is_burn = false
	enemy_is_cold = false
	enemy_is_paralized = false
	enemy_stack_count = 0
	enemy_burn_damage = 0

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

		println(values.Cyan + "enemy stop is calling the reset of status effects" + values.Reset)
		enemy_reset_effects()
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

func enemy_take_damage(player_damage int, effect spells.Effect) {
	if !enemy_is_alive {
		return
	}

	var defense float64 = float64(enemy_specs.Defense)
	if enemy_is_cold && enemy_stack_count > 0 {
		defense = defense / float64(enemy_stack_count+1)
		println(values.Blue+"Enemy is cold so it has a debuff on it's defense. Current debuf is:"+values.Reset, enemy_stack_count)
	}
	if defense < 1 {
		defense = 1
	}

	var raw float64 = float64(player_damage) / defense
	if raw < 1 {
		raw = 1
	}
	damage := utils1.CalcDamange(int(raw), int(raw)/2)

	print(fmt.Sprintf("%sPlayer is attacking with %d damage %s\n", values.Green, damage, values.Reset))

	enemy_health.TakeDamage(int(damage))
	enemy_scale(-1)

	enemy_show_damage_text(damage, render.White)
	time.AfterFunc(time.Millisecond*400, func() {
		enemy_scale(1)
	})

	if enemy_health.GetCurrent() <= 0 {
		enemy_stop()
		return
	}

	if effect.Stack > 0 {
		enemy_stack_count += effect.Stack

		switch effect.Type {
		case spells.Burn:
			enemy_apply_burn(damage)
		case spells.Cold:
			enemy_apply_cold()
		case spells.Paralysis:
			enemy_apply_paralysis()
		}
	}
}

func enemy_reset_effects() {
	if enemy_is_burn {
		enemy_is_burn = false
		enemy_burn_damage = 0
		enemy_burn_done <- true
	}

	if enemy_is_cold {
		enemy_is_cold = false
	}

	if enemy_is_paralized {
		enemy_is_paralized = false
	}

	enemy_stack_count = 0
	enemy_effect_icons = make([]spells.EffectIcon, 0)
}

func enemy_apply_burn(base_damage int) {
	var raw float64 = float64(base_damage) * 0.1
	if raw < 1 {
		raw = 1
	}

	enemy_burn_damage = int(raw) * enemy_stack_count
	print(fmt.Sprintf(values.Red+"Current Burning stack: %v\n"+values.Reset, enemy_stack_count))

	enemy_effect_icons = make([]spells.EffectIcon, 0)
	for range enemy_stack_count {
		enemy_effect_icons = append(enemy_effect_icons, spells.EffectIcon{
			Color: render.Orange,
		})
	}

	if !enemy_is_burn {
		enemy_is_burn = true

		go enemy_burn_task()
	}
}

func enemy_apply_cold() {
	if !enemy_is_cold {
		enemy_is_cold = true
	}

	enemy_effect_icons = make([]spells.EffectIcon, 0)
	for range enemy_stack_count {
		enemy_effect_icons = append(enemy_effect_icons, spells.EffectIcon{
			Color: render.Blue,
		})
	}
}

func enemy_apply_paralysis() {
	if !enemy_is_paralized {
		enemy_is_paralized = true
	}

	enemy_effect_icons = make([]spells.EffectIcon, 0)
	for range enemy_stack_count {
		enemy_effect_icons = append(enemy_effect_icons, spells.EffectIcon{
			Color: render.Purple,
		})
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
			skip := false
			if enemy_is_paralized {
				skip = true
				enemy_stack_count -= 1
				enemy_effect_icons = enemy_effect_icons[:len(enemy_effect_icons)-1]

				if enemy_stack_count <= 0 {
					events.Emit(events.Game, game_events.EnemyBreakParalysisEvent{})
					enemy_reset_effects()
				}
			}

			if skip {
				println(values.Yellow + "Enemy is paralyzed. Skipping attack" + values.Reset)
				enemy_shake()
			} else {
				damage := enemy_specs.Attack_Damage

				if enemy_is_cold && enemy_stack_count > 0 {
					damage = damage / (enemy_stack_count + 1)
					println(values.Blue+"Enemy is cold and have a attack debuff. Current debuff:"+values.Reset, enemy_stack_count)
				}
				if damage < 1 {
					damage = 1
				}

				enemy_scale(1)
				time.AfterFunc(time.Millisecond*400, func() {
					enemy_scale(-1)
				})

				time.AfterFunc(time.Millisecond*200, func() {
					game_events.Bus.Publish(game_events.EnemyAttackEvent{
						Damage: damage,
					})
				})
			}
		}
	}
}

func enemy_burn_task() {
	ticker := time.NewTicker(time.Millisecond * 1000)

	for {
		select {
		case <-enemy_burn_done:
			println("Enemy is no longer burning")
			ticker.Stop()
			return
		case <-ticker.C:
			print(fmt.Sprintf("Enemy is burning at %v damage\n", enemy_burn_damage))
			enemy_health.TakeDamage(enemy_burn_damage)
			enemy_show_damage_text(enemy_burn_damage, render.Orange)

			if enemy_health.GetCurrent() <= 0 {
				enemy_stop()
			}
		}
	}
}

func enemy_scale(direction int) {
	if !enemy_is_alive {
		return
	}

	enemy_anim_position.PosX -= 64 * direction
	enemy_anim_position.PosY -= 64 * direction
	enemy_anim_position.Width += 128 * direction
	enemy_anim_position.Height += 128 * direction
}

func enemy_shake() {
	if !enemy_is_alive {
		return
	}

	og_pos := enemy_anim_position.PosX
	enemy_anim_position.PosX += 64
	time.AfterFunc(time.Millisecond*100, func() {
		enemy_anim_position.PosX = og_pos - 64
	})
	time.AfterFunc(time.Millisecond*150, func() {
		enemy_anim_position.PosX = og_pos + 64
	})
	time.AfterFunc(time.Millisecond*200, func() {
		enemy_anim_position.PosX = og_pos - 64
	})
	time.AfterFunc(time.Millisecond*250, func() {
		enemy_anim_position.PosX = og_pos
	})
}

func enemy_show_damage_text(damage int, color render.Color) {
	enemy_damage_ui_text.UpdateText(fmt.Sprint(damage))
	enemy_damage_ui_text.UpdateColor(color)
	time.AfterFunc(time.Millisecond*1200, func() {
		enemy_damage_ui_text.UpdateColor(render.Transparent)
	})
}
