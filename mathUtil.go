package util

type IMathUtil interface {
	MinInt(a, b int) int
}

type mathUtil struct{}

func (mathUtil *mathUtil) MinInt(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
