// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/laws/laws.proto

package laws

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/html/htmlsanitizer"
)

func (m *Law) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := any(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Description

	if m.Description != nil {
		*m.Description = htmlsanitizer.Sanitize(*m.Description)
	}

	// Field: Hint

	if m.Hint != nil {
		*m.Hint = htmlsanitizer.Sanitize(*m.Hint)
	}

	// Field: Name
	m.Name = htmlsanitizer.Sanitize(m.Name)

	// Field: UpdatedAt
	if m.UpdatedAt != nil {
		if v, ok := any(m.GetUpdatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *LawBook) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := any(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Description

	if m.Description != nil {
		*m.Description = htmlsanitizer.Sanitize(*m.Description)
	}

	// Field: Laws
	for idx, item := range m.Laws {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Name
	m.Name = htmlsanitizer.Sanitize(m.Name)

	// Field: UpdatedAt
	if m.UpdatedAt != nil {
		if v, ok := any(m.GetUpdatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
