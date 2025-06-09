package stego

import (
	"image"
	"image/color"
	"math/bits"
)

const (
	Delta           = 1 << 0
	InfoPixelsCount = bits.UintSize / Delta
)

func GetBit[T int | rune | byte | uint | uint16](b T, idx int) T {
	return (b >> idx) & 0b01
}

func ClearLastBits[T byte | uint32 | uint16](b T, n int) T {
	return b >> n << n
}

func ReadPointColors(img image.Image) []PointColor {
	bounds := img.Bounds()
	pixelsCount := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	limit := pixelsCount * 3
	data := make([]PointColor, 0, limit)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data = append(data,
				PointColor{
					Point: image.Point{
						X: x,
						Y: y,
					},
					NRGBA64: color.NRGBA64{
						R: uint16(r),
						G: uint16(g),
						B: uint16(b),
						A: uint16(a),
					},
				},
			)
		}
	}

	return data
}
