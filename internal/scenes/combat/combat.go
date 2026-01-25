package combat

import (
	"fmt"
	"gates/config"
	"gates/internal/actors"
	data_enemies "gates/internal/data/enemies"
	"gates/internal/events"

	"github.com/Papiermond/eventbus"
	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
)

type GameState int

const (
	Running GameState = iota
	Preparing
	Waiting
	Stopped
)

var game_state GameState

var skills_fonts [4]*render.Font
var continue_font *render.Font
var rounds int

func Init() {

}

func Show() {

}

func Hide() {

}

func RunCombat() {
	listeners()
	game_state = Preparing
}

func listeners() {
	// Engine events

	gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_SPACE, func(data any) {
		sequence()
	})

	// Game events
	events.Bus.Subscribe(events.GAME_OVER_EVENT, func(e eventbus.Event) {
		println("Player is dead. Game over")
		game_state = Stopped
		lifecycle.Kill()
	})

	events.Bus.Subscribe(events.ENEMY_DEAD_EVENT, func(e eventbus.Event) {
		if game_state != Running {
			println(config.Yellow + "Trying to kill an enemy while game is not running. Current state:" + fmt.Sprint(game_state) + config.Reset)
			return
		}

		println("MAIN::" + e.(events.EnemyDeadEvent).Message)

		game_state = Waiting
		rounds += 1

		// Player Level Up
		if rounds == 1 ||
			rounds == 4 ||
			rounds == 7 ||
			rounds == 10 ||
			rounds == 13 ||
			rounds == 16 ||
			rounds == 19 ||
			rounds == 22 ||
			rounds == 25 {

			actors.PlayerLevelUp()
			show_level_up_text()
		} else {
			show_continue_message()
		}
	})
}

func sequence() {
	println("Sequence called")

	if game_state == Preparing {
		hide_ui_text()

		actors.Player()
		actors.LoadEnemy(data_enemies.Rat)
		actors.Enemy()

		game_state = Running
		println("Starting game...")
	}

	if game_state == Waiting {
		hide_ui_text()

		if rounds == 3 {
			actors.LoadEnemy(data_enemies.Wolf)
		}

		if rounds == 6 {
			actors.LoadEnemy(data_enemies.Zombie)
		}

		if rounds == 9 {
			actors.LoadEnemy(data_enemies.Goblin)
		}

		if rounds == 12 {
			actors.LoadEnemy(data_enemies.Skeleton)
		}

		if rounds == 15 {
			actors.LoadEnemy(data_enemies.Bandit)
		}

		if rounds == 18 {
			actors.LoadEnemy(data_enemies.Orc)
		}

		if rounds == 21 {
			actors.LoadEnemy(data_enemies.Werewolf)
		}

		if rounds == 24 {
			actors.LoadEnemy(data_enemies.Vampire)
		}

		game_state = Running
		events.Bus.Publish(events.GameRestartEvent{
			Message: "Restarting game!",
		})
	}
}

func show_level_up_text() {
	messages := []string{"LEVEL UP. Choose a skill to increase", "1 - STR", "2 - INT", "3 - SPD"}

	for i, m := range messages {
		if skills_fonts[i] == nil {
			skills_fonts[i] = render.NewFont(config.FONT_SPECS, config.SCREEN_SIZE)
		}

		skills_fonts[i].Init(m, render.Yellow, math.Vector2{X: 0, Y: 0})
		skills_fonts[i].AlignText(render.TopLeft, math.Vector2{X: 16, Y: 16 + (i * 32)})

		if !skills_fonts[i].IsEnable() {
			skills_fonts[i].Enable()
		}
	}
}

func hide_ui_text() {
	for _, f := range skills_fonts {
		if f == nil {
			continue
		}

		f.Disable()
	}

	if continue_font == nil {
		return
	}

	continue_font.Disable()
}

func show_continue_message() {
	if continue_font == nil {
		continue_font = render.NewFont(config.FONT_SPECS, config.SCREEN_SIZE)
		continue_font.Init("Enemy is dead. Press SPACE to continue...", render.Blue, math.Vector2{X: 0, Y: 0})
		continue_font.AlignText(render.TopLeft, math.Vector2{X: 16, Y: 16})
	}

	continue_font.Enable()
}
