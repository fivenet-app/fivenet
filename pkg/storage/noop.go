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

// Noop implements IStorage as a no-operation (noop) storage backend for testing or disabled storage scenarios.
type Noop struct{}

// NewNoop returns a new instance of Noop storage backend.
func NewNoop(p Params) (IStorage, error) {
	return &Noop{}, nil
}

// WithPrefix returns the same Noop instance, ignoring the prefix.
func (s *Noop) WithPrefix(prefix string) (IStorage, error) {
	return s, nil
}

// Get always returns nils, simulating a missing object.
func (s *Noop) Get(ctx context.Context, filePath string) (IObject, IObjectInfo, error) {
	return nil, nil, nil
}

// GetURL always returns nil, simulating no URL support.
func (s *Noop) GetURL(ctx context.Context, filePath string, expires time.Duration, reqParams url.Values) (*string, error) {
	return nil, nil
}

// Stat always returns nil, simulating a missing object.
func (s *Noop) Stat(ctx context.Context, filePath string) (IObjectInfo, error) {
	return nil, nil
}

// Put always returns an empty string, simulating a no-op upload.
func (s *Noop) Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error) {
	return "", nil
}

// Delete always returns nil, simulating a successful delete.
func (s *Noop) Delete(ctx context.Context, filePath string) error {
	return nil
}

// List always returns nil, simulating no files in storage.
func (s *Noop) List(ctx context.Context, filePath string, offset int, pageSize int) ([]*FileInfo, error) {
	return nil, nil
}

func (s *Noop) GetSpaceUsage(ctx context.Context) (int64, error) {
	return 0, nil
}
