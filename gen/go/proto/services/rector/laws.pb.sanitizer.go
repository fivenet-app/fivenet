// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/rector/laws.proto

package rector

func (m *CreateOrUpdateLawBookRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: LawBook
	if m.LawBook != nil {
		if v, ok := interface{}(m.GetLawBook()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *CreateOrUpdateLawBookResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: LawBook
	if m.LawBook != nil {
		if v, ok := interface{}(m.GetLawBook()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *CreateOrUpdateLawRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Law
	if m.Law != nil {
		if v, ok := interface{}(m.GetLaw()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *CreateOrUpdateLawResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Law
	if m.Law != nil {
		if v, ok := interface{}(m.GetLaw()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DeleteLawBookRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *DeleteLawBookResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *DeleteLawRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *DeleteLawResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}