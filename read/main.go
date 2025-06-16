package read

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/bits"
	"os"
	"unicode/utf8"

	"github.com/lunn06/stego"
)

func Read(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		panic(err)
	}

	nrgba, ok := img.(*image.NRGBA64)
	if !ok {
		nrgba = image.NewNRGBA64(img.Bounds())
		draw.Draw(nrgba, img.Bounds(), img, img.Bounds().Min, draw.Over)
	}

	info, err := GetInfo(nrgba)
	if err != nil {
		return "", err
	}

	secret, err := GetSecret(nrgba, info, stego.Delta, stego.OffsetLen)
	if err != nil {
		return "", err
	}

	return secret, nil
}

func GetInfo(data *image.NRGBA64) (uint, error) {
	bounds := data.Bounds()
	if bits.UintSize > bounds.Max.X*bounds.Max.Y*3 {
		return 0, errors.New("secret too large")
	}

	x := bounds.Min.X
	maxX := bounds.Max.X
	y := bounds.Min.Y
	var i, info uint
	for j := range bits.UintSize {
		c := data.NRGBA64At(x, y)

		var b uint16
		switch {
		case i == 0:
			b = stego.GetBit(c.R, 0)
		case i == 1:
			b = stego.GetBit(c.G, 0)
		case i == 2:
			b = stego.GetBit(c.B, 0)
		}

		info += uint(b * (1 << j))

		i++
		if i < 3 {
			continue
		}
		i = 0
		x++
		if x == maxX {
			x = 0
			y++
		}
	}

	return info, nil
}

func GetSecret(data *image.NRGBA64, info uint, delta, offset int) (string, error) {
	bounds := data.Bounds()
	if offset > bounds.Max.X*bounds.Max.Y*(1<<delta)*3 {
		return "", errors.New("secret too large")
	}

	var (
		secret = make([]byte, 0, info/8)
		maxX   = bounds.Max.X
		x      = bounds.Min.X + offset
		y      = bounds.Min.Y

		i int
	)
	if offset > maxX {
		x = bounds.Min.X + offset%maxX
		y += offset / maxX
	}

	changeB := func(c color.NRGBA64) uint16 {
		var b uint16
		switch {
		case i == 0:
			b = stego.GetBit(c.R, 0)
		case i == 1:
			b = stego.GetBit(c.G, 0)
		case i == 2:
			b = stego.GetBit(c.B, 0)
		}

		return b
	}

	incXYI := func() {
		i++
		if i < 3 {
			return
		}
		i = 0
		x++
		if x == maxX {
			x = 0
			y++
		}
	}

	for k := uint(0); k < info; {
		var runeLen int
		for j := range stego.UTF8Len {
			c := data.NRGBA64At(x, y)
			b := changeB(c)
			runeLen += int(b * (1 << j))

			incXYI()
			k++
		}

		var s rune
		for j := range runeLen * 8 {
			c := data.NRGBA64At(x, y)
			b := changeB(c)
			s += rune(b * (1 << j))

			incXYI()
			k++
		}

		if utf8.ValidRune(s) {
			secret = utf8.AppendRune(secret, s)
		} else {
			return "", errors.New("invalid UTF-8")
		}
	}

	return string(secret), nil
}
