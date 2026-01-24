package data

type EnemySpecs struct {
	Name            string
	Image_Path      string
	Size            int
	HP              int
	Attack_Interval int
	Attack_Damage   int
	Defense         int
}

const SPRITE_SIZE = 512
const BASE_INTERVAL = 3000

var Rat = EnemySpecs{
	Name:            "Rat",
	Image_Path:      "assets/images/sprites/enemies/rat.png",
	Size:            SPRITE_SIZE,
	HP:              10,
	Attack_Interval: BASE_INTERVAL,
	Attack_Damage:   4,
	Defense:         1,
}

var Wolf = EnemySpecs{
	Name:            "Wolf",
	Image_Path:      "assets/images/sprites/enemies/wolf.png",
	Size:            SPRITE_SIZE,
	HP:              16,
	Attack_Interval: BASE_INTERVAL,
	Attack_Damage:   6,
	Defense:         1,
}

var Zombie = EnemySpecs{
	Name:            "Zombie",
	Image_Path:      "assets/images/sprites/enemies/zombie.png",
	Size:            SPRITE_SIZE,
	HP:              28,
	Attack_Interval: 2500,
	Attack_Damage:   8,
	Defense:         2,
}

var Goblin = EnemySpecs{
	Name:            "Goblin",
	Image_Path:      "assets/images/sprites/enemies/goblin.png",
	Size:            SPRITE_SIZE,
	HP:              38,
	Attack_Interval: 2250,
	Attack_Damage:   12,
	Defense:         2,
}

var Skeleton = EnemySpecs{
	Name:            "Skeleton",
	Image_Path:      "assets/images/sprites/enemies/skeleton.png",
	Size:            SPRITE_SIZE,
	HP:              48,
	Attack_Interval: 2000,
	Attack_Damage:   16,
	Defense:         3,
}

var Bandit = EnemySpecs{
	Name:            "Bandit",
	Image_Path:      "assets/images/sprites/enemies/bandit.png",
	Size:            SPRITE_SIZE,
	HP:              58,
	Attack_Interval: 1700,
	Attack_Damage:   20,
	Defense:         3,
}

var Orc = EnemySpecs{
	Name:            "Orc",
	Image_Path:      "assets/images/sprites/enemies/orc.png",
	Size:            SPRITE_SIZE,
	HP:              60,
	Attack_Interval: 1500,
	Attack_Damage:   24,
	Defense:         4,
}

var Werewolf = EnemySpecs{
	Name:            "Werewolf",
	Image_Path:      "assets/images/sprites/enemies/werewolf.png",
	Size:            SPRITE_SIZE,
	HP:              76,
	Attack_Interval: 1000,
	Attack_Damage:   28,
	Defense:         5,
}

var Vampire = EnemySpecs{
	Name:            "Vampire",
	Image_Path:      "assets/images/sprites/enemies/vampire.png",
	Size:            SPRITE_SIZE,
	HP:              106,
	Attack_Interval: 750,
	Attack_Damage:   28,
	Defense:         5,
}

var Dragon = EnemySpecs{
	Name:            "Dragon",
	Image_Path:      "assets/images/placeholder/dragon.png",
	Size:            SPRITE_SIZE,
	HP:              180,
	Attack_Interval: BASE_INTERVAL,
	Attack_Damage:   32,
	Defense:         6,
}
