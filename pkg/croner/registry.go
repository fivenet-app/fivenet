package croner

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const DefaultCronjobTimeout = 15 * time.Second

const BucketName = "cron"

var RegistryModule = fx.Module("cron_registry",
	fx.Provide(
		NewRegistry,
	),
)

var ErrInvalidCronSyntax = errors.New("invalid cron syntax")

type IRegistry interface {
	RegisterCronjob(ctx context.Context, job *cron.Cronjob) error
	UnregisterCronjob(ctx context.Context, name string) error
}

type CronRegister interface {
	RegisterCronjobs(ctx context.Context, c IRegistry) error
	RegisterCronjobHandlers(h *Handlers) error
}

type RegistryParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
	Cfg    *config.Config

	Jobs []CronRegister `group:"cronjobregister"`
}

type Registry struct {
	logger *zap.Logger

	ctx   context.Context
	js    *events.JSWrapper
	store *store.Store[cron.Cronjob, *cron.Cronjob]
	kv    jetstream.KeyValue
}

type RegistryResult struct {
	fx.Out

	Registry  *Registry
	IRegistry IRegistry
}

func NewRegistry(p RegistryParams) (RegistryResult, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.WithOptions(zap.IncreaseLevel(p.Cfg.LogLevelOverrides.Get(config.LoggingComponentCron, p.Cfg.LogLevel))).
		Named("cron.registry")
	r := &Registry{
		logger: p.Logger,
		ctx:    ctxCancel,
		js:     p.JS,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxStartup, r.js); err != nil {
			return err
		}

		storeKV, err := r.js.CreateOrUpdateKeyValue(ctxStartup, jetstream.KeyValueConfig{
			Bucket:      BucketName,
			Description: BucketName + " Store",
			History:     2,
			Storage:     jetstream.MemoryStorage,
		})
		if err != nil {
			return fmt.Errorf("failed to create kv (bucket %s) for cron store. %w", BucketName, err)
		}
		r.kv = storeKV

		st, err := store.New[cron.Cronjob, *cron.Cronjob](ctxStartup, logger, p.JS, BucketName,
			store.WithJetstreamKV[cron.Cronjob, *cron.Cronjob](storeKV),
		)
		if err != nil {
			return err
		}
		r.store = st

		if err := st.Start(ctxCancel, true); err != nil {
			return err
		}

		for _, reg := range p.Jobs {
			if err := reg.RegisterCronjobs(ctxStartup, r); err != nil {
				return err
			}
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return nil
	}))

	return RegistryResult{
		Registry:  r,
		IRegistry: r,
	}, nil
}

func (r *Registry) ListCronjobs(ctx context.Context) []*cron.Cronjob {
	cj := []*cron.Cronjob{}

	r.store.Range(func(_ string, entry *cron.Cronjob) bool {
		cj = append(cj, entry)

		return true
	})

	return cj
}

func (r *Registry) RegisterCronjob(ctx context.Context, job *cron.Cronjob) error {
	if job.GetName() == "" {
		return fmt.Errorf("cron job name is required or uses reserved name: %s", job.GetName())
	}

	// Validate the cron schedule
	if !gronx.IsValid(job.GetSchedule()) {
		return ErrInvalidCronSyntax
	}

	r.logger.Debug("registering cronjob", zap.String("name", job.GetName()))

	if job.GetTimeout() == nil {
		job.Timeout = durationpb.New(DefaultCronjobTimeout)
	} else if job.GetTimeout().AsDuration() < 0 || job.GetTimeout().AsDuration() > 30*time.Minute {
		// Ensure the timeout is not negative and not bigger than 30 minutes
		return fmt.Errorf("cron job %s has negative timeout", job.GetName())
	}

	if job.GetState() == cron.CronjobState_CRONJOB_STATE_UNSPECIFIED {
		job.State = cron.CronjobState_CRONJOB_STATE_WAITING
	}

	nextTime, err := gronx.NextTick(job.GetSchedule(), false)
	if err != nil {
		return err
	}

	if job.GetNextScheduleTime() == nil || job.GetNextScheduleTime().AsTime() != nextTime {
		job.NextScheduleTime = timestamp.New(nextTime)
	}

	if err := r.store.ComputeUpdate(ctx, strings.ToLower(job.GetName()), func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
		if existing == nil {
			return job, true, nil
		}

		existing.Merge(job)

		return existing, true, nil
	}); err != nil {
		return fmt.Errorf("failed to register cron job %s in store. %w", job.GetName(), err)
	}

	return nil
}

func (r *Registry) UnregisterCronjob(ctx context.Context, name string) error {
	r.logger.Debug("unregistering cronjob", zap.String("name", name))
	if err := r.store.Delete(ctx, name); err != nil {
		return fmt.Errorf("failed to unregister cron job %s from store. %w", name, err)
	}

	return nil
}
