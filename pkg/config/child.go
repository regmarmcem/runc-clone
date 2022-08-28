package config

import (
	"fmt"
	"os"
	"os/exec"
	"regmarmcem/runc-clone/pkg/log"
)

const STACK_SIZE int = 1024 * 1024

func ChildProcess(config ContainerOpts) (cmd *exec.Cmd, err error) {
	log.Logger.Info("ChildProcess is started")
	cmd = exec.Command(config.path, config.argv...)
	if err = cmd.Start(); err != nil {
		fmt.Println("run process failed!!")
		os.Exit(1)
	}
	log.Logger.Info("process start!!")
	log.Logger.Infof("cmd is %v\n", cmd)
	return cmd, nil
}
