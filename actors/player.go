package actors

import (
	"fmt"
	"gates/actors/weapons"
	game_events "gates/events"
	"gates/spells"
	"gates/systems"
	utils1 "gates/utils"
	"gates/values"
	"math/rand/v2"
	"time"

	"github.com/Papiermond/eventbus"
	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var player_sprite *render.Sprite
var player_health *systems.Health
var player_skills *systems.Skill
var player_damage_ui_text *render.Font
var player_current_weapon weapons.Weapon
var player_weapon_rect utils.RectSpecs
var player_hp_rect utils.RectSpecs
var player_recovery_rect utils.RectSpecs
var player_anim_position utils.RectSpecs
var player_defense_rect utils.RectSpecs
var player_attack_start_time time.Time
var player_defense_start_time time.Time
var player_max_hp int
var player_max_hp_width int
var player_concentration_count = 0
var player_can_attack = true
var player_can_level_up = false
var player_is_attacking = false
var player_is_animating = false
var player_is_defending = false
var player_is_stunned = false
var player_is_concentrating = false

const PLAYER_HP_PER_LEVEL int = 5
const PLAYER_STUN_DELAY int = 6000
const PLAYER_SPELL_EFFECT_CHANCE = 50

func Player() {
	player_init()

	og_recovery_width := player_recovery_rect.Width

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			player_sprite.Init()

			// Input
			events.Subscribe(events.Input, events.INPUT_MOUSE_CLICK_DOWN, func(data any) {
				click := data.(events.InputMouseClickDownEvent)

				if player_is_animating || player_is_defending || player_is_stunned {
					println("Player cannot act now")
					return
				}

				if click.Index.Left == 1 {
					player_attack_listener()
				}

				if click.Index.Right == 1 {
					player_defend(true)
				}
			})

			events.Subscribe(events.Input, events.INPUT_MOUSE_CLICK_UP, func(data any) {
				player_defend(false)
			})

			// Level Up
			events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_1, func(data any) {
				player_level_up_listener(1)
			})

			events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_2, func(data any) {
				player_level_up_listener(2)
			})

			events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_3, func(data any) {
				player_level_up_listener(3)
			})

			// Take Damage
			game_events.Bus.Subscribe(game_events.ENEMY_ATTACK_EVENT, func(e eventbus.Event) {
				damage := e.(game_events.EnemyAttackEvent).Damage
				player_take_damage_listener(int(damage))
			})

			game_events.Bus.Subscribe(game_events.ENEMY_DEAD_EVENT, func(e eventbus.Event) {
				player_lost_concentration(false)
			})

			// Enemy
			events.Subscribe(events.Game, game_events.ENEMY_BREAK_PARALYSIS_EVENT, func(data any) {
				println("Enemy is trying to break free from Paralysis")
				player_lost_concentration(false)
			})
		},
		Update: func() {
			player_hp_rect.Width = (player_max_hp_width * player_health.GetCurrent()) / player_health.GetMax()

			if player_is_attacking {
				elapsed := time.Since(player_attack_start_time).Milliseconds()
				t := float64(elapsed) / float64(player_get_recovery())

				width := utils1.Lerp(float64(og_recovery_width), 0, t)
				player_recovery_rect.Width = int(width)

				if t > 1 {
					t = 1
					player_is_attacking = false
					player_recovery_rect.Width = og_recovery_width
				}
			}
		},
		Render: func() {
			player_sprite.UpdateRect(player_anim_position)

			if player_is_defending {
				render.DrawRect(player_defense_rect, render.White)
			}

			if player_is_stunned {
				render.DrawRect(player_hp_rect, render.Yellow)
				render.DrawRect(player_recovery_rect, render.Yellow)
				return
			}

			render.DrawRect(player_hp_rect, render.Green)
			render.DrawRect(player_recovery_rect, render.Blue)
		},
	})
}

func PlayerLevelUp() {
	player_can_level_up = true
}

func PlayerLoadWeapon(weapon weapons.Weapon) {
	player_current_weapon = weapon
	println(fmt.Sprintf("Player weapon loaded: %s", weapon.Name))
}

