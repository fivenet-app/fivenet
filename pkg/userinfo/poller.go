package userinfo

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var PollerModule = fx.Module(
	"userinfo.poller",
	fx.Provide(
		NewPoller,
	),
)

type userSnapshot struct {
	Job      string
	JobGrade int32
	Groups   []string
}

// Poller is responsible for polling user information via the Retriever.
// It is used to detect changes in user information and publish events accordingly.
// The events are mainly used by any streams.
type Poller struct {
	ctx context.Context //nolint:containedctx // Used as the service-lifecycle context for async KV operations.

	logger *zap.Logger
	jsCons jetstream.ConsumeContext

	db       *sql.DB
	cfg      *config.Config
	appCfg   appconfig.IConfig
	js       *events.JSWrapper
	enricher mstlystcdata.IEnricher
	kv       jetstream.KeyValue

	pendingMu sync.Mutex
	pending   map[string]*pb.PollReq
	snapMu    sync.Mutex
	lastSeen  map[int64]map[int32]*userSnapshot
	interval  time.Duration
	ttl       time.Duration
}

type PollerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	DB        *sql.DB
	Cfg       *config.Config
	Enricher  mstlystcdata.IEnricher
	JS        *events.JSWrapper
	AppConfig appconfig.IConfig
}

func NewPoller(p PollerParams) *Poller {
	ctxCancel, cancel := context.WithCancel(context.Background())

	poller := &Poller{
		logger: p.Logger.Named("userinfo.poller"),

		ctx:      ctxCancel,
		db:       p.DB,
		cfg:      p.Cfg,
		appCfg:   p.AppConfig,
		enricher: p.Enricher,
		js:       p.JS,

		pending:  make(map[string]*pb.PollReq),
		lastSeen: make(map[int64]map[int32]*userSnapshot),
		interval: 20 * time.Second,
		ttl:      20 * time.Second,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerStreams(ctxStartup, p.JS); err != nil {
			return fmt.Errorf("failed to register user info streams. %w", err)
		}

		kv, err := poller.js.CreateOrUpdateKeyValue(ctxStartup, jetstream.KeyValueConfig{
			Bucket:         KVBucketName,
			Description:    "User info polling ttl store",
			History:        2,
			LimitMarkerTTL: 2 * poller.interval,
		})
		if err != nil {
			return fmt.Errorf("failed to create user info kv store. %w", err)
		}
		poller.kv = kv

		if err := poller.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to register subscriptions for user info poller. %w", err)
		}

		go poller.start(ctxCancel)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return poller
}

func (p *Poller) registerSubscriptions(
	ctxStartup context.Context,
	ctxCancel context.Context,
) error {
	// Subscribe to poll requests
	consumer, err := p.js.CreateOrUpdateConsumer(
		ctxStartup,
		PollStreamName,
		jetstream.ConsumerConfig{
			Durable:           instance.ID() + "_ui_poller",
			AckPolicy:         jetstream.AckExplicitPolicy,
			FilterSubjects:    []string{PollSubject},
			InactiveThreshold: 1 * time.Minute, // Close consumer if inactive for 1 minute
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create/update consumer for %s. %w", PollStreamName, err)
	}

	if p.jsCons != nil {
		p.jsCons.Stop()
	}

	p.jsCons, err = consumer.Consume(p.handleMsg,
		p.js.ConsumeErrHandlerWithRestart(ctxCancel, p.logger, p.registerSubscriptions))
	if err != nil {
		return fmt.Errorf("failed to start message consumer for %s. %w", PollStreamName, err)
	}

	return nil
}

func (p *Poller) handleMsg(m jetstream.Msg) {
	var req pb.PollReq
	if err := protoutils.UnmarshalPartialJSON(m.Data(), &req); err == nil {
		key := fmt.Sprintf("%d:%d", req.GetAccountId(), req.GetUserId())

		// Try Create with TTL. ErrKeyExists means we skip.
		if _, err := p.kv.Create(p.ctx, key, []byte("1"), jetstream.KeyTTL(p.ttl)); err == nil {
			p.pendingMu.Lock()
			p.pending[key] = &req
			p.pendingMu.Unlock()
		}
	}
	m.Ack()
}

func (p *Poller) start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := p.doBatch(ctx); err != nil {
				// Log the error, but continue polling
				p.logger.Error("failed to process batch", zap.Error(err))
			}
		}
	}
}

