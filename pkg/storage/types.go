package storage

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"
)

var (
	// ErrNotFound is returned when a file is not found in storage.
	ErrNotFound = errors.New("file not found")
	// ErrInvalidPath is returned when a file path is invalid.
	ErrInvalidPath = errors.New("invalid file path")
)

// IStorage defines the interface for a storage backend.
// Storage defines the interface for a storage backend.
type IStorage interface {
	// WithPrefix returns a storage instance with the given prefix transparently added to all calls.
	WithPrefix(prefix string) (IStorage, error)

	// Get returns the object contents and info for the given key.
	Get(ctx context.Context, key string) (IObject, IObjectInfo, error)
	// Stat returns object info for the given key.
	Stat(ctx context.Context, key string) (IObjectInfo, error)
	// Put uploads a file; size and content type must be accurate.
	Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error)
	// Delete removes a file from storage.
	Delete(ctx context.Context, key string) error

	// List returns a list of files by offset and page size.
	List(ctx context.Context, key string, offset int, pageSize int) ([]*FileInfo, error)

	// GetSpaceUsage returns the total space used in storage.
	GetSpaceUsage(ctx context.Context) (int64, error)
}

// IObject defines the interface for a storage object, supporting read, seek, and random access.
type IObject interface {
	io.ReadSeekCloser
	io.ReaderAt
}

// IObjectInfo defines the interface for object metadata in storage.
type IObjectInfo interface {
	GetName() string
	GetExtension() string
	GetContentType() string
	GetSize() int64
	GetLastModified() time.Time
	GetExpiration() time.Time
}

// ObjectInfo implements IObjectInfo and holds metadata about a storage object.
type ObjectInfo struct {
	// The object's name
	name string
	// The object's file extension
	extension string
	// The object's MIME content type
	contentType string
	// The object's size in bytes
	size int64
	// The object's last modification time
	lastModified time.Time
	// The object's expiration time
	expiration time.Time
}

// GetName returns the object's name.
func (o *ObjectInfo) GetName() string {
	return o.name
}

// GetExtension returns the object's file extension.
func (o *ObjectInfo) GetExtension() string {
	return o.extension
}

// GetContentType returns the object's MIME content type.
func (o *ObjectInfo) GetContentType() string {
	return o.contentType
}

// GetSize returns the object's size in bytes.
func (o *ObjectInfo) GetSize() int64 {
	return o.size
}

// GetExpiration returns the object's expiration time.
func (o *ObjectInfo) GetExpiration() time.Time {
	return o.expiration
}

// GetLastModified returns the object's last modification time.
func (o *ObjectInfo) GetLastModified() time.Time {
	return o.lastModified
}

// GetFileInfo returns a FileInfo struct populated from the ObjectInfo.
func (o *ObjectInfo) GetFileInfo() *FileInfo {
	return &FileInfo{
		Name:         o.name,
		Size:         o.size,
		ContentType:  o.contentType,
		LastModified: o.lastModified,
	}
}

// FileInfo holds basic file metadata for listing operations.
type FileInfo struct {
	// The file's name
	Name string
	// The file's last modification time
	LastModified time.Time
	// The file's size in bytes
	Size int64
	// The file's MIME content type
	ContentType string
}

// GetFilename generates a deterministic storage filename based on user ID, file name, and extension.
// The file is sharded by a hash of the file name for better distribution.
func GetFilename(uid string, fileName string, fileExtension string) string {
	hasher := sha1.New()
	hasher.Write([]byte(fileName))
	fileHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	directory := fileHash[0:2]

	return fmt.Sprintf("%s/%s-%s.%s", directory, uid, fileHash, fileExtension)
}

// FileNameSplitter splits a file name for sharding/storage purposes.
// If the file name is less than 2 characters, it is padded with '0'.
func FileNameSplitter(fileName string) string {
	if len(fileName) < 2 {
		fileName = "0" + fileName
	}

	chars := fileName
	part1 := chars[0:1]
	return filepath.Join(part1, fileName)
}
