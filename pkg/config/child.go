package config

import (
	"os"
	"os/exec"
	"regmarmcem/runc-clone/pkg/log"
	"syscall"
)

const STACK_SIZE int = 1024 * 1024

func createDirsIfDontExist(dirs []string) error {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				log.Logger.Infof("Error creating directory: %v\n", err)
				return err
			}
		}
	}
	return nil
}

func ChildProcess(config ContainerOpts) (cmd *exec.Cmd, err error) {

	log.Logger.Debugf("config is %t", config)
	if err != nil {
		log.Logger.Infof("Unable to set containerConf %s", err)
		return nil, err
	}
	// close one of sockpair: config.fd is sockets[1]
	// err = syscall.Close(config.fd)
	if err != nil {
		log.Logger.Infof("Unable to close fd %s", err)
		return nil, err
	}

	log.Logger.Debugf("config.argv= %s", config.argv)
	args := append([]string{"init"}, config.path)
	args = append(args, config.argv...)
	cmd = exec.Command("/proc/self/exe", args...)
	log.Logger.Debugf("cmd %s", cmd)
	cmd.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS,
	}

	// user, err := user.LookupId(strconv.Itoa(int(config.uid)))
	// if err != nil {
	// log.Logger.Info("LookupId failed")
	// os.Exit(1)
	// }

	// var gidi int
	// gidi, err = strconv.Atoi(user.Gid)

	// if err != nil {
	// log.Logger.Info("gid is invalid")
	// os.Exit(1)
	// }

	// gid := uint32(gidi)
	// cmd.SysProcAttr = &syscall.SysProcAttr{}
	// cmd.SysProcAttr.Credential = &syscall.Credential{Uid: config.uid, Gid: gid}

	log.Logger.Infof("cmd is %v\n", cmd)
	log.Logger.Infof("cmd sysprocattr %v\n", cmd.SysProcAttr.Cloneflags)
	if err = cmd.Start(); err != nil {
		log.Logger.Infof("run process failed %s", err)
		os.Exit(1)
	}
	log.Logger.Info("process start!!")
	return cmd, nil
}

func ExecProcess(c ContainerOpts) (cmd *exec.Cmd, err error) {

	log.Logger.Debugf("config is %t", c)
	err = containerConf(c)
	if err != nil {
		log.Logger.Infof("Unable to set containerConf %s", err)
		return nil, err
	}
	// close one of sockpair: config.fd is sockets[1]
	// err = syscall.Close(config.fd)
	if err != nil {
		log.Logger.Infof("Unable to close fd %s", err)
		return nil, err
	}

	log.Logger.Info("ChildProcess is started")
	log.Logger.Debugf("config.path is %s", c.path)
	log.Logger.Debugf("config.argv is %s", c.argv)
	cmd = exec.Command(c.path, c.argv...)
	log.Logger.Debugf("cmd %s", cmd)
	cmd.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS,
	}

	// user, err := user.LookupId(strconv.Itoa(int(config.uid)))
	// if err != nil {
	// log.Logger.Info("LookupId failed")
	// os.Exit(1)
	// }

	// var gidi int
	// gidi, err = strconv.Atoi(user.Gid)

	// if err != nil {
	// log.Logger.Info("gid is invalid")
	// os.Exit(1)
	// }

	// gid := uint32(gidi)
	// cmd.SysProcAttr = &syscall.SysProcAttr{}
	// cmd.SysProcAttr.Credential = &syscall.Credential{Uid: config.uid, Gid: gid}

	log.Logger.Infof("cmd is %v\n", cmd)
	log.Logger.Infof("cmd sysprocattr %v\n", cmd.SysProcAttr.Cloneflags)
	if err = cmd.Start(); err != nil {
		log.Logger.Infof("run process failed %s", err)
		os.Exit(1)
	}
	log.Logger.Info("process start!!")
	return cmd, nil

}

func containerConf(config ContainerOpts) error {
	if err := syscall.Sethostname([]byte(config.Hostname)); err != nil {
		log.Logger.Infof("Unable to set hostname %s", err)
		return err
	}
	if err := SetMountPoint(config.MountDir); err != nil {
		log.Logger.Infof("Unable to set mount point %s", err)
		return err
	}
	log.Logger.Info("succeed in set hostname")
	// UserNs(config.fd, config.uid)
	return nil
}
