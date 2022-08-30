package config

import (
	"fmt"
	"math/rand"
	"os"
	"regmarmcem/runc-clone/pkg/log"
	"syscall"
)

func SetMountPoint(mountDir string) (err error) {
	fmt.Println("Pivoting root")
	oldRootTail := "oldroot." + RandomStr(6)
	newRoot := "/tmp/runc-clone." + RandomStr(12)
	putOld := newRoot + oldRootTail
	oldRoot := "/" + oldRootTail

	if err := os.Chdir("/"); err != nil {
		return err
	}

	UnmountPath(oldRoot)
	DeleteDir(oldRoot)

	err = CreateDir(putOld)
	if err != nil {
		log.Logger.Infof("Cannot create dir %s\n", err)
		return err
	}

	MountDir("", "/", syscall.MS_REC|syscall.MS_PRIVATE)
	if err != nil {
		log.Logger.Infof("Cannot mount dir %s\n", err)
		return err
	}

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
