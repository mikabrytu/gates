package actors

import (
	"fmt"
	"gates/config"
	data "gates/internal/data/enemies"
	"gates/internal/events"
	"gates/pkg/health"
	"gates/pkg/math"
	spells "gates/pkg/spell"
	"time"

	"github.com/Papiermond/eventbus"
	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	gomesmath "github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type Enemy struct {
	specs          data.EnemySpecs
	instance       *lifecycle.GameObject
	sprite         *render.Sprite
	health         *health.Health
	damage_ui_text *render.Font
	hp_rect        utils.RectSpecs
	anim_position  utils.RectSpecs
	effect_icons   []spells.EffectIcon
	hp_max_width   int
	stack_count    int
	burn_damage    int
	is_alive       bool
	is_burn        bool
	is_cold        bool
	is_paralized   bool
	attack_done    chan bool
	burn_done      chan bool
}

func NewEnemy() *Enemy {
	enemy := &Enemy{
		effect_icons: make([]spells.EffectIcon, 0),
		stack_count:  0,
		burn_damage:  0,
		is_alive:     false,
		is_burn:      false,
		is_cold:      false,
		is_paralized: false,
		attack_done:  make(chan bool),
		burn_done:    make(chan bool),
	}

	enemy.damage_ui_text = render.NewFont(config.FONT_SPECS, config.SCREEN_SIZE)
	enemy.damage_ui_text.Init("0", render.White, gomesmath.Vector2{X: 0, Y: 0})
	enemy.damage_ui_text.AlignText(render.TopCenter, gomesmath.Vector2{X: 0, Y: 32})
	enemy.damage_ui_text.Disable()

	events.Bus.Subscribe(events.GAME_RESTART_EVENT, func(e eventbus.Event) {
		if enemy.is_alive {
			println("Enemy is already alive. Skipping respawn")
			return
		}

		if enemy.instance != nil {
			go func() {
				enemy.attack_done = make(chan bool)
			}()

			enemy.Enable()
		}
	})

	events.Bus.Subscribe(events.GAME_OVER_EVENT, func(e eventbus.Event) {
		if !enemy.is_alive {
			return
		}

		println(e.(events.GameOverEvent).Message)
		enemy.stop()
	})

	enemy.instance = lifecycle.Register(&lifecycle.GameObject{
		Start:   enemy.start,
		Update:  enemy.update,
		Render:  enemy.render,
		Destroy: enemy.destroy,
	})

	return enemy
}

func (e *Enemy) start() {
	events.Bus.Subscribe(events.PLAYER_ATTACK_EVENT, func(data eventbus.Event) {
		damage := data.(events.PlayerAttackEvent).Damage
		effect := data.(events.PlayerAttackEvent).Effect

		e.take_damage(damage, effect)
	})

	gomesevents.Subscribe(gomesevents.Game, events.PLAYER_BREAK_SPELL_EVENT, func(data any) {
		println(config.Magenta + "Enemy knows the player is not concentrating anymore. All effects should stop now" + config.Reset)
		e.reset_effects()
	})
}

func (e *Enemy) update() {
	e.hp_rect.Width = (e.hp_max_width * e.health.GetCurrent()) / e.specs.HP
}

func (e *Enemy) render() {
	render.DrawRect(e.hp_rect, render.Red)
	e.sprite.UpdateRect(e.anim_position)

	if e.is_burn || e.is_cold || e.is_paralized {
		for i, icon := range e.effect_icons {
			icon.Rect.PosX = e.hp_rect.PosX + (i * (32 + 8))
			icon.Rect.PosY = e.hp_rect.PosY + e.hp_rect.Height + 8
			icon.Rect.Width = 32
			icon.Rect.Height = 32

			render.DrawRect(icon.Rect, icon.Color)
		}
	}
}

func (e *Enemy) destroy() {
	e.damage_ui_text.Clear()
	e.sprite.Clear()
	e.health = nil
}

func (e *Enemy) Enable() {
	if e.sprite != nil {
		e.sprite.Enable()
	}

	lifecycle.Enable(e.instance)
}

