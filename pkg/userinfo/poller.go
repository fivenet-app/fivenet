package userinfo

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pb "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type userSnapshot struct {
	Job       string
	JobGrade  int32
	Superuser bool
}

// Poller is responsible for polling user information via the Retriever.
type Poller struct {
	ctx context.Context

	logger *zap.Logger
	jsCons jetstream.ConsumeContext

	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	notifi   notifi.INotifi
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

	Logger   *zap.Logger
	DB       *sql.DB
	Enricher *mstlystcdata.Enricher
	Notifi   notifi.INotifi
	JS       *events.JSWrapper
}

func NewPoller(p PollerParams) *Poller {
	ctxCancel, cancel := context.WithCancel(context.Background())

	poller := &Poller{
		logger: p.Logger.Named("userinfo.poller"),

		ctx:      context.Background(),
		db:       p.DB,
		enricher: p.Enricher,
		notifi:   p.Notifi,
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

	userIds := []jet.Expression{}
	for _, req := range batch {
		userIds = append(userIds, jet.Int32(req.GetUserId()))
	}

	tUser := tables.User()
	tAccount := table.FivenetAccounts

	stmt := tUser.
		SELECT(
			tAccount.ID.AS("account_id"),
			tUser.ID.AS("user_id"),
			tUser.Job.AS("job"),
			tUser.JobGrade.AS("job_grade"),
			tAccount.Superuser.AS("superuser"),
			tUser.LastSeen.AS("last_seen"),
		).
		FROM(
			tUser.
				INNER_JOIN(
					tAccount,
					tAccount.License.LIKE(
						jet.RawString("SUBSTRING_INDEX(`users`.`identifier`, ':', -1)"),
					),
				),
		).
		WHERE(tUser.ID.IN(userIds...)).
		LIMIT(int64(len(userIds)))

	var dest []*struct {
		AccountId int64
		UserId    int32
		Job       string
		JobGrade  int32
		Superuser bool
		LastSeen  *timestamp.Timestamp
	}

	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return err
	}

	var errs error
	for _, row := range dest {
		if err := p.checkDiffAndPublish(ctx, row.AccountId, row.UserId, row.Job, row.JobGrade, row.Superuser, row.LastSeen); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to check diff. %w", err))
			continue
		}
	}

	return errs
}

func (p *Poller) checkDiffAndPublish(
	ctx context.Context,
	acct int64,
	uid int32,
	job string,
	grade int32,
	superuser bool,
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
		userMap[uid] = &userSnapshot{Job: job, JobGrade: grade}
		return nil
	}

	if old.Job != job || old.JobGrade != grade || old.Superuser != superuser {
		evt := &pb.UserInfoChanged{
			AccountId:   acct,
			UserId:      uid,
			OldJob:      old.Job,
			NewJob:      &job,
			OldJobGrade: old.JobGrade,
			NewJobGrade: &grade,
			Superuser:   &superuser,
			ChangedAt:   updatedAt,
		}
		p.enricher.EnrichJobInfo(evt)

		subj := fmt.Sprintf("userinfo.%d.changes", acct)
		if _, err := p.js.PublishAsyncProto(ctx, subj, evt); err != nil {
			p.logger.Error("failed to publish user info change event",
				zap.Int64("accountId", acct),
				zap.Int32("userId", uid),
				zap.String("job", job),
				zap.Int32("jobGrade", grade),
				zap.Bool("superuser", superuser),
				zap.Error(err),
			)
		}
		userMap[uid] = &userSnapshot{Job: job, JobGrade: grade, Superuser: superuser}
	}

	return nil
}
