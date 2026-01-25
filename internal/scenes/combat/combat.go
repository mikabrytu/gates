package combat

import (
	"gates/internal/actors"
	data_enemies "gates/internal/data/enemies"
	data_weapons "gates/internal/data/weapons"
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
var enemy *actors.Enemy

func Init() {
	player = actors.NewPlayer()
	player.Disable()

	enemy = actors.NewEnemy()
	enemy.Disable()
}

func Show() {
	player.Enable()
	enemy.Enable()
}

func Hide() {
	player.Disable()
	enemy.Disable()
}

func LoadPlayerData(skills skill.Skill, weapon data_weapons.Weapon) {
	player.LoadData(weapon, skills)
}

func LoadEnemy() {
	enemy.LoadData(data_enemies.Rat)
}