func (e *Enemy) Disable() {
	if e.damage_ui_text != nil {
		e.damage_ui_text.Disable()
	}

	if e.sprite != nil {
		e.sprite.Disable()
	}

	lifecycle.Disable(e.instance)
}

func (e *Enemy) LoadData(specs data.EnemySpecs) {
	message := fmt.Sprintf("Loading enemy %v", specs.Name)
	println(config.Red + message + config.Reset)

	rect := utils.RectSpecs{
		PosX:   (config.SCREEN_SIZE.X / 2) - (specs.Size / 2),
		PosY:   32,
		Width:  specs.Size,
		Height: specs.Size,
	}
	e.anim_position = rect

	e.hp_max_width = rect.Width
	e.hp_rect = rect
	e.hp_rect.PosY -= 24
	e.hp_rect.Height = 16

	if e.sprite == nil {
		e.sprite = render.NewSprite(
			specs.Name,
			specs.Image_Path,
			rect,
			render.White,
		)
		e.sprite.Init()
	} else {
		if e.specs.Name != specs.Name {
			e.sprite.UpdateImage(specs.Image_Path, render.Transparent)
			e.sprite.Disable()
		}
	}

	if e.health == nil {
		e.health = health.Init(specs.HP)
	} else {
		e.health.ChangeMax(specs.HP)
		e.health.Reset()
	}

	e.specs = specs
	e.is_alive = true

	go e.attack_task(e.specs.Attack_Interval)
}

func (e *Enemy) stop() {
	go func() {
		e.attack_done <- true

		println(config.Cyan + "enemy stop is calling the reset of status effects" + config.Reset)
		e.reset_effects()
	}()

	e.is_alive = false
	e.sprite.Disable()
	lifecycle.Disable(e.instance)

	go func() {
		message := config.Yellow + "Enemy " + e.specs.Name + " is dead" + config.Reset
		events.Bus.Publish(events.EnemyDeadEvent{
			XP:      e.specs.XP,
			Message: message,
		})
	}()
}

func (e *Enemy) take_damage(player_damage int, effect spells.Effect) {
	if !e.is_alive {
		return
	}

	var defense float64 = float64(e.specs.Defense)
	if e.is_cold && e.stack_count > 0 {
		defense = defense / float64(e.stack_count+1)
		println(config.Blue+"Enemy is cold so it has a debuff on it's defense. Current debuf is:"+config.Reset, e.stack_count)
	}
	if defense < 1 {
		defense = 1
	}

	var raw float64 = float64(player_damage) / defense
	if raw < 1 {
		raw = 1
	}
	damage := math.CalcDamange(int(raw), int(raw)/2)

	print(fmt.Sprintf("%sPlayer is attacking with %d damage %s\n", config.Green, damage, config.Reset))

	e.health.TakeDamage(int(damage))
	e.scale(-1)

	e.show_damage_text(damage, render.White)
	time.AfterFunc(time.Millisecond*400, func() {
		e.scale(1)
	})

	if e.health.GetCurrent() <= 0 {
		e.stop()
		return
	}

	if effect.Stack > 0 {
		e.stack_count += effect.Stack

		switch effect.Type {
		case spells.Burn:
			e.apply_burn(damage)
		case spells.Cold:
			e.apply_cold()
		case spells.Paralysis:
			e.apply_paralysis()
		}
	}
}

func (e *Enemy) reset_effects() {
	if e.is_burn {
		e.is_burn = false
		e.burn_damage = 0
		e.burn_done <- true
	}

	if e.is_cold {
		e.is_cold = false
	}

	if e.is_paralized {
		e.is_paralized = false
	}

	e.stack_count = 0
	e.effect_icons = make([]spells.EffectIcon, 0)
}

func (e *Enemy) apply_burn(base_damage int) {
	var raw float64 = float64(base_damage) * 0.1
	if raw < 1 {
		raw = 1
	}

	e.burn_damage = int(raw) * e.stack_count
	print(fmt.Sprintf(config.Red+"Current Burning stack: %v\n"+config.Reset, e.stack_count))

	e.effect_icons = make([]spells.EffectIcon, 0)
	for range e.stack_count {
		e.effect_icons = append(e.effect_icons, spells.EffectIcon{
			Color: render.Orange,
		})
	}

	if !e.is_burn {
		e.is_burn = true

		go e.burn_task()
	}
}

