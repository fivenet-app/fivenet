package filestore

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"path"
	"strings"

	"github.com/galexrt/fivenet/pkg/storage"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
)

const FilestoreURLPrefix = "/api/filestore/"

var validFilesMap = matchers.Map{
	matchers.TypeJpeg: matchers.Jpeg,
	matchers.TypePng:  matchers.Png,
}

var imagesMatchersMap = matchers.Map{
	matchers.TypeJpeg: matchers.Jpeg,
	matchers.TypePng:  matchers.Png,
}

type FilePrefix = string

const (
	Avatars  FilePrefix = "avatars"
	JobLogos FilePrefix = "job_logos"
	MugShots FilePrefix = "mugshots"
)

func (x *File) Scan(value any) error {
	switch t := value.(type) {
	case string:
		if t != "" {
			x.Url = &t
		}
	case []byte:
		url := string(t)
		if url != "" {
			x.Url = &url
		}
	}
	return nil
}

// Scan implements driver.Valuer for protobuf DocumentAccess.
func (x *File) Value() (driver.Value, error) {
	if x == nil || x.Url == nil {
		return nil, nil
	}

	return *x.Url, nil
}

func (x *File) Upload(ctx context.Context, st storage.IStorage, prefix FilePrefix, fileName string) error {
	if x.Data == nil {
		return fmt.Errorf("no file data given")
	}

	if !filetype.MatchesMap(x.Data, validFilesMap) {
		return fmt.Errorf("invalid file type")
	}

	if x.Url != nil {
		if err := st.Delete(ctx, strings.TrimPrefix(*x.Url, FilestoreURLPrefix)); err != nil {
			return err
		}
	}

	contentType, err := filetype.Match(x.Data)
	if err != nil {
		return err
	}

	fileName = path.Join(prefix, fmt.Sprintf("%s.%s", fileName, contentType.Extension))

	rd := bytes.NewReader(x.Data)
	url, err := st.Put(ctx, fileName, rd, int64(len(x.Data)), contentType.MIME.Value)
	if err != nil {
		return err
	}
	url = FilestoreURLPrefix + url

	x.Url = &url
	x.Data = nil

	return nil
}

func (x *File) IsImage() bool {
	return filetype.MatchesMap(x.Data, imagesMatchersMap)
}
