// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/centrum/access.proto

package centrum

func (m *UnitAccess) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Jobs
	for idx, item := range m.Jobs {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Qualifications
	for idx, item := range m.Qualifications {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *UnitJobAccess) Sanitize() error {
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

	return nil
}

func (m *UnitQualificationAccess) Sanitize() error {
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

	// Field: Qualification
	if m.Qualification != nil {
		if v, ok := interface{}(m.GetQualification()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *UnitUserAccess) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}