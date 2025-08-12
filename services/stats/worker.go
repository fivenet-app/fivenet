package stats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sync "sync"
	"sync/atomic"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var (
	tAccounts         = table.FivenetAccounts
	tDocuments        = table.FivenetDocuments
	tDispatches       = table.FivenetCentrumDispatches
	tCitizenActivity  = table.FivenetUserActivity
	tJobUserTimeclock = table.FivenetJobTimeclock
)

type worker struct {
	mu sync.Mutex

	logger *zap.Logger
	db     *sql.DB

	active atomic.Bool
	stats  atomic.Pointer[map[string]*stats.Stat]
}

func newWorker(logger *zap.Logger, db *sql.DB) *worker {
	return &worker{
		logger: logger,
		db:     db,
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
	data := Stats{}

	tUsers := tables.User()
	queries := map[string]jet.Statement{
		"users_registered": tAccounts.SELECT(jet.COUNT(tAccounts.ID).AS("value")).
			WHERE(tAccounts.Enabled.IS_TRUE()),
		"documents_created": tDocuments.SELECT(jet.COUNT(tDocuments.ID).AS("value")).
			WHERE(tDocuments.DeletedAt.IS_NULL()),
		"dispatches_created": tDispatches.SELECT(jet.MAX(tDispatches.ID).AS("value")),
		"citizen_activity":   tCitizenActivity.SELECT(jet.COUNT(tCitizenActivity.ID).AS("value")),
		"timeclock_tracked": tJobUserTimeclock.SELECT(
			jet.CAST(jet.SUM(tJobUserTimeclock.SpentTime)).AS_SIGNED().AS("value"),
		),
		"citizens_total": tUsers.SELECT(jet.COUNT(tUsers.ID).AS("value")),
	}

	errs := multierr.Combine()

	zero := int32(0)
	for key, query := range queries {
		dest := &struct {
			Value *int32 `alias:"value"`
		}{
			Value: &zero,
		}
		if err := query.QueryContext(ctx, s.db, dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				errs = multierr.Append(
					errs,
					fmt.Errorf("error during %q stats query. %w", key, err),
				)
				continue
			}
		}
		if (*dest.Value % 10) != 0 {
			*dest.Value = (10 - *dest.Value%10) + *dest.Value
		}

		data[key] = &stats.Stat{
			Value: dest.Value,
		}
	}

	s.stats.Store(&data)

	return errs
}

func (s *worker) GetStats() *Stats {
	return s.stats.Load()
}
