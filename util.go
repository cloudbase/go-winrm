package winrm

import (
	"crypto/rand"
	"fmt"
	"io"
)

func Uuid() (string, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}
