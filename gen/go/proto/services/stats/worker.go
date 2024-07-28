package stats

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

var (
	tAccounts        = table.FivenetAccounts
	tDocuments       = table.FivenetDocuments
	tDispatches      = table.FivenetCentrumDispatches
	tCitizenActivity = table.FivenetUserActivity
	tJobsTimeclock   = table.FivenetJobsTimeclock
	tUsers           = table.Users
)

func (s *Server) calculateStats(ctx context.Context) {
	for {
		if err := s.loadStats(ctx); err != nil {
			s.logger.Error("error while calculating stats", zap.Error(err))
		}

		select {
		case <-time.After(1 * time.Hour):
		case <-ctx.Done():
			return
		}
	}
}

func (s *Server) loadStats(ctx context.Context) error {
	data := Stats{}

	queries := map[string]jet.Statement{
		"users_registered":   tAccounts.SELECT(jet.COUNT(tAccounts.ID).AS("value")),
		"documents_created":  tDocuments.SELECT(jet.COUNT(tDocuments.ID).AS("value")),
		"dispatches_created": tDispatches.SELECT(jet.COUNT(tDispatches.ID).AS("value")),
		"citizen_activity":   tCitizenActivity.SELECT(jet.COUNT(tCitizenActivity.ID).AS("value")),
		"timeclock_tracked":  tJobsTimeclock.SELECT(jet.SUM(tJobsTimeclock.SpentTime).AS("value")),
		"citizens_total":     tUsers.SELECT(jet.COUNT(tUsers.ID).AS("value")),
	}

	zero := int32(0)
	for key, query := range queries {
		dest := &struct {
			Value *int32 `alias:"value"`
		}{
			Value: &zero,
		}
		if err := query.QueryContext(ctx, s.db, dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return err
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

	return nil
}
