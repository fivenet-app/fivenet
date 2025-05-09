package centrummanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Manager) ResolveUserById(ctx context.Context, u int32) (*users.User, error) {
	tUsers := tables.Users().AS("user")

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
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.NamePrefix,
			tJobsUserProps.NameSuffix,
			tUserProps.Avatar.AS("user.avatar"),
		).
		FROM(
			tUsers.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUsers.ID).
						AND(tJobsUserProps.Job.EQ(tUsers.Job)),
				),
		).
		WHERE(
			tUsers.ID.EQ(jet.Int32(u)),
		).
		LIMIT(1)

	dest := users.User{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to resolve user by id %d: %w", u, err)
		}

		return nil, nil
	}

	if dest.UserId == 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Manager) resolveUserShortById(ctx context.Context, u int32) (*jobs.Colleague, error) {
	us, err := s.resolveColleagueById(ctx, u)
	if err != nil {
		return nil, err
	}

	return us[0], nil
}

func (s *Manager) resolveColleagueById(ctx context.Context, u ...int32) ([]*jobs.Colleague, error) {
	if len(u) == 0 {
		return nil, nil
	}

	userIds := make([]jet.Expression, len(u))
	for i := range u {
		userIds[i] = jet.Int32(u[i])
	}

	tUsers := tables.Users().AS("colleague")

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
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.NamePrefix,
			tJobsUserProps.NameSuffix,
			tUserProps.Avatar.AS("colleague.avatar"),
		).
		FROM(
			tUsers.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUsers.ID).
						AND(tJobsUserProps.Job.EQ(tUsers.Job)),
				),
		).
		WHERE(
			tUsers.ID.IN(userIds...),
		).
		LIMIT(int64(len(u)))

	dest := []*jobs.Colleague{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to resolve colleagues by ids %+v: %w", u, err)
	}
	for i := range dest {
		if dest[i] != nil {
			s.enricher.EnrichJobInfo(dest[i])
		}
	}

	return dest, nil
}

func (s *Manager) resolveUsersForUnit(ctx context.Context, u *[]*centrum.UnitAssignment) error {
	userIds := make([]int32, len(*u))
	for i := 0; i < len(*u); i++ {
		userIds[i] = (*u)[i].UserId
	}

	if len(userIds) == 0 {
		return nil
	}

	us, err := s.resolveColleagueById(ctx, userIds...)
	if err != nil {
		return err
	}

	for i := 0; i < len(*u); i++ {
		(*u)[i].User = us[i]
	}

	return nil
}
