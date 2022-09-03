package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"regmarmcem/runc-clone/pkg/log"
	"regmarmcem/runc-clone/pkg/util"
	"strconv"
	"syscall"
)

const USERNS_OFFSET uint64 = 10000
const USERNS_COUNT uint64 = 2000

func UserNs(fd int, uid uint32) error {
	err := syscall.Unshare(syscall.CLONE_NEWUSER)
	if err != nil {
		log.Logger.Infof("Unable to unshare %s", err)
		util.SendBoolean(fd, false)
		return err
	}
	util.SendBoolean(fd, true)

	var b bool
	b, err = util.RecvBoolean(fd)
	if err != nil {
		log.Logger.Infof("Unable to recvboolean %s", err)
		return err
	}

	if b {
		return errors.New("create namespace")
	}

	user, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		log.Logger.Infof("Unable to lookupid %s", err)
		return err
	}
	g, err := user.GroupIds()
	if err != nil {
		log.Logger.Infof("Unable to get groupids %s", err)
		return err
	}
	var gids []int
	for _, n := range g {
		id, err := strconv.Atoi(n)
		if err != nil {
			log.Logger.Infof("failed to Atoi %s", n)
			continue
		}
		gids = append(gids, id)
	}
	u, err := strconv.Atoi(user.Uid)
	if err != nil {
		log.Logger.Infof("failed to Atoi %s", uid)
		return err
	}

	syscall.Setgroups(gids)
	syscall.Setresgid(gids[0], gids[1], gids[2])
	syscall.Setresuid(u, u, u)

	return nil
}

func HandleChildUidMap(pid int, fd int) error {
	r, err := util.RecvBoolean(fd)
	if err != nil {
		log.Logger.Infof("Unable to recv boolean %s", err)
		return err
	}
	if r {
		// UID/GID map
		uf, err := os.Create(fmt.Sprintf("/proc/%d/%s", pid, "uid_map"))
		if err != nil {
			return errors.New(("NamespaceError(4)"))
		}
		_, err = uf.WriteString(fmt.Sprintf("0 %d %d", USERNS_OFFSET, USERNS_COUNT))
		if err != nil {
			return errors.New("NamespaceError(5)")
		}

		gf, err := os.Create(fmt.Sprintf("/proc/%d/%s", pid, "gid_map"))
		if err != nil {
			return errors.New(("NamespaceError(6)"))
		}
		_, err = gf.WriteString(fmt.Sprintf("0 %d %d", USERNS_OFFSET, USERNS_COUNT))
		if err != nil {
			return errors.New("NamespaceError(7)")
		}

	} else {
		log.Logger.Debug("Child UID/GID map done, sending signal to child to continue...")
		util.SendBoolean(fd, false)
	}
	return nil
}
