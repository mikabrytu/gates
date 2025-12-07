package enemies

import "gates/actors"

var Skeleton = actors.EnemySpecs{
	Name:            "Skeleton",
	Image_Path:      "assets/images/placeholder/skeleton.png",
	Size:            512,
	HP:              60,
	Attack_Interval: 1500,
	Attack_Damage:   actors.D8,
}
