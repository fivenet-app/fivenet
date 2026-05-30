package filestore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadFilterValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		filter      *UploadFilter
		fileName    string
		contentType string
		expectError bool
	}{
		{
			name:        "Nil filter allows all",
			filter:      nil,
			fileName:    "documents/test.pdf",
			contentType: "application/pdf",
			expectError: false,
		},
		{
			name:        "Image filter allows png with mixed casing",
			filter:      NewImageUploadFilter(),
			fileName:    "user_profile_pictures/photo.PNG",
			contentType: "image/png",
			expectError: false,
		},
		{
			name:        "Image filter allows jpeg with parameters",
			filter:      NewImageUploadFilter(),
			fileName:    "user_profile_pictures/photo.jpeg",
			contentType: "image/jpeg; charset=utf-8",
			expectError: false,
		},
		{
			name:        "Image filter rejects unsupported extension",
			filter:      NewImageUploadFilter(),
			fileName:    "user_profile_pictures/photo.svg",
			contentType: "image/svg+xml",
			expectError: true,
		},
		{
			name:        "Image filter rejects unsupported content type",
			filter:      NewImageUploadFilter(),
			fileName:    "user_profile_pictures/photo.png",
			contentType: "application/octet-stream",
			expectError: true,
		},
		{
			name:        "Extension-only filter allows without content type",
			filter:      NewUploadFilter(nil, []string{"pdf"}),
			fileName:    "documents/file.pdf",
			contentType: "",
			expectError: false,
		},
		{
			name:        "Content-type-only filter allows without extension",
			filter:      NewUploadFilter([]string{"application/pdf"}, nil),
			fileName:    "documents/file",
			contentType: "application/pdf",
			expectError: false,
		},
		{
			name:        "Content-type-only filter rejects missing content type",
			filter:      NewUploadFilter([]string{"application/pdf"}, nil),
			fileName:    "documents/file.pdf",
			contentType: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.filter.Validate(tt.fileName, tt.contentType)
			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUploadFilterAllowedListsAreNormalizedAndSorted(t *testing.T) {
	t.Parallel()

	filter := NewUploadFilter(
		[]string{" IMAGE/JPEG ; charset=UTF-8 ", "image/png"},
		[]string{".PNG", " jpg "},
	)

	assert.Equal(t, []string{"image/jpeg", "image/png"}, filter.AllowedContentTypes())
	assert.Equal(t, []string{"jpg", "png"}, filter.AllowedExtensions())
}
