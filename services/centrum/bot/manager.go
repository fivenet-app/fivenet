package centrumbot

import (
	"context"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var metricBotActive = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "centrum_bot",
	Name:      "active",
	Help:      "If centrum bot is active or not.",
}, []string{"job_name"})

var Module = fx.Module("centrum_bot_manager",
	fx.Provide(
		NewManager,
	),
)

type Manager struct {
	logger *zap.Logger
	mu     *sync.RWMutex
	wg     sync.WaitGroup

	tracer trace.Tracer

	bots    *xsync.Map[string, *Bot]
	js      *events.JSWrapper
	tracker tracker.ITracker

	helpers    *helpers.Helpers
	settings   *settings.SettingsDB
	units      *units.UnitDB
	dispatches *dispatches.DispatchDB
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	JS      *events.JSWrapper
	Tracker tracker.ITracker

	Helpers    *helpers.Helpers
	Settings   *settings.SettingsDB
	Units      *units.UnitDB
	Dispatches *dispatches.DispatchDB
}

func NewManager(p Params) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	b := &Manager{
		logger: p.Logger.Named("centrum.bot.manager"),
		mu:     &sync.RWMutex{},
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum.cache"),

		bots:    xsync.NewMap[string, *Bot](),
		js:      p.JS,
		tracker: p.Tracker,

		helpers:    p.Helpers,
		settings:   p.Settings,
		units:      p.Units,
		dispatches: p.Dispatches,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			b.Run(ctx)
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		b.wg.Wait()

		return nil
	}))

	return b
}

func (s *Manager) Run(ctx context.Context) {
	s.logger.Info("started centrum bot manager")

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(3 * time.Second):
			func() {
				ctx, span := s.tracer.Start(ctx, "centrum.bots-check")
				defer span.End()

				s.checkIfBotsAreNeeded(ctx)
			}()
		}
	}
}

func (b *Manager) startBot(ctx context.Context, job string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Already a bot active
	if _, ok := b.bots.Load(job); ok {
		return nil
	}

	b.logger.Info("starting centrum dispatch bot", zap.String("job", job))
	bot := NewBot(ctx, b.logger, b.tracker, b.helpers, b.settings, b.units, b.dispatches, job)
	b.bots.Store(job, bot)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		bot.Run()
	}()

	metricBotActive.WithLabelValues(job).Set(1)

	return nil
}

func (b *Manager) stopBot(job string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	bot, ok := b.bots.Load(job)
	if !ok {
		return nil
	}

	b.logger.Info("stopping centrum dispatch bot", zap.String("job", job))

	bot.Stop()

	metricBotActive.WithLabelValues(job).Set(0)

	b.bots.Delete(job)

	return nil
}

func (s *Manager) checkIfBotsAreNeeded(ctx context.Context) {
	for _, settings := range s.settings.List(ctx) {
		if !s.helpers.CheckIfBotNeeded(ctx, settings.GetJob()) {
			if err := s.stopBot(settings.GetJob()); err != nil {
				s.logger.Error(
					"failed to stop dispatch center bot for job",
					zap.String("job", settings.GetJob()),
				)
			}

			continue
		}

		if err := s.startBot(ctx, settings.GetJob()); err != nil {
			s.logger.Error(
				"failed to start dispatch center bot for job",
				zap.String("job", settings.GetJob()),
			)
		}
	}
}
