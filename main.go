package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "runc_clone",
		Usage: "runc",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "mount",
			},
			&cli.IntFlag{
				Name: "uid",
			},
			&cli.BoolFlag{
				Name: "debug",
			},
			&cli.StringFlag{
				Name: "command",
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Printf("Options are %s and %d and %t and %s", ctx.String("mount"), ctx.Int("uid"), ctx.Bool("debug"), ctx.String("command"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
