package filestore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/galexrt/fivenet/pkg/images"
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
	JobLogos FilePrefix = "joblogos"
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

func (x *File) GetHash() string {
	hasher := sha256.New()
	hasher.Write(x.Data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (x *File) IsImage() bool {
	return filetype.MatchesMap(x.Data, imagesMatchersMap)
}

func (x *File) Optimize(ctx context.Context) error {
	var data io.Reader
	data = bytes.NewReader(x.Data)

	contentType, err := filetype.Match(x.Data)
	if err != nil {
		return err
	}
	x.Extension = &contentType.Extension
	x.ContentType = &contentType.MIME.Value

	if x.IsImage() {
		data, err = images.ResizeImage(contentType, data, 850, 850)
		if err != nil {
			return err
		}

		x.Data, _ = io.ReadAll(data)
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
		if err := st.Delete(ctx, strings.TrimPrefix(*x.Url, FilestoreURLPrefix)); err != nil {
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
