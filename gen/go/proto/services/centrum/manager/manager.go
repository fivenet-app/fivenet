package manager

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adrg/strutil/metrics"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/tracker/postals"
	"github.com/nats-io/nats.go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("centrum_manager", fx.Provide(
	New,
))

type Manager struct {
	ctx    context.Context
	logger *zap.Logger

	tracer   trace.Tracer
	db       *sql.DB
	js       nats.JetStreamContext
	enricher *mstlystcdata.Enricher
	tracker  *tracker.Tracker
	postals  *postals.Postals

	trackedJobs []string
	publicJobs  []string

	*state.State
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	TP       *tracesdk.TracerProvider
	DB       *sql.DB
	JS       nats.JetStreamContext
	Enricher *mstlystcdata.Enricher
	Postals  *postals.Postals
	Tracker  *tracker.Tracker
	Config   *config.Config

	State *state.State
}

func New(p Params) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	jw := metrics.NewJaroWinkler()
	jw.CaseSensitive = false

	s := &Manager{
		ctx:    ctx,
		logger: p.Logger.Named("centrum_state"),

		tracer:   p.TP.Tracer("centrum-manager"),
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
		postals:  p.Postals,
		tracker:  p.Tracker,

		trackedJobs: p.Config.Game.Livemap.Jobs,
		publicJobs:  p.Config.Game.PublicJobs,

		State: p.State,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := eventscentrum.RegisterStreams(ctx, s.js); err != nil {
			return fmt.Errorf("failed to register events: %w", err)
		}

		if err := s.loadData(); err != nil {
			return err
		}

		if err := s.registerSubscriptions(); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s
}

func (s *Manager) loadData() error {
	ctx, span := s.tracer.Start(s.ctx, "centrum-loaddata")
	defer span.End()

	s.logger.Debug("loading settings")
	if err := s.LoadSettings(ctx, ""); err != nil {
		return fmt.Errorf("failed to load centrum settings: %w", err)
	}

	s.logger.Debug("loading disponents")
	if err := s.LoadDisponents(ctx, ""); err != nil {
		return fmt.Errorf("failed to load centrum disponents: %w", err)
	}

	s.logger.Debug("loading units")
	if err := s.LoadUnits(ctx, 0); err != nil {
		return fmt.Errorf("failed to load centrum units: %w", err)
	}

	s.logger.Debug("loading dispatches")
	if err := s.LoadDispatches(ctx, 0); err != nil {
		return fmt.Errorf("failed to load centrum dispatches: %w", err)
	}

	s.logger.Debug("loaded all centrum data")
	return nil
}
