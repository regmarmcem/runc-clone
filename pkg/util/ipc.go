package util

import (
	"os"
	"regmarmcem/runc-clone/pkg/log"
	"strconv"
)

func SendBoolean(fd *os.File, data bool) (err error) {

	log.Logger.Debugf("sent data is %s", data)
	if _, err = fd.WriteString(strconv.FormatBool(data)); err != nil {
		log.Logger.Infof("Failed to send data: %s", err)
		return err
	}

	return nil
}

func RecvBoolean(fd *os.File) bool {
	buf := make([]byte, 1)

	if _, err := fd.Read(buf); err != nil {
		log.Logger.Infof("Failed to read data: %s", err)
		return false
	}

	b, err := strconv.ParseBool(string(buf))
	log.Logger.Debugf("received value is %s", b)
	if err != nil {
		log.Logger.Infof("Failed to read data: %s", err)
		return false
	}
	return b
}
