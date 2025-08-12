package query

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// Copied and adapted from https://github.com/golang-migrate/migrate/pull/793#issuecomment-1954565692
// Provides template-based file and FS wrappers for migration logic.

// templateFile wraps an fs.File and applies Go text/template processing on Read.
type templateFile struct {
	// File is the underlying file being wrapped.
	File fs.File
	// data is the template data context used for template execution.
	data map[string]any
}

// Read reads the file, applies template processing, and returns the rendered bytes.
// It also supports migration-specific logic for ESX compatibility and migration number.
func (t *templateFile) Read(p []byte) (int, error) {
	st, err := t.File.Stat()
	if err != nil {
		return 0, err
	}

	fileName := filepath.Base(st.Name())
	split := strings.Split(fileName, "_")
	migrationNumber, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse migration number from file name %s: %v", fileName, err))
	}

	// Migration-specific logic: set UsersTableName for pre-big_rename migrations if not ESXCompat.
	esxCompat, ok := t.data["ESXCompat"].(bool)
	if migrationNumber <= 1747999113 && (ok && !esxCompat) { // big_rename migration
		t.data["UsersTableName"] = "fivenet_users"
	}

	out, err := io.ReadAll(t.File)
	if err != nil {
		return 0, err
	}
	if len(out) == 0 {
		return 0, io.EOF
	}

	tmpl, err := template.New("").Parse(string(out))
	if err != nil {
		return 0, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, t.data); err != nil {
		return 0, err
	}

	return copy(p, buf.Bytes()), nil
}

// Close closes the underlying file. Satisfies fs.File interface.
func (t *templateFile) Close() error {
	return t.File.Close()
}

// Stat returns the FileInfo structure describing file. Satisfies fs.File interface.
func (t *templateFile) Stat() (fs.FileInfo, error) {
	return t.File.Stat()
}

// templateFS wraps an embed.FS and injects template data for use in templateFile.
type templateFS struct {
	// FS is the embedded filesystem being wrapped.
	embed.FS

	// data is the template data context used for all opened files.
	data map[string]any
}

// Open opens a file from the embedded FS and returns a templateFile with template data context.
func (t *templateFS) Open(name string) (fs.File, error) {
	file, err := t.FS.Open(name)
	if err != nil {
		return nil, err
	}

	// Copy template data to avoid mutation between files.
	dataCopy := make(map[string]any)
	maps.Copy(dataCopy, t.data)

	return &templateFile{
		File: file,
		data: dataCopy,
	}, nil
}
