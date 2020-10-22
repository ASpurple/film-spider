package tools

import (
	"crypto/md5"
	"fmt"
)

// Md ..
func Md(txt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(txt)))
}
