package spells

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
