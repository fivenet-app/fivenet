package storage

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestFilesystem_Get(t *testing.T) {
	lc := fxtest.NewLifecycle(t)

	basePath := t.TempDir()
	dirPrefix := "test-prefix"

	fs, err := NewFilesystem(Params{
		LC: lc,
		Cfg: &config.Config{
			Storage: config.Storage{
				Type: config.StorageTypeFilesystem,
				Filesystem: config.FilesystemStorage{
					Path:   basePath,
					Prefix: dirPrefix,
				},
			},
		},
	})
	require.NoError(t, err)

	err = lc.Start(t.Context())
	require.NoError(t, err)

	// Create a test file
	testFilePath := filepath.Join(basePath, dirPrefix, "testfile.txt")
	err = os.MkdirAll(filepath.Dir(testFilePath), 0o770)
	require.NoError(t, err)

	content := []byte("test content")
	err = os.WriteFile(testFilePath, content, 0o600)
	require.NoError(t, err)

	t.Run("valid file", func(t *testing.T) {
		ctx := t.Context()
		obj, info, err := fs.Get(ctx, "testfile.txt")
		require.NoError(t, err)
		require.NotNil(t, obj)
		require.NotNil(t, info)

		defer obj.Close()

		// Verify file content
		data, err := io.ReadAll(obj)
		require.NoError(t, err)
		require.Equal(t, content, data)

		// Verify file info
		require.Equal(t, "testfile.txt", info.GetName())
		require.Equal(t, int64(len(content)), info.GetSize())
	})

	t.Run("file not found", func(t *testing.T) {
		ctx := t.Context()
		_, _, err := fs.Get(ctx, "nonexistent.txt")
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("invalid path", func(t *testing.T) {
		ctx := t.Context()
		_, _, err := fs.Get(ctx, "../invalidpath.txt")
		require.ErrorIs(t, err, ErrInvalidPath)
	})

	t.Run("put and get file", func(t *testing.T) {
		ctx := t.Context()
		contentToPut := []byte("content to put")
		key := "putfile.txt"

		// Put the file
		_, err := fs.Put(
			ctx,
			key,
			io.NopCloser(
				io.LimitReader(
					io.NewSectionReader(bytes.NewReader(contentToPut), 0, int64(len(contentToPut))),
					int64(len(contentToPut)),
				),
			),
			int64(len(contentToPut)),
			"text/plain",
		)
		require.NoError(t, err)

		// Get the file
		obj, info, err := fs.Get(ctx, key)
		require.NoError(t, err)
		require.NotNil(t, obj)
		require.NotNil(t, info)

		defer obj.Close()

		// Verify file content
		data, err := io.ReadAll(obj)
		require.NoError(t, err)
		require.Equal(t, contentToPut, data)

		// Verify file info
		require.Equal(t, key, info.GetName())
		require.Equal(t, int64(len(contentToPut)), info.GetSize())
	})
}
