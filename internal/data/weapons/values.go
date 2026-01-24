package data

import "github.com/mikabrytu/gomes-engine/math"

type WeaponType int

const (
	Physical WeaponType = iota
	Fire
	Ice
	Shock
)

type Weapon struct {
	Name         string
	Type         WeaponType
	SpritePath   string
	SpriteSize   math.Vector2
	SpriteOffset math.Vector2
	Damage       int
	Recovery     int
}

const SPRITE_SIZE_MULTIPLIER int = 8
