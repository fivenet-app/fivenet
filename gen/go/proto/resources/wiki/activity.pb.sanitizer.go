// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/wiki/activity.proto

package wiki

func (m *PageAccessJobsDiff) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: ToCreate
	for idx, item := range m.ToCreate {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: ToDelete
	for idx, item := range m.ToDelete {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: ToUpdate
	for idx, item := range m.ToUpdate {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *PageAccessUpdated) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Jobs
	if m.Jobs != nil {
		if v, ok := any(m.GetJobs()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Users
	if m.Users != nil {
		if v, ok := any(m.GetUsers()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *PageAccessUsersDiff) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: ToCreate
	for idx, item := range m.ToCreate {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: ToDelete
	for idx, item := range m.ToDelete {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: ToUpdate
	for idx, item := range m.ToUpdate {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *PageActivity) Sanitize() error {
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

	// Field: Creator
	if m.Creator != nil {
		if v, ok := any(m.GetCreator()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Data
	if m.Data != nil {
		if v, ok := any(m.GetData()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *PageActivityData) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: AccessUpdated
	switch v := m.Data.(type) {

	case *PageActivityData_AccessUpdated:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: Updated
	case *PageActivityData_Updated:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *PageFilesChange) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *PageUpdated) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: FilesChange
	if m.FilesChange != nil {
		if v, ok := any(m.GetFilesChange()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
