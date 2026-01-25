package combat

import (
	"gates/internal/actors"
	data "gates/internal/data/weapons"
	"gates/pkg/skill"
)

type GameState int

const (
	Running GameState = iota
	Preparing
	Waiting
	Stopped
)

var player *actors.Player

func Init() {
	player = actors.NewPlayer()
	player.Disable()
}

func Show() {
	player.Enable()
}

func Hide() {
	player.Disable()
}

func LoadPlayerData(skills skill.Skill, weapon data.Weapon) {
	player.LoadData(weapon, skills)
}
