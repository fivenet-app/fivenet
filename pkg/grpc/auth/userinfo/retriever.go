package userinfo

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
)

const cacheTTL = 20 * time.Second

var ErrAccountError = fmt.Errorf("failed to retrieve account data")

var tFivenetAccounts = table.FivenetAccounts

type UserInfoRetriever interface {
	GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error)
	GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*UserInfo, error)
	SetUserInfo(ctx context.Context, accountId uint64, superuser bool, job *string, jobGrade *int32) error
}

// UIRetriever implements UserInfoRetriever and provides user info retrieval with caching.
type userAccountKey struct {
	UserID    int32
	AccountID uint64
}

type UIRetriever struct {
	db *sql.DB

	// userCache caches user info by userAccountKey (userId+accountId)
	userCache    *cache.Cache[userAccountKey, *UserInfo]
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

// NewUIRetriever creates a new UIRetriever with LRU cache and configures lifecycle hooks.
func NewUIRetriever(p Params) UserInfoRetriever {
	ctx, cancel := context.WithCancel(context.Background())

	userCache := cache.NewContext(
		ctx,
		cache.AsLRU[userAccountKey, *UserInfo](lru.WithCapacity(350)),
		cache.WithJanitorInterval[userAccountKey, *UserInfo](cacheTTL),
	)

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		return nil
	}))

	return &UIRetriever{
		db:              p.DB,
		userCache:       userCache,
		userCacheTTL:    cacheTTL,
		superuserGroups: p.Config.Auth.SuperuserGroups,
		superuserUsers:  p.Config.Auth.SuperuserUsers,
	}
}

// GetUserInfo retrieves user info for a given userId and accountId, using cache for performance.
func (ui *UIRetriever) GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error) {
	key := userAccountKey{UserID: userId, AccountID: accountId}
	if dest, ok := ui.userCache.Get(key); ok {
		return dest, nil
	}

	dest := &UserInfo{}
	tUsers := tables.User().AS("user_info")

	stmt := tUsers.
		SELECT(
			tFivenetAccounts.ID.AS("user_info.account_id"),
			tFivenetAccounts.Enabled.AS("user_info.enabled"),
			tFivenetAccounts.License.AS("user_info.license"),
			tFivenetAccounts.OverrideJob.AS("user_info.override_job"),
			tFivenetAccounts.OverrideJobGrade.AS("user_info.override_job_grade"),
			tFivenetAccounts.Superuser.AS("user_info.superuser"),
			tFivenetAccounts.LastChar.AS("user_info.last_char"),
			tUsers.ID.AS("user_info.userid"),
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

	// Set superuser status and override job/grade if applicable
	ui.setSuperuserStatus(dest)

	ui.userCache.Set(key, dest, cache.WithExpiration(ui.userCacheTTL))

	return dest, nil
}

// GetUserInfoWithoutAccountId retrieves user info for a userId without account context.
func (ui *UIRetriever) GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*UserInfo, error) {
	tUsers := tables.User().AS("user_info")

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("user_info.userid"),
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Group,
		).
		FROM(tUsers).
		WHERE(jet.AND(
			tUsers.ID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	dest := &UserInfo{}
	if err := stmt.QueryContext(ctx, ui.db, dest); err != nil {
		return nil, errswrap.NewError(err, ErrAccountError)
	}

	// Set superuser status (without license check)
	ui.setSuperuserStatus(dest)

	return dest, nil
}

// setSuperuserStatus sets CanBeSuper and applies override job/grade if user is a superuser.
func (ui *UIRetriever) setSuperuserStatus(dest *UserInfo) {
	// Defensive nil check for dest
	if dest == nil {
		return
	}
	// Check if user is superuser by group or license (license may be nil)
	isSuperGroup := slices.Contains(ui.superuserGroups, dest.Group)
	if isSuperGroup || slices.Contains(ui.superuserUsers, dest.License) {
		dest.CanBeSuperuser = true
		// Only override if both are non-nil and OverrideJob is not empty
		if dest.OverrideJob != nil && *dest.OverrideJob != "" && dest.OverrideJobGrade != nil {
			dest.Job = *dest.OverrideJob
			dest.JobGrade = *dest.OverrideJobGrade
		}
	}
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
