package nrgba64

import (
	"image/color"
	"io"

	"github.com/lunn06/stego/internal/bit"
	"github.com/lunn06/stego/internal/rgb"
)

func WriteStringToNRGBA64(img Image, s string) (int, error) {
	return NewNRGBA64StringWriter(img).WriteString(s)
}

func NewNRGBA64StringWriter(img Image) *NRGBA64StringWriter {
	bounds := img.Bounds()
	return &NRGBA64StringWriter{
		img: img,
		x:   bounds.Min.X,
		y:   bounds.Min.Y,
		c:   rgb.Red,
	}
}

var _ io.Writer = (*NRGBA64StringWriter)(nil)

type NRGBA64StringWriter struct {
	img     Image
	x, y, n int
	c       rgb.RGB
}

func (writer *NRGBA64StringWriter) Write(p []byte) (int, error) {
	lastN := writer.n
	for _, b := range p {
		if err := writer.WriteByte(b); err != nil {
			return writer.n - lastN, err
		}
	}

	return writer.n - lastN, nil
}

func (writer *NRGBA64StringWriter) WriteString(s string) (int, error) {
	return writer.Write([]byte(s))
}

func (writer *NRGBA64StringWriter) WriteByte(b byte) error {
	for j := range byteSize {
		color := writer.img.NRGBA64At(writer.x, writer.y)
		value := writer.c.NRGBA64Value(color)

		value = bit.ClearFirstN(value, 1)
		value += uint16(bit.ByIndex(b, byteSize-1-j))

		writer.setColorValue(value, &color)
		writer.img.SetNRGBA64(writer.x, writer.y, color)

		writer.n++
		writer.c = writer.c.Next()

		bounds := writer.img.Bounds()
		if writer.c == rgb.Red {
			writer.x++
		}
		if writer.x >= bounds.Max.X {
			writer.x = bounds.Min.X
			writer.y++
		}
		if writer.y >= bounds.Max.Y {
			return io.EOF
		}
	}

	return nil
}

func (writer *NRGBA64StringWriter) setColorValue(value uint16, color *color.NRGBA64) {
	switch writer.c {
	case rgb.Red:
		color.R = value
	case rgb.Green:
		color.G = value
	case rgb.Blue:
		color.B = value
	case rgb.Undefined:
		fallthrough
	default:
		panic("undefined color value")
	}
}
