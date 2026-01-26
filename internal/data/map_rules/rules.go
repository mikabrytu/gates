package data

import (
	"gates/internal/data/items"
	"gates/pkg/tilemap"

	"github.com/mikabrytu/gomes-engine/render"
)

var RULES []tilemap.TileRules = []tilemap.TileRules{
	{
		Chan:       tilemap.B,
		ChanValue:  255,
		SpritePath: "",
		Walkable:   true,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  0,
		SpritePath: "assets/images/sprites/tiles/wall_base.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  1,
		SpritePath: "assets/images/sprites/tiles/wall_n.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  2,
		SpritePath: "assets/images/sprites/tiles/wall_s.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  3,
		SpritePath: "assets/images/sprites/tiles/wall_e.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  4,
		SpritePath: "assets/images/sprites/tiles/wall_w.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  5,
		SpritePath: "assets/images/sprites/tiles/wall_ne.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  6,
		SpritePath: "assets/images/sprites/tiles/wall_ns.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  7,
		SpritePath: "assets/images/sprites/tiles/wall_nw.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  8,
		SpritePath: "assets/images/sprites/tiles/wall_se.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  9,
		SpritePath: "assets/images/sprites/tiles/wall_sw.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  10,
		SpritePath: "assets/images/sprites/tiles/wall_ew.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  11,
		SpritePath: "assets/images/sprites/tiles/wall_new.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  12,
		SpritePath: "assets/images/sprites/tiles/wall_sew.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  13,
		SpritePath: "assets/images/sprites/tiles/wall_nes.png",
		Walkable:   false,
	},
	{
		Chan:       tilemap.B,
		ChanValue:  14,
		SpritePath: "assets/images/sprites/tiles/wall_nws.png",
		Walkable:   false,
	},
	{
		Chan:      tilemap.R,
		ChanValue: 255,
		Color:     render.Red,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 255,
		Color:     render.Yellow,
		Walkable:  true,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 1,
		Item:      items.DOORS[0],
		Color:     render.Brown,
		Walkable:  false,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 2,
		Item:      items.DOORS[1],
		Color:     render.Brown,
		Walkable:  false,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 3,
		Item:      items.DOORS[2],
		Color:     render.Brown,
		Walkable:  false,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 4,
		Item:      items.DOORS[3],
		Color:     render.Brown,
		Walkable:  false,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 5,
		Item:      items.KEYS[0],
		Color:     render.Yellow,
		Walkable:  true,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 6,
		Item:      items.KEYS[1],
		Color:     render.Yellow,
		Walkable:  true,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 7,
		Item:      items.KEYS[2],
		Color:     render.Yellow,
		Walkable:  true,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 8,
		Item:      items.KEYS[3],
		Color:     render.Yellow,
		Walkable:  true,
	},
}
