package main

import (
	"fmt"
	l "log"
	"regmarmcem/runc-clone/pkg/config"
	"regmarmcem/runc-clone/pkg/log"

	"github.com/urfave/cli/v2"
)

var runCommand = &cli.Command{
	Name:        "run",
	Usage:       "run a container",
	Description: `Run a container.`,
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
		fmt.Printf("Options are %s and %d and %t and %s\n", ctx.String("mount"), ctx.Int("uid"), ctx.Bool("debug"), ctx.String("command"))
		if err := log.InitLogger(ctx.Bool("debug")); err != nil {
			l.Fatal(err)
		}
		fmt.Printf("DebugOption is %t\n", log.DebugOption)
		config.Start(ctx)
		return nil
	},
}
