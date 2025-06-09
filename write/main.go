package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/bits"
	"os"
	"stego"
	"unicode/utf8"
)

func main() {
	//fileBytes, err := os.ReadFile("image.jpg")
	fileBytes, err := os.ReadFile("image2.png")
	if err != nil {
		panic(err)
	}

	//img, err := jpeg.Decode(bytes.NewReader(fileBytes))
	img, err := png.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		panic(err)
	}
	//img, err := png.Decode(bytes.NewReader(fileBytes))
	//if err != nil {
	//	panic(err)
	//}

	//secret, err := os.ReadFile("1")
	//if err != nil {
	//	panic(err)
	//}
	secret := "test1234567890"

	nrgba := image.NewNRGBA64(img.Bounds())
	draw.Draw(nrgba, img.Bounds(), img, img.Bounds().Min, draw.Src)

	c, err := InsertSecret(secret, nrgba, stego.Delta, stego.OffsetLen)
	if err != nil {
		panic(err)
	}

	_, err = InsertInfo(uint(c), nrgba)
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	//if err = jpeg.Encode(outFile, rgba, &jpeg.Options{Quality: 100}); err != nil {
	//	panic(err)
	//}
	if err = png.Encode(outFile, nrgba); err != nil {
		panic(err)
	}
}

func InsertInfo(info uint, data *image.NRGBA64) (int, error) {
	bounds := data.Bounds()
	if bits.UintSize > bounds.Max.X*bounds.Max.Y*3 {
		return 0, errors.New("secret too large")
	}

	x := bounds.Min.X
	maxX := bounds.Max.X
	y := bounds.Min.Y
	var i, counter int
	for j := range bits.UintSize {
		b := uint16(stego.GetBit(info, j))

		c := data.NRGBA64At(x, y)
		switch {
		case i == 0:
			c.R = stego.ClearLastBits(c.R, 1) + b
		case i == 1:
			c.G = stego.ClearLastBits(c.G, 1) + b
		case i == 2:
			c.B = stego.ClearLastBits(c.B, 1) + b
		}

		data.SetNRGBA64(x, y, c)

		i++
		if i < 3 {
			continue
		}
		i = 0
		x++
		counter++
		if x == maxX {
			x = 0
			y++
		}
	}

	return counter, nil
}

func InsertSecret(secret string, data *image.NRGBA64, delta, offset int) (int, error) {
	bounds := data.Bounds()
	if len(secret)*8/delta > bounds.Max.X*bounds.Max.Y*(1<<delta)*3 {
		return 0, errors.New("secret too large")
	}

	var (
		x    = bounds.Min.X + offset
		maxX = bounds.Max.X
		y    = bounds.Min.Y

		i, counter int
	)

	changeC := func(b uint16) color.NRGBA64 {
		c := data.NRGBA64At(x, y)
		switch {
		case i == 0:
			c.R = stego.ClearLastBits(c.R, 1) + b
		case i == 1:
			c.G = stego.ClearLastBits(c.G, 1) + b
		case i == 2:
			c.B = stego.ClearLastBits(c.B, 1) + b
		}
		return c
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

	for _, s := range secret {
		runeLen := utf8.RuneLen(s)
		if runeLen == -1 {
			return 0, errors.New("invalid secret")
		}
		for j := range stego.UT8Len {
			b := uint16(stego.GetBit(runeLen, j))
			c := changeC(b)

			data.SetNRGBA64(x, y, c)

			counter++
			incXYI()
		}

		for j := range runeLen * 8 {
			b := uint16(stego.GetBit(s, j))

			c := changeC(b)

			data.SetNRGBA64(x, y, c)

			counter++
			incXYI()
		}
	}

	return counter, nil
}
