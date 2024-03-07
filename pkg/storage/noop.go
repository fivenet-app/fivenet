package storage

import (
	"context"
	"io"

	"github.com/galexrt/fivenet/pkg/config"
	"go.uber.org/fx"
)

func init() {
	storageFactories["noop"] = NewNoop
}

type Noop struct{}

func NewNoop(lc fx.Lifecycle, cfg *config.Config) (IStorage, error) {
	return &Noop{}, nil
}

func (s *Noop) WithPrefix(prefix string) (IStorage, error) {
	return s, nil
}

func (s *Noop) Get(ctx context.Context, filePath string) (IObject, IObjectInfo, error) {
	return nil, nil, nil
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
