package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanFilePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		valid    bool
	}{
		{
			name:     "Absolute path",
			input:    "/absolute/path/file.txt",
			expected: "/absolute/path/file.txt",
			valid:    true,
		},
		{
			name:     "Absolute path to a directory",
			input:    "/absolute/path/to/directory/",
			expected: "",
			valid:    false,
		},
		{
			name:     "Absolute path with trailing slash",
			input:    "/absolute/path/file.txt/",
			expected: "",
			valid:    false,
		},
		{
			name:     "Invalid local path",
			input:    "./test/file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Empty path",
			input:    "",
			expected: "",
			valid:    false,
		},
		{
			name:     "Current directory",
			input:    ".",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path with redundant elements 1",
			input:    "./test/../file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path with redundant elements 2",
			input:    "./test/../../../../../file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path with redundant elements 3",
			input:    "./test/dir/..//../file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path beginning with traversal",
			input:    "../outside/file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path beginning with traversal (Windows style)",
			input:    "c:/../outside/file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path with redundant elements 2",
			input:    "../../../test/../../../../../file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path with multi slashes and traversal",
			input:    "//////../../../test//../..////../../../file.txt",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path with invalid utf-8 rune",
			input:    "/absolute/path/with/invalid/\xff/file.txt",
			expected: "",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, valid := CleanFilePath(tt.input)
			assert.Equal(t, tt.expected, result, tt.name)
			assert.Equal(t, tt.valid, valid, tt.name)
		})
	}
}
