package main

import (
	"os"

	"regmarmcem/runc-clone/pkg/log"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "runc_clone",
		Usage: "runc",
		Commands: []*cli.Command{
			runCommand,
			initCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
