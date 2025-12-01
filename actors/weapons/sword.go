package weapons

import (
	"gates/actors"
	"gates/systems"
)

var Sword = actors.Weapon{
	Name:     "Sword",
	Damage:   6,
	Recovery: 3000,
	Modifier: systems.STR,
}
