// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/sync/activity.proto

package sync

func (m *ColleagueProps) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Props
	if m.Props != nil {
		if v, ok := any(m.GetProps()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *TimeclockUpdate) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserOAuth2Conn) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserProps) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Props
	if m.Props != nil {
		if v, ok := any(m.GetProps()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *UserUpdate) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}
