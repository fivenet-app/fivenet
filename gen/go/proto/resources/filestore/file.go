package filestore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"path"
	"strings"

	"github.com/fivenet-app/fivenet/pkg/images"
	"github.com/fivenet-app/fivenet/pkg/storage"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
)

const FilestoreURLPrefix = "/api/filestore/"

func StripURLPrefix(in string) string {
	return strings.TrimPrefix(in, FilestoreURLPrefix)
}

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
	Avatars   FilePrefix = "avatars"
	JobLogos  FilePrefix = "joblogos"
	MugShots  FilePrefix = "mugshots"
	JobAssets FilePrefix = "jobassets"

	QualificationExamAssets FilePrefix = "qualifications_exam"
)

// Scan implements driver.Valuer for protobuf DocumentAccess.
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

// Value marshals the value into driver.Valuer.
func (x *File) Value() (driver.Value, error) {
	if x == nil || x.Url == nil {
		return nil, nil
	}

	return *x.Url, nil
}

func (x *File) GetHash() string {
	hasher := sha256.New()
	hasher.Write(x.Data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (x *File) IsImage() bool {
	return filetype.MatchesMap(x.Data, imagesMatchersMap)
}

func (x *File) MatchFileExt() error {
	contentType, err := filetype.Match(x.Data)
	if err != nil {
		return err
	}
	x.Extension = &contentType.Extension
	x.ContentType = &contentType.MIME.Value

	return nil
}

func (x *File) Optimize(ctx context.Context) error {
	if x.Extension == nil {
		if err := x.MatchFileExt(); err != nil {
			return err
		}
	}

	if x.IsImage() {
		data, err := images.ResizeImage(*x.Extension, bytes.NewReader(x.Data), 850, 850)
		if err != nil {
			return err
		}
		if data != nil {
			x.Data = data
		}
	}

	return nil
}

func (x *File) Upload(ctx context.Context, st storage.IStorage, prefix FilePrefix, fileName string) error {
	if x.Data == nil {
		return fmt.Errorf("no file data given")
	}

	if !filetype.MatchesMap(x.Data, validFilesMap) {
		return fmt.Errorf("invalid file type")
	}

	if x.Url != nil {
		if err := st.Delete(ctx, StripURLPrefix(*x.Url)); err != nil {
			return err
		}
	}

	if x.Extension == nil {
		if err := x.MatchFileExt(); err != nil {
			return err
		}
	}

	fileName = path.Join(prefix, fmt.Sprintf("%s.%s", fileName, *x.Extension))
	url, err := st.Put(ctx, fileName, bytes.NewReader(x.Data), int64(len(x.Data)), *x.ContentType)
	if err != nil {
		return err
	}
	url = FilestoreURLPrefix + url

	x.Url = &url
	x.Data = nil

	return nil
}
