package config

import "os/exec"

const STACK_SIZE int = 1024 * 1024

func ChildProcess(config ContainerOpts) (pid int, err error) {
	cmd := exec.Command(config.path, config.argv...)
	return cmd.Process.Pid, nil
}
