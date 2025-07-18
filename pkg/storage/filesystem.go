package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/h2non/filetype"
	"go.uber.org/fx"
)

func init() {
	storageFactories[config.StorageTypeFilesystem] = NewFilesystem
}

// Filesystem implements IStorage for local filesystem-based storage.
type Filesystem struct {
	IStorage

	// basePath is the root directory for all stored files.
	basePath string
	// prefix is an optional subdirectory prefix for namespacing.
	prefix string
}

// NewFilesystem creates a new Filesystem storage backend using the provided parameters.
// It ensures the base path exists and registers a start hook for lifecycle management.
func NewFilesystem(p Params) (IStorage, error) {
	f := &Filesystem{
		basePath: p.Cfg.Storage.Filesystem.Path,
		prefix:   p.Cfg.Storage.Filesystem.Prefix,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := os.MkdirAll(f.basePath, 0o770); err != nil {
			return err
		}

		return nil
	}))

	return f, nil
}

// WithPrefix returns a new Filesystem instance with the given prefix, ensuring the directory exists.
func (s *Filesystem) WithPrefix(prefix string) (IStorage, error) {
	if err := os.MkdirAll(filepath.Join(s.basePath, prefix), 0o770); err != nil {
		return nil, err
	}

	return &Filesystem{
		basePath: s.basePath,
		prefix:   prefix,
	}, nil
}

// Get retrieves a file and its metadata from the filesystem.
// Returns an open file and ObjectInfo, or an error if not found or invalid.
func (s *Filesystem) Get(ctx context.Context, keyIn string) (IObject, IObjectInfo, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return nil, nil, ErrInvalidPath
	}
	key = filepath.Join(s.basePath, s.prefix, key)

	stat, err := os.Stat(key)
	if os.IsNotExist(err) {
		return nil, nil, ErrNotFound
	}

	f, err := os.OpenFile(key, os.O_RDONLY, 0o600)
	if err != nil {
		return nil, nil, err
	}

	name := stat.Name()

	return f, &ObjectInfo{
		name:         strings.TrimPrefix(name, s.prefix),
		extension:    strings.TrimPrefix(filepath.Ext(name), "."),
		size:         stat.Size(),
		lastModified: stat.ModTime(),
	}, nil
}

// Stat returns metadata for a file in the filesystem, or an error if not found or invalid.
func (s *Filesystem) Stat(ctx context.Context, keyIn string) (IObjectInfo, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return nil, ErrInvalidPath
	}
	key = filepath.Join(s.basePath, s.prefix, key)

	stat, err := os.Stat(key)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}

	return &ObjectInfo{
		name: strings.TrimPrefix(key, s.prefix),
		size: stat.Size(),
	}, nil
}

// Put writes a file to the filesystem, creating directories as needed.
// Returns the relative path or an error.
func (s *Filesystem) Put(ctx context.Context, keyIn string, reader io.Reader, size int64, contentType string) (string, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return "", ErrInvalidPath
	}
	key = filepath.Join(s.basePath, s.prefix, key)

	dir := filepath.Dir(key)
	if err := os.MkdirAll(dir, 0o770); err != nil {
		return "", err
	}

	f, err := os.OpenFile(key, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return filepath.Join(s.prefix, keyIn), nil
}

// Delete removes a file from the filesystem. Returns nil if the file does not exist.
func (s *Filesystem) Delete(ctx context.Context, keyIn string) error {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return ErrInvalidPath
	}
	keyIn = filepath.Join(s.basePath, s.prefix, key)

	if err := os.Remove(keyIn); err != nil {
		e, ok := err.(*os.PathError)
		if ok && e.Err != syscall.ENOENT {
			return err
		}
	}

	return nil
}

// List returns a list of files and their metadata in the given directory.
// Returns an error if the directory is invalid or cannot be read.
func (s *Filesystem) List(ctx context.Context, keyIn string, offset int, pageSize int) ([]*FileInfo, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return nil, ErrInvalidPath
	}
	key = filepath.Join(s.basePath, s.prefix, key)

	entries, err := os.ReadDir(key)
	if err != nil {
		return nil, err
	}

	files := []*FileInfo{}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			return nil, err
		}

		name := e.Name()
		contentType := filetype.GetType(strings.TrimPrefix(filepath.Ext(name), "."))

		files = append(files, &FileInfo{
			Name:         strings.TrimPrefix(name, s.prefix),
			LastModified: info.ModTime(),
			Size:         info.Size(),
			ContentType:  contentType.MIME.Value,
		})
	}

	return files, nil
}

// GetSpaceUsage calculates the total space used by files in the filesystem.
// It walks the directory tree and sums the sizes of all files.
// Returns the total size in bytes or an error if the directory cannot be read.
func (s *Filesystem) GetSpaceUsage(ctx context.Context) (int64, error) {
	var totalSize int64

	err := filepath.Walk(filepath.Join(s.basePath, s.prefix), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return totalSize, nil
}
