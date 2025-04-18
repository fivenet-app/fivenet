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
			name:     "Valid local path",
			input:    "./test/file.txt",
			expected: "test/file.txt",
			valid:    true,
		},
		{
			name:     "Absolute path",
			input:    "/absolute/path/file.txt",
			expected: "/absolute/path/file.txt",
			valid:    true,
		},
		{
			name:     "Empty path",
			input:    "",
			expected: "",
			valid:    false,
		},
		{
			name:     "Path with redundant elements 1",
			input:    "./test/../file.txt",
			expected: "file.txt",
			valid:    true,
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
			expected: "file.txt",
			valid:    true,
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
