package math

func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}
