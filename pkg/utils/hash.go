package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5HashFromString(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
