package stego

import (
	"image/color"
	"math"
	"reflect"
	"testing"
)

func Test_getBit(t *testing.T) {
	type args struct {
		b   byte
		idx int
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "getBit first",
			args: args{
				b:   1,
				idx: 0,
			},
			want: 1,
		},
		{
			name: "getBit last",
			args: args{
				b:   255,
				idx: 7,
			},
			want: 1,
		},
		{
			name: "getBit middle",
			args: args{
				b:   16,
				idx: 4,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBit(tt.args.b, tt.args.idx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertSecret(t *testing.T) {
	type args struct {
		secret []byte
		data   []PointColor
		delta  int
	}
	tests := []struct {
		name    string
		args    args
		want    []PointColor
		wantErr bool
	}{
		{
			name: "InsertSecret",
			args: args{
				secret: []byte{1},
				data: []PointColor{
					{RGBA64: color.RGBA64{R: 2}},
					{},
					{},
				},
				delta: 1,
			},
			want: []PointColor{
				{RGBA64: color.RGBA64{R: 2}},
				{},
				{RGBA64: color.RGBA64{G: 1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertSecret(tt.args.secret, tt.args.data, tt.args.delta)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clearLastBits(t *testing.T) {
	type args struct {
		b byte
		n int
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "clearLastBits one",
			args: args{
				b: 0b0000_0001,
				n: 1,
			},
			want: 0b0000_0000,
		},
		{
			name: "clearLastBits some",
			args: args{
				b: 0b0000_0111,
				n: 2,
			},
			want: 0b0000_0100,
		},
		{
			name: "clearLastBits some2",
			args: args{
				b: 0b0000_0011,
				n: 2,
			},
			want: 0b0000_0000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clearLastBits(tt.args.b, tt.args.n); got != tt.want {
				t.Errorf("clearLastBits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertInfo(t *testing.T) {
	info2data := make([]PointColor, 256)
	info2data[0] = PointColor{RGBA64: color.RGBA64{A: 0}}
	info2data[15] = PointColor{RGBA64: color.RGBA64{A: 1}}

	info3data := make([]PointColor, 256*256*256*256)
	info2data[255*255] = PointColor{RGBA64: color.RGBA64{A: 1}}
	type args struct {
		info uint
		data []PointColor
		maxX int
		maxY int
	}
	tests := []struct {
		name    string
		args    args
		want    []PointColor
		wantErr bool
	}{
		{
			name: "InsertInfo1",
			args: args{
				info: 123,
				data: make([]PointColor, 4),
				maxX: 1,
				maxY: 1,
			},
			want: []PointColor{
				{RGBA64: color.RGBA64{A: 123}},
				{},
				{},
				{},
			},
		},
		{
			name: "InsertInfo2",
			args: args{
				info: math.MaxUint16 + 1,
				data: make([]PointColor, 256),
				maxX: 15,
				maxY: 15,
			},
			want: info2data,
		},
		{
			name: "InsertInfo3",
			args: args{
				info: 4294967296,
				data: make([]PointColor, 256*256*256*256),
				maxX: 255,
				maxY: 255,
			},
			want: info3data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertInfo(tt.args.info, tt.args.data, tt.args.maxX, tt.args.maxY, 1, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInfo(t *testing.T) {
	info2data := make([]PointColor, 256)
	info2data[0] = PointColor{RGBA64: color.RGBA64{A: 2}}
	info2data[15] = PointColor{RGBA64: color.RGBA64{A: 1}}
	type args struct {
		data []PointColor
		maxX int
		maxY int
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "GetInfo2",
			args: args{
				data: info2data,
				maxX: 15,
				maxY: 15,
			},
			want: math.MaxUint16 + 3,
		},
		{
			name: "GetInfo1",
			args: args{
				data: []PointColor{
					{RGBA64: color.RGBA64{A: 123}},
					{},
					{},
					{},
				},
				maxX: 1,
				maxY: 1,
			},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInfo(tt.args.data, tt.args.maxX, tt.args.maxY); got != tt.want {
				t.Errorf("GetInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
