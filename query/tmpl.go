package query

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// Copied from https://github.com/golang-migrate/migrate/pull/793#issuecomment-1954565692

type templateFile struct {
	fs.File

	data map[string]any
}

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

	if migrationNumber < 1748173399 && !t.data["ESXCompat"].(bool) { // big_rename migration
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

type templateFS struct {
	embed.FS

	data map[string]any
}

func (t *templateFS) Open(name string) (fs.File, error) {
	file, err := t.FS.Open(name)
	if err != nil {
		return nil, err
	}

	dataCopy := make(map[string]any)
	for k, v := range t.data {
		dataCopy[k] = v
	}

	return &templateFile{
		File: file,
		data: dataCopy,
	}, nil
}
