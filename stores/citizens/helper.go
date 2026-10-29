package citizensstore

import (
	"context"
	"errors"

	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) GetUserAccess(ctx context.Context, userId int32) (*users.User, error) {
	tUser := table.FivenetUser.AS("user")

	stmt := tUser.
		SELECT(
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(tUser).
		WHERE(tUser.ID.EQ(mysql.Int32(userId))).
		LIMIT(1)

	u := &users.User{}
	if err := stmt.QueryContext(ctx, s.db, u); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return u, nil
}

func (s *Store) ListExpiredWantedUserProps(
	ctx context.Context,
	maxDays int64,
	limit int64,
) ([]int32, error) {
	tValues := table.FivenetUserPropertyValues.AS("value")
	tState := table.FivenetUserPropertyState.AS("state")

	stmt := tValues.
		SELECT(tValues.UserID).
		FROM(tValues.INNER_JOIN(tState, mysql.AND(
			tState.UserID.EQ(tValues.UserID),
			tState.PropertyKey.EQ(tValues.PropertyKey),
		))).
		WHERE(mysql.AND(
			tValues.PropertyKey.EQ(mysql.String("user.wanted")),
			tValues.ValueType.EQ(mysql.Uint8(uint8(userPropertyValueTypeBool))),
			tValues.ValueBool.IS_TRUE(),
			mysql.OR(
				tState.ActiveAt.LT(
					mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(maxDays, mysql.DAY)),
				),
				tState.ExpiresAt.LT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		LIMIT(limit)

	dest := []int32{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	return dest, nil
}

func (s *Store) GetAvatarFileID(ctx context.Context, userId int32) (*int64, error) {
	tUser := table.FivenetUser

	stmt := tUser.
		SELECT(tUser.AvatarFileID.AS("avatar_file_id")).
		WHERE(tUser.ID.EQ(mysql.Int32(userId))).
		LIMIT(1)

	var props struct {
		AvatarFileID *int64
	}
	if err := stmt.QueryContext(ctx, s.db, &props); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return props.AvatarFileID, nil
}