func (p *Poller) doBatch(ctx context.Context) error {
	p.pendingMu.Lock()
	batch := p.pending
	p.pending = make(map[string]*pb.PollReq)
	p.pendingMu.Unlock()

	if len(batch) == 0 {
		return nil
	}

	userIds := []mysql.Expression{}
	for _, req := range batch {
		userIds = append(userIds, mysql.Int32(req.GetUserId()))
	}

	tUser := table.FivenetUser
	tUserAccounts := table.FivenetUserAccounts
	tAccount := table.FivenetAccounts

	stmt := tUser.
		SELECT(
			tAccount.ID.AS("account_id"),
			tAccount.License.AS("license"),
			tUser.ID.AS("user_id"),
			tUser.Job.AS("job"),
			tUser.JobGrade.AS("job_grade"),
			tAccount.Groups.AS("groups"),
			tUser.UpdatedAt.AS("updated_at"),
		).
		FROM(
			tUser.
				LEFT_JOIN(
					tUserAccounts,
					tUserAccounts.UserID.EQ(tUser.ID),
				).
				INNER_JOIN(
					tAccount,
					mysql.OR(
						tAccount.ID.EQ(tUserAccounts.AccountID),
						tAccount.License.EQ(tUser.License),
					),
				),
		).
		WHERE(tUser.ID.IN(userIds...)).
		LIMIT(int64(len(userIds)))

	var dest []*struct {
		AccountId int64
		License   string
		UserId    int32
		Job       string
		JobGrade  int32
		Groups    *accounts.AccountGroups
		UpdatedAt *timestamp.Timestamp
	}

	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return err
	}

	var errs error
	for _, row := range dest {
		if err := p.checkDiffAndPublish(
			ctx,
			row.AccountId,
			row.License,
			row.UserId,
			row.Job,
			row.JobGrade,
			row.Groups,
			row.UpdatedAt,
		); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to check diff. %w", err))
			continue
		}
	}

	return errs
}

func (p *Poller) checkDiffAndPublish(
	ctx context.Context,
	acct int64,
	license string,
	uid int32,
	job string,
	grade int32,
	groups *accounts.AccountGroups,
	updatedAt *timestamp.Timestamp,
) error {
	p.snapMu.Lock()
	defer p.snapMu.Unlock()

	userMap, ok := p.lastSeen[acct]
	if !ok {
		userMap = make(map[int32]*userSnapshot)
		p.lastSeen[acct] = userMap
	}
	old, exists := userMap[uid]
	if !exists {
		// First-seen: record snapshot, no event
		userMap[uid] = &userSnapshot{
			Job:      job,
			JobGrade: grade,
			Groups:   groupsSlice(groups),
		}
		return nil
	}

	newGroups := groupsSlice(groups)
	jobChanged := old.Job != job || old.JobGrade != grade
	groupsChanged := !slices.Equal(old.Groups, newGroups)

	if jobChanged {
		evt := BuildUserInfoChangedEvent(acct, uid, updatedAt, job, grade, p.enricher)

		if _, err := p.js.PublishAsyncProto(
			ctx,
			fmt.Sprintf("userinfo.%d.changes", acct),
			evt,
		); err != nil {
			p.logger.Error("failed to publish user info change event",
				zap.Int64("accountId", acct),
				zap.Int32("userId", uid),
				zap.String("job", job),
				zap.Int32("jobGrade", grade),
				zap.Strings("groups", newGroups),
				zap.Error(err),
			)
		}
	}

	if groupsChanged {
		var jobAdminGroups, jobAdminUsers, configAdminGroups, configAdminUsers []string
		if p.cfg != nil {
			jobAdminGroups = p.cfg.Auth.GetJobAdminGroups()
			jobAdminUsers = p.cfg.Auth.GetJobAdminUsers()
			configAdminGroups = p.cfg.Auth.GetConfigAdminGroups()
			configAdminUsers = p.cfg.Auth.GetConfigAdminUsers()
		}
		jobAdminGroups, jobAdminUsers = EffectiveJobAdminLists(
			jobAdminGroups,
			jobAdminUsers,
			configAdminGroups,
			configAdminUsers,
			p.appCfg,
		)
		evt := BuildAccountGroupsChangedEvent(
			acct,
			updatedAt,
			groups,
			CanBeSuperuser(groups, license, jobAdminGroups, jobAdminUsers),
		)

		if _, err := p.js.PublishAsyncProto(
			ctx,
			fmt.Sprintf("userinfo.%d.groups", acct),
			evt,
		); err != nil {
			p.logger.Error("failed to publish user groups change event",
				zap.Int64("accountId", acct),
				zap.Int32("userId", uid),
				zap.Strings("groups", newGroups),
				zap.Error(err),
			)
		}
	}

	if jobChanged || groupsChanged {
		userMap[uid] = &userSnapshot{
			Job:      job,
			JobGrade: grade,
			Groups:   newGroups,
		}
	}

	return nil
}

func groupsSlice(groups *accounts.AccountGroups) []string {
	if groups == nil {
		return nil
	}
	return slices.Clone(groups.GetGroups())
}
