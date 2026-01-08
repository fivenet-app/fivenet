package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanStorageKey(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "Valid path",
			input:     "valid/path",
			expected:  "valid/path",
			expectErr: false,
		},
		{
			name:      "Path with leading slash",
			input:     "/leading/slash",
			expected:  "leading/slash",
			expectErr: false,
		},
		{
			name:      "Path with trailing slash",
			input:     "trailing/slash/",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Path with dot segments",
			input:     "path/./to/./file",
			expected:  "path/to/file",
			expectErr: false,
		},
		{
			name:      "Path traversal attempt",
			input:     "../traversal",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Path traversal within path",
			input:     "path/../traversal",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Empty path",
			input:     "",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid UTF-8 string",
			input:     string([]byte{0xff, 0xfe, 0xfd}),
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Path with multiple slashes",
			input:     "path//to///file",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Path with only dots",
			input:     ".",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Path with dotdot segment",
			input:     "path/..",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Path with empty segment",
			input:     "path//segment",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Windows path with mixed slashes",
			input:     "path\\to/../file",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Windows drive path rejected",
			input:     "c:/path/to/file",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Windows drive path with backslashes rejected",
			input:     "C:\\path\\to\\file",
			expectErr: true,
		},
		{
			name:      "UNC path rejected forward slashes",
			input:     "//server/share/file.txt",
			expectErr: true,
		},
		{
			name:      "UNC path rejected backslashes",
			input:     `\\server\share\file.txt`,
			expectErr: true,
		},
		{
			name:      "Real World #1",
			input:     "fivenet/fivenet/documents/20260106/3b56a375-fce0-40c9-b481-d60455233282-you-guys-are-getting-blueprints-v0-9ijtqacktz9g1.webp",
			expected:  "fivenet/fivenet/documents/20260106/3b56a375-fce0-40c9-b481-d60455233282-you-guys-are-getting-blueprints-v0-9ijtqacktz9g1.webp",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CleanStorageKey(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFSRootPath(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		key       string
		expected  string
		expectErr bool
	}{
		{
			name:      "Join prefix and key",
			prefix:    "images",
			key:       "subdir/file.txt",
			expected:  filepath.Join("images", "subdir", "file.txt"),
			expectErr: false,
		},
		{
			name:      "Empty prefix returns key",
			prefix:    "",
			key:       "subdir/file.txt",
			expected:  filepath.Join("subdir", "file.txt"),
			expectErr: false,
		},
		{
			name:      "Prefix with trailing slash",
			prefix:    "images/",
			key:       "subdir/file.txt",
			expected:  filepath.Join("images", "subdir", "file.txt"),
			expectErr: false,
		},
		{
			name:      "Nested prefix",
			prefix:    "images/nested",
			key:       "subdir/file.txt",
			expected:  filepath.Join("images", "nested", "subdir", "file.txt"),
			expectErr: false,
		},
		{
			name:      "Unicode key",
			prefix:    "images",
			key:       "ümlaut/猫.png",
			expected:  filepath.Join("images", "ümlaut", "猫.png"),
			expectErr: false,
		},

		// Minimal invariants
		{name: "Empty key rejected", prefix: "images", key: "", expectErr: true},
		{name: "Dot key rejected", prefix: "images", key: ".", expectErr: true},
		{name: "Dotdot key rejected", prefix: "images", key: "..", expectErr: true},
		{name: "NUL rejected", prefix: "images", key: "a\x00b", expectErr: true},
		{name: "Backslash rejected", prefix: "images", key: `subdir\file.txt`, expectErr: true},

		// Prefix misconfig defense-in-depth
		{name: "Absolute prefix rejected", prefix: "/images", key: "file.txt", expectErr: true},
		{name: "Traversal prefix rejected", prefix: "../images", key: "file.txt", expectErr: true},
		{name: "Traversal key rejected", prefix: "../images", key: "../file.txt", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FSRootPath(tt.prefix, tt.key)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
			assert.False(t, filepath.IsAbs(got))
			assert.Empty(t, filepath.VolumeName(got))
		})
	}
}
