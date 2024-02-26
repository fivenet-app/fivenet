package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/galexrt/fivenet/pkg/config"
)

func init() {
	storageFactories["filesystem"] = NewFilesystem
}

type Filesystem struct {
	IStorage

	basePath string
	prefix   string
}

func NewFilesystem(cfg *config.Config) (IStorage, error) {
	f := &Filesystem{
		basePath: cfg.Storage.Filesystem.Path,
	}

	if err := os.MkdirAll(f.basePath, 0770); err != nil {
		return nil, err
	}

	return f, nil
}

func (s *Filesystem) WithPrefix(prefix string) (IStorage, error) {
	if err := os.MkdirAll(path.Join(s.basePath, prefix), 0770); err != nil {
		return nil, err
	}

	return &Filesystem{
		basePath: s.basePath,
		prefix:   prefix,
	}, nil
}

func (s *Filesystem) Get(ctx context.Context, filePath string) (IObject, IObjectInfo, error) {
	filePath = path.Join(s.basePath, s.prefix, filePath)

	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, nil, err
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, nil, err
	}

	return f, &ObjectInfo{
		size: stat.Size(),
	}, nil
}

func (s *Filesystem) Stat(ctx context.Context, filePath string) (IObjectInfo, error) {
	filePath = path.Join(s.basePath, s.prefix, filePath)

	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, err
	}

	return &ObjectInfo{
		size: stat.Size(),
	}, nil
}

func (s *Filesystem) Put(ctx context.Context, filePathIn string, reader io.Reader, size int64, contentType string) (string, error) {
	filePath := path.Join(s.basePath, s.prefix, filePathIn)

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		if err := os.MkdirAll(dir, 0770); err != nil {
			return "", err
		}
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return path.Join(s.prefix, filePathIn), nil
}

func (s *Filesystem) Delete(ctx context.Context, filePath string) error {
	filePath = path.Join(s.basePath, s.prefix, filePath)

	if err := os.Remove(filePath); err != nil {
		e, ok := err.(*os.PathError)
		if ok && e.Err != syscall.ENOENT {
			return err
		}
	}

	return nil
}