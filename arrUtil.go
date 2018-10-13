package util

type IArrayUtil interface {
	IsEmptyStringArray(arr []string) bool
}

type ArrayUtil struct{}

func (arrUtil *ArrayUtil) IsEmptyStringArray(arr []string) bool {
	return arr == nil || len(arr) == 0
}
