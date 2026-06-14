package stats

import (
	"context"
	sync "sync"
	"sync/atomic"
	"time"

	statsstore "github.com/fivenet-app/fivenet/v2026/stores/stats"
	"go.uber.org/zap"
)

type Stats = statsstore.Stats

type worker struct {
	mu sync.Mutex

	logger *zap.Logger
	store  statsstore.IStore

	active atomic.Bool
	stats  atomic.Pointer[Stats]
}

func newWorker(logger *zap.Logger, store statsstore.IStore) *worker {
	return &worker{
		logger: logger,
		store:  store,
		active: atomic.Bool{},
	}
}

func (s *worker) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Enable stats calculation when enabled and not active yet
	if s.active.Load() {
		return
	}

	go s.calculateStats(ctx)

	s.active.Store(true)
}

func (s *worker) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active.Load() {
		return
	}

	s.active.Store(false)
}

func (s *worker) calculateStats(ctx context.Context) {
	for {
		if err := s.loadStats(ctx); err != nil {
			s.logger.Error("error while calculating stats", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			s.active.Store(false)
			return

		case <-time.After(30 * time.Minute):
		}
	}
}

func (s *worker) loadStats(ctx context.Context) error {
	data, err := s.store.LoadPublicStats(ctx)
	if data == nil {
		data = Stats{}
	}

	for _, stat := range data {
		if stat == nil || stat.Value == nil {
			continue
		}

		if *stat.Value%10 != 0 {
			*stat.Value = (10 - *stat.Value%10) + *stat.Value
		}
	}

	s.stats.Store(&data)

	return err
}

func (s *worker) GetStats() *Stats {
	return s.stats.Load()
}
