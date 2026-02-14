package utils

import (
	"errors"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

var (
	ErrEmptyPathSegment = errors.New("empty path segment not allowed")
	ErrEmptyPath        = errors.New("empty path")
)

func CleanStoragePath(input string, emptyOk bool) (string, error) {
	// Forbid invalid UTF-8
	if !utf8.ValidString(input) {
		return "", errors.New("invalid utf-8")
	}
	// Forbid NUL byte (defensive check)
	if strings.ContainsRune(input, 0) {
		return "", errors.New("nul byte not allowed")
	}
	// Forbid ambiguous separators
	if strings.Contains(input, `\`) {
		return "", errors.New("backslash not allowed")
	}
	// Forbid UNC paths
	if strings.HasPrefix(input, `\\`) || strings.HasPrefix(input, "//") {
		return "", errors.New("unc paths not allowed")
	}
	// Forbid Windows drive paths (e.g., C:\ or C:/)
	if hasWindowsDrivePrefix(input) {
		return "", errors.New("windows drive paths not allowed")
	}

	// Trim leading slash to treat as relative path
	input = strings.TrimPrefix(input, "/")

	// Disallow any empty and ".." path segments
	parts := strings.SplitSeq(input, "/")
	for p := range parts {
		if p == "" && !emptyOk {
			return "", ErrEmptyPathSegment
		}
		if p == ".." {
			return "", errors.New("dotdot segment not allowed")
		}
	}

	// Normalize as URL path always uses '/'
	clean := path.Clean("/" + input)
	clean = strings.TrimPrefix(clean, "/")

	if clean == "" || clean == "." {
		if !emptyOk {
			return "", ErrEmptyPath
		}
	}

	// Prevent escaping: Anything starting with ".." is traversal
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return "", errors.New("path traversal detected")
	}
	// Re-add trailing slash if it was present in the input and empty segments are allowed
	if emptyOk && strings.HasSuffix(input, "/") {
		clean += "/"
	}

	return clean, nil
}

func CleanStorageKey(input string) (string, error) {
	return CleanStoragePath(input, false)
}

func hasWindowsDrivePrefix(s string) bool {
	if len(s) < 2 {
		return false
	}
	c := s[0]
	if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
		return false
	}
	return s[1] == ':'
}

func FSRootFile(prefix string, key string) (string, error) {
	return FSRootPath(prefix, key, false)
}

func FSRootPath(prefix string, key string, emptyOk bool) (string, error) {
	if key == "" || key == "." || key == ".." {
		if emptyOk {
			return "", nil
		}
		return "", errors.New("empty key")
	}
	// Forbid NUL byte (defensive check)
	if strings.ContainsRune(key, 0) {
		return "", errors.New("nul byte not allowed")
	}
	// Forbid ambiguous separators
	if strings.Contains(key, `\`) {
		return "", errors.New("backslash not allowed")
	}

	// Disallow any empty and ".." path segments
	parts := strings.SplitSeq(key, "/")
	for p := range parts {
		if p == "" && !emptyOk {
			return "", ErrEmptyPathSegment
		}
		if p == ".." {
			return "", errors.New("dotdot segment not allowed")
		}
	}

	// key is canonical slash-based; convert to OS separators and join with prefix
	rel := filepath.Join(prefix, filepath.FromSlash(key))

	// Ensure rel is still relative (no absolute, no volume).
	if filepath.IsAbs(rel) || filepath.VolumeName(rel) != "" {
		return "", errors.New("result must be relative")
	}

	// Disallow leading ".." segments in the final relative path
	clean := filepath.Clean(rel)
	if clean == "." && !emptyOk {
		return "", errors.New("invalid path")
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("path traversal detected")
	}

	// Re-add trailing slash if it was present in the input and empty segments are allowed
	if emptyOk && strings.HasSuffix(key, "/") {
		clean += "/"
	}

	return clean, nil
}
