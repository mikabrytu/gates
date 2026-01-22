package weapons

import (
	"github.com/mikabrytu/gomes-engine/math"
)

var Bow = Weapon{
	Name:         "Bow",
	Type:         Physical,
	SpritePath:   "assets/images/sprites/weapons/bow.png",
	SpriteSize:   math.Vector2{X: 16 * SPRITE_SIZE_MULTIPLIER, Y: 106 * SPRITE_SIZE_MULTIPLIER},
	SpriteOffset: math.Vector2{X: -32, Y: 32 * SPRITE_SIZE_MULTIPLIER},
	Damage:       3,
	Recovery:     1500,
}
