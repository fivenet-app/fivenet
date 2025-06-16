package utils

import (
	"path/filepath"
)

// CleanFilePath returns a cleaned version of the given file path and a boolean indicating validity.
// The path is considered valid if it is local or absolute. Returns false if not.
func CleanFilePath(filePath string) (string, bool) {
	if !filepath.IsLocal(filePath) && !filepath.IsAbs(filePath) {
		return "", false
	}

	return filepath.Clean(filePath), true
}
