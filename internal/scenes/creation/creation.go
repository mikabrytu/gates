package creation

import (
	"gates/config"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
)

type CreationState int

const (
	Race CreationState = iota
	Class
	Weapon
)

var state CreationState
var fonts [6]*render.Font

func Init() {
	state = Race

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
	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_1, func(data any) {
		option_listener(1)
	})

	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_2, func(data any) {
		option_listener(2)
	})

	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_3, func(data any) {
		option_listener(3)
	})

	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_4, func(data any) {
		option_listener(4)
	})

	events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_5, func(data any) {
		option_listener(5)
	})
}

func init_fonts() {
	for i := 0; i < len(fonts); i++ {
		fonts[i] = render.NewFont(config.FONT_SPECS, config.SCREEN_SIZE)
		fonts[i].Init("0", render.White, math.Vector2{X: 0, Y: 0})
		fonts[i].Disable()
	}
}

func option_listener(option int) {
	switch state {
	case Race:
		state = Class
		prepare_class_text()
	case Class:
		state = Weapon
		prepare_weapon_text()
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
