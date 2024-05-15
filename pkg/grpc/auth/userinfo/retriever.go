package userinfo

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
)

const cacheTTL = 20 * time.Second

var (
	ErrAccountError = fmt.Errorf("failed to retrieve account data")
)

type UserInfoRetriever interface {
	GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error)
	GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*UserInfo, error)
	SetUserInfo(ctx context.Context, accountId uint64, superuser bool, job *string, jobGrade *int32) error
}

type UIRetriever struct {
	db *sql.DB

	userCache    *cache.Cache[int32, *UserInfo]
	userCacheTTL time.Duration

	superuserGroups []string
	superuserUsers  []string
}

type Params struct {
	fx.In

	LC     fx.Lifecycle
	DB     *sql.DB
	Config *config.Config
}

func NewUIRetriever(p Params) UserInfoRetriever {
	ctx, cancel := context.WithCancel(context.Background())

	userCache := cache.NewContext(
		ctx,
		cache.AsLRU[int32, *UserInfo](lru.WithCapacity(350)),
		cache.WithJanitorInterval[int32, *UserInfo](cacheTTL),
	)

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return &UIRetriever{
		db: p.DB,

		userCache:    userCache,
		userCacheTTL: cacheTTL,

		superuserGroups: p.Config.Auth.SuperuserGroups,
		superuserUsers:  p.Config.Auth.SuperuserUsers,
	}
}

func (ui *UIRetriever) GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error) {
	var dest *UserInfo
	var ok bool
	if dest, ok = ui.userCache.Get(userId); ok {
		return dest, nil
	}
	if dest == nil {
		dest = &UserInfo{}
	}

	stmt := tUsers.
		SELECT(
			tFivenetAccounts.ID.AS("userinfo.account_id"),
			tFivenetAccounts.Enabled.AS("userinfo.enabled"),
			tFivenetAccounts.License.AS("userinfo.license"),
			tFivenetAccounts.OverrideJob.AS("userinfo.override_job"),
			tFivenetAccounts.OverrideJobGrade.AS("userinfo.override_job_grade"),
			tFivenetAccounts.Superuser.AS("userinfo.superuser"),
			tFivenetAccounts.LastChar.AS("userinfo.last_char"),
			tUsers.ID.AS("userinfo.userid"),
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Group,
		).
		FROM(
			tFivenetAccounts,
			tUsers,
		).
		WHERE(jet.AND(
			tFivenetAccounts.ID.EQ(jet.Uint64(accountId)),
			tUsers.ID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, ui.db, dest); err != nil {
		return nil, errswrap.NewError(err, ErrAccountError)
	}

	// If account is not enabled, fail here
	if !dest.Enabled {
		return nil, ErrAccountError
	}

	// Check if user is superuser
	if slices.Contains(ui.superuserGroups, dest.Group) || slices.Contains(ui.superuserUsers, dest.License) {
		dest.CanBeSuper = true

		if dest.OverrideJob != nil && *dest.OverrideJob != "" {
			dest.Job = *dest.OverrideJob
			dest.JobGrade = *dest.OverrideJobGrade
		}
	}

	ui.userCache.Set(userId, dest, cache.WithExpiration(ui.userCacheTTL))

	return dest, nil
}

func (ui *UIRetriever) GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*UserInfo, error) {
	dest := &UserInfo{}

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("userinfo.userid"),
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Group,
		).
		FROM(tUsers).
		WHERE(jet.AND(
			tUsers.ID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, ui.db, dest); err != nil {
		return nil, errswrap.NewError(err, ErrAccountError)
	}

	// Check if user is superuser
	if slices.Contains(ui.superuserGroups, dest.Group) {
		dest.CanBeSuper = true
	}

	return dest, nil
}

func (ui *UIRetriever) SetUserInfo(ctx context.Context, accountId uint64, superuser bool, job *string, jobGrade *int32) error {
	stmt := tFivenetAccounts.
		UPDATE(
			tFivenetAccounts.Superuser,
			tFivenetAccounts.OverrideJob,
			tFivenetAccounts.OverrideJobGrade,
		).
		SET(
			superuser,
			job,
			jobGrade,
		).
		WHERE(
			tFivenetAccounts.ID.EQ(jet.Uint64(accountId)),
		)

	if _, err := stmt.ExecContext(ctx, ui.db); err != nil {
		return err
	}

	return nil
}
