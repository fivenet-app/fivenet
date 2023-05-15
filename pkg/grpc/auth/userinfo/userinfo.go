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
	tUsers = table.Users.AS("userinfo")
)

type UserInfo struct {
	UserId    int32
	Job       string
	JobGrade  int32
	Group     string
	SuperUser bool
}

type UserInfoRetriever interface {
	GetUserInfo(userId int32) (*UserInfo, error)
}

type UIRetriever struct {
	ctx context.Context
	db  *sql.DB

	userCache *cache.Cache[int32, *UserInfo]
}

func NewUIRetriever(ctx context.Context, db *sql.DB) *UIRetriever {
	userCache := cache.NewContext(
		ctx,
		cache.AsLRU[int32, *UserInfo](lru.WithCapacity(32)),
		cache.WithJanitorInterval[int32, *UserInfo](60*time.Second),
	)

	return &UIRetriever{
		ctx: ctx,
		db:  db,

		userCache: userCache,
	}
}

func (ui *UIRetriever) GetUserInfo(userId int32) (*UserInfo, error) {
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
		).
		FROM(tUsers).
		WHERE(
			tUsers.ID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	dest := &UserInfo{}
	if err := stmt.QueryContext(ui.ctx, ui.db, dest); err != nil {
		return nil, err
	}

	// Check if user is superuser
	if utils.InStringSlice(config.C.Game.SuperuserGroups, dest.Group) {
		dest.SuperUser = true
	}

	ui.userCache.Set(userId, dest)

	return dest, nil
}
