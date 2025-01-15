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

	// Field: Data
	if m.Data != nil {
		if v, ok := interface{}(m.GetData()).(interface{ Sanitize() error }); ok {
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

func (m *UserActivityData) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: DocumentRelation
	switch v := m.Data.(type) {

	case *UserActivityData_DocumentRelation:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: FineChange
	case *UserActivityData_FineChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: JailChange
	case *UserActivityData_JailChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: JobChange
	case *UserActivityData_JobChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: LabelsChange
	case *UserActivityData_LabelsChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: LicensesChange
	case *UserActivityData_LicensesChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: MugshotChange
	case *UserActivityData_MugshotChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: NameChange
	case *UserActivityData_NameChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: TrafficInfractionPointsChange
	case *UserActivityData_TrafficInfractionPointsChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: WantedChange
	case *UserActivityData_WantedChange:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *UserDocumentRelation) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserFineChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserJailChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserJobChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserLabelsChange) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Added
	for idx, item := range m.Added {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Removed
	for idx, item := range m.Removed {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *UserLicenseChange) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Added
	for idx, item := range m.Added {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Removed
	for idx, item := range m.Removed {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *UserMugshotChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserNameChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserTrafficInfractionPointsChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserWantedChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}
