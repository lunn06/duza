package nrgba64_test

import (
	"testing"

	"github.com/lunn06/duza/nrgba64"
)

func TestReadStringFromNRGBA64(t *testing.T) {
	tests := []struct {
		name    string
		img     nrgba64.Image
		want    string
		wantErr bool
	}{
		{
			// 01001000 -> H
			name: "Success With One",
			img: nrgba64.NewMockImage(
				2, 2,
				[]uint16{
					0, 1, 0 /**/, 0, 1, 0,
					0, 0, 0 /**/, 0, 0, 0,
				}...,
			),
			want: "H",
		},
		{
			// 01001000 -> H
			name: "Success With One Reverted",
			img: nrgba64.NewMockImage(
				2, 2,
				[]uint16{
					254, 255, 254 /**/, 254, 255, 254,
					254, 254, 254 /**/, 254, 254, 254,
				}...,
			),
			want: "H",
		},
		{
			// 01001000 01100101 01111001 11110000 10011111 10010001 10001011  -> Hey👋
			name: "Success With Many",
			img: nrgba64.NewMockImage(
				5, 4,
				[]uint16{
					0, 1, 0 /**/, 0, 1, 0 /**/, 0, 0, 0 /**/, 1, 1, 0 /**/, 0, 1, 0,
					1, 0, 1 /**/, 1, 1, 1 /**/, 0, 0, 1 /**/, 1, 1, 1 /**/, 1, 0, 0,
					0, 0, 1 /**/, 0, 0, 1 /**/, 1, 1, 1 /**/, 1, 1, 0 /**/, 0, 1, 0,
					0, 0, 1 /**/, 1, 0, 0 /**/, 0, 1, 0 /**/, 1, 1, 0 /**/, 1, 1, 1,
				}...,
			),
			want: "Hey👋",
		},
		{
			name: "Invalid chars",
			img: nrgba64.NewMockImage(
				5, 1,
				[]uint16{
					255, 255, 255 /**/, 0, 0, 0 /**/, 0, 0, 0 /**/, 1, 1, 0 /**/, 0, 1, 0,
				}...,
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nrgba64.ReadString(tt.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
