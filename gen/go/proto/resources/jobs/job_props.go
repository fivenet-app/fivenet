package jobs

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
	tJobProps := table.FivenetJobProps.AS("job_props")
	tFiles := table.FivenetFiles.AS("logo_file")

	stmt := tJobProps.
		SELECT(
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.DeletedAt,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordLastSync,
			tJobProps.DiscordSyncSettings,
			tJobProps.DiscordSyncChanges,
			tJobProps.LogoFileID,
			tFiles.ID,
			tFiles.FilePath,
		).
		FROM(
			tJobProps.
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
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
	if x.GetJob() == "" {
		x.Job = job
	}

	if x.GetQuickButtons() == nil {
		x.QuickButtons = &QuickButtons{
			PenaltyCalculator: false,
			MathCalculator:    false,
		}
	}

	if x.GetLivemapMarkerColor() == "" {
		x.LivemapMarkerColor = DefaultLivemapMarkerColor
	}

	// Discord Sync Settings
	if x.GetDiscordSyncSettings() == nil {
		x.DiscordSyncSettings = &DiscordSyncSettings{
			DryRun:                   false,
			UserInfoSync:             false,
			QualificationsRoleFormat: DefaultQualificationsRoleFormat,
		}
	}

	if x.GetDiscordSyncSettings().GetUserInfoSyncSettings() == nil {
		x.DiscordSyncSettings.UserInfoSyncSettings = &UserInfoSyncSettings{
			EmployeeRoleEnabled: true,
			UnemployedEnabled:   false,
			UnemployedMode:      UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE,
			SyncNicknames:       true,
			GroupMapping:        []*GroupMapping{},
		}
	}

	employeeRoleFormat := DefaultEmployeeRoleFormat
	if x.GetDiscordSyncSettings().GetUserInfoSyncSettings().GetEmployeeRoleFormat() == "" {
		x.DiscordSyncSettings.UserInfoSyncSettings.EmployeeRoleFormat = employeeRoleFormat
	}

	gradeRoleFormat := DefaultGradeRoleFormat
	if x.GetDiscordSyncSettings().GetUserInfoSyncSettings().GetGradeRoleFormat() == "" {
		x.DiscordSyncSettings.UserInfoSyncSettings.GradeRoleFormat = gradeRoleFormat
	}

	unemployedRoleName := DefaultUnemployedRoleName
	if x.GetDiscordSyncSettings().GetUserInfoSyncSettings().GetUnemployedRoleName() == "" {
		x.DiscordSyncSettings.UserInfoSyncSettings.UnemployedRoleName = unemployedRoleName
	}

	// Status Log Settings
	if x.GetDiscordSyncSettings().GetStatusLogSettings() == nil {
		x.DiscordSyncSettings.StatusLogSettings = &StatusLogSettings{}
	}

	// Jobs Abscene Role
	if x.GetDiscordSyncSettings().GetJobsAbsenceSettings() == nil {
		x.DiscordSyncSettings.JobsAbsenceSettings = &JobsAbsenceSettings{
			AbsenceRole: DefaultJobsAbsenceRoleName,
		}
	}

	// Group Sync Settings
	if x.GetDiscordSyncSettings().GetGroupSyncSettings() == nil {
		x.DiscordSyncSettings.GroupSyncSettings = &GroupSyncSettings{
			IgnoredRoleIds: []string{},
		}
	}

	if x.GetDiscordSyncSettings().GetQualificationsRoleFormat() == "" {
		x.DiscordSyncSettings.QualificationsRoleFormat = DefaultQualificationsRoleFormat
	}

	// Job Settings
	if x.GetSettings() == nil {
		x.Settings = &JobSettings{}
	}
	x.GetSettings().Default()
}
