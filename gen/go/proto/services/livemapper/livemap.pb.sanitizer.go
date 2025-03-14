// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/livemapper/livemap.proto

package livemapper

func (m *CreateOrUpdateMarkerRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Marker
	if m.Marker != nil {
		if v, ok := any(m.GetMarker()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *CreateOrUpdateMarkerResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Marker
	if m.Marker != nil {
		if v, ok := any(m.GetMarker()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DeleteMarkerRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *DeleteMarkerResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *JobsList) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Markers
	for idx, item := range m.Markers {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Users
	for idx, item := range m.Users {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *MarkerMarkersUpdates) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Updated
	for idx, item := range m.Updated {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *StreamRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: MarkersUpdatedAt
	if m.MarkersUpdatedAt != nil {
		if v, ok := any(m.GetMarkersUpdatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *StreamResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Jobs
	switch v := m.Data.(type) {

	case *StreamResponse_Jobs:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: Markers
	case *StreamResponse_Markers:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: Users
	case *StreamResponse_Users:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *UserMarkersUpdates) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Updated
	for idx, item := range m.Updated {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}
