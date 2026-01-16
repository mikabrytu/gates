package utils

import "math/rand/v2"

func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func CalcDamange(max, min int) int {
	return (rand.IntN(max-min) + min) + 1
}
