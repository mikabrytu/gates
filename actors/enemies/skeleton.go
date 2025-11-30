package enemies

import "gates/actors"

var Skeleton = actors.EnemySpecs{
	Name:            "Skeleton",
	Image_Path:      "assets/images/skeleton.png",
	Size:            512,
	HP:              100,
	Attack_Interval: 1500,
	Attack_Damage:   15,
}
