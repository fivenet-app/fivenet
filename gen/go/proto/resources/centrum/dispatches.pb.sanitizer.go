// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/centrum/dispatches.proto

package centrum

import (
	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
)

func (m *Dispatch) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Attributes
	if m.Attributes != nil {
		if v, ok := interface{}(m.GetAttributes()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := interface{}(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Creator
	if m.Creator != nil {
		if v, ok := interface{}(m.GetCreator()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Description

	if m.Description != nil {
		*m.Description = htmlsanitizer.Sanitize(*m.Description)
	}

	// Field: Message
	m.Message = htmlsanitizer.Sanitize(m.Message)

	// Field: Postal

	if m.Postal != nil {
		*m.Postal = htmlsanitizer.Sanitize(*m.Postal)
	}

	// Field: References
	if m.References != nil {
		if v, ok := interface{}(m.GetReferences()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Status
	if m.Status != nil {
		if v, ok := interface{}(m.GetStatus()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Units
	for idx, item := range m.Units {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: UpdatedAt
	if m.UpdatedAt != nil {
		if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DispatchAssignment) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := interface{}(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: ExpiresAt
	if m.ExpiresAt != nil {
		if v, ok := interface{}(m.GetExpiresAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Unit
	if m.Unit != nil {
		if v, ok := interface{}(m.GetUnit()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DispatchAssignments) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Units
	for idx, item := range m.Units {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *DispatchReference) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *DispatchReferences) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: References
	for idx, item := range m.References {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *DispatchStatus) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Code

	if m.Code != nil {
		*m.Code = htmlsanitizer.Sanitize(*m.Code)
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := interface{}(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Postal

	if m.Postal != nil {
		*m.Postal = htmlsanitizer.Sanitize(*m.Postal)
	}

	// Field: Reason

	if m.Reason != nil {
		*m.Reason = htmlsanitizer.Sanitize(*m.Reason)
	}

	// Field: Unit
	if m.Unit != nil {
		if v, ok := interface{}(m.GetUnit()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: User
	if m.User != nil {
		if v, ok := interface{}(m.GetUser()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
