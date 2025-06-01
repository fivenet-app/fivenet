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

const FilestoreURLPrefix = "/api/filestore/"

type countingReader struct {
	io.Reader
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.Reader.Read(p)
	c.n += int64(n)
	return n, err
}

// Utility: sanitise + key builder

var sanitizeRegex = regexp.MustCompile(`[^0-9A-Za-z._-]`)

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

func sniff(userCType, name string) string {
	if userCType != "" {
		return userCType
	}
	if typ := mime.TypeByExtension(filepath.Ext(name)); typ != "" {
		return typ
	}
	return "application/octet-stream"
}

func buildKey(ns, name string) string {
	return fmt.Sprintf("%s/%s/%s-%s",
		ns, time.Now().UTC().Format("20060102"), uuid.NewString(), name)
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
