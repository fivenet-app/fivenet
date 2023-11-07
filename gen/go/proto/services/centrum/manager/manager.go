package manager

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/adrg/strutil/metrics"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/tracker/postals"
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
	wg     sync.WaitGroup

	tracer   trace.Tracer
	db       *sql.DB
	events   *events.Eventus
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
	Events   *events.Eventus
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
		wg:     sync.WaitGroup{},

		tracer:   p.TP.Tracer("centrum-state"),
		db:       p.DB,
		events:   p.Events,
		enricher: p.Enricher,
		postals:  p.Postals,
		tracker:  p.Tracker,

		trackedJobs: p.Config.Game.Livemap.Jobs,
		publicJobs:  p.Config.Game.PublicJobs,

		State: p.State,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := eventscentrum.RegisterEvents(ctx, s.events); err != nil {
			return fmt.Errorf("failed to register events: %w", err)
		}

		if err := s.loadData(); err != nil {
			return err
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.watchEvents()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.watchUserChanges()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.housekeeper()
		}()

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

	if err := s.LoadSettings(ctx, ""); err != nil {
		return fmt.Errorf("failed to load centrum settings: %w", err)
	}

	if err := s.LoadDisponents(ctx, ""); err != nil {
		return fmt.Errorf("failed to load centrum disponents: %w", err)
	}

	if err := s.LoadUnits(ctx, 0); err != nil {
		return fmt.Errorf("failed to load centrum units: %w", err)
	}

	if err := s.LoadDispatches(ctx, 0); err != nil {
		return fmt.Errorf("failed to load centrum dispatches: %w", err)
	}

	return nil
}
