package data

import (
	"gates/pkg/tilemap"

	"github.com/mikabrytu/gomes-engine/render"
)

var RULES []tilemap.TileRules = []tilemap.TileRules{
	{
		Chan:      tilemap.R,
		ChanValue: 255,
		Color:     render.Red,
	},
	{
		Chan:      tilemap.G,
		ChanValue: 255,
		Color:     render.Yellow,
	},
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
}
