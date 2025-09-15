package userinfo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/cache"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const cacheTTL = 20 * time.Second

var ErrAccountError = errors.New("failed to retrieve account data")

type UserInfoRetriever interface {
	GetUserInfo(ctx context.Context, userId int32, accountId int64) (*pbuserinfo.UserInfo, error)
	GetUserInfoWithoutAccountId(ctx context.Context, userId int32) (*pbuserinfo.UserInfo, error)
	SetUserInfo(
		ctx context.Context,
		accountId int64,
		userId int32,
		superuser bool,
		job *string,
		jobGrade *int32,
	) error
	RefreshUserInfo(ctx context.Context, userId int32, accountId int64) error
}

// UIRetriever implements UserInfoRetriever and provides user info retrieval with caching.
type userAccountKey struct {
	UserID    int32
	AccountID int64
}

type Retriever struct {
	logger *zap.Logger
	jsCons jetstream.ConsumeContext

	ctx      context.Context
	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	notifi   notifi.INotifi

	// userCache caches user info by userAccountKey (userId+accountId)
	userCache    *cache.LRUCache[userAccountKey, *pbuserinfo.UserInfo]
	userCacheTTL time.Duration

	superuserGroups []string
	superuserUsers  []string
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	DB       *sql.DB
	JS       *events.JSWrapper
	Enricher *mstlystcdata.Enricher
	Config   *config.Config
	Notifi   notifi.INotifi
}

