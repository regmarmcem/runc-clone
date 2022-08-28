package util

import (
	"net"
	"os"
	"regmarmcem/runc-clone/pkg/log"
)

func SendBoolean(f *os.File, boolean bool) (err error) {
	var data []byte
	var conn net.Conn

	if conn, err = net.FileConn(f); err != nil {
		log.Logger.Infof("Failed to FileConn: %s", err)
		return err
	}

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

func RecvBoolean(f *os.File) (_ bool, err error) {
	data := []byte("0")
	var conn net.Conn

	if conn, err = net.FileConn(f); err != nil {
		log.Logger.Infof("Failed to FileConn: %s", err)
		return false, err
	}

	if _, err = conn.Read(data); err != nil {
		log.Logger.Infof("Failed to read data: %s", err)
		return false, err
	}

	return true, nil
}
