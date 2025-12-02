package weapons

import (
	"gates/actors"
	"gates/systems"
)

var SpellFire = actors.Weapon{
	Name:     "Fire Spell",
	Damage:   10,
	Recovery: 5000,
	Modifier: systems.INT,
}
