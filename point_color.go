package stego

import (
	"image"
	"image/color"
)

type PointColor struct {
	image.Point
	color.NRGBA64
}

func (pc *PointColor) Get(i int) uint16 {
	switch i {
	case 0:
		return pc.R
	case 1:
		return pc.G
	case 2:
		return pc.B
	case 3:
		return pc.A
	default:
		panic("invalid point index")
	}
}

func (pc *PointColor) Set(v uint16, i int) {
	switch i {
	case 0:
		pc.R = v
	case 1:
		pc.G = v
	case 2:
		pc.B = v
	case 3:
		pc.A = v
	default:
		panic("invalid point index")
	}
}
