package combat

import (
	"gates/internal/actors"
	data_enemies "gates/internal/data/enemies"
	data_weapons "gates/internal/data/weapons"
	"gates/pkg/skill"
	"math/rand/v2"
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
var enemy_pool []data_enemies.EnemySpecs

func Init() {
	player = actors.NewPlayer()
	player.Disable()

	enemy = actors.NewEnemy()
	enemy.Disable()

	enemy_pool = []data_enemies.EnemySpecs{
		data_enemies.Rat,
		data_enemies.Wolf,
		data_enemies.Zombie,
		data_enemies.Goblin,
		data_enemies.Skeleton,
		data_enemies.Bandit,
		data_enemies.Orc,
		data_enemies.Werewolf,
		data_enemies.Vampire,
	}
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
	pool := []data_enemies.EnemySpecs{}
	for _, e := range enemy_pool {
		if e.MinLevel <= player.GetLevel() {
			pool = append(pool, e)
		}
	}

	index := rand.IntN(len(pool))
	enemy.LoadData(pool[index])
}
