package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilesystem_Get(t *testing.T) {
	basePath := t.TempDir()
	fs := &Filesystem{
		basePath: basePath,
		prefix:   "test-prefix",
	}

	// Create a test file
	testFilePath := filepath.Join(basePath, "test-prefix", "testfile.txt")
	err := os.MkdirAll(filepath.Dir(testFilePath), 0o770)
	require.NoError(t, err)

	content := []byte("test content")
	err = os.WriteFile(testFilePath, content, 0o600)
	require.NoError(t, err)

	t.Run("valid file", func(t *testing.T) {
		ctx := context.Background()
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
		ctx := context.Background()
		_, _, err := fs.Get(ctx, "nonexistent.txt")
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("invalid path", func(t *testing.T) {
		ctx := context.Background()
		_, _, err := fs.Get(ctx, "../invalidpath.txt")
		require.ErrorIs(t, err, ErrInvalidPath)
	})
}