func (e *Enemy) apply_cold() {
	if !e.is_cold {
		e.is_cold = true
	}

	e.effect_icons = make([]spells.EffectIcon, 0)
	for range e.stack_count {
		e.effect_icons = append(e.effect_icons, spells.EffectIcon{
			Color: render.Blue,
		})
	}
}

func (e *Enemy) apply_paralysis() {
	if !e.is_paralized {
		e.is_paralized = true
	}

	e.effect_icons = make([]spells.EffectIcon, 0)
	for range e.stack_count {
		e.effect_icons = append(e.effect_icons, spells.EffectIcon{
			Color: render.Purple,
		})
	}
}

func (e *Enemy) attack_task(interval int) {
	println("Starting enemy attack task...")
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))

	for {
		select {
		case <-e.attack_done:
			println("Stopping enemy attack")
			ticker.Stop()
			return
		case <-ticker.C:
			skip := false
			if e.is_paralized {
				skip = true
				e.stack_count -= 1
				e.effect_icons = e.effect_icons[:len(e.effect_icons)-1]

				if e.stack_count <= 0 {
					gomesevents.Emit(gomesevents.Game, events.EnemyBreakParalysisEvent{})
					e.reset_effects()
				}
			}

			if skip {
				println(config.Yellow + "Enemy is paralyzed. Skipping attack" + config.Reset)
				e.shake()
			} else {
				damage := math.CalcDamange(e.specs.Attack_Damage, e.specs.Attack_Damage/2)

				if e.is_cold && e.stack_count > 0 {
					damage = damage / (e.stack_count + 1)
					println(config.Blue+"Enemy is cold and have a attack debuff. Current debuff:"+config.Reset, e.stack_count)
				}

				e.scale(1)
				time.AfterFunc(time.Millisecond*400, func() {
					e.scale(-1)
				})

				time.AfterFunc(time.Millisecond*200, func() {
					events.Bus.Publish(events.EnemyAttackEvent{
						Damage: damage,
					})
				})
			}
		}
	}
}

func (e *Enemy) burn_task() {
	ticker := time.NewTicker(time.Millisecond * 1000)

	for {
		select {
		case <-e.burn_done:
			println("Enemy is no longer burning")
			ticker.Stop()
			return
		case <-ticker.C:
			print(fmt.Sprintf("Enemy is burning at %v damage\n", e.burn_damage))
			e.health.TakeDamage(e.burn_damage)
			e.show_damage_text(e.burn_damage, render.Orange)

			if e.health.GetCurrent() <= 0 {
				e.stop()
			}
		}
	}
}

func (e *Enemy) scale(direction int) {
	if !e.is_alive {
		return
	}

	e.anim_position.PosX -= 64 * direction
	e.anim_position.PosY -= 64 * direction
	e.anim_position.Width += 128 * direction
	e.anim_position.Height += 128 * direction
}

func (e *Enemy) shake() {
	if !e.is_alive {
		return
	}

	og_pos := e.anim_position.PosX
	e.anim_position.PosX += 64
	time.AfterFunc(time.Millisecond*100, func() {
		e.anim_position.PosX = og_pos - 64
	})
	time.AfterFunc(time.Millisecond*150, func() {
		e.anim_position.PosX = og_pos + 64
	})
	time.AfterFunc(time.Millisecond*200, func() {
		e.anim_position.PosX = og_pos - 64
	})
	time.AfterFunc(time.Millisecond*250, func() {
		e.anim_position.PosX = og_pos
	})
}

func (e *Enemy) show_damage_text(damage int, color render.Color) {
	e.damage_ui_text.UpdateText(fmt.Sprint(damage))
	e.damage_ui_text.UpdateColor(color)
	e.damage_ui_text.Enable()
	time.AfterFunc(time.Millisecond*1200, func() {
		e.damage_ui_text.Disable()
	})
}
