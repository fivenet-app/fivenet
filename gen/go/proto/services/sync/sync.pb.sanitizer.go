// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/sync/sync.proto

package sync

func (m *AddActivityRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Dispatch
	switch v := m.Activity.(type) {

	case *AddActivityRequest_Dispatch:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: JobsTimeclock
	case *AddActivityRequest_JobsTimeclock:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: JobsUserActivity
	case *AddActivityRequest_JobsUserActivity:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: JobsUserProps
	case *AddActivityRequest_JobsUserProps:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: UserActivity
	case *AddActivityRequest_UserActivity:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: UserOauth2
	case *AddActivityRequest_UserOauth2:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: UserProps
	case *AddActivityRequest_UserProps:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: UserUpdate
	case *AddActivityRequest_UserUpdate:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *AddActivityResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *GetStatusRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *GetStatusResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Jobs
	if m.Jobs != nil {
		if v, ok := interface{}(m.GetJobs()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Licenses
	if m.Licenses != nil {
		if v, ok := interface{}(m.GetLicenses()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Users
	if m.Users != nil {
		if v, ok := interface{}(m.GetUsers()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Vehicles
	if m.Vehicles != nil {
		if v, ok := interface{}(m.GetVehicles()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *SendDataRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Jobs
	switch v := m.Data.(type) {

	case *SendDataRequest_Jobs:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: Licenses
	case *SendDataRequest_Licenses:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: UserLocations
	case *SendDataRequest_UserLocations:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: Users
	case *SendDataRequest_Users:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

		// Field: Vehicles
	case *SendDataRequest_Vehicles:
		if v, ok := interface{}(v).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *SendDataResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *StreamRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *StreamResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}
