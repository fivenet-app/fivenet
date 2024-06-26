package jobs

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getTimeclockStats(ctx context.Context, condition jet.BoolExpression) (*jobs.TimeclockStats, error) {
	stmt := tTimeClock.
		SELECT(
			tTimeClock.Job.AS("timeclock_stats.job"),
			jet.SUM(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_sum"),
			jet.AVG(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_avg"),
			jet.MAX(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_max"),
		).
		FROM(tTimeClock).
		WHERE(jet.AND(
			condition,
			tTimeClock.Date.BETWEEN(jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(7, jet.DAY))), jet.CURRENT_DATE()),
		))

	var dest jobs.TimeclockStats
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) getTimeclockWeeklyStats(ctx context.Context, condition jet.BoolExpression) ([]*jobs.TimeclockWeeklyStats, error) {
	stmt := tTimeClock.
		SELECT(
			jet.RawString("YEAR(timeclock_entry.`date`)").AS("timeclock_weekly_stats.year"),
			jet.RawString("WEEK(timeclock_entry.`date`)").AS("timeclock_weekly_stats.calendar_week"),
			jet.SUM(tTimeClock.SpentTime).AS("timeclock_weekly_stats.sum"),
			jet.AVG(tTimeClock.SpentTime).AS("timeclock_weekly_stats.avg"),
			jet.MAX(tTimeClock.SpentTime).AS("timeclock_weekly_stats.max"),
		).
		FROM(tTimeClock).
		WHERE(jet.AND(
			condition,
		)).
		GROUP_BY(
			jet.RawString("YEAR(timeclock_entry.`date`)"),
			jet.RawString("WEEK(timeclock_entry.`date`)"),
		).
		ORDER_BY(
			jet.RawString("`timeclock_weekly_stats.year` DESC"),
			jet.RawString("`timeclock_weekly_stats.calendar_week` DESC"),
		).
		LIMIT(12)

	var dest []*jobs.TimeclockWeeklyStats
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	slices.Reverse(dest)

	return dest, nil
}
