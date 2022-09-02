package config

import (
	"errors"
	"regmarmcem/runc-clone/pkg/log"
	"regmarmcem/runc-clone/pkg/util"
	"syscall"
)

func UserNs(fd int, uid uint32) error {
	err := syscall.Unshare(syscall.CLONE_NEWUSER)
	if err != nil {
		log.Logger.Infof("Unable to unshare %s", err)
		return err
	}
	util.SendBoolean(fd, true)

	var b bool
	b, err = util.RecvBoolean(fd)
	if b {
		return errors.New("create namespace")
	}
	return nil
}

func HandleChildUidMap(pid int, fd int) error {
	r, err := util.RecvBoolean(fd)
	if err != nil {
		log.Logger.Infof("Unable to recv boolean %s", err)
		return err
	}
	if r {
		// TODO: UID/GID map
	} else {
		log.Logger.Debug("Child UID/GID map done, sending signal to child to continue...")
		util.SendBoolean(fd, false)
	}
	return nil
}
