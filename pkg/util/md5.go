package util

import (
	"crypto/md5"
	"fmt"
)

// Get md5 string
func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
