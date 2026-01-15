package weapons

import (
	"github.com/mikabrytu/gomes-engine/math"
)

var SpellFire = Weapon{
	Name:         "Fire Spell",
	Type:         Fire,
	SpritePath:   "assets/images/sprites/spell.png",
	SpriteSize:   math.Vector2{X: 76 * SPRITE_SIZE_MULTIPLIER, Y: 48 * SPRITE_SIZE_MULTIPLIER},
	SpriteOffset: math.Vector2{X: 0, Y: 0},
	Damage:       6,
	Recovery:     5000,
}

var SpellIce = Weapon{
	Name:         "Ice Spell",
	Type:         Ice,
	SpritePath:   "assets/images/sprites/spell.png",
	SpriteSize:   math.Vector2{X: 76 * SPRITE_SIZE_MULTIPLIER, Y: 48 * SPRITE_SIZE_MULTIPLIER},
	SpriteOffset: math.Vector2{X: 0, Y: 0},
	Damage:       6,
	Recovery:     5000,
}

var SpellShock = Weapon{
	Name:         "Shock Spell",
	Type:         Shock,
	SpritePath:   "assets/images/sprites/spell.png",
	SpriteSize:   math.Vector2{X: 76 * SPRITE_SIZE_MULTIPLIER, Y: 48 * SPRITE_SIZE_MULTIPLIER},
	SpriteOffset: math.Vector2{X: 0, Y: 0},
	Damage:       6,
	Recovery:     5000,
}
