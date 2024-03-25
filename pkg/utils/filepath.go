package utils

import (
	"path/filepath"
)

func CleanFilePath(filePath string) (string, bool) {
	if !filepath.IsLocal(filePath) {
		return "", false
	}

	return filepath.Clean(filePath), true
}
