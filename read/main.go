package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/bits"
	"os"
	"stego"
	"unicode/utf8"
)

func main() {
	fileBytes, err := os.ReadFile("output.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(bytes.NewReader(fileBytes))
	//img, err := jpeg.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		panic(err)
	}

	nrgba, ok := img.(*image.RGBA64)
	if !ok {
		panic(fmt.Sprintf("expected *image.RGBA64, got = %T", img))
	}

	info, err := GetInfo(nrgba)
	if err != nil {
		panic(err)
	}

	secret, err := GetSecret(nrgba, info, 1, 21)
	if err != nil {
		panic(err)
	}

	fmt.Println(secret)
	//err = os.WriteFile("2", secret, 0644)
	//if err != nil {
	//	panic(err)
	//}
}

func GetInfo(data *image.RGBA64) (uint, error) {
	bounds := data.Bounds()
	if bits.UintSize > bounds.Max.X*bounds.Max.Y*3 {
		return 0, errors.New("secret too large")
	}

	x := bounds.Min.X
	maxX := bounds.Max.X
	y := bounds.Min.Y
	var i, info uint
	for j := range bits.UintSize {
		c := data.RGBA64At(x, y)

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

func GetSecret(data *image.RGBA64, info uint, delta, offset int) (string, error) {
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

	changeB := func(c color.RGBA64) uint16 {
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
		for j := range 3 {
			c := data.RGBA64At(x, y)
			b := changeB(c)
			runeLen += int(b * (1 << j))

			incXYI()
			k++
		}
		//secret = append(secret, s)

		var s rune
		for j := range runeLen * 8 {
			c := data.RGBA64At(x, y)
			b := changeB(c)
			s += rune(b * (1 << j))

			incXYI()
			k++

		}
		secret = utf8.AppendRune(secret, s)
	}

	return string(secret), nil
}
