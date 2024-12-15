// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/jobs/conduct.proto

package jobs

func (m *CreateConductEntryRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Entry
	if m.Entry != nil {
		if v, ok := interface{}(m.GetEntry()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *CreateConductEntryResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Entry
	if m.Entry != nil {
		if v, ok := interface{}(m.GetEntry()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DeleteConductEntryRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *DeleteConductEntryResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *ListConductEntriesRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Pagination
	if m.Pagination != nil {
		if v, ok := interface{}(m.GetPagination()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Sort
	if m.Sort != nil {
		if v, ok := interface{}(m.GetSort()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Types
	for idx, item := range m.Types {
		_, _ = idx, item

	}

	return nil
}

func (m *ListConductEntriesResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Entries
	for idx, item := range m.Entries {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Pagination
	if m.Pagination != nil {
		if v, ok := interface{}(m.GetPagination()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *UpdateConductEntryRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Entry
	if m.Entry != nil {
		if v, ok := interface{}(m.GetEntry()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *UpdateConductEntryResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Entry
	if m.Entry != nil {
		if v, ok := interface{}(m.GetEntry()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}