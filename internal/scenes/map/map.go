package gamemap

import (
	"gates/config"
	"gates/internal/data/items"
	data "gates/internal/data/map_rules"
	"gates/internal/events"
	"gates/pkg/tilemap"

	gomesevents "github.com/mikabrytu/gomes-engine/events"
	"github.com/mikabrytu/gomes-engine/lifecycle"
	gomesmath "github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

var tmap *tilemap.TileMap
var wall_instance *lifecycle.GameObject
var player_instance *lifecycle.GameObject
var player_rect utils.RectSpecs
var player_coord gomesmath.Vector2
var player_items []items.Item
var enabled bool

var MAP_SIZE = gomesmath.Vector2{X: 28, Y: 22}

const SCALE int = 42

func Init() {
	player_items = make([]items.Item, 0)

	drawMap()
	init_player()
}

func Show() {
	enabled = true
	tmap.Enable()
	lifecycle.Enable(wall_instance)
	lifecycle.Enable(player_instance)
}

func Hide() {
	enabled = false
	tmap.Disable()
	lifecycle.Disable(wall_instance)
	lifecycle.Disable(player_instance)
}

func drawMap() {
	map_file := config.MAP_DEMO_FILE

	offset := 0
	rect := utils.RectSpecs{
		PosX:   (config.SCREEN_SIZE.X / 2) - ((SCALE * MAP_SIZE.X) / 2) - ((offset * MAP_SIZE.X) / 2),
		PosY:   (config.SCREEN_SIZE.Y / 2) - ((SCALE * MAP_SIZE.Y) / 2) - ((offset * MAP_SIZE.Y) / 2),
		Width:  SCALE,
		Height: SCALE,
	}

	tmap = tilemap.NewTileMap(MAP_SIZE, rect, offset)
	tmap.DrawMapAssetsFromFile(data.DEMO_RULES, map_file)
	tmap.Disable()

	wall_color := render.Color{R: 155, G: 173, B: 183, A: 255}
	map_walls := []utils.RectSpecs{
		{
			PosX:   (config.SCREEN_SIZE.X / 2) - ((SCALE * MAP_SIZE.X) / 2) - ((offset * MAP_SIZE.X) / 2) - 4,
			PosY:   (config.SCREEN_SIZE.Y / 2) - ((SCALE * MAP_SIZE.Y) / 2) - ((offset * MAP_SIZE.Y) / 2),
			Width:  4,
			Height: SCALE * (MAP_SIZE.Y + offset),
		},
		{
			PosX:   (config.SCREEN_SIZE.X / 2) + ((SCALE * MAP_SIZE.X) / 2) + ((offset * MAP_SIZE.X) / 2) + 4,
			PosY:   (config.SCREEN_SIZE.Y / 2) - ((SCALE * MAP_SIZE.Y) / 2) - ((offset * MAP_SIZE.Y) / 2),
			Width:  4,
			Height: SCALE * (MAP_SIZE.Y + offset),
		},
		{
			PosX:   (config.SCREEN_SIZE.X / 2) - ((SCALE * MAP_SIZE.X) / 2) - ((offset * MAP_SIZE.X) / 2) - 4,
			PosY:   (config.SCREEN_SIZE.Y / 2) - ((SCALE * MAP_SIZE.Y) / 2) - ((offset * MAP_SIZE.Y) / 2) - 4,
			Width:  SCALE*(MAP_SIZE.X+offset) + 12,
			Height: 4,
		},
		{
			PosX:   (config.SCREEN_SIZE.X / 2) - ((SCALE * MAP_SIZE.X) / 2) - ((offset * MAP_SIZE.X) / 2) - 4,
			PosY:   (config.SCREEN_SIZE.Y / 2) + ((SCALE * MAP_SIZE.Y) / 2) + ((offset * MAP_SIZE.Y) / 2),
			Width:  SCALE*(MAP_SIZE.X+offset) + 12,
			Height: 4,
		},
	}

	wall_instance = lifecycle.Register(&lifecycle.GameObject{
		Render: func() {
			for _, wall := range map_walls {
				render.DrawRect(wall, wall_color)
			}
		},
	})
	lifecycle.Disable(wall_instance)
}

func init_player() {
	var tile *tilemap.Tile = nil
	for _, row := range tmap.Tiles {
		for _, t := range row {
			if t.IsSpawn {
				tile = t
			}
		}
	}

	if tile == nil {
		panic("Player spawn point not found")
	}

	player_coord = tile.Coord
	player_rect = utils.RectSpecs{
		PosX:   tile.Rect.PosX + (SCALE / 4),
		PosY:   tile.Rect.PosY + (SCALE / 4),
		Width:  SCALE / 2,
		Height: SCALE / 2,
	}

	show_adjacent()

	player_instance = lifecycle.Register(&lifecycle.GameObject{
		Start: func() {
			gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_W, func(data any) {
				move_player(gomesmath.Vector2{X: 0, Y: -1})
			})
			gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_S, func(data any) {
				move_player(gomesmath.Vector2{X: 0, Y: 1})
			})
			gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_D, func(data any) {
				move_player(gomesmath.Vector2{X: 1, Y: 0})
			})
			gomesevents.Subscribe(gomesevents.Input, gomesevents.INPUT_KEYBOARD_PRESSED_A, func(data any) {
				move_player(gomesmath.Vector2{X: -1, Y: 0})
			})
		},
		Render: func() {
			render.DrawRect(player_rect, render.Green)

			tile := tmap.Tiles[player_coord.Y][player_coord.X]
			if tile.HasEnemy || tile.HasItem {
				tile.Enabled = false
			}
		},
	})

	lifecycle.Disable(player_instance)
}

