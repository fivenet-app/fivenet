package centrumutils

import (
	"fmt"
	"strconv"
	"strings"
)

func IdKey(id uint64) string {
	return strconv.FormatUint(id, 10)
}

func JobIdKey(job string, id uint64) string {
	return job + "." + strconv.FormatUint(id, 10)
}

// ExtractID takes a key like "police.123"  ➜  123
func ExtractID(key string) (uint64, error) {
	idx := strings.LastIndexByte(key, '.')
	if idx < 0 || idx+1 >= len(key) {
		return 0, fmt.Errorf("key %q does not contain a numeric suffix", key)
	}

	return strconv.ParseUint(key[idx+1:], 10, 64)
}

// ExtractIDString takes a key like "police.123"  ➜  123
func ExtractIDString(key string) (string, error) {
	idx := strings.LastIndexByte(key, '.')
	if idx < 0 || idx+1 >= len(key) {
		return "", fmt.Errorf("key %q does not contain a numeric suffix", key)
	}

	return key[idx+1:], nil
}
