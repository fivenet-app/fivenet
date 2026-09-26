package croner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
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

	logger := p.Logger.WithOptions(zap.IncreaseLevel(p.Cfg.Log.LevelOverrides.Get(config.LoggingComponentCron, p.Cfg.LogLevel))).
		Named("cron.registry")
	r := &Registry{
		logger: logger,
		js:     p.JS,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxStartup, r.js); err != nil {
			return err
		}

		storeKV, err := r.js.CreateOrUpdateKeyValue(ctxStartup, jetstream.KeyValueConfig{
			Bucket:      BucketName,
			Description: BucketName + " Store",
			History:     1,
			Storage:     jetstream.FileStorage,
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

func (r *Registry) ListCronjobs(_ context.Context) []*cron.Cronjob {
	cj := []*cron.Cronjob{}

	r.store.Range(func(_ string, entry *cron.Cronjob) bool {
		cj = append(cj, entry)

		return true
	})

	return cj
}

func (r *Registry) GetCronjob(ctx context.Context, name string) (*cron.Cronjob, error) {
	return r.store.Get(normalizeCronjobName(name))
}

func (r *Registry) RegisterCronjob(ctx context.Context, cronjob *cron.Cronjob) error {
	cronjob.Name = normalizeCronjobName(cronjob.GetName())
	if cronjob.GetName() == "" {
		return fmt.Errorf("cron job name is required or uses reserved name: %s", cronjob.GetName())
	}

	// Validate the cron schedule
	if !gronx.IsValid(cronjob.GetSchedule()) {
		return ErrInvalidCronSyntax
	}

	r.logger.Debug("registering cronjob", zap.String("name", cronjob.GetName()))

	if cronjob.GetTimeout() == nil {
		cronjob.Timeout = durationpb.New(DefaultCronjobTimeout)
	} else if cronjob.GetTimeout().AsDuration() <= 0 || cronjob.GetTimeout().AsDuration() > 30*time.Minute {
		// Ensure the timeout is positive and not bigger than 30 minutes.
		return fmt.Errorf("cron job %s has invalid timeout", cronjob.GetName())
	}

	if cronjob.GetState() == cron.CronjobState_CRONJOB_STATE_UNSPECIFIED {
		cronjob.State = cron.CronjobState_CRONJOB_STATE_WAITING
	}

	nextTime, err := gronx.NextTick(cronjob.GetSchedule(), false)
	if err != nil {
		return err
	}

	if cronjob.GetNextScheduleTime() == nil || cronjob.GetNextScheduleTime().AsTime() != nextTime {
		cronjob.NextScheduleTime = timestamp.New(nextTime)
	}

	if err := r.store.ComputeUpdate(
		ctx,
		cronjob.GetName(),
		func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
			if existing == nil {
				return cronjob, true, nil
			}

			wasRunning := existing.GetState() == cron.CronjobState_CRONJOB_STATE_RUNNING
			runtimeState := existing.GetState()
			runtimeRunID := existing.GetRunId()
			runtimeStartedTime := existing.GetStartedTime()
			runtimeLastAttemptTime := existing.GetLastAttemptTime()
			runtimeLastCompletedEvent := existing.GetLastCompletedEvent()

			existing.Merge(cronjob)
			if wasRunning {
				existing.State = runtimeState
				existing.SetRunId(runtimeRunID)
				existing.StartedTime = runtimeStartedTime
				existing.LastAttemptTime = runtimeLastAttemptTime
				existing.LastCompletedEvent = runtimeLastCompletedEvent
			}

			return existing, true, nil
		},
	); err != nil {
		return fmt.Errorf("failed to register cron job %s in store. %w", cronjob.GetName(), err)
	}

	return nil
}

func (r *Registry) UnregisterCronjob(ctx context.Context, name string) error {
	r.logger.Debug("unregistering cronjob", zap.String("name", name))
	if err := r.store.Delete(ctx, normalizeCronjobName(name)); err != nil {
		return fmt.Errorf("failed to unregister cron job %s from store. %w", name, err)
	}

	return nil
}
