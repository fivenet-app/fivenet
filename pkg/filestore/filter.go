package filestore

import (
	"maps"
	"mime"
	"path/filepath"
	"slices"
	"strings"
)

// UploadFilter validates file type information before storing uploads.
type UploadFilter struct {
	allowedContentTypes map[string]struct{}
	allowedExtensions   map[string]struct{}
}

// NewUploadFilter creates a filter that allow-lists content types and/or file extensions.
func NewUploadFilter(allowedContentTypes []string, allowedExtensions []string) *UploadFilter {
	f := &UploadFilter{
		allowedContentTypes: map[string]struct{}{},
		allowedExtensions:   map[string]struct{}{},
	}

	for _, ctype := range allowedContentTypes {
		ctype = normalizeContentType(ctype)
		if ctype == "" {
			continue
		}
		f.allowedContentTypes[ctype] = struct{}{}
	}

	for _, ext := range allowedExtensions {
		ext = normalizeExtension(ext)
		if ext == "" {
			continue
		}
		f.allowedExtensions[ext] = struct{}{}
	}

	return f
}

// NewImageUploadFilter creates an allow-list for common web image formats.
func NewImageUploadFilter() *UploadFilter {
	return NewUploadFilter(
		[]string{
			"image/gif",
			"image/jpeg",
			"image/png",
			"image/webp",
		},
		[]string{
			"gif",
			"jpg",
			"jpeg",
			"png",
			"webp",
		},
	)
}

// Validate checks if a file key and MIME type match the filter.
func (f *UploadFilter) Validate(fileName, contentType string) error {
	if f == nil {
		return nil
	}

	ext := normalizeExtension(fileName)
	ctype := normalizeContentType(contentType)

	if len(f.allowedExtensions) > 0 {
		if _, ok := f.allowedExtensions[ext]; !ok {
			return ErrUploadFileTypeNotAllowed(map[string]any{
				"extension": ext,
			})
		}
	}

	if len(f.allowedContentTypes) > 0 {
		if _, ok := f.allowedContentTypes[ctype]; !ok {
			return ErrUploadFileTypeNotAllowed(map[string]any{
				"contentType": ctype,
			})
		}
	}

	return nil
}

// AllowedContentTypes returns a copy of the allow-listed content types.
func (f *UploadFilter) AllowedContentTypes() []string {
	if f == nil {
		return nil
	}

	values := make([]string, 0, len(f.allowedContentTypes))
	for contentType := range f.allowedContentTypes {
		values = append(values, contentType)
	}
	slices.Sort(values)
	return values
}

// AllowedExtensions returns a copy of the allow-listed extensions.
func (f *UploadFilter) AllowedExtensions() []string {
	if f == nil {
		return nil
	}

	values := make([]string, 0, len(f.allowedExtensions))
	for ext := range f.allowedExtensions {
		values = append(values, ext)
	}
	slices.Sort(values)
	return values
}

func normalizeContentType(contentType string) string {
	contentType = strings.TrimSpace(strings.ToLower(contentType))
	if contentType == "" {
		return ""
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}

	return strings.ToLower(mediaType)
}

func normalizeExtension(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return ""
	}

	ext := strings.TrimPrefix(filepath.Ext(name), ".")
	if ext != "" {
		return ext
	}

	return strings.TrimPrefix(name, ".")
}

func (f *UploadFilter) clone() *UploadFilter {
	if f == nil {
		return nil
	}

	return &UploadFilter{
		allowedContentTypes: maps.Clone(f.allowedContentTypes),
		allowedExtensions:   maps.Clone(f.allowedExtensions),
	}
}
