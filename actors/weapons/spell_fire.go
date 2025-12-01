package weapons

import (
	"gates/actors"
	"gates/systems"
)

var SpellFire = actors.Weapon{
	Name:     "Fire Spell",
	Damage:   10000,
	Recovery: 5,
	Modifier: systems.INT,
}
