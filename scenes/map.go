package scenes

import (
	"gates/systems"
	"gates/values"

	"github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var tilemap *systems.TileMap
var player_rect utils.RectSpecs
var player_coord math.Vector2

var MAP_SIZE = math.Vector2{X: 3, Y: 3}

const SCALE int = 128

func RunMap() {
	drawMap()
	init_player()
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
			Walkable:   true,
		},
		{
			Chan:       systems.B,
			ChanValue:  0,
			SpritePath: "assets/images/sprites/tiles/wall_simple.png",
			Walkable:   false,
		},
	}

	tilemap = systems.NewTileMap(MAP_SIZE, rect, offset)
	tilemap.DrawMapAssetsFromFile(rules, map_file)
}

func init_player() {
	tile := tilemap.Tiles[MAP_SIZE.Y-1][MAP_SIZE.X/2]
	player_coord = tile.Coord
	player_rect = utils.RectSpecs{
		PosX:   tile.Rect.PosX + (SCALE / 4),
		PosY:   tile.Rect.PosY + (SCALE / 4),
		Width:  SCALE / 2,
		Height: SCALE / 2,
	}

	lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_W, func(data any) {
				move_player(math.Vector2{X: 0, Y: -1})
			})
			events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_S, func(data any) {
				move_player(math.Vector2{X: 0, Y: 1})
			})
			events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_D, func(data any) {
				move_player(math.Vector2{X: 1, Y: 0})
			})
			events.Subscribe(events.Input, events.INPUT_KEYBOARD_PRESSED_A, func(data any) {
				move_player(math.Vector2{X: -1, Y: 0})
			})
		},
		Render: func() {
			render.DrawRect(player_rect, render.Green)

			tile := tilemap.Tiles[player_coord.Y][player_coord.X]
			if tile.HasEnemy || tile.HasItem {
				tile.Enabled = false
			}
		},
	})
}

func move_player(coord math.Vector2) {
	if (player_coord.X+coord.X) >= MAP_SIZE.X ||
		(player_coord.X+coord.X) < 0 ||
		(player_coord.Y+coord.Y) >= MAP_SIZE.Y ||
		(player_coord.Y+coord.Y) < 0 {
		println("Trying to move out of bounds. Stopping player movement")
		return
	}

	tile := tilemap.Tiles[player_coord.Y+coord.Y][player_coord.X+coord.X]

	if !tile.IsWalkable {
		println("Player is trying to move to a wall. Stopping movement")
		return
	}

	player_coord.X += coord.X
	player_coord.Y += coord.Y

	player_rect.PosX = tile.Rect.PosX + (SCALE / 4)
	player_rect.PosY = tile.Rect.PosY + (SCALE / 4)
}
