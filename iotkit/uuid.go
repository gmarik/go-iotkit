package iotkit

import (
	"fmt"
	"os"
)

// uuid implements uuid generation
// http://stackoverflow.com/a/15130965
func uuid() (string, error) {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		return "", err
	}

	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}
