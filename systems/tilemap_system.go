package systems

import (
	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type Tile struct {
	Coord math.Vector2
	Rect  utils.RectSpecs
}

type TileMap struct {
	Tiles      []Tile
	instance   *lifecycle.GameObject
	can_render bool
}

func NewTileMap(map_size math.Vector2, tile_rect utils.RectSpecs, offset int) *TileMap {
	tiles := make([]Tile, 0)
	for x := range map_size.X {
		for y := range map_size.Y {
			rect := tile_rect
			rect.PosX += x * (tile_rect.Width + offset)
			rect.PosY += y * (tile_rect.Height + offset)

			tiles = append(tiles, Tile{
				Coord: math.Vector2{X: x, Y: y},
				Rect:  rect,
			})
		}
	}

	tilemap := &TileMap{
		Tiles:      tiles,
		can_render: false,
	}

	tilemap.instance = lifecycle.Register(&lifecycle.GameObject{
		Render: tilemap.render,
	})

	return tilemap
}

func (m *TileMap) DrawMap() {
	m.can_render = true
}

func (m *TileMap) render() {
	if !m.can_render {
		return
	}

	for _, t := range m.Tiles {
		render.DrawRect(t.Rect, render.Green)
	}
}
