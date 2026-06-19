package authstore

import (
	"context"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

var (
	tUserProps = table.FivenetUserProps.AS("user_props")
	tJobProps  = table.FivenetJobProps.AS("job_props")
)

func (s *Store) ListCharacters(
	ctx context.Context,
	accountID int64,
	license string,
) ([]*accounts.Character, error) {
	tUsers := table.FivenetUser.AS("user")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			dbutils.Columns{
				tUsers.ID.AS("user.user_id"),
				tUsers.AccountID,
				tUsers.Identifier,
				tUsers.Job,
				tUsers.JobGrade,
				tUsers.Firstname,
				tUsers.Lastname,
				tUsers.Dateofbirth,
				tUsers.Sex,
				tUsers.Height,
				tUsers.PhoneNumber,
				tUserProps.AvatarFileID.AS("user.profile_picture_file_id"),
				tAvatar.FilePath.AS("user.profile_picture"),
				tUsers.Group.AS("character.group"),
				s.customDB.Columns.User.GetVisum(tUsers.Alias()),
				s.customDB.Columns.User.GetPlaytime(tUsers.Alias()),
			}.Get()...,
		).
		FROM(tUsers.
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUsers.ID),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(mysql.OR(
			tUsers.AccountID.EQ(mysql.Int64(accountID)),
			tUsers.License.EQ(mysql.String(license)),
		)).
		ORDER_BY(tUsers.ID).
		LIMIT(10)

	chars := []*accounts.Character{}
	if err := stmt.QueryContext(ctx, s.db, &chars); err != nil {
		return nil, err
	}

	return chars, nil
}

func (s *Store) GetCharacter(
	ctx context.Context,
	charID int32,
) (*users.User, *jobsprops.JobProps, error) {
	tUsers := table.FivenetUser.AS("user")
	tLogo := table.FivenetFiles.AS("logo_file")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.ID.AS("user.user_id"),
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUserProps.AvatarFileID.AS("user.profile_picture_file_id"),
			tAvatar.FilePath.AS("user.profile_picture"),
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.DeletedAt,
			tJobProps.DiscordGuildID,
			tJobProps.LivemapMarkerColor,
			tJobProps.QuickButtons,
			tJobProps.RadioFrequency,
			tJobProps.LogoFileID,
			tLogo.ID,
			tLogo.FilePath,
		).
		FROM(
			tUsers.
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tUsers.Job),
				).
				LEFT_JOIN(tLogo,
					tLogo.ID.EQ(tJobProps.LogoFileID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(tUsers.ID.EQ(mysql.Int32(charID))).
		LIMIT(1)

	var dest struct {
		*users.User

		JobProps *jobsprops.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, nil, err
	}

	return dest.User, dest.JobProps, nil
}

func (s *Store) GetJobWithProps(
	ctx context.Context,
	jobName string,
) (*jobs.Job, int32, *jobsprops.JobProps, error) {
	tJobs := table.FivenetJobs.AS("job")
	tJobsGrades := table.FivenetJobsGrades.AS("jg")
	tFiles := table.FivenetFiles.AS("logo_file")

	stmt := tJobs.
		SELECT(
			tJobs.Name,
			tJobs.Label,
			tJobsGrades.Grade.AS("job_grade"),
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.DeletedAt,
			tJobProps.DiscordGuildID,
			tJobProps.LivemapMarkerColor,
			tJobProps.QuickButtons,
			tJobProps.RadioFrequency,
			tJobProps.LogoFileID,
			tFiles.ID,
			tFiles.FilePath,
		).
		FROM(
			tJobs.
				INNER_JOIN(tJobsGrades,
					tJobsGrades.JobName.EQ(tJobs.Name),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tJobs.Name),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(
			tJobs.Name.EQ(mysql.String(jobName)),
		).
		ORDER_BY(tJobsGrades.Grade.DESC()).
		LIMIT(1)

	var dest struct {
		Job      *jobs.Job
		JobGrade int32 `alias:"job_grade"`
		JobProps *jobsprops.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, 0, nil, err
	}

	return dest.Job, dest.JobGrade, dest.JobProps, nil
}
