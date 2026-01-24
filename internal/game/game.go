package game

import (
	"gates/internal/scenes"

	"github.com/mikabrytu/gomes-engine/render"
)

func Init() {
	render.SetBackgroundColor(render.Color{R: 25, G: 20, B: 43, A: 255})
	scenes.RunCombat()
}
