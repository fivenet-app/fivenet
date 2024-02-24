package storage

import (
	"context"
	"io"
	"time"
)

type IStorage interface {
	WithPrefix(prefix string) IStorage

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
