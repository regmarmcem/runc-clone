package config

import (
	"fmt"
	"os"
	"os/exec"
)

const STACK_SIZE int = 1024 * 1024

func ChildProcess(config ContainerOpts) (cmd *exec.Cmd, err error) {
	fmt.Println("ChildProcess is started")
	cmd = exec.Command(config.path, config.argv...)
	if err = cmd.Start(); err != nil {
		fmt.Println("run process failed!!")
		os.Exit(1)
	}
	fmt.Println("start process succeed!!")
	fmt.Printf("cmd is %v\n", cmd)
	return cmd, nil
}
