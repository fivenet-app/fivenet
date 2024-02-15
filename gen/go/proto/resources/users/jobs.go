package users

import (
	"database/sql/driver"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	DefaultLivemapMarkerColor = "#5c7aff"

	DefaultEmployeeRoleFormat = "%s Personal"
	DefaultGradeRoleFormat    = "[%grade%] %grade_label%"
)

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}

func (x *JobProps) Default(job string) {
	if x.Job == "" {
		x.Job = job
	}

	if x.Theme == "" {
		x.Theme = "defaultTheme"
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

	if x.DiscordSyncSettings == nil {
		x.DiscordSyncSettings = &DiscordSyncSettings{
			UserInfoSync: false,
		}
	}

	if x.DiscordSyncSettings.UserInfoSyncSettings == nil {
		x.DiscordSyncSettings.UserInfoSyncSettings = &UserInfoSyncSettings{
			EmployeeRoleEnabled: true,
			UnemployedEnabled:   false,
			UnemployedMode:      UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE,
		}
	}

	employeeRoleFormat := DefaultEmployeeRoleFormat
	if x.DiscordSyncSettings.UserInfoSyncSettings.EmployeeRoleFormat == nil {
		x.DiscordSyncSettings.UserInfoSyncSettings.EmployeeRoleFormat = &employeeRoleFormat
	}

	gradeRoleFormat := DefaultGradeRoleFormat
	if x.DiscordSyncSettings.UserInfoSyncSettings.GradeRoleFormat == nil {
		x.DiscordSyncSettings.UserInfoSyncSettings.GradeRoleFormat = &gradeRoleFormat
	}

	unemployedRoleName := "Citizen"
	if x.DiscordSyncSettings.UserInfoSyncSettings.UnemployedRoleName == nil {
		x.DiscordSyncSettings.UserInfoSyncSettings.UnemployedRoleName = &unemployedRoleName
	}
}

func (x *QuickButtons) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.UnmarshalFromString(t, x)
	case []byte:
		return json.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf QuickButtons.
func (x *QuickButtons) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.MarshalToString(x)
	return out, err
}

func (x *DiscordSyncSettings) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf DiscordSyncSettings.
func (x *DiscordSyncSettings) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protojson.Marshal(x)
	return string(out), err
}
