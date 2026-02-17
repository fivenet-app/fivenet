package userinfo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ErrAccountError = errors.New("failed to retrieve account data")

var RetrieverModule = fx.Module(
	"userinfo.retriever",
	fx.Provide(
		NewRetriever,
	),
)

type UserInfoRetriever interface {
	GetUserInfo(ctx context.Context, userId int32) (*pbuserinfo.UserInfo, error)
	GetUserInfoFromClaims(
		ctx context.Context,
		userClaims *authclaims.UserInfoClaims,
		accClaims *authclaims.AccountInfoClaims,
	) (*pbuserinfo.UserInfo, error)
	RefreshUserInfo(ctx context.Context, userId int32) error
}

// UIRetriever implements UserInfoRetriever and provides user info retrieval with caching.
type Retriever struct {
	UserInfoRetriever

	logger *zap.Logger
	jsCons jetstream.ConsumeContext

	ctx      context.Context
	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	notifi   notifi.INotifi

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

	retriever := &Retriever{
		logger:   p.Logger.Named("userinfo.retriever"),
		ctx:      ctxCancel,
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
		notifi:   p.Notifi,

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
	if err := protoutils.UnmarshalPartialJSON(m.Data(), &evt); err != nil {
		r.logger.Error("failed to unmarshal user info changed event",
			zap.Error(err),
			zap.String("subject", m.Subject()),
		)
		return
	}

	// Notify UI about the user info change
	r.logger.Debug(
		"User info changed, notifying user",
		zap.Int32("userId", evt.GetUserId()),
		zap.Int64("accountId", evt.GetAccountId()),
	)

	r.notifi.SendUserEvent(r.ctx, evt.GetUserId(), &notificationsevents.UserEvent{
		Data: &notificationsevents.UserEvent_UserInfoChanged{
			UserInfoChanged: &evt,
		},
	})
}

// GetUserInfo retrieves user info for a given userId and accountId, using cache for performance.
func (r *Retriever) GetUserInfo(
	ctx context.Context,
	userId int32,
) (*pbuserinfo.UserInfo, error) {
	dbUser, err := r.getUserInfoFromDB(ctx, userId)
	if err != nil {
		return nil, errswrap.NewError(
			fmt.Errorf("failed to get user info from db. %w", err),
			ErrAccountError,
		)
	}

	// If account is not enabled, fail here
	if !dbUser.GetEnabled() {
		return nil, ErrAccountError
	}

	return dbUser, nil
}

func (r *Retriever) getUserInfoFromDB(
	ctx context.Context,
	userId int32,
) (*pbuserinfo.UserInfo, error) {
	tAccount := table.FivenetAccounts
	tUsers := table.FivenetUser.AS("user_info")

	stmt := tUsers.
		SELECT(
			tAccount.ID.AS("user_info.account_id"),
			tAccount.Enabled.AS("user_info.enabled"),
			tAccount.License.AS("user_info.license"),
			tAccount.Groups.AS("user_info.groups"),
			tAccount.LastChar.AS("user_info.last_char"),
			tUsers.ID.AS("user_info.userid"),
			tUsers.Job,
			tUsers.JobGrade,
		).
		FROM(
			tAccount.
				INNER_JOIN(tUsers,
					tAccount.ID.EQ(tUsers.AccountID),
				),
		).
		WHERE(tUsers.ID.EQ(mysql.Int32(userId))).
		LIMIT(1)

	user := &pbuserinfo.UserInfo{}
	if err := stmt.QueryContext(ctx, r.db, user); err != nil {
		return nil, err
	}

	return user, nil
}

// checkAndSetSuperuser check if user is superuser by group or license.
func (r *Retriever) checkAndSetSuperuser(userInfo *pbuserinfo.UserInfo) {
	if userInfo.GetGroups().ContainsAnyGroup(r.superuserGroups) ||
		slices.Contains(r.superuserUsers, userInfo.GetLicense()) {
		userInfo.CanBeSuperuser = true
	} else {
		userInfo.CanBeSuperuser = false
		userInfo.Superuser = false
	}
}

// GetUserInfoFromClaims retrieves user info based on the claims.
func (r *Retriever) GetUserInfoFromClaims(
	ctx context.Context,
	userClaims *authclaims.UserInfoClaims,
	accClaims *authclaims.AccountInfoClaims,
) (*pbuserinfo.UserInfo, error) {
	userInfo, err := r.GetUserInfo(ctx, userClaims.UserID)
	if err != nil {
		return nil, err
	}

	// Set superuser status and override job/grade if applicable
	r.checkAndSetSuperuser(userInfo)

	if userInfo.CanBeSuperuser && userClaims.Superuser != nil && *userClaims.Superuser {
		userInfo.Superuser = true
	}

	// If the claims contain job info, override the user info with it
	if userClaims.OriginalJob != nil {
		if (userInfo.Job == userClaims.OriginalJob.Job ||
			userInfo.JobGrade == userClaims.OriginalJob.JobGrade) && !userInfo.Superuser {
			userInfo.Job = userClaims.OriginalJob.Job
			userInfo.JobGrade = userClaims.OriginalJob.JobGrade
		}
	}

	return userInfo, nil
}

func (r *Retriever) RefreshUserInfo(ctx context.Context, userId int32) error {
	dest, err := r.getUserInfoFromDB(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get user info from db. %w", err)
	}

	r.checkAndSetSuperuser(dest)

	return nil
}
