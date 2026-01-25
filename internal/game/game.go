package game

import (
	"gates/config"
	"gates/internal/events"
	"gates/internal/scenes/combat"
	"gates/internal/scenes/creation"
	gamemap "gates/internal/scenes/map"
	"time"

	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/render"
)

func Init() {
	register_events()

	render.SetBackgroundColor(render.Color{R: 25, G: 20, B: 43, A: 255})

	creation.Init()
	gamemap.Init()
	combat.Init()

	time.AfterFunc(time.Millisecond*200, func() {
		creation.Show()
		gamemap.Hide()
		combat.Hide()
	})
}

func register_events() {
	gomesevents.Subscribe(gomesevents.Game, events.SCENE_CHANGE_EVENT, func(data any) {
		to_close := data.(events.SceneChangeEvent)
		change_scene(to_close.Scene)
	})
}

func change_scene(current string) {
	switch current {
	case config.SCENE_CREATION:
		creation.Hide()
		gamemap.Show()

	case config.SCENE_MAP:
		gamemap.Hide()
		combat.Show()

	case config.SCENE_COMBAT:
		combat.Hide()
	}
}
