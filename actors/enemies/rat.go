package enemies

import "gates/actors"

var Rat = actors.EnemySpecs{
	Name:            "Rat",
	Image_Path:      "assets/images/rat.png",
	Size:            516,
	HP:              50,
	Attack_Interval: 3000,
	Attack_Damage:   1,
}
