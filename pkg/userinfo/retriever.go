package userinfo

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/cache"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const cacheTTL = 20 * time.Second

var ErrAccountError = fmt.Errorf("failed to retrieve account data")

var tFivenetAccounts = table.FivenetAccounts

type UserInfoRetriever interface {
	GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*pbuserinfo.UserInfo, error)
	GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*pbuserinfo.UserInfo, error)
	SetUserInfo(ctx context.Context, accountId uint64, superuser bool, job *string, jobGrade *int32) error
}

// UIRetriever implements UserInfoRetriever and provides user info retrieval with caching.
type userAccountKey struct {
	UserID    int32
	AccountID uint64
}

type Retriever struct {
	logger *zap.Logger
	jsCons jetstream.ConsumeContext

	ctx    context.Context
	db     *sql.DB
	js     *events.JSWrapper
	notifi notifi.INotifi

	// userCache caches user info by userAccountKey (userId+accountId)
	userCache    *cache.LRUCache[userAccountKey, *pbuserinfo.UserInfo]
	userCacheTTL time.Duration

	superuserGroups []string
	superuserUsers  []string
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	JS     *events.JSWrapper
	Config *config.Config
	Notifi notifi.INotifi
}

// NewRetriever creates a new Retriever with LRU cache and configures lifecycle hooks.
func NewRetriever(p Params) UserInfoRetriever {
	ctxCancel, cancel := context.WithCancel(context.Background())

	userCache := cache.NewLRUCache[userAccountKey, *pbuserinfo.UserInfo](350)

	retriever := &Retriever{
		logger: p.Logger.Named("userinfo.retriever"),
		ctx:    ctxCancel,
		db:     p.DB,
		js:     p.JS,
		notifi: p.Notifi,

		userCache:    userCache,
		userCacheTTL: cacheTTL,

		superuserGroups: p.Config.Auth.SuperuserGroups,
		superuserUsers:  p.Config.Auth.SuperuserUsers,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerStreams(ctxCancel, p.JS); err != nil {
			return fmt.Errorf("failed to register user info stream. %w", err)
		}

		if err := retriever.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to register subscriptions for user info retriever. %w", err)
		}

		// Run cache janitor to clean up old entries every 41 seconds
		go userCache.StartJanitor(ctxCancel, 41*time.Second)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return retriever
}

func (r *Retriever) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	// Subscribe to the userinfo diffs stream
	consumer, err := r.js.CreateOrUpdateConsumer(ctxStartup, UserInfoStreamName, jetstream.ConsumerConfig{
		Durable:           instance.ID() + "_ui_retriever",
		AckPolicy:         jetstream.AckExplicitPolicy,
		FilterSubjects:    []string{UserInfoSubject},
		InactiveThreshold: 1 * time.Minute, // Close consumer if inactive for 1 minute
	})
	if err != nil {
		return fmt.Errorf("failed to create/update consumer for stream %q. %w", UserInfoStreamName, err)
	}

	if r.jsCons != nil {
		r.jsCons.Stop()
	}

	r.jsCons, err = consumer.Consume(r.handleMsg,
		r.js.ConsumeErrHandlerWithRestart(ctxCancel, r.logger, r.registerSubscriptions))
	if err != nil {
		return fmt.Errorf("failed to start message consumer for %s. %w", PollStreamName, err)
	}

	return nil
}

func (r *Retriever) handleMsg(m jetstream.Msg) {
	var evt pbuserinfo.UserInfoChanged
	if err := proto.Unmarshal(m.Data(), &evt); err == nil {
		key := userAccountKey{UserID: evt.UserId, AccountID: evt.AccountId}
		if userInfo, ok := r.userCache.Get(key); ok {
			// If we already have this user in cache, clone the current user info and update it
			clone := userInfo.Clone()
			clone.Job = evt.NewJob
			clone.JobGrade = evt.NewJobGrade
			r.userCache.Put(key, clone, r.userCacheTTL)
		}

		// Notify UI about the user info change
		r.logger.Debug("User info changed, notifying user", zap.Int32("userId", evt.UserId), zap.Uint64("accountId", evt.AccountId))

		r.notifi.SendUserEvent(r.ctx, evt.UserId, &notifications.UserEvent{
			Data: &notifications.UserEvent_UserInfoChanged{
				UserInfoChanged: &evt,
			},
		})
	}

	m.Ack()
}

// GetUserInfo retrieves user info for a given userId and accountId, using cache for performance.
func (r *Retriever) GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*pbuserinfo.UserInfo, error) {
	key := userAccountKey{UserID: userId, AccountID: accountId}
	if dest, ok := r.userCache.Get(key); ok {
		return dest, nil
	}

	dest := &pbuserinfo.UserInfo{}
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

	if err := stmt.QueryContext(ctx, r.db, dest); err != nil {
		return nil, errswrap.NewError(err, ErrAccountError)
	}

	// If account is not enabled, fail here
	if !dest.Enabled {
		return nil, ErrAccountError
	}

	// Set superuser status and override job/grade if applicable
	r.setSuperuserStatus(dest)

	r.userCache.Put(key, dest, r.userCacheTTL)

	return dest, nil
}

// GetUserInfoWithoutAccountId retrieves user info for a userId without account context.
func (r *Retriever) GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*pbuserinfo.UserInfo, error) {
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

	dest := &pbuserinfo.UserInfo{}
	if err := stmt.QueryContext(ctx, r.db, dest); err != nil {
		return nil, errswrap.NewError(err, ErrAccountError)
	}

	// Set superuser status (without license check)
	r.setSuperuserStatus(dest)

	return dest, nil
}

// setSuperuserStatus sets CanBeSuper and applies override job/grade if user is a superuser.
func (r *Retriever) setSuperuserStatus(dest *pbuserinfo.UserInfo) {
	// Defensive nil check for dest
	if dest == nil {
		return
	}

	// Check if user is superuser by group or license
	isSuperGroup := slices.Contains(r.superuserGroups, dest.Group)
	if isSuperGroup || slices.Contains(r.superuserUsers, dest.License) {
		dest.CanBeSuperuser = true
		// Only override if both are non-nil and OverrideJob is not empty
		if dest.OverrideJob != nil && *dest.OverrideJob != "" && dest.OverrideJobGrade != nil {
			dest.Job = *dest.OverrideJob
			dest.JobGrade = *dest.OverrideJobGrade
		}
	}
}

func (r *Retriever) SetUserInfo(ctx context.Context, accountId uint64, superuser bool, job *string, jobGrade *int32) error {
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

	if _, err := stmt.ExecContext(ctx, r.db); err != nil {
		return err
	}

	return nil
}
