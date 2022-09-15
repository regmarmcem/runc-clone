package main

import (
	"flag"
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
	Action: func(ctx *cli.Context) error {
		fs := flag.FlagSet{}
		fmt.Printf("args %s\n", os.Args)
		if err := fs.Parse(os.Args[2:]); err != nil {
			fmt.Println("Error parsing: ", err)
		}
		fmt.Printf("printf runc-clone init")
		if err := log.InitLogger(true); err != nil {
			l.Fatal(err)
		}
		config.Initialize(os.Args)
		return nil
	},
}
