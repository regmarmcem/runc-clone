package util

import (
	"net"
	"os"
	"regmarmcem/runc-clone/pkg/log"
)

func SendBoolean(fd int, boolean bool) (err error) {
	var data []byte
	var conn net.Conn

	f := os.NewFile(uintptr(fd), "")
	if conn, err = net.FileConn(f); err != nil {
		log.Logger.Infof("Failed to FileConn: %s", err)
		return err
	}
	defer conn.Close()

	if boolean {
		data = []byte("1")
	} else {
		data = []byte("0")
	}

	if _, err = conn.Write(data); err != nil {
		log.Logger.Infof("Failed to send data: %s", err)
		return err
	}

	return nil
}

func RecvBoolean(fd int) (_ bool, err error) {
	data := []byte("0")
	var conn net.Conn

	f := os.NewFile(uintptr(fd), "")
	if conn, err = net.FileConn(f); err != nil {
		log.Logger.Infof("Failed to FileConn: %s", err)
		return false, err
	}
	defer conn.Close()

	if _, err = conn.Read(data); err != nil {
		log.Logger.Infof("Failed to read data: %s", err)
		return false, err
	}

	return true, nil
}
