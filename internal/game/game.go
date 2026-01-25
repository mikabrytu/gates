package game

import (
	"gates/config"
	weapons "gates/internal/data/weapons"
	"gates/internal/events"
	"gates/internal/scenes/combat"
	"gates/internal/scenes/creation"
	levelup "gates/internal/scenes/level_up"
	gamemap "gates/internal/scenes/map"
	"gates/pkg/skill"
	"time"

	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/render"
)

var player_skills skill.Skill
var player_weapon weapons.Weapon

func Init() {
	register_events()

	render.SetBackgroundColor(render.Color{R: 25, G: 20, B: 43, A: 255})

	creation.Init()
	gamemap.Init()
	combat.Init()
	levelup.Init()

	time.AfterFunc(time.Millisecond*200, func() {
		creation.Show()
		gamemap.Hide()
		combat.Hide()
		levelup.Hide()
	})
}

func register_events() {
	gomesevents.Subscribe(gomesevents.Game, events.SCENE_CHANGE_EVENT, func(data any) {
		to_close := data.(events.SceneChangeEvent)

		if to_close.Scene == config.SCENE_CREATION {
			player_skills = to_close.Data[0].(skill.Skill)
			player_weapon = to_close.Data[1].(weapons.Weapon)
		}

		change_scene(to_close.Scene)
	})
}

func change_scene(to_close string) {
	switch to_close {
	case config.SCENE_CREATION:
		creation.Hide()
		gamemap.Show()

	case config.SCENE_MAP:
		gamemap.Hide()

		combat.LoadPlayerData(player_skills, player_weapon)
		combat.Show()

	case config.SCENE_COMBAT:
		combat.Hide()

	case config.SCENE_LEVEL_UP:
		levelup.Hide()
	}
}
