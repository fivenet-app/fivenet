package jobs

import (
	"context"
	"errors"
	"slices"
	"time"

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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
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
		LIMIT(10)

	var dest []*jobs.TimeclockWeeklyStats
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	slices.Reverse(dest)

	// Add "null" values at the begin of the stats for better UX
	if len(dest) > 0 {
		last := dest[len(dest)-1]
		lastCalendarWeek := last.CalendarWeek

		for i, s := range slices.Backward(dest) {
			if last.Year != s.Year {
				continue
			}

			if len(dest) >= i {
				if dest[i].CalendarWeek == lastCalendarWeek-1 {
					lastCalendarWeek--
					continue
				}
			}

			for lastCalendarWeek-s.CalendarWeek > 1 {
				lastCalendarWeek--
				dest = append([]*jobs.TimeclockWeeklyStats{
					{
						Year:         last.Year,
						CalendarWeek: lastCalendarWeek,
						Sum:          0,
						Avg:          0,
						Max:          0,
					},
				}, dest...)
			}
		}

		slices.SortFunc(dest, func(a, b *jobs.TimeclockWeeklyStats) int {
			return int(a.Year - b.Year + a.CalendarWeek - b.CalendarWeek)
		})

		if dest[0].Year == last.Year && dest[0].CalendarWeek > 1 {
			dest = append([]*jobs.TimeclockWeeklyStats{
				{
					Year:         last.Year,
					CalendarWeek: dest[0].CalendarWeek - 1,
					Sum:          0,
					Avg:          0,
					Max:          0,
				},
			}, dest...)
		}
	} else {
		// No stats? Add two empty ones so the graph doesn't break
		year, week := time.Now().ISOWeek()
		dest = append(dest,
			&jobs.TimeclockWeeklyStats{
				Year:         int32(year),
				CalendarWeek: int32(week),
				Sum:          0,
				Avg:          0,
				Max:          0,
			},
			&jobs.TimeclockWeeklyStats{
				Year:         int32(year),
				CalendarWeek: int32(week + 1),
				Sum:          0,
				Avg:          0,
				Max:          0,
			})
	}

	return dest, nil
}
