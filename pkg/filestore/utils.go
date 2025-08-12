// Package filestore provides file storage utilities and helpers.
package filestore

import (
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// FilestoreURLPrefix is the URL prefix for accessing files in the filestore API.
const FilestoreURLPrefix = "/api/filestore/"

// countingReader wraps an io.Reader and counts the number of bytes read.
type countingReader struct {
	// Reader is the underlying io.Reader.
	io.Reader

	// n is the total number of bytes read so far.
	n int64
}

// Read reads from the underlying reader and increments the byte count.
func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.Reader.Read(p)
	c.n += int64(n)
	return n, err
}

// sanitizeRegex matches any character not allowed in sanitized file names.
var sanitizeRegex = regexp.MustCompile(`[^0-9A-Za-z._-]`)

// SanitizeFileName returns a sanitized version of the file name, replacing disallowed characters
// and truncating the base name to a maximum of 240 characters.
func SanitizeFileName(name string) string {
	name = filepath.Base(name)
	ext := filepath.Ext(name)
	base := name[:len(name)-len(ext)]
	base = sanitizeRegex.ReplaceAllString(base, "_")
	if len(base) > 240 {
		base = base[:240]
	}
	return base + ext
}

// sniff determines the MIME type for a file, preferring the user-provided type,
// then falling back to the file extension, and finally defaulting to application/octet-stream.
func sniff(userCType, name string) string {
	if userCType != "" {
		return userCType
	}
	if typ := mime.TypeByExtension(filepath.Ext(name)); typ != "" {
		return typ
	}
	return "application/octet-stream"
}

// buildKey generates a unique storage key for a file using the namespace, current UTC date,
// a new UUID, and the sanitized file name.
func buildKey(ns, name string) string {
	return fmt.Sprintf("%s/%s/%s-%s",
		ns, time.Now().UTC().Format("20060102"), uuid.NewString(), name)
}
