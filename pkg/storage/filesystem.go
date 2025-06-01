package storage

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/h2non/filetype"
	"go.uber.org/fx"
)

func init() {
	storageFactories[config.StorageTypeFilesystem] = NewFilesystem
}

type Filesystem struct {
	IStorage

	basePath string
	prefix   string
}

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

func (s *Filesystem) WithPrefix(prefix string) (IStorage, error) {
	if err := os.MkdirAll(filepath.Join(s.basePath, prefix), 0o770); err != nil {
		return nil, err
	}

	return &Filesystem{
		basePath: s.basePath,
		prefix:   prefix,
	}, nil
}

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

func (s *Filesystem) GetURL(ctx context.Context, key string, expires time.Duration, reqParams url.Values) (*string, error) {
	return nil, nil
}

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
