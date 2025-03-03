// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/users/job_settings.proto

package users

import (
	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
)

func (m *DiscordSyncChange) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Time
	if m.Time != nil {
		if v, ok := any(m.GetTime()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DiscordSyncChanges) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Changes
	for idx, item := range m.Changes {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *DiscordSyncSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: GroupSyncSettings
	if m.GroupSyncSettings != nil {
		if v, ok := any(m.GetGroupSyncSettings()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: JobsAbsenceSettings
	if m.JobsAbsenceSettings != nil {
		if v, ok := any(m.GetJobsAbsenceSettings()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: StatusLogSettings
	if m.StatusLogSettings != nil {
		if v, ok := any(m.GetStatusLogSettings()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: UserInfoSyncSettings
	if m.UserInfoSyncSettings != nil {
		if v, ok := any(m.GetUserInfoSyncSettings()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *GroupMapping) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *GroupSyncSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: IgnoredRoleIds
	for idx, item := range m.IgnoredRoleIds {
		_, _ = idx, item

		m.IgnoredRoleIds[idx] = htmlsanitizer.StripTags(m.IgnoredRoleIds[idx])

	}

	return nil
}

func (m *JobSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *JobsAbsenceSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *StatusLogSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *UserInfoSyncSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: GroupMapping
	for idx, item := range m.GroupMapping {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}
