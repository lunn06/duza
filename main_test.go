package stego

import "testing"

func TestGetBit(t *testing.T) {
	type args[T interface {
		int | rune | byte | uint | uint16
	}] struct {
		b   T
		idx int
	}
	type testCase[T interface {
		int | rune | byte | uint | uint16
	}] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[uint16]{
		{
			name: "GetBit",
			args: args[uint16]{
				b:   uint16(0b00001),
				idx: 0,
			},
			want: uint16(1),
		},
		{
			name: "GetBit",
			args: args[uint16]{
				b:   uint16(0b01001),
				idx: 15,
			},
			want: uint16(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBit(tt.args.b, tt.args.idx); got != tt.want {
				t.Errorf("GetBit() = %v, want %v", got, tt.want)
			}
		})
	}
}
