package citizens

import (
	"context"
	"database/sql"

	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type Store struct {
	db       *sql.DB
	customDB config.CustomDB
}

func New(db *sql.DB, customDB config.CustomDB) *Store {
	return &Store{
		db:       db,
		customDB: customDB,
	}
}

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
		if err != qrm.ErrNoRows {
			return nil, err
		}
	}

	return u, nil
}

func (s *Store) ListExpiredWantedUserProps(ctx context.Context, maxDays int64, limit int64) ([]int32, error) {
	tUserProps := table.FivenetUserProps

	stmt := tUserProps.
		SELECT(tUserProps.UserID).
		FROM(tUserProps).
		WHERE(mysql.AND(
			tUserProps.Wanted.IS_TRUE(),
			mysql.OR(
				tUserProps.WantedAt.LT(mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(maxDays, "DAY"))),
				tUserProps.WantedTill.LT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		LIMIT(limit)

	var dest []int32
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if err != qrm.ErrNoRows {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) GetAvatarFileID(ctx context.Context, userId int32) (*int64, error) {
	tUserProps := table.FivenetUserProps

	stmt := tUserProps.
		SELECT(tUserProps.AvatarFileID).
		WHERE(tUserProps.UserID.EQ(mysql.Int32(userId))).
		LIMIT(1)

	var props struct {
		AvatarFileID *int64
	}
	if err := stmt.QueryContext(ctx, s.db, &props); err != nil {
		if err != qrm.ErrNoRows {
			return nil, err
		}
	}

	return props.AvatarFileID, nil
}

func (s *Store) GetMugshotFileID(ctx context.Context, userId int32) (*int64, error) {
	props, err := s.GetUserProps(ctx, s.db, userId)
	if err != nil {
		return nil, err
	}
	return props.MugshotFileId, nil
}
