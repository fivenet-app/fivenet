// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/qualifications/access.proto

package qualifications

func (m *QualificationAccess) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Jobs
	for idx, item := range m.Jobs {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *QualificationJobAccess) Sanitize() error {
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

	return nil
}

func (m *QualificationUserAccess) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}
