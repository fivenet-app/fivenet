package croner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/locks"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ErrInvalidCronSyntax = errors.New("invalid cron syntax")

var Module = fx.Module("cron",
	fx.Provide(
		New,
	),
)

type ICron interface {
	RegisterCronjob(ctx context.Context, job *cron.Cronjob) error
	UnregisterCronjob(ctx context.Context, name string) error
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	Cfg       *config.Config
	JS        *events.JSWrapper
	State     *State
	Scheduler *Scheduler
}

type Cron struct {
	name string

	ctx       context.Context
	logger    *zap.Logger
	js        *events.JSWrapper
	ownerKv   jetstream.KeyValue
	ownerLock *locks.Locks
	state     *State
	scheduler *Scheduler
}

func New(p Params) (ICron, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	_, port, err := net.SplitHostPort(p.Cfg.HTTP.AdminListen)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	cr := &Cron{
		name: fmt.Sprintf("%s-%s", hostname, port),

		ctx:       ctx,
		logger:    p.Logger.Named("cron"),
		js:        p.JS,
		state:     p.State,
		scheduler: p.Scheduler,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := registerCronStreams(ctx, cr.js); err != nil {
			return err
		}

		ownerKv, err := p.JS.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:  "cron_locks",
			Storage: jetstream.MemoryStorage,
			History: 5,
		})
		if err != nil {
			return err
		}
		cr.ownerKv = ownerKv

		ownerLock, err := locks.New(p.Logger, ownerKv, ownerKv.Bucket(), 20*time.Second)
		if err != nil {
			return err
		}
		cr.ownerLock = ownerLock

		go cr.lockLoop()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return nil
	}))

	return cr, nil
}

func (c *Cron) RegisterCronjob(ctx context.Context, job *cron.Cronjob) error {
	// Validate the cron schedule
	if !gronx.IsValid(job.Schedule) {
		return ErrInvalidCronSyntax
	}

	c.logger.Debug("registering cronjob", zap.String("name", job.Name))
	cj, err := c.state.store.GetOrLoad(ctx, job.Name)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return fmt.Errorf("failed to load existing cron job %s. %w", job.Name, err)
	}
	if cj != nil {
		// Merge with existing cronjob data
		cj.Merge(job)
	} else {
		cj = job
	}

	if cj.State == cron.CronjobState_CRONJOB_STATE_UNSPECIFIED {
		cj.State = cron.CronjobState_CRONJOB_STATE_PENDING
	}

	nextTime, err := gronx.NextTick(cj.Schedule, false)
	if err != nil {
		return err
	}

	if cj.NextScheduleTime == nil || cj.NextScheduleTime.AsTime() != nextTime {
		cj.NextScheduleTime = timestamp.New(nextTime)
	}

	if err := c.state.store.ComputeUpdate(ctx, strings.ToLower(job.Name), true, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
		if existing == nil {
			return cj, true, nil
		}

		existing.Merge(cj)

		return existing, true, nil
	}); err != nil {
		return fmt.Errorf("failed to register cron job %s in store. %w", job.Name, err)
	}

	return nil
}

func (c *Cron) UnregisterCronjob(ctx context.Context, name string) error {
	c.logger.Debug("unregistering cronjob", zap.String("name", name))
	if err := c.state.store.Delete(ctx, name); err != nil {
		return fmt.Errorf("failed to unregister cron job %s from store. %w", name, err)
	}

	return nil
}