func player_init() {
	player_max_hp_width = 512

	player_weapon_rect = utils.RectSpecs{
		PosX:   values.SCREEN_SIZE.X - player_current_weapon.SpriteSize.X + player_current_weapon.SpriteOffset.X,
		PosY:   values.SCREEN_SIZE.Y - player_current_weapon.SpriteSize.Y + player_current_weapon.SpriteOffset.Y,
		Width:  player_current_weapon.SpriteSize.X,
		Height: player_current_weapon.SpriteSize.Y,
	}
	player_anim_position = player_weapon_rect

	player_hp_rect = utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - (player_max_hp_width / 2),
		PosY:   values.SCREEN_SIZE.Y - 40,
		Width:  player_max_hp_width,
		Height: 24,
	}

	player_defense_rect = utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - 256,
		PosY:   values.SCREEN_SIZE.Y - (values.SCREEN_SIZE.Y / 3),
		Width:  512,
		Height: 512,
	}

	player_recovery_rect = player_hp_rect
	player_recovery_rect.PosY -= 24
	player_recovery_rect.Height = 12

	player_sprite = render.NewSprite(
		"Player Weapon",
		player_current_weapon.SpritePath,
		player_weapon_rect,
		render.Transparent,
	)

	//player_max_hp = 1000
	player_skills = systems.NewSkill()
	player_max_hp = PLAYER_HP_PER_LEVEL * player_skills.GetTotalSkillPoints()
	player_health = systems.InitHealth(player_max_hp)

	print(fmt.Sprintf("Player initialized with %d health\n", player_health.GetCurrent()))

	if player_damage_ui_text == nil {
		player_damage_ui_text = render.NewFont(values.FONT_SPECS, values.SCREEN_SIZE)
		player_damage_ui_text.Init("10", render.Transparent, math.Vector2{X: 0, Y: 0})
		player_damage_ui_text.AlignText(render.BottomCenter, math.Vector2{X: 0, Y: 96})
	}
}

func player_damage() int {
	return player_current_weapon.Damage * player_skills.STR
}

func player_spell_effect() spells.Effect {
	if player_current_weapon.Type == weapons.Physical {
		return spells.Effect{}
	}

	if player_concentration_count >= player_skills.INT {
		println("Player cannot active more spells until concentration break")
		return spells.Effect{}
	}

	roll := rand.IntN(100)

	if roll <= PLAYER_SPELL_EFFECT_CHANCE {
		player_is_concentrating = true
		player_concentration_count += 1

		var effect_type spells.EffectType
		name := ""
		switch player_current_weapon.Type {
		case weapons.Fire:
			effect_type = spells.Burn
			name = "Burn"
		case weapons.Ice:
			effect_type = spells.Cold
			name = "Cold"
		case weapons.Shock:
			effect_type = spells.Paralysis
			name = "Paralysis"
		}

		print(fmt.Sprintf(values.Yellow+"Player triggered spell effect: %v. Max Stack: %v\n"+values.Reset, name, player_skills.INT))

		return spells.Effect{
			Type:  effect_type,
			Stack: 1,
		}
	}

	return spells.Effect{}
}

func player_attack_listener() {
	if !player_can_attack {
		return
	}

	player_can_attack = false
	player_is_attacking = true
	player_is_animating = true
	player_attack_start_time = time.Now()

	game_events.Bus.Publish(game_events.PlayerAttackEvent{
		Damage: player_damage(),
		Effect: player_spell_effect(),
	})

	temp_rect := player_weapon_rect
	temp_rect.PosX = (values.SCREEN_SIZE.X / 2) - (player_current_weapon.SpriteSize.X / 2)
	player_anim_position = temp_rect

	// Recovery delay
	time.AfterFunc(time.Millisecond*time.Duration(player_get_recovery()), func() {
		player_can_attack = true
	})

	// Weapon animation
	time.AfterFunc(time.Millisecond*350, func() {
		player_is_animating = false
		player_anim_position = player_weapon_rect
	})
}

