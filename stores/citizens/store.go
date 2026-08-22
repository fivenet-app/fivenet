package citizensstore

import (
	"context"
	"database/sql"

	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
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
		userInfo *userinfo.UserInfo,
		search string,
		ownJobOnly bool,
		canCreateLabel bool,
		minAccess int32,
		includeDeleted bool,
	) (*citizenslabels.Labels, error)
	GetUserLabelsForUser(
		ctx context.Context,
		userInfo *userinfo.UserInfo,
		userId int32,
	) (*citizenslabels.Labels, error)
	NextLabelSortOrder(ctx context.Context, q qrm.Queryable, job string) (int32, error)
	GetLabel(
		ctx context.Context,
		q qrm.Queryable,
		job string,
		labelId int64,
		includeDeleted bool,
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
		ctes []mysql.CommonTableExpression,
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
		opts ListUserActivityOptions,
	) ([]*usersactivity.UserActivity, error)
	CountUserActivity(ctx context.Context, opts CountUserActivityOptions) (int64, error)
}

type Store struct {
	db                 *sql.DB
	customDB           *config.CustomDB
	labelsAccess       *access.CitizenLabelsObjectAccess
	userActivitySorter *database.SorterBuilder
}

type Params struct {
	fx.In

	DB           *sql.DB
	CustomDB     *config.CustomDB
	LabelsAccess *access.CitizenLabelsObjectAccess
}

func New(p Params) IStore {
	return &Store{
		db:           p.DB,
		customDB:     p.CustomDB,
		labelsAccess: p.LabelsAccess,
		userActivitySorter: database.New(
			database.SpecMap{
				"createdAt": database.Column{
					Col:       table.FivenetUserActivity.AS("user_activity").CreatedAt,
					NullsLast: true,
				},
			},
			[]mysql.OrderByClause{
				table.FivenetUserActivity.AS("user_activity").CreatedAt.DESC().NULLS_LAST(),
			},
			[]mysql.OrderByClause{table.FivenetUserActivity.AS("user_activity").ID.DESC()},
			"createdAt",
			3,
		),
	}
}
