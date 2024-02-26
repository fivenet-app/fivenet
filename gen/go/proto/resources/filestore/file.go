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
	"golang.org/x/exp/slices"
)

const FilestoreURLPrefix = "/api/filestore/"

var AllowedFileExtensions = []string{"jpg", "jpeg", "png", "webp"}

type FilePrefix = string

const (
	Avatars  FilePrefix = "avatars"
	JobLogos FilePrefix = "job_logos"
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
	if x.Type == nil {
		return fmt.Errorf("no file type given")
	}
	if x.Data == nil {
		return fmt.Errorf("no file data given")
	}

	if !slices.Contains(AllowedFileExtensions, *x.Type) {
		return fmt.Errorf("disallowed file extension")
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
	return filetype.IsImage(x.Data)
}
