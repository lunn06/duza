package nrgba64

import (
	"image"
	"image/color"
)

var _ Image = (*MockImage)(nil)

func NewMockImage(maxX, maxY int, pix ...uint16) MockImage {
	return MockImage{
		pix:  pix,
		rect: image.Rect(0, 0, maxX, maxY),
	}
}

type MockImage struct {
	pix  []uint16
	rect image.Rectangle
}

func (m MockImage) NRGBA64At(x, y int) color.NRGBA64 {
	// p: 1 1 1 0 0 0 1 0 1 0 0 0
	// x: 0 0 0 1 1 1 0 0 0 1 1 1
	// y: 0 0 0 0 0 0 1 1 1 0 0 0
	// i: 0 1 2 3 4 5 6 7 8 9 10 11
	bounds := m.rect.Bounds()
	idx := x*3 + bounds.Max.X*y*3

	return color.NRGBA64{
		B: m.pix[idx+2],
		G: m.pix[idx+1],
		R: m.pix[idx],
		A: 0,
	}
}

func (m MockImage) SetNRGBA64(x, y int, c color.NRGBA64) {
	bounds := m.rect.Bounds()
	idx := x*3 + bounds.Max.X*y*3

	m.pix[idx+2] = c.B
	m.pix[idx+1] = c.G
	m.pix[idx] = c.R
}

func (m MockImage) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (m MockImage) Bounds() image.Rectangle {
	return m.rect
}

func (m MockImage) At(x, y int) color.Color {
	return m.NRGBA64At(x, y)
}
