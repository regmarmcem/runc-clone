package main

import (
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
	Action: func(ctx *cli.Context) error {
		if err := log.InitLogger(true); err != nil {
			l.Fatal(err)
		}
		log.Logger.Debugf("args %s\n", os.Args)
		log.Logger.Debugf("printf runc-clone init")
		config.Initialize(os.Args[2:])
		return nil
	},
}
