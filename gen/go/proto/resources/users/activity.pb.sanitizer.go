// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/users/activity.proto

package users

import (
	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
)

func (m *UserActivity) Sanitize() error {
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

	// Field: Key
	m.Key = htmlsanitizer.Sanitize(m.Key)

	// Field: Reason
	m.Reason = htmlsanitizer.Sanitize(m.Reason)

	// Field: SourceUser
	if m.SourceUser != nil {
		if v, ok := interface{}(m.GetSourceUser()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: TargetUser
	if m.TargetUser != nil {
		if v, ok := interface{}(m.GetTargetUser()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
