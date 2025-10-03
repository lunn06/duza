package bit

import (
	"testing"

	"golang.org/x/exp/constraints"
)

func TestByIndex(t *testing.T) {
	const b byte = 72
	type testCase[T constraints.Integer] struct {
		name string
		idx  int
		want T
	}
	tests := []testCase[byte]{
		{name: "7", idx: 7, want: 0},
		{name: "6", idx: 6, want: 1},
		{name: "5", idx: 5, want: 0},
		{name: "4", idx: 4, want: 0},
		{name: "3", idx: 3, want: 1},
		{name: "2", idx: 2, want: 0},
		{name: "1", idx: 1, want: 0},
		{name: "0", idx: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByIndex(b, tt.idx); got != tt.want {
				t.Errorf("ByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
