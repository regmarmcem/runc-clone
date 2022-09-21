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

func UserNs(fd *os.File, uid int) error {
	err := syscall.Unshare(syscall.CLONE_NEWUSER)
	if err != nil {
		log.Logger.Infof("Unable to unshare %s", err)
		util.SendBoolean(fd, false)
		return err
	}
	log.Logger.Debug("UserNs SendBoolean")
	util.SendBoolean(fd, true)

	log.Logger.Debug("UserNs RecvBoolean")
	r := util.RecvBoolean(fd)
	if !r {
		return errors.New("create namespace")
	}

	var u *user.User
	u, err = user.LookupId(strconv.Itoa(uid))
	if err != nil {
		log.Logger.Infof("Unable to lookupid %s", err)
		return err
	}
	g, err := u.GroupIds()
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

	var uid_i int
	uid_i, err = strconv.Atoi(u.Uid)
	if err != nil {
		log.Logger.Infof("failed to Atoi %s", uid)
		return err
	}

	syscall.Setgroups(gids)
	syscall.Setresgid(gids[0], gids[1], gids[2])
	syscall.Setresuid(uid_i, uid_i, uid_i)

	return nil
}

func HandleChildUidMap(pid int, fd *os.File) error {

	log.Logger.Debug("HandleChildUidMap")
	r := util.RecvBoolean(fd)
	if r {
		// UID/GID map
		uf, err := os.Create(fmt.Sprintf("/proc/%d/%s", pid, "uid_map"))
		if err != nil {
			log.Logger.Debug("uf create pid, pid_map failed")
			return errors.New(("NamespaceError(4)"))
		}
		_, err = uf.WriteString(fmt.Sprintf("0 %d %d", USERNS_OFFSET, USERNS_COUNT))
		if err != nil {
			log.Logger.Debug("uf USERNS_OFFSET, USERNS_COUNT failed")
			return errors.New("NamespaceError(5)")
		}

		gf, err := os.Create(fmt.Sprintf("/proc/%d/%s", pid, "gid_map"))
		if err != nil {
			log.Logger.Debug("gf create pid, pid_map failed")
			return errors.New(("NamespaceError(6)"))
		}
		_, err = gf.WriteString(fmt.Sprintf("0 %d %d", USERNS_OFFSET, USERNS_COUNT))
		if err != nil {
			log.Logger.Debug("gf USERNS_OFFSET, USERNS_COUNT failed")
			return errors.New("NamespaceError(7)")
		}
		log.Logger.Debug("Namespace creation succeed")
	} else {
		log.Logger.Debug("Child UID/GID map done, sending signal to child to continue...")
		util.SendBoolean(fd, false)
	}
	return nil
}
