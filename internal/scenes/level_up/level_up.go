package levelup

import (
	"gates/config"
	"gates/internal/events"

	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
)

var fonts [4]*render.Font
var enabled bool

func Init() {
	register_events()
	init_fonts()
}

func Show() {
	enabled = true
	show_level_up_text()
}

func Hide() {
	enabled = false
}

func register_events() {
	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_1, func(data any) {
		choice_listener(1)
	})

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_2, func(data any) {
		choice_listener(2)
	})

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_3, func(data any) {
		choice_listener(3)
	})
}

func init_fonts() {
	for i := range len(fonts) {
		fonts[i] = render.NewFont(config.FONT_SPECS, config.SCREEN_SIZE)
		fonts[i].Init("0", render.Yellow, math.Vector2{X: 0, Y: 0})
		fonts[i].Disable()
	}
}

func show_level_up_text() {
	messages := []string{
		"LEVEL UP. Choose a skill to increase",
		"1 - STR",
		"2 - INT",
		"3 - SPD",
	}

	for i, m := range messages {
		fonts[i].Enable()
		fonts[i].UpdateText(m)
		fonts[i].AlignText(render.TopLeft, math.Vector2{X: 16, Y: 16 + (i * 32)})
	}
}

func choice_listener(choice int) {
	// TODO: Save this on a player data container

	gomesevents.Emit(gomesevents.Game, events.SceneChangeEvent{
		Scene: config.SCENE_LEVEL_UP,
	})
}
