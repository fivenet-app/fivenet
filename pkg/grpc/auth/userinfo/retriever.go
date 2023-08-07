package userinfo

import (
	"context"
	"database/sql"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
)

type UserInfoRetriever interface {
	GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error)
	GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*UserInfo, error)
	SetUserInfo(ctx context.Context, accountId uint64, job string, jobGrade int32) error
}

type UIRetriever struct {
	ctx context.Context
	db  *sql.DB

	userCache    *cache.Cache[int32, *UserInfo]
	userCacheTTL time.Duration

	superuserGroups []string
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
		cache.AsLRU[int32, *UserInfo](lru.WithCapacity(300)),
		cache.WithJanitorInterval[int32, *UserInfo](40*time.Second),
	)

	p.LC.Append(fx.StopHook(func(_ context.Context) {
		cancel()
	}))

	return &UIRetriever{
		ctx: ctx,
		db:  p.DB,

		userCache:    userCache,
		userCacheTTL: 30 * time.Second,

		superuserGroups: p.Config.Game.SuperuserGroups,
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
			tUsers.ID.AS("userinfo.userid"),
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Group,
			tFivenetAccounts.ID.AS("userinfo.acc_id"),
			tFivenetAccounts.OverrideJob.AS("userinfo.orig_job"),
			tFivenetAccounts.OverrideJobGrade.AS("userinfo.orig_job_grade"),
		).
		FROM(
			tUsers,
			tFivenetAccounts,
		).
		WHERE(jet.AND(
			tUsers.ID.EQ(jet.Int32(userId)),
			tFivenetAccounts.ID.EQ(jet.Uint64(accountId)),
		)).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, ui.db, dest); err != nil {
		return nil, err
	}

	// Check if user is superuser
	if utils.InSlice(ui.superuserGroups, dest.Group) {
		dest.SuperUser = true
		if dest.OrigJob != "" {
			dest.Job = dest.OrigJob
			dest.JobGrade = dest.OrigJobGrade
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
		FROM(
			tUsers,
			tFivenetAccounts,
		).
		WHERE(jet.AND(
			tUsers.ID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, ui.db, dest); err != nil {
		return nil, err
	}

	// Check if user is superuser
	if utils.InSlice(ui.superuserGroups, dest.Group) {
		dest.SuperUser = true
	}

	return dest, nil
}

func (ui *UIRetriever) SetUserInfo(ctx context.Context, accountId uint64, job string, jobGrade int32) error {
	jobVal := jet.NULL
	jobGradeVal := jet.NULL
	if job != "" && jobGrade > 0 {
		jobVal = jet.String(job)
		jobGradeVal = jet.Int32(jobGrade)
	}

	stmt := tFivenetAccounts.
		UPDATE(
			tFivenetAccounts.OverrideJob,
			tFivenetAccounts.OverrideJobGrade,
		).
		SET(
			tFivenetAccounts.OverrideJob.SET(jet.StringExp(jobVal)),
			tFivenetAccounts.OverrideJobGrade.SET(jet.IntExp(jobGradeVal)),
		).
		WHERE(
			tFivenetAccounts.ID.EQ(jet.Uint64(accountId)),
		)

	if _, err := stmt.ExecContext(ctx, ui.db); err != nil {
		return err
	}

	return nil
}
