package util

type IArrayUtil interface {
	IsEmptyStringArray(arr []string) bool
}

type arrayUtil struct{}

func (arrUtil *arrayUtil) IsEmptyStringArray(arr []string) bool {
	return arr == nil || len(arr) == 0
}
