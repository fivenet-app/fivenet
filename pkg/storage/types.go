package storage

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"time"
)

var (
	ErrNotFound    = errors.New("file not found")
	ErrInvalidPath = errors.New("invalid file path")
)

type IStorage interface {
	// Return storage with prefix transparently added to the calls
	WithPrefix(prefix string) (IStorage, error)

	// Get return object contents and info
	Get(ctx context.Context, filePath string) (IObject, IObjectInfo, error)
	// Get URL of object (not every storage adapter supports this)
	GetURL(ctx context.Context, filePath string, expires time.Duration, reqParams url.Values) (*string, error)
	// Return object info
	Stat(ctx context.Context, filePath string) (IObjectInfo, error)
	// Upload file, size and content type must be accurate
	Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error)
	// Delete file
	Delete(ctx context.Context, filePath string) error

	// List files by offset and page size
	List(ctx context.Context, filePath string, offset int, pageSize int) ([]*FileInfo, error)
}

type IObject interface {
	io.ReadSeekCloser
	io.ReaderAt
}

type IObjectInfo interface {
	GetName() string
	GetExtension() string
	GetContentType() string
	GetSize() int64
	GetLastModified() time.Time
	GetExpiration() time.Time
}

type ObjectInfo struct {
	name         string
	extension    string
	contentType  string
	size         int64
	lastModified time.Time
	expiration   time.Time
}

func (o *ObjectInfo) GetName() string {
	return o.name
}

func (o *ObjectInfo) GetExtension() string {
	return o.extension
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

func (o *ObjectInfo) GetLastModified() time.Time {
	return o.lastModified
}

func (o *ObjectInfo) GetFileInfo() *FileInfo {
	return &FileInfo{
		Name:         o.name,
		Size:         o.size,
		ContentType:  o.contentType,
		LastModified: o.lastModified,
	}
}

type FileInfo struct {
	Name         string
	LastModified time.Time
	Size         int64
	ContentType  string
}

func GetFilename(uid string, fileName string, fileExtension string) string {
	hasher := sha1.New()
	hasher.Write([]byte(fileName))
	fileHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	directory := fileHash[0:2]

	return fmt.Sprintf("%s/%s-%s.%s", directory, uid, fileHash, fileExtension)
}

func FileNameSplitter(fileName string) string {
	if len(fileName) < 2 {
		fileName = "0" + fileName
	}

	chars := fileName
	part1 := chars[0:1]
	return filepath.Join(part1, fileName)
}
