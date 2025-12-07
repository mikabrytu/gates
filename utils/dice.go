package utils

type Dice int

const (
	D4 Dice = iota
	D6
	D8
	D12
	D20
)

var DiceValue = map[Dice]int{
	D4:  4,
	D6:  6,
	D8:  8,
	D12: 12,
	D20: 20,
}
