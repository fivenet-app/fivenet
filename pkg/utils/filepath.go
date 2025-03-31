package utils

import (
	"path/filepath"
)

func CleanFilePath(filePath string) (string, bool) {
	if !filepath.IsLocal(filePath) && !filepath.IsAbs(filePath) {
		return "", false
	}

	return filepath.Clean(filePath), true
}
