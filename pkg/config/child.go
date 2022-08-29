package config

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regmarmcem/runc-clone/pkg/log"
	"strconv"
	"syscall"
)

const STACK_SIZE int = 1024 * 1024

func ChildProcess(config ContainerOpts) (cmd *exec.Cmd, err error) {
	err = containerConf(&config)
	if err != nil {
		fmt.Println("Unable to containerConf")
		log.Logger.Infof("Unable to set containerConf %s", err)
		return nil, err
	}
	log.Logger.Info("ChildProcess is started")
	cmd = exec.Command(config.path, config.argv...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	user, err := user.LookupId(strconv.Itoa(int(config.uid)))
	if err != nil {
		fmt.Println("LookupId failed")
		os.Exit(1)
	}

	var gidi int
	gidi, err = strconv.Atoi(user.Gid)

	if err != nil {
		fmt.Println("gid is invalid")
		os.Exit(1)
	}

	gid := uint32(gidi)
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: config.uid, Gid: gid}
	if err = cmd.Start(); err != nil {
		fmt.Println("run process failed!!")
		os.Exit(1)
	}
	log.Logger.Info("process start!!")
	log.Logger.Infof("cmd is %v\n", cmd)
	return cmd, nil
}

func containerConf(config *ContainerOpts) error {
	if err := syscall.Sethostname([]byte(config.Hostname)); err != nil {
		fmt.Printf("Unable to set hostname %s\n", err)
		log.Logger.Infof("Unable to set hostname %s", err)
		return err
	}
	fmt.Println("Succeed in set hostname")
	return nil
}
