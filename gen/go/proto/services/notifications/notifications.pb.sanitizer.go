// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/notifications/notifications.proto

package notifications

func (m *GetNotificationsRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Categories
	for idx, item := range m.Categories {
		_, _ = idx, item

	}

	// Field: Pagination
	if m.Pagination != nil {
		if v, ok := any(m.GetPagination()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *GetNotificationsResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Notifications
	for idx, item := range m.Notifications {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Pagination
	if m.Pagination != nil {
		if v, ok := any(m.GetPagination()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *MarkNotificationsRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *MarkNotificationsResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *StreamRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: ClientView
	switch v := m.Data.(type) {

	case *StreamRequest_ClientView:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
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

	// Field: JobEvent
	switch v := m.Data.(type) {

	case *StreamResponse_JobEvent:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: JobGradeEvent
	case *StreamResponse_JobGradeEvent:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: MailerEvent
	case *StreamResponse_MailerEvent:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: ObjectEvent
	case *StreamResponse_ObjectEvent:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: SystemEvent
	case *StreamResponse_SystemEvent:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: UserEvent
	case *StreamResponse_UserEvent:
		if v, ok := any(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}
