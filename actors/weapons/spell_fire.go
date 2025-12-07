package weapons

import (
	"gates/actors"
	"gates/systems"
)

var SpellFire = actors.Weapon{
	Name:       "Fire Spell",
	SpritePath: "assets/images/sprites/sword.png",
	Damage:     10,
	Recovery:   5000,
	Modifier:   systems.INT,
}
