package game

import (
	"gates/internal/scenes/combat"
	"gates/internal/scenes/creation"
	gamemap "gates/internal/scenes/map"
	"time"

	"github.com/mikabrytu/gomes-engine/render"
)

func Init() {
	render.SetBackgroundColor(render.Color{R: 25, G: 20, B: 43, A: 255})

	creation.Init()
	combat.Init()
	gamemap.Init()

	time.AfterFunc(time.Millisecond*200, func() {
		creation.Show()
	})
}
