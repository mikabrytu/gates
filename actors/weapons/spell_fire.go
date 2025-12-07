package weapons

import (
	"gates/actors"
	"gates/systems"

	"github.com/mikabrytu/gomes-engine/math"
)

var SpellFire = actors.Weapon{
	Name:         "Fire Spell",
	SpritePath:   "assets/images/sprites/spell.png",
	SpriteSize:   math.Vector2{X: 76 * SPRITE_SIZE_MULTIPLIER, Y: 48 * SPRITE_SIZE_MULTIPLIER},
	SpriteOffset: math.Vector2{X: 0, Y: 0},
	Damage:       10,
	Recovery:     5000,
	Modifier:     systems.INT,
}
