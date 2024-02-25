package storage

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

type IStorage interface {
	WithPrefix(prefix string) (IStorage, error)

	Get(ctx context.Context, filePath string) (IObject, IObjectInfo, error)
	Stat(ctx context.Context, filePath string) (IObjectInfo, error)
	Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, filePath string) error
}

type IObject interface {
	Read(p []byte) (n int, err error)
	ReadAt(p []byte, off int64) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}

type IObjectInfo interface {
	GetContentType() string
	GetSize() int64
	GetExpiration() time.Time
}

type ObjectInfo struct {
	contentType string
	size        int64
	expiration  time.Time
}

func (o *ObjectInfo) GetContentType() string {
	return o.contentType
}

func (o *ObjectInfo) GetSize() int64 {
	return o.size
}

func (o *ObjectInfo) GetExpiration() time.Time {
	return o.expiration
}

func GetFilename(uid string, fileName string, fileExtension string) string {
	hasher := sha1.New()
	hasher.Write([]byte(fileName))
	fileHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	directory := fileHash[0:2]

	return fmt.Sprintf("%s/%s-%s.%s", directory, uid, fileHash, fileExtension)
}
