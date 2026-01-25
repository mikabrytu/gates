package actors

import (
	"fmt"
	"gates/config"
	data "gates/internal/data/weapons"
	"gates/internal/events"
	"gates/pkg/health"
	"gates/pkg/math"
	"gates/pkg/skill"
	"gates/pkg/spell"
	"math/rand/v2"
	"time"

	"github.com/Papiermond/eventbus"
	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	gomesmath "github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type Player struct {
	instance            *lifecycle.GameObject
	sprite              *render.Sprite
	health              *health.Health
	skills              *skill.Skill
	damage_ui_text      *render.Font
	current_weapon      data.Weapon
	weapon_rect         utils.RectSpecs
	hp_rect             utils.RectSpecs
	recovery_rect       utils.RectSpecs
	anim_position       utils.RectSpecs
	defense_rect        utils.RectSpecs
	attack_start_time   time.Time
	defense_start_time  time.Time
	max_hp              int
	max_hp_width        int
	og_recovery_width   int
	concentration_count int
	can_attack          bool
	can_level_up        bool
	is_attacking        bool
	is_animating        bool
	is_defending        bool
	is_stunned          bool
	is_concentrating    bool
	enabled             bool
}

const PLAYER_HP_PER_LEVEL int = 5
const PLAYER_STUN_DELAY int = 6000
const PLAYER_SPELL_EFFECT_CHANCE = 50

func NewPlayer() *Player {
	player := &Player{
		concentration_count: 0,
		can_attack:          true,
		can_level_up:        false,
		is_attacking:        false,
		is_animating:        false,
		is_defending:        false,
		is_stunned:          false,
		is_concentrating:    false,
		enabled:             false,
	}

	player.max_hp_width = 512

	player.hp_rect = utils.RectSpecs{
		PosX:   (config.SCREEN_SIZE.X / 2) - (player.max_hp_width / 2),
		PosY:   config.SCREEN_SIZE.Y - 40,
		Width:  player.max_hp_width,
		Height: 24,
	}

	player.defense_rect = utils.RectSpecs{
		PosX:   (config.SCREEN_SIZE.X / 2) - 256,
		PosY:   config.SCREEN_SIZE.Y - (config.SCREEN_SIZE.Y / 3),
		Width:  512,
		Height: 512,
	}

	player.recovery_rect = player.hp_rect
	player.recovery_rect.PosY -= 24
	player.recovery_rect.Height = 12

	if player.damage_ui_text == nil {
		player.damage_ui_text = render.NewFont(config.FONT_SPECS, config.SCREEN_SIZE)
		player.damage_ui_text.Init("0", render.White, gomesmath.Vector2{X: 0, Y: 0})
		player.damage_ui_text.AlignText(render.BottomCenter, gomesmath.Vector2{X: 0, Y: 96})
		player.damage_ui_text.Disable()
	}

	player.og_recovery_width = player.recovery_rect.Width

	player.instance = lifecycle.Register(&lifecycle.GameObject{
		Start:   player.start,
		Update:  player.update,
		Render:  player.render,
		Destroy: player.destroy,
	})

	return player
}

func (p *Player) start() {
	// Input
	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_MOUSE_CLICK_DOWN, func(data any) {
		if !p.enabled {
			return
		}

		click := data.(gomesevents.InputMouseClickDownEvent)

		if p.is_animating || p.is_defending || p.is_stunned {
			println("Player cannot act now")
			return
		}

		if click.Index.Left == 1 {
			p.attack_listener()
		}

		if click.Index.Right == 1 {
			p.defend(true)
		}
	})

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_MOUSE_CLICK_UP, func(data any) {
		if !p.enabled {
			return
		}

		p.defend(false)
	})

	// Take Damage
	events.Bus.Subscribe(events.ENEMY_ATTACK_EVENT, func(e eventbus.Event) {
		damage := e.(events.EnemyAttackEvent).Damage
		p.take_damage_listener(int(damage))
	})

	events.Bus.Subscribe(events.ENEMY_DEAD_EVENT, func(e eventbus.Event) {
		p.lost_concentration(false)
	})

	// Enemy
	gomesevents.Subscribe(gomesevents.Game, events.ENEMY_BREAK_PARALYSIS_EVENT, func(data any) {
		println("Enemy is trying to break free from Paralysis")
		p.lost_concentration(false)
	})
}

