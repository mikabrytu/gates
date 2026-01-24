package spell

import (
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type EffectType int

const (
	Burn EffectType = iota
	Cold
	Paralysis
)

type Effect struct {
	Type  EffectType
	Stack int
}

type EffectIcon struct {
	Rect  utils.RectSpecs
	Color render.Color
}
