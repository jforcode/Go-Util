package util

import "testing"

func TestMinInt(t *testing.T) {
	mathUtil := &MathUtil{}

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"first less", -1, 20, -1},
		{"second less", 25, 20, 20},
		{"equal", 20, 20, 20},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := mathUtil.MinInt(test.a, test.b)
			if actual != test.expected {
				t.FailNow()
			}
		})
	}

}