func move_player(coord gomesmath.Vector2) {
	if !enabled {
		return
	}

	if (player_coord.X+coord.X) >= MAP_SIZE.X ||
		(player_coord.X+coord.X) < 0 ||
		(player_coord.Y+coord.Y) >= MAP_SIZE.Y ||
		(player_coord.Y+coord.Y) < 0 {
		return
	}

	tile := tmap.Tiles[player_coord.Y+coord.Y][player_coord.X+coord.X]

	if !tile.IsWalkable {
		unlock := false
		if tile.Item.Type == items.Lock {
			for _, i := range player_items {
				if tile.Item.Link == i.ID {
					unlock = true
					tile.IsWalkable = true
					break
				}
			}
		}

		if !unlock {
			return
		}
	}

	player_coord.X += coord.X
	player_coord.Y += coord.Y

	player_rect.PosX = tile.Rect.PosX + (SCALE / 4)
	player_rect.PosY = tile.Rect.PosY + (SCALE / 4)

	show_adjacent()

	if tile.HasEnemy {
		tile.HasEnemy = false
		gomesevents.Emit(gomesevents.Game, events.SceneChangeEvent{
			Scene: config.SCENE_MAP,
		})
	}

	if tile.HasItem {
		tile.HasItem = false

		if tile.Item.Type != items.Lock {
			player_items = append(player_items, tile.Item)
		}
	}
}

func show_adjacent() {
	adjacents := []gomesmath.Vector2{
		{X: player_coord.X, Y: player_coord.Y - 1},
		{X: player_coord.X + 1, Y: player_coord.Y - 1},
		{X: player_coord.X - 1, Y: player_coord.Y - 1},
		{X: player_coord.X - 1, Y: player_coord.Y},
		{X: player_coord.X + 1, Y: player_coord.Y},
		{X: player_coord.X, Y: player_coord.Y + 1},
		{X: player_coord.X + 1, Y: player_coord.Y + 1},
		{X: player_coord.X - 1, Y: player_coord.Y + 1},
	}

	for _, coord := range adjacents {
		if coord.X >= MAP_SIZE.X ||
			coord.X < 0 ||
			coord.Y >= MAP_SIZE.Y ||
			coord.Y < 0 {
			continue
		}

		tile := tmap.Tiles[coord.Y][coord.X]
		tile.Enabled = true
		if tile.Sprite != nil {
			tile.Sprite.Enable()
		}
	}
}
