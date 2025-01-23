package query

import (
	"bytes"
	"embed"
	"io/fs"
	"text/template"
)

// Copied from https://github.com/golang-migrate/migrate/pull/793#issuecomment-1954565692

type templateFile struct {
	fs.File

	data any
}

func (t *templateFile) Read(name []byte) (int, error) {
	if _, err := t.File.Read(name); err != nil {
		return 0, err
	}

	tmpl, err := template.New("").Parse(string(name))
	if err != nil {
		return 0, err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t.data)
	if err != nil {
		return 0, err
	}

	return copy(name, buf.Bytes()), nil
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
