package centrumutils

import (
	"fmt"
	"strconv"
	"strings"
)

func IdKey(id int64) string {
	return strconv.FormatInt(id, 10)
}

func JobIdKey(job string, id int64) string {
	return job + "." + strconv.FormatInt(id, 10)
}

// ExtractID takes a key like "police.123" ➜ 123.
func ExtractID(key string) (int64, error) {
	idx := strings.LastIndexByte(key, '.')
	if idx < 0 || idx+1 >= len(key) {
		return 0, fmt.Errorf("key %q does not contain a numeric suffix", key)
	}

	return strconv.ParseInt(key[idx+1:], 10, 64)
}

// ExtractIDString takes a key like "police.123" ➜ 123.
func ExtractIDString(key string) (string, error) {
	idx := strings.LastIndexByte(key, '.')
	if idx < 0 || idx+1 >= len(key) {
		return "", fmt.Errorf("key %q does not contain a numeric suffix", key)
	}

	return key[idx+1:], nil
}
