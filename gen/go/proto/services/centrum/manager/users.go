package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Manager) ResolveUserById(ctx context.Context, u int32) (*users.User, error) {
	tUsers := tUsers.AS("user")
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
			tUserProps.Avatar.AS("usershort.avatar"),
		).
		FROM(
			tUsers.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
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

	return &dest, nil
}

func (s *Manager) resolveUserShortById(ctx context.Context, u int32) (*users.UserShort, error) {
	us, err := s.resolveUserShortsByIds(ctx, u)
	if err != nil {
		return nil, err
	}

	return us[0], nil
}

func (s *Manager) resolveUserShortsByIds(ctx context.Context, u ...int32) ([]*users.UserShort, error) {
	if len(u) == 0 {
		return nil, nil
	}

	userIds := make([]jet.Expression, len(u))
	for i := 0; i < len(u); i++ {
		userIds[i] = jet.Int32(u[i])
	}

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
			tUserProps.Avatar.AS("usershort.avatar"),
		).
		FROM(
			tUsers.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				),
		).
		WHERE(
			tUsers.ID.IN(userIds...),
		).
		LIMIT(int64(len(u)))

	dest := []*users.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to resolve usershorts by ids %+v: %w", u, err)
	}

	for i := 0; i < len(dest); i++ {
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

	us, err := s.resolveUserShortsByIds(ctx, userIds...)
	if err != nil {
		return err
	}

	for i := 0; i < len(*u); i++ {
		(*u)[i].User = us[i]
	}

	return nil
}
