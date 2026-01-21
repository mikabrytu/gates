package scenes

import (
	"gates/systems"
	"gates/values"

	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var tilemap *systems.TileMap

var MAP_SIZE = math.Vector2{X: 3, Y: 3}

const SCALE int = 128

func RunMap() {
	drawMap()
	initPlayer()
}

func drawMap() {
	map_file := "assets/images/maps/map3x3.png"

	offset := 1
	rect := utils.RectSpecs{
		PosX:   (values.SCREEN_SIZE.X / 2) - ((SCALE * MAP_SIZE.X) / 2) - ((offset * MAP_SIZE.X) / 2),
		PosY:   (values.SCREEN_SIZE.Y / 2) - ((SCALE * MAP_SIZE.Y) / 2) - ((offset * MAP_SIZE.Y) / 2),
		Width:  SCALE,
		Height: SCALE,
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

	tilemap = systems.NewTileMap(MAP_SIZE, rect, offset)
	tilemap.DrawMapAssetsFromFile(rules, map_file)
}

func initPlayer() {
	start_tile_rect := tilemap.Tiles[MAP_SIZE.Y-1][MAP_SIZE.X/2].Rect
	rect := utils.RectSpecs{
		PosX:   start_tile_rect.PosX + (SCALE / 4),
		PosY:   start_tile_rect.PosY + (SCALE / 4),
		Width:  SCALE / 2,
		Height: SCALE / 2,
	}

	lifecycle.Register(&lifecycle.GameObject{
		Render: func() {
			render.DrawRect(rect, render.Green)
		},
	})
}
