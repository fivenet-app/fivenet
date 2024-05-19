package users

import (
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	DefaultTheme              = "defaultTheme"
	DefaultLivemapMarkerColor = "#5c7aff"

	DefaultEmployeeRoleFormat  = "%s Personal"
	DefaultGradeRoleFormat     = "[%grade%] %grade_label%"
	DefaultUnemployedRoleName  = "Citizen"
	DefaultJobsAbsenceRoleName = "Absent"
)

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}

func (x *JobProps) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *JobProps) Default(job string) {
	if x.Job == "" {
		x.Job = job
	}

	if x.Theme == "" {
		x.Theme = DefaultTheme
	}

	if x.QuickButtons == nil {
		x.QuickButtons = &QuickButtons{
			BodyCheckup:       false,
			PenaltyCalculator: false,
		}
	}

	if x.LivemapMarkerColor == "" {
		x.LivemapMarkerColor = DefaultLivemapMarkerColor
	}

	// Discord Sync Settings
	if x.DiscordSyncSettings == nil {
		x.DiscordSyncSettings = &DiscordSyncSettings{
			DryRun:       false,
			UserInfoSync: false,
		}
	}

	if x.DiscordSyncSettings.UserInfoSyncSettings == nil {
		x.DiscordSyncSettings.UserInfoSyncSettings = &UserInfoSyncSettings{
			EmployeeRoleEnabled: true,
			UnemployedEnabled:   false,
			UnemployedMode:      UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE,
			SyncNicknames:       true,
			GroupMapping:        []*GroupMapping{},
		}
	}

	employeeRoleFormat := DefaultEmployeeRoleFormat
	if x.DiscordSyncSettings.UserInfoSyncSettings.EmployeeRoleFormat == "" {
		x.DiscordSyncSettings.UserInfoSyncSettings.EmployeeRoleFormat = employeeRoleFormat
	}

	gradeRoleFormat := DefaultGradeRoleFormat
	if x.DiscordSyncSettings.UserInfoSyncSettings.GradeRoleFormat == "" {
		x.DiscordSyncSettings.UserInfoSyncSettings.GradeRoleFormat = gradeRoleFormat
	}

	unemployedRoleName := DefaultUnemployedRoleName
	if x.DiscordSyncSettings.UserInfoSyncSettings.UnemployedRoleName == "" {
		x.DiscordSyncSettings.UserInfoSyncSettings.UnemployedRoleName = unemployedRoleName
	}

	// Status Log Settings
	if x.DiscordSyncSettings.StatusLogSettings == nil {
		x.DiscordSyncSettings.StatusLogSettings = &StatusLogSettings{}
	}

	// Jobs Abscene Role
	if x.DiscordSyncSettings.JobsAbsenceSettings == nil {
		x.DiscordSyncSettings.JobsAbsenceSettings = &JobsAbsenceSettings{
			AbsenceRole: DefaultJobsAbsenceRoleName,
		}
	}

	// Group Sync Settings
	if x.DiscordSyncSettings.GroupSyncSettings == nil {
		x.DiscordSyncSettings.GroupSyncSettings = &GroupSyncSettings{
			IgnoredRoleIds: []string{},
		}
	}
}

// Scan implements driver.Valuer for protobuf QuickButtons.
func (x *QuickButtons) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *QuickButtons) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

// Scan implements driver.Valuer for protobuf DiscordSyncSettings.
func (x *DiscordSyncSettings) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *DiscordSyncSettings) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func (x *DiscordSyncSettings) IsStatusLogEnabled() bool {
	return x.StatusLog && x.StatusLogSettings != nil && x.StatusLogSettings.ChannelId != ""
}

// Scan implements driver.Valuer for protobuf CitizenAttributes.
func (x *CitizenAttributes) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *CitizenAttributes) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func (x *CitizenAttribute) Equal(a *CitizenAttribute) bool {
	return x.Name == a.Name
}

// Scan implements driver.Valuer for protobuf JobSettings.
func (x *JobSettings) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *JobSettings) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

// Scan implements driver.Valuer for protobuf DiscordSyncDiff.
func (x *DiscordSyncDiff) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *DiscordSyncDiff) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
