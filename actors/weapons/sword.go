package weapons

import (
	"gates/actors"
	"gates/systems"

	"github.com/mikabrytu/gomes-engine/math"
)

var Sword = actors.Weapon{
	Name:         "Sword",
	SpritePath:   "assets/images/sprites/sword.png",
	SpriteSize:   math.Vector2{X: 18 * 7, Y: 106 * 7},
	SpriteOffset: math.Vector2{X: -32, Y: 132},
	Damage:       6,
	Recovery:     3000,
	Modifier:     systems.STR,
}
