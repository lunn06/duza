package nrgba64_test

import (
	"reflect"
	"testing"

	"github.com/lunn06/duza/nrgba64"
)

func TestWriteStringToNRGBA64(t *testing.T) {
	type args struct {
		img nrgba64.Image
		s   string
	}
	tests := []struct {
		name    string
		args    args
		want    nrgba64.Image
		wantErr bool
	}{
		{
			name: "Succes With One",
			args: args{
				s: "H",
				img: nrgba64.NewMockImage(
					2, 2,
					[]uint16{
						0, 0, 0 /**/, 0, 0, 0,
						0, 0, 0 /**/, 0, 0, 0,
					}...,
				),
			},
			want: nrgba64.NewMockImage(
				2, 2,
				[]uint16{
					0, 1, 0 /**/, 0, 1, 0,
					0, 0, 0 /**/, 0, 0, 0,
				}...,
			),
		},
		{
			name: "Succes With Many",
			args: args{
				s: "Hey👋",
				img: nrgba64.NewMockImage(
					5, 4,
					[]uint16{
						0, 0, 0 /**/, 0, 0, 0 /**/, 0, 0, 0 /**/, 0, 0, 0 /**/, 0, 1, 0,
						1, 0, 1 /**/, 1, 0, 1 /**/, 0, 0, 1 /**/, 1, 1, 1 /**/, 1, 0, 0,
						0, 0, 1 /**/, 0, 0, 0 /**/, 1, 1, 1 /**/, 1, 1, 0 /**/, 0, 1, 0,
						0, 0, 1 /**/, 0, 0, 0 /**/, 0, 1, 0 /**/, 0, 1, 0 /**/, 1, 1, 1,
					}...,
				),
			},
			want: nrgba64.NewMockImage(
				5, 4,
				[]uint16{
					0, 1, 0 /**/, 0, 1, 0 /**/, 0, 0, 0 /**/, 1, 1, 0 /**/, 0, 1, 0,
					1, 0, 1 /**/, 1, 1, 1 /**/, 0, 0, 1 /**/, 1, 1, 1 /**/, 1, 0, 0,
					0, 0, 1 /**/, 0, 0, 1 /**/, 1, 1, 1 /**/, 1, 1, 0 /**/, 0, 1, 0,
					0, 0, 1 /**/, 1, 0, 0 /**/, 0, 1, 0 /**/, 1, 1, 0 /**/, 1, 1, 1,
				}...,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nrgba64.WriteStringToNRGBA64(tt.args.img, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteStringToNRGBA64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.want, tt.args.img) {
				t.Errorf("WriteStringToNRGBA64() got = %v, want %v", tt.args.img, tt.want)
			}
		})
	}
}
