package weapons

import (
	"gates/actors"
	"gates/systems"

	"github.com/mikabrytu/gomes-engine/math"
)

var Bow = actors.Weapon{
	Name:         "Bow",
	SpritePath:   "assets/images/sprites/bow.png",
	SpriteSize:   math.Vector2{X: 16 * SPRITE_SIZE_MULTIPLIER, Y: 106 * SPRITE_SIZE_MULTIPLIER},
	SpriteOffset: math.Vector2{X: -32, Y: 32 * SPRITE_SIZE_MULTIPLIER},
	Damage:       4,
	Recovery:     1500,
	Modifier:     systems.SPD,
}
