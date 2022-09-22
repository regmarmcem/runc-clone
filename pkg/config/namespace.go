package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"regmarmcem/runc-clone/pkg/log"
	"strconv"
	"syscall"
)

const USERNS_OFFSET uint64 = 10000
const USERNS_COUNT uint64 = 2000

func UserNs(fd *os.File, uid int) error {

	log.Logger.Debug("UserNs RecvBoolean")
	// r := util.RecvBoolean(fd)
	// if !r {
	// return errors.New("create namespace")
	// }

	var u *user.User
	u, err := user.LookupId(strconv.Itoa(uid))
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
	syscall.Setresgid(gids[0], gids[0], gids[0])
	syscall.Setresuid(uid_i, uid_i, uid_i)

	return nil
}

func HandleChildUidMap(pid int, fd *os.File) error {

	log.Logger.Debug("HandleChildUidMap")
	// util.SendBoolean(fd, true)
	// r := util.RecvBoolean(fd)
	// log.Logger.Debugf("received handle is: %s", r)
	// UID/GID map
	err := os.Mkdir(fmt.Sprintf("/proc/%d", pid), 0555)
	if err != nil {
		log.Logger.Debugf("mkdir /proc/$$ failed: %s", err)
		return errors.New(("NamespaceError(4)"))
	}

	uf, err := os.Create(fmt.Sprintf("/proc/%d/%s", pid, "uid_map"))
	if err != nil {
		log.Logger.Debugf("uf create pid, pid_map failed: %s", err)
		return errors.New(("NamespaceError(4)"))
	}
	_, err = uf.WriteString(fmt.Sprintf("0 %d %d", USERNS_OFFSET, USERNS_COUNT))
	if err != nil {
		log.Logger.Debugf("uf USERNS_OFFSET, USERNS_COUNT failed %s", err)
		return errors.New("NamespaceError(5)")
	}

	gf, err := os.Create(fmt.Sprintf("/proc/%d/%s", pid, "gid_map"))
	if err != nil {
		log.Logger.Debugf("gf create pid, pid_map failed %s", err)
		return errors.New(("NamespaceError(6)"))
	}
	_, err = gf.WriteString(fmt.Sprintf("0 %d %d", USERNS_OFFSET, USERNS_COUNT))
	if err != nil {
		log.Logger.Debugf("gf USERNS_OFFSET, USERNS_COUNT failed %s", err)
		return errors.New("NamespaceError(7)")
	}
	log.Logger.Debug("Namespace creation succeed")
	// util.SendBoolean(fd, true)
	return nil
}
