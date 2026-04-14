package utils_test

import (
	"testing"

	"github.com/shaolei/cubism-go/renderer/utils"
)

func TestNormalize(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		x      float32
		n      float32
		m      float32
		expect float32
	}{
		{"center", 5, 0, 10, 0},
		{"min", 0, 0, 10, -1},
		{"max", 10, 0, 10, 1},
		{"below range", -5, 0, 10, -2},
		{"above range", 15, 0, 10, 2},
		{"quarter", 2.5, 0, 10, -0.5},
		{"three quarters", 7.5, 0, 10, 0.5},
		{"same n and m returns zero", 5, 3, 3, 0},
		{"same n and m with x equals n", 3, 3, 3, 0},
		{"negative range", 0, -10, 10, 0},
		{"small range", 0.5, 0, 1, 0},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := utils.Normalize(tc.x, tc.n, tc.m)
			if got != tc.expect {
				t.Errorf("Normalize(%v, %v, %v) = %v, want %v", tc.x, tc.n, tc.m, got, tc.expect)
			}
		})
	}
}