func (p *Player) update() {
	p.hp_rect.Width = (p.max_hp_width * p.health.GetCurrent()) / p.health.GetMax()

	if p.is_attacking {
		elapsed := time.Since(p.attack_start_time).Milliseconds()
		t := float64(elapsed) / float64(p.get_recovery())

		width := math.Lerp(float64(p.og_recovery_width), 0, t)
		p.recovery_rect.Width = int(width)

		if t > 1 {
			t = 1
			p.is_attacking = false
			p.recovery_rect.Width = p.og_recovery_width
		}
	}
}

func (p *Player) render() {
	p.sprite.UpdateRect(p.anim_position)

	if p.is_defending {
		render.DrawRect(p.defense_rect, render.White)
	}

	if p.is_stunned {
		render.DrawRect(p.hp_rect, render.Yellow)
		render.DrawRect(p.recovery_rect, render.Yellow)
		return
	}

	render.DrawRect(p.hp_rect, render.Green)
	render.DrawRect(p.recovery_rect, render.Blue)
}

func (p *Player) destroy() {
	p.damage_ui_text.Clear()
	p.sprite.Clear()
	p.health = nil
	p.skills = nil
	p = nil
}

func (p *Player) Enable() {
	p.enabled = true

	if p.sprite != nil {
		p.sprite.Enable()
	}

	lifecycle.Enable(p.instance)
}

func (p *Player) Disable() {
	p.enabled = false

	if p.damage_ui_text != nil {
		p.damage_ui_text.Disable()
	}

	if p.sprite != nil {
		p.sprite.Disable()
		println("Player sprite should be disabled")
	}

	lifecycle.Disable(p.instance)
}

func (p *Player) LoadData(weapon data.Weapon, skills skill.Skill) {
	p.current_weapon = weapon
	p.weapon_rect = utils.RectSpecs{
		PosX:   config.SCREEN_SIZE.X - p.current_weapon.SpriteSize.X + p.current_weapon.SpriteOffset.X,
		PosY:   config.SCREEN_SIZE.Y - p.current_weapon.SpriteSize.Y + p.current_weapon.SpriteOffset.Y,
		Width:  p.current_weapon.SpriteSize.X,
		Height: p.current_weapon.SpriteSize.Y,
	}
	p.anim_position = p.weapon_rect

	p.sprite = render.NewSprite(
		"Player Weapon",
		p.current_weapon.SpritePath,
		p.weapon_rect,
		render.Transparent,
	)
	p.sprite.Init()
	p.sprite.Disable()

	p.skills = skill.NewSkill()
	p.skills.IncreaseSTR(skills.STR)
	p.skills.IncreaseINT(skills.INT)
	p.skills.IncreaseSPD(skills.SPD)
	p.max_hp = PLAYER_HP_PER_LEVEL * p.skills.GetTotalSkillPoints()
	p.health = health.Init(p.max_hp)

	message := fmt.Sprintf(
		"Player data loaded. Level %v: { STR: %v, INT: %v SPD: %v } | Weapon: %v\n",
		p.skills.GetLevel(),
		p.skills.STR,
		p.skills.INT,
		p.skills.SPD,
		p.current_weapon.Name,
	)
	print(config.Green + message + config.Reset)
}

func (p *Player) GetLevel() int {
	return p.skills.GetLevel()
}

func (p *Player) calc_damage() int {
	return p.current_weapon.Damage * p.skills.STR
}

func (p *Player) spell_effect() spell.Effect {
	if p.current_weapon.Type == data.Physical {
		return spell.Effect{}
	}

	if p.concentration_count >= p.skills.INT {
		println("Player cannot active more spells until concentration break")
		return spell.Effect{}
	}

	roll := rand.IntN(100)

	if roll <= PLAYER_SPELL_EFFECT_CHANCE {
		p.is_concentrating = true
		p.concentration_count += 1

		var effect_type spell.EffectType
		name := ""
		switch p.current_weapon.Type {
		case data.Fire:
			effect_type = spell.Burn
			name = "Burn"
		case data.Ice:
			effect_type = spell.Cold
			name = "Cold"
		case data.Shock:
			effect_type = spell.Paralysis
			name = "Paralysis"
		}

		print(fmt.Sprintf(config.Yellow+"Player triggered spell effect: %v. Max Stack: %v\n"+config.Reset, name, p.skills.INT))

		return spell.Effect{
			Type:  effect_type,
			Stack: 1,
		}
	}

	return spell.Effect{}
}

