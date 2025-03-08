package storage

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/utils"
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

func (s *Filesystem) Get(ctx context.Context, filePathIn string) (IObject, IObjectInfo, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return nil, nil, ErrInvalidPath
	}
	filePath = filepath.Join(s.basePath, s.prefix, filePath)

	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, nil, ErrNotFound
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0o600)
	if err != nil {
		return nil, nil, err
	}

	name := stat.Name()

	return f, &ObjectInfo{
		name:         name,
		extension:    strings.TrimPrefix(filepath.Ext(name), "."),
		size:         stat.Size(),
		lastModified: stat.ModTime(),
	}, nil
}

func (s *Filesystem) GetURL(ctx context.Context, filePath string, expires time.Duration, reqParams url.Values) (*string, error) {
	return nil, nil
}

func (s *Filesystem) Stat(ctx context.Context, filePathIn string) (IObjectInfo, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return nil, ErrInvalidPath
	}
	filePath = filepath.Join(s.basePath, s.prefix, filePath)

	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, err
	}

	return &ObjectInfo{
		size: stat.Size(),
	}, nil
}

func (s *Filesystem) Put(ctx context.Context, filePathIn string, reader io.Reader, size int64, contentType string) (string, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return "", ErrInvalidPath
	}
	filePath = filepath.Join(s.basePath, s.prefix, filePath)

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		if err := os.MkdirAll(dir, 0o770); err != nil {
			return "", err
		}
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return filepath.Join(s.prefix, filePathIn), nil
}

func (s *Filesystem) Delete(ctx context.Context, filePathIn string) error {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return ErrInvalidPath
	}
	filePathIn = filepath.Join(s.basePath, s.prefix, filePath)

	if err := os.Remove(filePathIn); err != nil {
		e, ok := err.(*os.PathError)
		if ok && e.Err != syscall.ENOENT {
			return err
		}
	}

	return nil
}

func (s *Filesystem) List(ctx context.Context, filePathIn string, offset int, pageSize int) ([]*FileInfo, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return nil, ErrInvalidPath
	}
	filePath = filepath.Join(s.basePath, s.prefix, filePath)

	entries, err := os.ReadDir(filePath)
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
			Name:         name,
			LastModified: info.ModTime(),
			Size:         info.Size(),
			ContentType:  contentType.MIME.Value,
		})
	}

	return files, nil
}
