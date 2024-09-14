package bot

import (
	"context"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/server/admin"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v3"
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

var Module = fx.Module("centrum_bot_manager", fx.Provide(
	NewManager,
))

type Manager struct {
	logger *zap.Logger
	mutex  sync.RWMutex
	wg     sync.WaitGroup

	tracer trace.Tracer

	bots *xsync.MapOf[string, *Bot]
	js   *events.JSWrapper

	state   *manager.Manager
	tracker tracker.ITracker
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	State   *manager.Manager
	JS      *events.JSWrapper
	Tracker tracker.ITracker
}

func NewManager(p Params) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	b := &Manager{
		logger: p.Logger.Named("centrum.bot.manager"),
		mutex:  sync.RWMutex{},
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		bots:    xsync.NewMapOf[string, *Bot](),
		js:      p.JS,
		state:   p.State,
		tracker: p.Tracker,
	}

	p.LC.Append(fx.StartHook(func(c context.Context) error {
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
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(3 * time.Second):
			func() {
				ctx, span := s.tracer.Start(ctx, "centrum-bots-check")
				defer span.End()

				if err := s.checkIfBotsAreNeeded(ctx); err != nil {
					s.logger.Error("failed to check if bots need to be (de-)activated", zap.Error(err))
				}
			}()
		}
	}
}

func (b *Manager) Start(ctx context.Context, job string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Already a bot active
	if _, ok := b.bots.Load(job); ok {
		return nil
	}

	b.logger.Info("starting centrum dispatch bot", zap.String("job", job))
	bot := NewBot(ctx, b.logger.With(zap.String("job", job)), job, b.state, b.tracker)
	b.bots.Store(job, bot)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		bot.Run()
	}()

	metricBotActive.WithLabelValues(job).Set(1)

	return nil
}

func (b *Manager) Stop(job string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

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

func (s *Manager) checkIfBotsAreNeeded(ctx context.Context) error {
	for _, settings := range s.state.ListSettings(ctx) {
		if s.state.CheckIfBotNeeded(ctx, settings.Job) {
			if err := s.Start(ctx, settings.Job); err != nil {
				s.logger.Error("failed to start dispatch center bot for job", zap.String("job", settings.Job))
			}
		} else {
			if err := s.Stop(settings.Job); err != nil {
				s.logger.Error("failed to stop dispatch center bot for job", zap.String("job", settings.Job))
			}
		}
	}

	return nil
}
