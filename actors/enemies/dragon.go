package enemies

import "gates/actors"

var Dragon = actors.EnemySpecs{
	Name:            "Dragon",
	Image_Path:      "assets/images/dragon.png",
	Size:            640,
	HP:              300,
	Attack_Interval: 3000,
	Attack_Damage:   15,
}
