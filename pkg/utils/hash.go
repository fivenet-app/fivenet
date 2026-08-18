package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetSHA256HashFromString returns the SHA256 hash of the input string as a hexadecimal string.
func GetSHA256HashFromString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
