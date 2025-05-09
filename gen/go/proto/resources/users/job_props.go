package users

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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
			tJobProps.DeletedAt,
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
