package userinfo

import (
	"context"
	"database/sql"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tUsers           = table.Users.AS("userinfo")
	tFivenetAccounts = table.FivenetAccounts
)

type UserInfo struct {
	AccId  uint64
	UserId int32

	Job          string
	JobGrade     int32
	OrigJob      string
	OrigJobGrade int32

	Group     string
	SuperUser bool
}

type UserInfoRetriever interface {
	GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error)
	SetUserInfo(ctx context.Context, accountId uint64, job string, jobGrade int32) error
}

type UIRetriever struct {
	ctx context.Context
	db  *sql.DB

	userCache *cache.Cache[int32, *UserInfo]
}

func NewUIRetriever(ctx context.Context, db *sql.DB) *UIRetriever {
	userCache := cache.NewContext(
		ctx,
		cache.AsLRU[int32, *UserInfo](lru.WithCapacity(300)),
		cache.WithJanitorInterval[int32, *UserInfo](360*time.Second),
	)

	return &UIRetriever{
		ctx: ctx,
		db:  db,

		userCache: userCache,
	}
}

func (ui *UIRetriever) GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error) {
	if ui.userCache.Contains(userId) {
		if userInfo, ok := ui.userCache.Get(userId); ok {
			return userInfo, nil
		}
	}

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("userinfo.userid"),
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Group,
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

	dest := &UserInfo{
		AccId: accountId,
	}
	if err := stmt.QueryContext(ui.ctx, ui.db, dest); err != nil {
		return nil, err
	}

	// Check if user is superuser
	if utils.InStringSlice(config.C.Game.SuperuserGroups, dest.Group) {
		dest.SuperUser = true
		if dest.OrigJob != "" {
			dest.Job = dest.OrigJob
			dest.JobGrade = dest.OrigJobGrade
		}
	}

	ui.userCache.Set(userId, dest, cache.WithExpiration(60*time.Second))

	return dest, nil
}

func (ui *UIRetriever) SetUserInfo(ctx context.Context, accountId uint64, job string, jobGrade int32) error {
	stmt := tFivenetAccounts.
		UPDATE(
			tFivenetAccounts.OverrideJob,
			tFivenetAccounts.OverrideJobGrade,
		).
		SET(
			tFivenetAccounts.OverrideJob.SET(jet.String(job)),
			tFivenetAccounts.OverrideJobGrade.SET(jet.Int32(jobGrade)),
		).
		WHERE(
			tFivenetAccounts.ID.EQ(jet.Uint64(accountId)),
		)

	if _, err := stmt.ExecContext(ctx, ui.db); err != nil {
		return err
	}

	return nil
}
