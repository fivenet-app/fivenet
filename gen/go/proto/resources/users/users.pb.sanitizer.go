// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/users/users.proto

package users

func (m *License) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *User) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Avatar
	if m.Avatar != nil {
		if v, ok := interface{}(m.GetAvatar()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Licenses
	for idx, item := range m.Licenses {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Props
	if m.Props != nil {
		if v, ok := interface{}(m.GetProps()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *UserLicenses) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Licenses
	for idx, item := range m.Licenses {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *UserShort) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Avatar
	if m.Avatar != nil {
		if v, ok := interface{}(m.GetAvatar()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
