package util

import (
	"syscall"
)

func GenerateSocketPair() (fds [2]int, err error) {
	fds, err = syscall.Socketpair(syscall.AF_LOCAL, syscall.FD_CLOEXEC, 0)
	return fds, err
}
