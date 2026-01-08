package storage

import (
	"context"
	"io"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
)

type Wrapper struct {
	IStorage

	backend IStorage
}

// Get returns the object contents and info for the given key.
func (s *Wrapper) Get(ctx context.Context, key string) (IObject, IObjectInfo, error) {
	key, err := utils.CleanStorageKey(key)
	if err != nil {
		return nil, nil, ErrInvalidPath
	}

	return s.backend.Get(ctx, key)
}

// Stat returns object info for the given key.
func (s *Wrapper) Stat(ctx context.Context, key string) (IObjectInfo, error) {
	key, err := utils.CleanStorageKey(key)
	if err != nil {
		return nil, ErrInvalidPath
	}

	return s.backend.Stat(ctx, key)
}

// Put uploads a file; size and content type must be accurate.
func (s *Wrapper) Put(
	ctx context.Context,
	key string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {
	key, err := utils.CleanStorageKey(key)
	if err != nil {
		return "", ErrInvalidPath
	}

	return s.backend.Put(ctx, key, reader, size, contentType)
}

// Delete removes a file from storage.
func (s *Wrapper) Delete(ctx context.Context, key string) error {
	key, err := utils.CleanStorageKey(key)
	if err != nil {
		return ErrInvalidPath
	}

	return s.backend.Delete(ctx, key)
}

// List returns a list of files by offset and page size.
func (s *Wrapper) List(
	ctx context.Context,
	key string,
	offset int,
	pageSize int,
) ([]*FileInfo, error) {
	key, err := utils.CleanStorageKey(key)
	if err != nil {
		return nil, ErrInvalidPath
	}
	return s.backend.List(ctx, key, offset, pageSize)
}

// GetSpaceUsage returns the total space used in storage.
func (s *Wrapper) GetSpaceUsage(ctx context.Context) (int64, error) {
	return s.backend.GetSpaceUsage(ctx)
}
