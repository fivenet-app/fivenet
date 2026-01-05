package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
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

	// root is the "root" directory of the storage.
	root *os.Root
	// prefix is an optional subdirectory prefix for namespacing.
	prefix string
}

// NewFilesystem creates a new Filesystem storage backend using the provided parameters.
// It ensures the base path exists and registers a start hook for lifecycle management.
func NewFilesystem(p Params) (IStorage, error) {
	f := &Filesystem{
		prefix: filepath.Clean(p.Cfg.Storage.Filesystem.Prefix),
	}

	basePath := p.Cfg.Storage.Filesystem.Path
	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := os.MkdirAll(basePath, 0o770); err != nil {
			return err
		}

		root, err := os.OpenRoot(basePath)
		if err != nil {
			return fmt.Errorf("failed to open filesystem storage path. %w", err)
		}
		f.root = root

		return nil
	}))

	return f, nil
}

// Get retrieves a file and its metadata from the filesystem.
// Returns an open file and ObjectInfo, or an error if not found or invalid.
func (s *Filesystem) Get(ctx context.Context, key string) (IObject, IObjectInfo, error) {
	key, err := utils.FSRootPath(filepath.Join(s.root.Name(), s.prefix), key)
	if err != nil {
		return nil, nil, ErrInvalidPath
	}

	stat, err := s.root.Stat(key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, ErrNotFound
		}

		return nil, nil, ErrNotFound
	}
	if stat.IsDir() {
		return nil, nil, ErrNotFound
	}

	f, err := s.root.OpenFile(key, os.O_RDONLY, 0o600)
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
func (s *Filesystem) Stat(ctx context.Context, key string) (IObjectInfo, error) {
	key, err := utils.FSRootPath(filepath.Join(s.root.Name(), s.prefix), key)
	if err != nil {
		return nil, ErrInvalidPath
	}

	stat, err := s.root.Stat(key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}

		return nil, ErrNotFound
	}
	if stat.IsDir() {
		return nil, ErrNotFound
	}

	return &ObjectInfo{
		name: strings.TrimPrefix(key, s.prefix),
		size: stat.Size(),
	}, nil
}

// Put writes a file to the filesystem, creating directories as needed.
// Returns the relative path or an error.
func (s *Filesystem) Put(
	ctx context.Context,
	key string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {
	key, err := utils.FSRootPath(filepath.Join(s.root.Name(), s.prefix), key)
	if err != nil {
		return "", ErrInvalidPath
	}

	dir := filepath.Dir(key)
	if err := s.root.MkdirAll(dir, 0o770); err != nil {
		return "", err
	}

	f, err := s.root.OpenFile(key, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return filepath.Join(s.prefix, key), nil
}

// Delete removes a file from the filesystem. Returns nil if the file does not exist.
func (s *Filesystem) Delete(ctx context.Context, key string) error {
	key, err := utils.FSRootPath(filepath.Join(s.root.Name(), s.prefix), key)
	if err != nil {
		return ErrInvalidPath
	}

	if err := s.root.Remove(key); err != nil {
		if !errors.Is(err, syscall.ENOENT) {
			return err
		}
	}

	return nil
}

// List returns a list of files and their metadata in the given directory.
// Returns an error if the directory is invalid or cannot be read.
func (s *Filesystem) List(
	ctx context.Context,
	key string,
	offset int,
	pageSize int,
) ([]*FileInfo, error) {
	key, err := utils.FSRootPath(filepath.Join(s.root.Name(), s.prefix), key)
	if err != nil {
		return nil, ErrInvalidPath
	}

	// os.Root doesn't have a `ReadDir` method, so open the directory and read from it
	f, err := s.root.Open(key)
	if err != nil {
		return nil, err
	}
	// combine offset + pageSize to improve performance
	entries, err := f.ReadDir(offset + pageSize)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	files := []*FileInfo{}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			return nil, err
		}

		name := filepath.Join(key, e.Name())
		contentType := filetype.GetType(strings.TrimPrefix(filepath.Ext(name), "."))

		files = append(files, &FileInfo{
			Name:         strings.TrimPrefix(name, s.prefix),
			LastModified: info.ModTime(),
			Size:         info.Size(),
			ContentType:  contentType.MIME.Value,
			IsDir:        e.IsDir(),
		})
	}

	return files, nil
}

// GetSpaceUsage calculates the total space used by files in the filesystem.
// It walks the directory tree and sums the sizes of all files.
// Returns the total size in bytes or an error if the directory cannot be read.
func (s *Filesystem) GetSpaceUsage(ctx context.Context) (int64, error) {
	var totalSize int64

	err := filepath.WalkDir(
		filepath.Join(s.root.Name(), s.prefix),
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if d.IsDir() {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return err
			}
			totalSize += info.Size()

			return nil
		},
	)
	if err != nil {
		return 0, err
	}

	return totalSize, nil
}
