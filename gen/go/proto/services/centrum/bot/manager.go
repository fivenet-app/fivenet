package bot

import (
	"context"
	"sync"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/galexrt/fivenet/pkg/server/admin"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	metricBotActive = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "centrum_bot",
		Name:      "active",
		Help:      "If centrum bot is active or not.",
	}, []string{"job"})
)

var Module = fx.Module("centrum_bot_manager", fx.Provide(
	NewManager,
))

type Manager struct {
	ctx    context.Context
	logger *zap.Logger
	mutex  sync.RWMutex
	wg     sync.WaitGroup

	tracer trace.Tracer

	bots *xsync.MapOf[string, context.CancelFunc]
	js   nats.JetStreamContext

	state *manager.Manager
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	State  *manager.Manager
	JS     nats.JetStreamContext
}

func NewManager(p Params) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	b := &Manager{
		ctx:    ctx,
		logger: p.Logger.Named("centrum.bot.manager"),
		mutex:  sync.RWMutex{},
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		bots:  xsync.NewMapOf[string, context.CancelFunc](),
		js:    p.JS,
		state: p.State,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			b.Run()
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

func (s *Manager) Run() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(3 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-bots-check")
				defer span.End()

				if err := s.checkIfBotsAreNeeded(ctx); err != nil {
					s.logger.Error("failed to check if bots need to be (de-)activated", zap.Error(err))
				}
			}()
		}
	}
}

func (b *Manager) Start(job string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Already a bot active
	if _, ok := b.bots.Load(job); ok {
		return nil
	}

	b.logger.Info("starting centrum dispatch bot", zap.String("job", job))
	bot := NewBot(b.logger.With(zap.String("job", job)), job, b.state)
	ctx, cancel := context.WithCancel(b.ctx)
	b.bots.Store(job, cancel)
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		bot.Run(ctx)
	}()

	metricBotActive.WithLabelValues(job).Set(1)

	return nil
}

func (b *Manager) Stop(job string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	cancel, ok := b.bots.Load(job)
	if !ok {
		return nil
	}

	b.logger.Info("stopping centrum dispatch bot", zap.String("job", job))

	cancel()

	metricBotActive.WithLabelValues(job).Set(0)

	b.bots.Delete(job)

	return nil
}

func (s *Manager) checkIfBotsAreNeeded(ctx context.Context) error {
	for _, settings := range s.state.ListSettings() {
		if s.state.CheckIfBotNeeded(settings.Job) {
			if err := s.Start(settings.Job); err != nil {
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