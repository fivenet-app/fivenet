package users

import (
	"context"
	"database/sql/driver"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	DefaultTheme              = "defaultTheme"
	DefaultLivemapMarkerColor = "#5c7aff"

	DefaultEmployeeRoleFormat  = "%s Personal"
	DefaultGradeRoleFormat     = "[%grade%] %grade_label%"
	DefaultUnemployedRoleName  = "Citizen"
	DefaultJobsAbsenceRoleName = "Absent"

	DefaultQualificationsRoleFormat = "%name% Qualification"
)

func GetJobProps(ctx context.Context, tx qrm.DB, job string) (*JobProps, error) {
	tJobProps := table.FivenetJobProps.AS("jobprops")
	stmt := tJobProps.
		SELECT(
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.Theme,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordLastSync,
			tJobProps.DiscordSyncSettings,
			tJobProps.DiscordSyncChanges,
			tJobProps.LogoURL,
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		).
		LIMIT(1)

	dest := &JobProps{}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	dest.Default(job)

	return dest, nil
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
			MathCalculator:    false,
		}
	}

	if x.LivemapMarkerColor == "" {
		x.LivemapMarkerColor = DefaultLivemapMarkerColor
	}

	// Discord Sync Settings
	if x.DiscordSyncSettings == nil {
		x.DiscordSyncSettings = &DiscordSyncSettings{
			DryRun:                   false,
			UserInfoSync:             false,
			QualificationsRoleFormat: DefaultQualificationsRoleFormat,
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

	if x.DiscordSyncSettings.QualificationsRoleFormat == "" {
		x.DiscordSyncSettings.QualificationsRoleFormat = DefaultQualificationsRoleFormat
	}

	// Job Settings
	if x.Settings == nil {
		x.Settings = &JobSettings{}
	}
	x.Settings.Default()
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

// Scan implements driver.Valuer for protobuf CitizenLabels.
func (x *CitizenLabels) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *CitizenLabels) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func (x *CitizenLabel) Equal(a *CitizenLabel) bool {
	return x.Name == a.Name
}

// Scan implements driver.Valuer for protobuf DiscordSyncChanges.
func (x *DiscordSyncChanges) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *DiscordSyncChanges) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func (x *DiscordSyncChanges) Add(change *DiscordSyncChange) {
	if x.Changes == nil {
		x.Changes = []*DiscordSyncChange{}
	}

	if len(x.Changes) > 0 {
		lastChange := x.Changes[len(x.Changes)-1]

		if lastChange.Plan == change.Plan {
			return
		}
	}

	x.Changes = append(x.Changes, change)

	if len(x.Changes) > 12 {
		x.Changes = slices.Delete(x.Changes, 0, 1)
	}
}
