package util

import "testing"

func TestIsEmptyStringArray(t *testing.T) {
	arrUtil := &arrayUtil{}

	tests := []struct {
		name     string
		arr      []string
		expected bool
	}{
		{"nil arr", nil, true},
		{"empty arr", make([]string, 0), true},
		{"1 element arr", []string{"asdf"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := arrUtil.IsEmptyStringArray(test.arr)
			if test.expected != actual {
				t.FailNow()
			}
		})
	}
}
