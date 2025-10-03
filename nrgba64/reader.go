package nrgba64

import (
	"errors"
	"image"
	"image/color"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/lunn06/stego/internal/bit"
	"github.com/lunn06/stego/internal/rgb"
)

const byteSize = 8

var ErrInvalidUTF8 = errors.New("invalid UTF-8")

type Image interface {
	image.Image
	NRGBA64At(x, y int) color.NRGBA64
	SetNRGBA64(x, y int, c color.NRGBA64)
}

func ReadString(img Image) (string, error) {
	return NewNRGBA64StringReader(img).ReadString()
}

func NewNRGBA64StringReader(img Image) *NRGBA64StringReader {
	bounds := img.Bounds()
	return &NRGBA64StringReader{
		c:   rgb.Red,
		img: img,
		x:   bounds.Min.X,
		y:   bounds.Min.Y,
	}
}

type NRGBA64StringReader struct {
	img     Image
	x, y, n int
	c       rgb.RGB
}

func (reader *NRGBA64StringReader) ReadByte() (byte, error) {
	var b byte
	for j := range byteSize {
		colorValue, ok := reader.nextColor()
		if !ok {
			return b, io.EOF
		}
		bit := bit.ByIndex(colorValue, 0)

		b += byte(bit * (1 << (byteSize - 1 - j)))
	}

	return b, nil
}

func (reader *NRGBA64StringReader) Read(p []byte) (int, error) {
	lastN := reader.n
	for i := range len(p) {
		b, err := reader.ReadByte()
		if err != nil {
			return reader.n - lastN, err
		}

		p[i] = b
	}

	return reader.n - lastN, nil
}

func (reader *NRGBA64StringReader) ReadString() (string, error) {
	var b strings.Builder
	for {
		r, n, err := reader.ReadRune()
		b.Grow(n)
		switch {
		case err == nil:
			b.WriteRune(r)
		case errors.Is(err, io.EOF):
			return b.String(), nil
		case b.Len() > 0 && errors.Is(err, ErrInvalidUTF8):
			return b.String(), nil
		default:
			return "", err
		}
	}
}

func (reader *NRGBA64StringReader) ReadRune() (rune, int, error) {
	runeBytes, bytesCount, err := reader.readRuneBytes()
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, bytesCount, err
	}

	r, size := utf8.DecodeRune(runeBytes)
	if r == utf8.RuneError {
		return 0, bytesCount, ErrInvalidUTF8
	}

	return r, size, err
}

func (reader *NRGBA64StringReader) readRuneBytes() ([]byte, int, error) {
	runeBytes := make([]byte, 0, utf8.UTFMax)
	for i := range utf8.UTFMax {
		b, err := reader.ReadByte()
		if err != nil {
			return runeBytes, i, err
		}

		runeBytes = append(runeBytes, b)

		if utf8.Valid(runeBytes) {
			return runeBytes, i + 1, nil
		}
	}

	return runeBytes, len(runeBytes), nil
}

func (reader *NRGBA64StringReader) nextColor() (uint16, bool) {
	color := reader.img.NRGBA64At(reader.x, reader.y)
	value := reader.c.NRGBA64Value(color)

	reader.n++
	reader.c = reader.c.Next()

	bounds := reader.img.Bounds()
	if reader.c == rgb.Red {
		reader.x++
	}
	if reader.x >= bounds.Max.X {
		reader.x = bounds.Min.X
		reader.y++
	}
	if reader.y >= bounds.Max.Y {
		return 0, false
	}

	return value, true
}
