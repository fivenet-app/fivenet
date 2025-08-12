//nolint:gosec // MD5 is not used for security purposes but for simple hashing of strings.
package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5HashFromString returns the MD5 hash of the input string as a hexadecimal string.
func GetMD5HashFromString(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
