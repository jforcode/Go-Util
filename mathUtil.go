package util

type IMathUtil interface {
	MinInt(a, b int) int
}

type MathUtil struct{}

func (mathUtil *MathUtil) MinInt(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
