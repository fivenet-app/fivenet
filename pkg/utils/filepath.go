package utils

import (
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// CleanFilePath returns a cleaned version of the given file path and a boolean indicating validity.
// The path is considered valid if it is absolute. Returns false if not.
// Is for file paths and not directories.
func CleanFilePath(filePath string) (string, bool) {
	// Reject directory paths or multi slashes
	if strings.HasPrefix(filePath, "//") || strings.HasSuffix(filePath, "/") {
		return "", false
	}

	if !filepath.IsAbs(filePath) || !utf8.ValidString(filePath) {
		return "", false
	}

	return filepath.Clean(filePath), true
}