// NewRetriever creates a new Retriever with LRU cache and configures lifecycle hooks.
func NewRetriever(p Params) UserInfoRetriever {
	ctxCancel, cancel := context.WithCancel(context.Background())

	userCache := cache.NewLRUCache[userAccountKey, *pbuserinfo.UserInfo](350)

	retriever := &Retriever{
		logger:   p.Logger.Named("userinfo.retriever"),
		ctx:      ctxCancel,
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
		notifi:   p.Notifi,

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

func (r *Retriever) registerSubscriptions(
	ctxStartup context.Context,
	ctxCancel context.Context,
) error {
	// Subscribe to the userinfo diffs stream
	consumer, err := r.js.CreateOrUpdateConsumer(
		ctxStartup,
		UserInfoStreamName,
		jetstream.ConsumerConfig{
			Durable:           instance.ID() + "_ui_retriever",
			AckPolicy:         jetstream.AckExplicitPolicy,
			FilterSubjects:    []string{UserInfoSubject},
			InactiveThreshold: 1 * time.Minute, // Close consumer if inactive for 1 minute
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create/update consumer for stream %q. %w",
			UserInfoStreamName,
			err,
		)
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
	if err := m.Ack(); err != nil {
		r.logger.Error("failed to ack message", zap.Error(err), zap.String("subject", m.Subject()))
	}

	var evt pbuserinfo.UserInfoChanged
	if err := proto.Unmarshal(m.Data(), &evt); err != nil {
		r.logger.Error("failed to unmarshal user info changed event",
			zap.Error(err),
			zap.String("subject", m.Subject()),
		)
		return
	}

	key := userAccountKey{UserID: evt.GetUserId(), AccountID: evt.GetAccountId()}
	if userInfo, ok := r.userCache.Get(key); ok {
		// If we already have this user in cache, clone the current user info and update it
		clone := userInfo.Clone()
		clone.Job = evt.GetNewJob()
		clone.JobGrade = evt.GetNewJobGrade()
		r.userCache.Put(key, clone, r.userCacheTTL)
	}

	// Notify UI about the user info change
	r.logger.Debug(
		"User info changed, notifying user",
		zap.Int32("userId", evt.GetUserId()),
		zap.Int64("accountId", evt.GetAccountId()),
	)

	r.notifi.SendUserEvent(r.ctx, evt.GetUserId(), &notifications.UserEvent{
		Data: &notifications.UserEvent_UserInfoChanged{
			UserInfoChanged: &evt,
		},
	})
}

// GetUserInfo retrieves user info for a given userId and accountId, using cache for performance.
func (r *Retriever) GetUserInfo(
	ctx context.Context,
	userId int32,
	accountId int64,
) (*pbuserinfo.UserInfo, error) {
	key := userAccountKey{UserID: userId, AccountID: accountId}
	if dest, ok := r.userCache.Get(key); ok {
		if !dest.GetEnabled() {
			return nil, ErrAccountError
		}

		return dest, nil
	}

	dest, err := r.getUserInfoFromDB(ctx, userId, accountId)
	if err != nil {
		return nil, err
	}

	// If account is not enabled, fail here
	if !dest.GetEnabled() {
		return nil, ErrAccountError
	}

	// Set superuser status and override job/grade if applicable
	r.checkAndSetSuperuser(dest)

	r.userCache.Put(key, dest, r.userCacheTTL)

	return dest, nil
}

func (r *Retriever) getUserInfoFromDB(
	ctx context.Context,
	userId int32,
	accountId int64,
) (*pbuserinfo.UserInfo, error) {
	tUsers := tables.User().AS("user_info")
	tAccount := table.FivenetAccounts

	stmt := tUsers.
		SELECT(
			tAccount.ID.AS("user_info.account_id"),
			tAccount.Enabled.AS("user_info.enabled"),
			tAccount.License.AS("user_info.license"),
			tAccount.OverrideJob.AS("user_info.override_job"),
			tAccount.OverrideJobGrade.AS("user_info.override_job_grade"),
			tAccount.Superuser.AS("user_info.superuser"),
			tAccount.LastChar.AS("user_info.last_char"),
			tUsers.ID.AS("user_info.userid"),
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Group,
		).
		FROM(
			tAccount,
			tUsers,
		).
		WHERE(mysql.AND(
			tAccount.ID.EQ(mysql.Int64(accountId)),
			tUsers.ID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	dest := &pbuserinfo.UserInfo{}
	if err := stmt.QueryContext(ctx, r.db, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

// GetUserInfoWithoutAccountId retrieves user info for a userId without account context.
func (r *Retriever) GetUserInfoWithoutAccountId(
	ctx context.Context,
	userId int32,
) (*pbuserinfo.UserInfo, error) {
	tUsers := tables.User().AS("user_info")

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("user_info.userid"),
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Group,
		).
		FROM(tUsers).
		WHERE(mysql.AND(
			tUsers.ID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	dest := &pbuserinfo.UserInfo{}
	if err := stmt.QueryContext(ctx, r.db, dest); err != nil {
		return nil, errswrap.NewError(err, ErrAccountError)
	}

	// Set superuser status if applicable
	r.checkAndSetSuperuser(dest)

	return dest, nil
}

// checkAndSetSuperuser check if user is superuser by group or license.
func (r *Retriever) checkAndSetSuperuser(dest *pbuserinfo.UserInfo) {
	if slices.Contains(r.superuserGroups, dest.GetGroup()) ||
		slices.Contains(r.superuserUsers, dest.GetLicense()) {
		dest.CanBeSuperuser = true

		// Only override if both are non-nil and OverrideJob is not empty
		if dest.OverrideJob != nil && dest.GetOverrideJob() != "" && dest.OverrideJobGrade != nil {
			dest.Job = dest.GetOverrideJob()
			dest.JobGrade = dest.GetOverrideJobGrade()
		}
	}
}

// SetUserInfo updates user info for a given accountId, setting superuser status and job/grade in the account table.
func (r *Retriever) SetUserInfo(
	ctx context.Context,
	accountId int64,
	userId int32,
	superuser bool,
	overrideJob *string,
	overrideJobGrade *int32,
) error {
	tAccount := table.FivenetAccounts

	stmt := tAccount.
		UPDATE(
			tAccount.Superuser,
			tAccount.OverrideJob,
			tAccount.OverrideJobGrade,
		).
		SET(
			superuser,
			overrideJob,
			overrideJobGrade,
		).
		WHERE(tAccount.ID.EQ(mysql.Int64(accountId))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, r.db); err != nil {
		return err
	}

	// If userId is not provided, we only update the account info
	if userId <= 0 {
		return nil
	}

	// If userId is provided, we need to update the user info in cache by publishing the update
	key := userAccountKey{UserID: userId, AccountID: accountId}
	dest, ok := r.userCache.Get(key)
	if !ok {
		return nil
	}

	ui := dest.Clone()
	ui.Superuser = superuser
	if ui.Superuser {
		ui.CanBeSuperuser = true
	}
	if overrideJob != nil {
		ui.Job = *overrideJob
	}
	if overrideJobGrade != nil {
		ui.JobGrade = *overrideJobGrade
	}

	r.userCache.Put(key, ui, r.userCacheTTL)

	evt := &pbuserinfo.UserInfoChanged{
		AccountId:      ui.AccountId,
		UserId:         ui.UserId,
		OldJob:         dest.Job,
		NewJob:         overrideJob,
		OldJobGrade:    dest.JobGrade,
		NewJobGrade:    overrideJobGrade,
		CanBeSuperuser: &superuser,
		Superuser:      &superuser,
		ChangedAt:      timestamp.Now(),
	}
	r.enricher.EnrichJobInfo(evt)

	// Publish the user info change to the NATS JetStream
	subj := fmt.Sprintf("userinfo.%d.changes", ui.AccountId)
	if _, err := r.js.PublishAsyncProto(ctx, subj, evt); err != nil {
		return fmt.Errorf(
			"failed to publish user info change for userId %d, accountId %d. %w",
			userId,
			accountId,
			err,
		)
	}

	return nil
}

func (r *Retriever) RefreshUserInfo(ctx context.Context, userId int32, accountId int64) error {
	dest, err := r.getUserInfoFromDB(ctx, userId, accountId)
	if err != nil {
		return err
	}

	r.checkAndSetSuperuser(dest)

	key := userAccountKey{UserID: userId, AccountID: accountId}
	r.userCache.Put(key, dest, r.userCacheTTL)

	return nil
}
