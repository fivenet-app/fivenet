package query

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"text/template"
)

// Copied from https://github.com/golang-migrate/migrate/pull/793#issuecomment-1954565692

type templateFile struct {
	fs.File

	data any
}

func (t *templateFile) Read(p []byte) (int, error) {
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

	data any
}

func (t *templateFS) Open(name string) (fs.File, error) {
	file, err := t.FS.Open(name)
	if err != nil {
		return nil, err
	}

	return &templateFile{
		File: file,
		data: t.data,
	}, nil
}
