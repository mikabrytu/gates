package data

import (
	"github.com/mikabrytu/gomes-engine/math"
)

var Sword = Weapon{
	Name:         "Sword",
	Type:         Physical,
	SpritePath:   "assets/images/sprites/weapons/sword.png",
	SpriteSize:   math.Vector2{X: 18 * SPRITE_SIZE_MULTIPLIER, Y: 106 * SPRITE_SIZE_MULTIPLIER},
	SpriteOffset: math.Vector2{X: -32, Y: 132},
	Damage:       4,
	Recovery:     3000,
}
