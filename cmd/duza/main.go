package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/lunn06/duza/nrgba64"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

var cmd = &cli.Command{
	Commands: []*cli.Command{
		{
			Name:  "encode",
			Usage: "Скрывает информацию в изображении",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "sourcePath",
				},
				&cli.StringArg{
					Name: "secret",
				},
				&cli.StringArg{
					Name: "outPath",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				file, err := os.Open(cmd.StringArg("sourcePath"))
				if err != nil {
					return err
				}
				defer func() { _ = file.Close() }()

				img, err := png.Decode(file)
				if err != nil {
					return err
				}

				nrgba := image.NewNRGBA64(img.Bounds())
				draw.Draw(nrgba, img.Bounds(), img, img.Bounds().Min, draw.Src)

				_, err = nrgba64.WriteStringToNRGBA64(nrgba, cmd.StringArg("secret"))
				if err != nil {
					return err
				}

				outFile, err := os.Create(cmd.StringArg("outPath"))
				if err != nil {
					return err
				}
				defer func() { _ = outFile.Close() }()

				if err = png.Encode(outFile, nrgba); err != nil {
					return err
				}

				return nil
			},
		},

		{
			Name:  "decode",
			Usage: "Раскрывает информацию из изображения",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "sourcePath",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				file, err := os.Open(cmd.StringArg("sourcePath"))
				if err != nil {
					return err
				}
				defer func() { _ = file.Close() }()

				img, err := png.Decode(file)
				if err != nil {
					return err
				}

				nrgba, ok := img.(*image.NRGBA64)
				if !ok {
					nrgba = image.NewNRGBA64(img.Bounds())
					draw.Draw(nrgba, img.Bounds(), img, img.Bounds().Min, draw.Over)
				}

				s, err := nrgba64.ReadString(nrgba)
				if err != nil {
					return err
				}

				fmt.Println(s)

				return nil
			},
		},
	},
}
