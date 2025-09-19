//nolint:gosec // G115 year and week should never be larger than an int32
package jobs

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getTimeclockStats(
	ctx context.Context,
	condition mysql.BoolExpression,
) (*jobs.TimeclockStats, error) {
	stmt := tTimeClock.
		SELECT(
			tTimeClock.Job.AS("timeclock_stats.job"),
			mysql.SUM(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_sum"),
			mysql.AVG(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_avg"),
			mysql.MAX(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_max"),
		).
		FROM(tTimeClock).
		WHERE(mysql.AND(
			condition,
			tTimeClock.Date.BETWEEN(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(7, mysql.DAY))),
				mysql.CURRENT_DATE(),
			),
		)).
		GROUP_BY(tTimeClock.Job)

	var dest jobs.TimeclockStats
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) getTimeclockWeeklyStats(
	ctx context.Context,
	condition mysql.BoolExpression,
) ([]*jobs.TimeclockWeeklyStats, error) {
	yearExpr := dbutils.YEAR(tTimeClock.Date)
	weekExpr := dbutils.WEEK(tTimeClock.Date)

	stmt := tTimeClock.
		SELECT(
			yearExpr.AS("timeclock_weekly_stats.year"),
			weekExpr.AS("timeclock_weekly_stats.calendar_week"),
			mysql.SUM(tTimeClock.SpentTime).AS("timeclock_weekly_stats.sum"),
			mysql.AVG(tTimeClock.SpentTime).AS("timeclock_weekly_stats.avg"),
			mysql.MAX(tTimeClock.SpentTime).AS("timeclock_weekly_stats.max"),
		).
		FROM(tTimeClock).
		WHERE(condition).
		GROUP_BY(
			yearExpr,
			weekExpr,
		).
		ORDER_BY(
			yearExpr.DESC(),
			weekExpr.DESC(),
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
		lastCalendarWeek := last.GetCalendarWeek()

		for i, s := range slices.Backward(dest) {
			if last.GetYear() != s.GetYear() {
				continue
			}

			if len(dest) >= i {
				if dest[i].GetCalendarWeek() == lastCalendarWeek-1 {
					lastCalendarWeek--
					continue
				}
			}

			for lastCalendarWeek-s.GetCalendarWeek() > 1 {
				lastCalendarWeek--
				dest = append([]*jobs.TimeclockWeeklyStats{
					{
						Year:         last.GetYear(),
						CalendarWeek: lastCalendarWeek,
					},
				}, dest...)
			}
		}

		slices.SortFunc(dest, func(a, b *jobs.TimeclockWeeklyStats) int {
			return int(a.GetYear() - b.GetYear() + a.GetCalendarWeek() - b.GetCalendarWeek())
		})

		if dest[0].GetYear() == last.GetYear() && dest[0].GetCalendarWeek() > 1 {
			dest = append([]*jobs.TimeclockWeeklyStats{
				{
					Year:         last.GetYear(),
					CalendarWeek: dest[0].GetCalendarWeek() - 1,
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
			},
			&jobs.TimeclockWeeklyStats{
				Year:         int32(year),
				CalendarWeek: int32(week + 1),
			})
	}

	return dest, nil
}
