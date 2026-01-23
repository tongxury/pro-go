package helper

import (
	"crypto/md5"
	"fmt"
)

func MD5(content []byte) string {
	return fmt.Sprintf("%x", md5.Sum(content))
}

func MD52(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}
