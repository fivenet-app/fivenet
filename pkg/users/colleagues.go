package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func RetrieveColleagueById(
	ctx context.Context,
	db *sql.DB,
	enricher *mstlystcdata.Enricher,
	u ...int32,
) ([]*jobs.Colleague, error) {
	if len(u) == 0 {
		return nil, nil
	}

	userIds := make([]jet.Expression, len(u))
	for i := range u {
		userIds[i] = jet.Int32(u[i])
	}

	tUsers := tables.User().AS("colleague")
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")
	tUserProps := table.FivenetUserProps.AS("user_props")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
		).
		FROM(
			tUsers.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tUsers.ID).
						AND(tColleagueProps.Job.EQ(tUsers.Job)),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tUsers.ID.IN(userIds...),
		).
		LIMIT(int64(len(u)))

	dest := []*jobs.Colleague{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		return nil, fmt.Errorf("failed to retrieve colleagues by ids %+v. %w", u, err)
	}
	for i := range dest {
		if dest[i] != nil {
			enricher.EnrichJobInfo(dest[i])
		}
	}

	return dest, nil
}

func RetrieveUserShortById(
	ctx context.Context,
	db *sql.DB,
	enricher *mstlystcdata.Enricher,
	u int32,
) (*jobs.Colleague, error) {
	us, err := RetrieveColleagueById(ctx, db, enricher, u)
	if err != nil {
		return nil, err
	}

	return us[0], nil
}

func RetrieveUsersForUnit(
	ctx context.Context,
	db *sql.DB,
	enricher *mstlystcdata.Enricher,
	u *[]*centrum.UnitAssignment,
) error {
	userIds := make([]int32, len(*u))
	for i := range *u {
		userIds[i] = (*u)[i].GetUserId()
	}

	if len(userIds) == 0 {
		return nil
	}

	us, err := RetrieveColleagueById(ctx, db, enricher, userIds...)
	if err != nil {
		return err
	}

	for i := range *u {
		(*u)[i].User = us[i]
	}

	return nil
}

func RetrieveUserById(ctx context.Context, db *sql.DB, u int32) (*users.User, error) {
	tUsers := tables.User().AS("user")
	tUserProps := table.FivenetUserProps.AS("user_props")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
		).
		FROM(
			tUsers.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tUsers.ID.EQ(jet.Int32(u)),
		).
		LIMIT(1)

	dest := users.User{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to retrieve user by id %d. %w", u, err)
		}

		return nil, nil
	}

	if dest.GetUserId() == 0 {
		return nil, nil
	}

	return &dest, nil
}
