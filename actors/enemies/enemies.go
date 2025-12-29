package enemies

import "gates/actors"

const SPRITE_SIZE = 512
const BASE_INTERVAL = 3000

var Rat = actors.EnemySpecs{
	Name:            "Rat",
	Image_Path:      "assets/images/sprites/rat.png",
	Size:            SPRITE_SIZE,
	HP:              8,
	Attack_Interval: BASE_INTERVAL,
	Attack_Damage:   4,
}

var Wolf = actors.EnemySpecs{
	Name:            "Wolf",
	Image_Path:      "assets/images/sprites/wolf.png",
	Size:            SPRITE_SIZE,
	HP:              16,
	Attack_Interval: 2000,
	Attack_Damage:   6,
}

var Zombie = actors.EnemySpecs{
	Name:            "Zombie",
	Image_Path:      "assets/images/sprites/zombie.png",
	Size:            SPRITE_SIZE,
	HP:              28,
	Attack_Interval: 2000,
	Attack_Damage:   12,
}

var Goblin = actors.EnemySpecs{
	Name:            "Goblin",
	Image_Path:      "assets/images/sprites/goblin.png",
	Size:            SPRITE_SIZE,
	HP:              38,
	Attack_Interval: 2000,
	Attack_Damage:   16,
}

var Skeleton = actors.EnemySpecs{
	Name:            "Skeleton",
	Image_Path:      "assets/images/sprites/skeleton.png",
	Size:            SPRITE_SIZE,
	HP:              48,
	Attack_Interval: 1750,
	Attack_Damage:   20,
}

var Bandit = actors.EnemySpecs{
	Name:            "Bandit",
	Image_Path:      "assets/images/sprites/bandit.png",
	Size:            SPRITE_SIZE,
	HP:              58,
	Attack_Interval: 1500,
	Attack_Damage:   22,
}

var Orc = actors.EnemySpecs{
	Name:            "Orc",
	Image_Path:      "assets/images/sprites/orc.png",
	Size:            SPRITE_SIZE,
	HP:              60,
	Attack_Interval: 1200,
	Attack_Damage:   26,
}

var Werewolf = actors.EnemySpecs{
	Name:            "Werewolf",
	Image_Path:      "assets/images/sprites/werewolf.png",
	Size:            SPRITE_SIZE,
	HP:              76,
	Attack_Interval: 1000,
	Attack_Damage:   28,
}

var Vampire = actors.EnemySpecs{
	Name:            "Vampire",
	Image_Path:      "assets/images/sprites/vampire.png",
	Size:            SPRITE_SIZE,
	HP:              106,
	Attack_Interval: 750,
	Attack_Damage:   32,
}

var Dragon = actors.EnemySpecs{
	Name:            "Dragon",
	Image_Path:      "assets/images/placeholder/dragon.png",
	Size:            SPRITE_SIZE,
	HP:              180,
	Attack_Interval: BASE_INTERVAL,
	Attack_Damage:   52,
}