func (p *Player) attack_listener() {
	if !p.can_attack {
		return
	}

	p.can_attack = false
	p.is_attacking = true
	p.is_animating = true
	p.attack_start_time = time.Now()

	events.Bus.Publish(events.PlayerAttackEvent{
		Damage: p.calc_damage(),
		Effect: p.spell_effect(),
	})

	temp_rect := p.weapon_rect
	temp_rect.PosX = (config.SCREEN_SIZE.X / 2) - (p.current_weapon.SpriteSize.X / 2)
	p.anim_position = temp_rect

	// Recovery delay
	time.AfterFunc(time.Millisecond*time.Duration(p.get_recovery()), func() {
		p.can_attack = true
	})

	// Weapon animation
	time.AfterFunc(time.Millisecond*350, func() {
		p.is_animating = false
		p.anim_position = p.weapon_rect
	})
}

func (p *Player) take_damage_listener(base_damage int) {
	damage := max(base_damage/p.skills.STR, 1)

	if p.is_defending {
		p.negate_damage(base_damage)
		return
	}

	message := config.Red + fmt.Sprintf("Enemy attacks with %d damage", damage) + config.Reset
	println(message)

	p.health.TakeDamage(damage)

	p.damage_ui_text.UpdateText(fmt.Sprint(damage))
	p.damage_ui_text.Enable()
	time.AfterFunc(time.Millisecond*1200, func() {
		p.damage_ui_text.Disable()
	})

	if p.health.GetCurrent() <= 0 {
		// events.Bus.Publish(events.GameOverEvent{
		// 	Message: "Player is dead",
		// })

		lifecycle.Kill()
	}
}

func (p *Player) defend(enable bool) {
	if enable {
		p.is_defending = true
		p.sprite.Disable()
		p.defense_start_time = time.Now()
	} else {
		p.is_defending = false
		p.sprite.Enable()
	}
}

func (p *Player) negate_damage(damage int) {
	hold_time := int(time.Since(p.defense_start_time).Milliseconds())

	break_roll := rand.IntN(100)
	break_chance := (damage*100)/p.health.GetMax() + (max(hold_time, 1) / 100)

	print(fmt.Sprintf(
		config.Blue+"Player absorbed %v damange\nDefense held for %v.\nBreak roll: %v | Break chance: %v\n"+config.Reset,
		damage,
		hold_time,
		break_roll,
		break_chance,
	))

	if break_roll <= break_chance {

		if p.is_concentrating {
			p.lost_concentration(true)
		}

		if p.current_weapon.Type != data.Fire {
			p.is_stunned = true
			print(fmt.Sprintf("%s", config.Yellow+"Player is Stunned\n"+config.Reset))

			time.AfterFunc(time.Millisecond*time.Duration(PLAYER_STUN_DELAY), func() {
				p.is_stunned = false
			})
		}
	}
}

func (p *Player) level_up_listener(s int) {
	if !p.can_level_up {
		return
	}

	level_up := skill.Skill{}
	if s == 1 {
		level_up.STR = 1
	}
	if s == 2 {
		level_up.INT = 1
	}
	if s == 3 {
		level_up.SPD = 1
	}

	p.skills.LevelUp(level_up)
	p.can_level_up = false
	p.health.ChangeMax(p.health.GetMax() + (p.max_hp / 2))
	p.health.Reset()

	println(fmt.Sprintf("Player Current Level: %d {STR: %d, INT: %d, SPD: %v} \n",
		p.skills.GetLevel(),
		p.skills.STR,
		p.skills.INT,
		p.skills.SPD))
}

func (p *Player) get_recovery() int {
	return p.current_weapon.Recovery / p.skills.SPD
}

func (p *Player) lost_concentration(fire_event bool) {
	p.concentration_count = 0
	p.is_concentrating = false

	if fire_event {
		println(config.Yellow + "Player is no longer concentrating and will publish a break event" + config.Reset)
		gomesevents.Emit(gomesevents.Game, events.PlayerBreakSpellEvent{})
	}
}
