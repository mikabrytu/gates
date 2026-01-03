package actors

import (
	"gates/values"

	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type EnemySpecs struct {
	Name            string
	Image_Path      string
	Size            int
	HP              int
	Attack_Interval int
	Attack_Damage   int
	Defense         int
}

var enemy_specs EnemySpecs
var enemy_hp_rect utils.RectSpecs

func Enemy() {
	enemy_init()

	lifecycle.Register(&lifecycle.GameObject{
		Render: func() {
			render.DrawRect(enemy_hp_rect, render.Red)
		},
	})
}

func LoadEnemy(specs EnemySpecs) {
	enemy_specs = specs
}

func enemy_init() {
	message := "Initializing enemy"
	println(message)

	rect := utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - (enemy_specs.Size / 2),
		PosY:   32,
		Width:  enemy_specs.Size,
		Height: enemy_specs.Size,
	}

	enemy_hp_rect = rect
	enemy_hp_rect.PosY -= 24
	enemy_hp_rect.Height = 16
}
