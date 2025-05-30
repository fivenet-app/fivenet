// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/internet/domain.proto

package internet

func (m *Domain) Sanitize() error {
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

	// Field: DeletedAt
	if m.DeletedAt != nil {
		if v, ok := any(m.GetDeletedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Tld
	if m.Tld != nil {
		if v, ok := any(m.GetTld()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

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

func (m *TLD) Sanitize() error {
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

	// Field: DeletedAt
	if m.DeletedAt != nil {
		if v, ok := any(m.GetDeletedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

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