func player_take_damage_listener(base_damage int) {
	damage := max(base_damage/player_skills.STR, 1)

	if player_is_defending {
		player_negate_damage(base_damage)
		return
	}

	message := values.Red + fmt.Sprintf("Enemy attacks with %d damage", damage) + values.Reset
	println(message)

	player_health.TakeDamage(damage)

	player_damage_ui_text.UpdateText(fmt.Sprint(damage))
	player_damage_ui_text.UpdateColor(render.White)
	time.AfterFunc(time.Millisecond*1200, func() {
		player_damage_ui_text.UpdateColor(render.Transparent)
	})

	if player_health.GetCurrent() <= 0 {
		// game_events.Bus.Publish(game_events.GameOverEvent{
		// 	Message: "Player is dead",
		// })

		lifecycle.Kill()
	}
}

func player_defend(enable bool) {
	if enable {
		player_is_defending = true
		player_sprite.Disable()
		player_defense_start_time = time.Now()
	} else {
		player_is_defending = false
		player_sprite.Enable()
	}
}

func player_negate_damage(damage int) {
	hold_time := int(time.Since(player_defense_start_time).Milliseconds())

	break_roll := rand.IntN(100)
	break_chance := (damage*100)/player_health.GetMax() + (max(hold_time, 1) / 100)

	print(fmt.Sprintf(
		values.Blue+"Player absorbed %v damange\nDefense held for %v.\nBreak roll: %v | Break chance: %v\n"+values.Reset,
		damage,
		hold_time,
		break_roll,
		break_chance,
	))

	if break_roll <= break_chance {

		if player_is_concentrating {
			player_lost_concentration(true)
		}

		if player_current_weapon.Type != weapons.Fire {
			player_is_stunned = true
			print(fmt.Sprintf("%s", values.Yellow+"Player is Stunned\n"+values.Reset))

			time.AfterFunc(time.Millisecond*time.Duration(PLAYER_STUN_DELAY), func() {
				player_is_stunned = false
			})
		}
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
	player_health.ChangeMax(player_health.GetMax() + (player_max_hp / 2))
	player_health.Reset()

	println(fmt.Sprintf("Player Current Level: %d {STR: %d, INT: %d, SPD: %v} \n",
		player_skills.GetLevel(),
		player_skills.STR,
		player_skills.INT,
		player_skills.SPD))
}

func player_get_recovery() int {
	return player_current_weapon.Recovery / player_skills.SPD
}

func player_lost_concentration(fire_event bool) {
	player_concentration_count = 0
	player_is_concentrating = false

	if fire_event {
		println(values.Yellow + "Player is no longer concentrating and will publish a break event" + values.Reset)
		events.Emit(events.Game, game_events.PlayerBreakSpellEvent{})
	}
}

func old_math() {
	// damage := 0
	// base := player_skills.STR * player_current_weapon.Damage
	// raw_damage := utils1.CalcDamange(base, base/2)

	// //print(fmt.Sprintf("%sPlayer Raw Damage of %d %s\n", values.Green, raw_damage, values.Reset))

	// crit_chance := player_skills.SPD * 10 // TODO: Set this multiplier somewhere else
	// crit_index := rand.IntN(100)
	// crit_hit := crit_index <= crit_chance

	// //print(fmt.Sprintf("%sCritical Chance of %d. Dice rolled %d. Was a crit? %v %s\n", values.Green, crit_chance, crit_index, crit_hit, values.Reset))

	// crit_damage := 0
	// if crit_hit {
	// 	crit_damage = (raw_damage * (player_skills.INT * 25)) / 100 // TODO: Set this multiplier somewhere else

	// 	print(fmt.Sprintf("%sCRITICAL HIT! Player is dealing additional %d damage%s\n", values.Green, crit_damage, values.Reset))
	// }

	// damage = raw_damage + crit_damage
	// print(fmt.Sprintf("%sPlayer is attacking with %d damage %s\n", values.Green, damage, values.Reset))
	// return damage
}
