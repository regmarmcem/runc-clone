package main

import (
	"fmt"
	l "log"
	"os"
	"regmarmcem/runc-clone/pkg/config"
	"regmarmcem/runc-clone/pkg/log"

	"github.com/urfave/cli/v2"
)

var initCommand = &cli.Command{
	Name:        "init",
	Usage:       "initialize a container",
	Description: "initialize a container",
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
		if err := log.InitLogger(true); err != nil {
			l.Fatal(err)
		}
		log.Logger.Debugf("args %s\n", os.Args)
		log.Logger.Debugf("printf runc-clone init")
		config.Initialize(ctx)
		return nil
	},
}
