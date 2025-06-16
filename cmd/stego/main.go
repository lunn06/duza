package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/lunn06/stego/read"
	"github.com/lunn06/stego/write"
)

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
				return write.Write(
					cmd.StringArg("sourcePath"),
					cmd.StringArg("secret"),
					cmd.StringArg("outPath"),
				)
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
				secret, err := read.Read(
					cmd.StringArg("sourcePath"),
				)
				if err != nil {
					return err
				}

				fmt.Print(secret)

				return nil
			},
		},
	},
}

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
