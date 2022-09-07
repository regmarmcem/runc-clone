package config

import (
	"math/rand"
	"os"
	"regmarmcem/runc-clone/pkg/log"
	"syscall"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func SetMountPoint(mountDir string) (err error) {
	log.Logger.Debug("Investigating...")
	f := mountDir + "/sbin/pivot_root"
	if !Exists(f) {
		log.Logger.Debug("pivot_root not found")
		os.Exit(1)
	}

	log.Logger.Debug("Setting mount points..")
	newRoot := "/tmp/runc-clone." + RandomStr(12)
	log.Logger.Debugf("Mounting temp directory %s", newRoot)
	CreateDir(newRoot)
	err = MountDir(mountDir, newRoot, uintptr(syscall.MS_BIND|syscall.MS_REC|syscall.MS_PRIVATE))
	if err != nil {
		log.Logger.Infof("Cannot mount dir %s\n", err)
		os.Exit(1)
		return err
	}

	log.Logger.Debug("Pivoting root")
	oldRootTail := "oldroot." + RandomStr(6)
	putOld := newRoot + "/" + oldRootTail

	err = CreateDir(putOld)
	if err != nil {
		log.Logger.Infof("Cannot create dir %s\n", err)
		os.Exit(1)
		return err
	}

	log.Logger.Debugf("newRoot: %s, putOld: %s\n", newRoot, putOld)
	err = syscall.PivotRoot(newRoot, putOld)
	if err != nil {
		log.Logger.Infof("Cannot pivot root %s\n", err)
		os.Exit(1)
		return err
	}

	log.Logger.Debug("Unmounting old root")
	oldRoot := "/" + oldRootTail
	err = syscall.Chdir("/")
	if err != nil {
		log.Logger.Info("MountError")
		os.Exit(1)
		return err
	}

	UnmountPath(oldRoot)
	DeleteDir(oldRoot)

	log.Logger.Debug("Unmount finished")

	return nil
}

func CleanMounts(mountDir string) (err error) {
	return nil
}

func MountDir(path string, mountPoint string, flags uintptr) (err error) {
	if err = syscall.Mount(path, mountPoint, "", flags, ""); err != nil {
		log.Logger.Infof("Cannot mount %s to %s: %s\n", path, mountPoint, err)
		return err
	}
	log.Logger.Debug("Mount succeess!!")
	return nil
}

func RandomStr(n uint64) string {
	charSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

func CreateDir(path string) (err error) {
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Logger.Infof("Cannot create dir %s\n", err)
		return err
	}
	return nil
}

func UnmountPath(path string) (err error) {
	err = syscall.Unmount(path, syscall.MNT_DETACH)
	if err != nil {
		log.Logger.Infof("Unable to unmount %s: %s\n", path, err)
		return err
	}
	return nil
}

func DeleteDir(path string) (err error) {
	err = os.Remove(path)
	if err != nil {
		log.Logger.Infof("Unable to delete %s: %s\n", path, err)
		return err
	}
	return nil
}
