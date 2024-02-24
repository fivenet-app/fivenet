package storage

import (
	"context"
	"io"
	"os"
	"path"

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

	return f, nil
}

func (s *Filesystem) WithPrefix(prefix string) IStorage {
	return &Filesystem{
		basePath: s.basePath,
		prefix:   prefix,
	}
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

func (s *Filesystem) Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error) {
	filePath = path.Join(s.basePath, s.prefix, filePath)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return path.Join(s.prefix, filePath), nil
}

func (s *Filesystem) Delete(ctx context.Context, filePath string) error {
	filePath = path.Join(s.basePath, s.prefix, filePath)

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
