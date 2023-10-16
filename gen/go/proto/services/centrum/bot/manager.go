package bot

import (
	"context"
	"sync"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/galexrt/fivenet/pkg/events"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

	bots   map[string]context.CancelFunc
	events *events.Eventus

	state *manager.Manager
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	State   *manager.Manager
	Eventus *events.Eventus
}

func NewManager(p Params) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	b := &Manager{
		ctx:    ctx,
		logger: p.Logger.Named("centrum_bots"),
		mutex:  sync.RWMutex{},
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		bots:   map[string]context.CancelFunc{},
		events: p.Eventus,
		state:  p.State,
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
	if _, ok := b.bots[job]; ok {
		return nil
	}

	bot := &Bot{
		state: b.state,
		job:   job,
	}
	ctx, cancel := context.WithCancel(b.ctx)
	b.bots[job] = cancel
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		bot.Run(ctx)
	}()

	return nil
}

func (b *Manager) Stop(job string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	cancel, ok := b.bots[job]
	if !ok {
		return nil
	}

	cancel()

	return nil
}

func (s *Manager) checkIfBotsAreNeeded(ctx context.Context) error {
	s.state.Settings.Range(func(job string, value *dispatch.Settings) bool {
		if s.state.CheckIfBotNeeded(job) {
			if err := s.Start(job); err != nil {
				s.logger.Error("failed to start dispatch center bot for job", zap.String("job", job))
			}
		} else {
			if err := s.Stop(job); err != nil {
				s.logger.Error("failed to stop dispatch center bot for job", zap.String("job", job))
			}
		}

		return true
	})

	return nil
}
