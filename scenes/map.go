package scenes

import (
	"gates/systems"
	"gates/values"

	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/utils"
)

func RunMap() {
	offset := 2
	size := math.Vector2{X: 10, Y: 10}
	rect := utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - ((64 * size.X) / 2) - ((offset * size.X) / 2),
		PosY:   (values.SCREEN_SIZE.Y / 2) - ((64 * size.Y) / 2) - ((offset * size.Y) / 2),
		Width:  64,
		Height: 64,
	}

	tilemap := systems.NewTileMap(size, rect, offset)
	tilemap.DrawFromFile("assets/images/sprites/rgb_10.png")
}
