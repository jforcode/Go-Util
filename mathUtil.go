package util

type MathUtil struct{}

func (mathUtil *MathUtil) MinInt(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
