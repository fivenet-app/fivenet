package citizensstore

import (
	"context"
	"database/sql"
	"errors"

	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	ListCitizens(
		ctx context.Context,
		req *pbcitizens.ListCitizensRequest,
		opts ListCitizensOptions,
	) (*pbcitizens.ListCitizensResponse, error)
	GetUser(
		ctx context.Context,
		req *pbcitizens.GetUserRequest,
		opts GetUserOptions,
	) (*pbcitizens.GetUserResponse, error)
	ListLabels(
		ctx context.Context,
		q qrm.Queryable,
		condition mysql.BoolExpression,
		includeDeleted bool,
	) (*citizenslabels.Labels, error)
	NextLabelSortOrder(ctx context.Context, q qrm.Queryable, job string) (int32, error)
	GetLabel(
		ctx context.Context,
		q qrm.Queryable,
		job string,
		labelId int64,
	) (*citizenslabels.Label, error)
	UpdateLabel(ctx context.Context, tx qrm.DB, label *citizenslabels.Label, job string) error
	InsertLabel(ctx context.Context, tx qrm.DB, label *citizenslabels.Label) (int64, error)
	DeleteLabel(
		ctx context.Context,
		tx qrm.DB,
		job string,
		labelId int64,
		deletedAt *timestamp.Timestamp,
	) error
	ReorderLabels(ctx context.Context, job string, labelIds []int64) error
	GetUserLabels(
		ctx context.Context,
		q qrm.Queryable,
		condition mysql.BoolExpression,
	) (*citizenslabels.Labels, error)
	ValidateLabels(
		ctx context.Context,
		userJob string,
		labels []*citizenslabels.Label,
	) (bool, error)
	GetUserAccess(ctx context.Context, userId int32) (*users.User, error)
	ListExpiredWantedUserProps(ctx context.Context, maxDays int64, limit int64) ([]int32, error)
	GetAvatarFileID(ctx context.Context, userId int32) (*int64, error)
	GetMugshotFileID(ctx context.Context, userId int32) (*int64, error)
	GetUserProps(ctx context.Context, tx qrm.DB, userId int32) (*usersprops.UserProps, error)
	HandleUserPropsChanges(
		ctx context.Context,
		tx qrm.DB,
		x *usersprops.UserProps,
		in *usersprops.UserProps,
		sourceUserId *int32,
		reason string,
	) ([]*usersactivity.UserActivity, error)
	ListUserActivity(
		ctx context.Context,
		req *pbcitizens.ListUserActivityRequest,
		limit int64,
	) ([]*usersactivity.UserActivity, error)
	CountUserActivity(ctx context.Context, req *pbcitizens.ListUserActivityRequest) (int64, error)
}

type Store struct {
	db       *sql.DB
	customDB *config.CustomDB
}

func New(db *sql.DB, customDB *config.CustomDB) IStore {
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
	tUserProps := table.FivenetUserProps

	stmt := tUserProps.
		SELECT(tUserProps.UserID).
		FROM(tUserProps).
		WHERE(mysql.AND(
			tUserProps.Wanted.IS_TRUE(),
			mysql.OR(
				tUserProps.WantedAt.LT(
					mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(maxDays, "DAY")),
				),
				tUserProps.WantedTill.LT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		LIMIT(limit)

	var dest []int32
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
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
		if !errors.Is(err, qrm.ErrNoRows) {
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
