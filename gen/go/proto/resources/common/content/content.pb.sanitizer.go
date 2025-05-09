// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/common/content/content.proto

package content

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/html/htmlsanitizer"
)

func (m *Content) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Content
	if m.Content != nil {
		if v, ok := any(m.GetContent()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: RawContent

	if m.RawContent != nil {
		*m.RawContent = htmlsanitizer.Sanitize(*m.RawContent)
	}

	return nil
}

func (m *JSONNode) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Attrs
	for idx, item := range m.Attrs {
		_, _ = idx, item

		m.Attrs[idx] = htmlsanitizer.StripTags(m.Attrs[idx])

	}

	// Field: Content
	for idx, item := range m.Content {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Id

	if m.Id != nil {
		*m.Id = htmlsanitizer.StripTags(*m.Id)
	}

	// Field: Tag
	m.Tag = htmlsanitizer.StripTags(m.Tag)

	// Field: Text

	if m.Text != nil {
		*m.Text = htmlsanitizer.StripTags(*m.Text)
	}

	return nil
}
