package scenes

import (
	"gates/systems"
	"gates/values"

	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

func RunMap() {
	map_file := "assets/images/maps/map3x3.png"

	scale := 128
	offset := 0
	size := math.Vector2{X: 3, Y: 3}
	rect := utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - ((scale * size.X) / 2) - ((offset * size.X) / 2),
		PosY:   (values.SCREEN_SIZE.Y / 2) - ((scale * size.Y) / 2) - ((offset * size.Y) / 2),
		Width:  scale,
		Height: scale,
	}

	rules := []systems.TileRules{
		{
			Chan:      systems.R,
			ChanValue: 255,
			Color:     render.Red,
		},
		{
			Chan:      systems.G,
			ChanValue: 255,
			Color:     render.Yellow,
		},
		{
			Chan:       systems.B,
			ChanValue:  255,
			SpritePath: "assets/images/sprites/tiles/tile_simple.png",
		},
		{
			Chan:       systems.B,
			ChanValue:  0,
			SpritePath: "assets/images/sprites/tiles/wall_simple.png",
		},
	}

	tilemap := systems.NewTileMap(size, rect, offset)
	tilemap.DrawMapAssetsFromFile(rules, map_file)
}
