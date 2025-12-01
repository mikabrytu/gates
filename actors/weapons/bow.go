package weapons

import (
	"gates/actors"
	"gates/systems"
)

var Bow = actors.Weapon{
	Name:     "Bow",
	Damage:   3,
	Recovery: 1500,
	Modifier: systems.SPD,
}
