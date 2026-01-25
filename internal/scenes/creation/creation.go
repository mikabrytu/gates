package creation

import (
	"fmt"
	"gates/config"
	data "gates/internal/data/weapons"
	"gates/internal/events"
	"gates/pkg/skill"

	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
)

type CreationState int

const (
	Race CreationState = iota
	Class
	Equipment
)

var state CreationState
var fonts [6]*render.Font
var character skill.Skill

func Init() {
	state = Race
	character = skill.Skill{}

	register_events()
	init_fonts()
}

func Show() {
	if state == Race {
		prepare_race_text()
	} else {
		show_text(nil)
	}
}

func Hide() {
	hide_text()
}

func register_events() {
	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_1, func(data any) {
		option_listener(1)
	})

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_2, func(data any) {
		option_listener(2)
	})

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_3, func(data any) {
		option_listener(3)
	})

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_4, func(data any) {
		option_listener(4)
	})

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_5, func(data any) {
		option_listener(5)
	})
}

func init_fonts() {
	for i := range len(fonts) {
		fonts[i] = render.NewFont(config.FONT_SPECS, config.SCREEN_SIZE)
		fonts[i].Init("0", render.White, math.Vector2{X: 0, Y: 0})
		fonts[i].Disable()
	}
}

func option_listener(option int) {
	if state != Equipment {
		switch option {
		case 1:
			character.STR += 1
		case 2:
			character.INT += 1
		case 3:
			character.SPD += 1
		}
	}

	switch state {
	case Race:
		state = Class
		prepare_class_text()
	case Class:
		state = Equipment
		prepare_weapon_text()
	case Equipment:
		message := fmt.Sprintf(
			"Character Created! Attributes: { STR: %v, INT: %v, SPD: %v }\n",
			character.STR,
			character.INT,
			character.SPD,
		)
		print(config.Yellow + message + config.Reset)

		var weapon data.Weapon = data.Weapon{}
		switch option {
		case 1:
			weapon = data.Sword
		case 2:
			weapon = data.Bow
		case 3:
			weapon = data.FireSpell
		case 4:
			weapon = data.IceSpell
		case 5:
			weapon = data.ShockSpell
		}

		gomesevents.Emit(gomesevents.Game, events.SceneChangeEvent{
			Scene: config.SCENE_CREATION,
			Data: []any{
				character,
				weapon,
			},
		})
	}
}

func prepare_race_text() {
	messages := []string{
		"Choose your race",
		"1 - Dwarf",
		"2 - Human",
		"3 - Elf",
	}

	show_text(messages)
}

func prepare_class_text() {
	messages := []string{
		"Choose your class",
		"1 - Warrior",
		"2 - Wizard",
		"3 - Rogue",
	}

	show_text(messages)
}

func prepare_weapon_text() {
	messages := []string{
		"Choose your weapon",
		"1 - Sword",
		"2 - Bow",
		"3 - Fire Spell",
		"4 - Ice Spell",
		"5 - Shock Spell",
	}

	show_text(messages)
}

func show_text(messages []string) {
	if messages == nil {
		for i := range len(fonts) {
			fonts[i].Enable()
		}

		return
	}

	if len(messages) > len(fonts) {
		panic("More text than fonts available to show")
	}

	for i, m := range messages {
		fonts[i].Enable()
		fonts[i].UpdateText(m)
		fonts[i].AlignText(render.TopLeft, math.Vector2{X: 16, Y: 16 + (i * 32)})
	}
}

func hide_text() {
	for _, f := range fonts {
		f.Disable()
	}
}
