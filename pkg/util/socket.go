package util

import (
	"os"
	"syscall"
)

func GenerateSocketPair() (fdfiles [2]*os.File, err error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.FD_CLOEXEC, 0)
	fdfiles[0] = os.NewFile(uintptr(fds[0]), "parent-fd")
	fdfiles[1] = os.NewFile(uintptr(fds[1]), "child-fd")
	return fdfiles, err
}
