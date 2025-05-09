package storage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
)

func init() {
	storageFactories[config.StorageTypeNoop] = NewNoop
}

type Noop struct{}

func NewNoop(p Params) (IStorage, error) {
	return &Noop{}, nil
}

func (s *Noop) WithPrefix(prefix string) (IStorage, error) {
	return s, nil
}

func (s *Noop) Get(ctx context.Context, filePath string) (IObject, IObjectInfo, error) {
	return nil, nil, nil
}

func (s *Noop) GetURL(ctx context.Context, filePath string, expires time.Duration, reqParams url.Values) (*string, error) {
	return nil, nil
}

func (s *Noop) Stat(ctx context.Context, filePath string) (IObjectInfo, error) {
	return nil, nil
}

func (s *Noop) Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error) {
	return "", nil
}

func (s *Noop) Delete(ctx context.Context, filePath string) error {
	return nil
}

func (s *Noop) List(ctx context.Context, filePath string, offset int, pageSize int) ([]*FileInfo, error) {
	return nil, nil
}
