//nolint:gosec // G115 year and week should never be larger than an int32
package jobs

import (
	"context"
	"errors"
	"slices"
	"time"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const TimeclockMaxDays = (365 / 2) * 24 * time.Hour

func (s *Store) CleanupTimeclock(ctx context.Context, db qrm.DB) error {
	stmt := tTimeClock.
		UPDATE().
		SET(
			tTimeClock.StartTime.SET(mysql.TimestampExp(mysql.NULL)),
		).
		WHERE(mysql.AND(
			tTimeClock.Date.BETWEEN(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(14, mysql.DAY))),
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(2, mysql.DAY))),
			),
			tTimeClock.StartTime.IS_NOT_NULL(),
			tTimeClock.EndTime.IS_NULL(),
		)).
		LIMIT(1000)

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (s *Store) GetTimeclockStats(
	ctx context.Context,
	db qrm.DB,
	q TimeclockQuery,
) (*jobstimeclock.TimeclockStats, error) {
	condition := mysql.AND(tTimeClock.Job.EQ(mysql.String(q.Job)))
	if len(q.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			ids[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tTimeClock.UserID.IN(ids...))
	} else if q.UserID > 0 {
		condition = condition.AND(tTimeClock.UserID.EQ(mysql.Int32(q.UserID)))
	}

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

	var dest jobstimeclock.TimeclockStats
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Store) GetTimeclockWeeklyStats(
	ctx context.Context,
	db qrm.DB,
	q TimeclockQuery,
) ([]*jobstimeclock.TimeclockWeeklyStats, error) {
	condition := mysql.AND(tTimeClock.Job.EQ(mysql.String(q.Job)))
	if len(q.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			ids[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tTimeClock.UserID.IN(ids...))
	} else if q.UserID > 0 {
		condition = condition.AND(tTimeClock.UserID.EQ(mysql.Int32(q.UserID)))
	}

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
		GROUP_BY(yearExpr, weekExpr).
		ORDER_BY(yearExpr.DESC(), weekExpr.DESC()).
		LIMIT(10)

	var dest []*jobstimeclock.TimeclockWeeklyStats
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	slices.Reverse(dest)

	if len(dest) > 0 {
		last := dest[len(dest)-1]
		lastCalendarWeek := last.GetCalendarWeek()

		for i, s := range slices.Backward(dest) {
			if last.GetYear() != s.GetYear() {
				continue
			}

			if len(dest) > i {
				if s.GetCalendarWeek() == lastCalendarWeek-1 {
					lastCalendarWeek--
					continue
				}
			}

			for lastCalendarWeek-s.GetCalendarWeek() > 1 {
				lastCalendarWeek--
				dest = append(
					[]*jobstimeclock.TimeclockWeeklyStats{
						{Year: last.GetYear(), CalendarWeek: lastCalendarWeek},
					},
					dest...)
			}
		}

		slices.SortFunc(dest, func(a, b *jobstimeclock.TimeclockWeeklyStats) int {
			return int(a.GetYear() - b.GetYear() + a.GetCalendarWeek() - b.GetCalendarWeek())
		})

		if dest[0].GetYear() == last.GetYear() && dest[0].GetCalendarWeek() > 1 {
			dest = append(
				[]*jobstimeclock.TimeclockWeeklyStats{
					{Year: last.GetYear(), CalendarWeek: dest[0].GetCalendarWeek() - 1},
				},
				dest...)
		}
	} else {
		year, week := time.Now().ISOWeek()
		dest = append(dest,
			&jobstimeclock.TimeclockWeeklyStats{Year: int32(year), CalendarWeek: int32(week)},
			&jobstimeclock.TimeclockWeeklyStats{Year: int32(year), CalendarWeek: int32(week + 1)},
		)
	}

	return dest, nil
}

func (s *Store) CountInactiveEmployees(
	ctx context.Context,
	db qrm.DB,
	q InactiveEmployeesQuery,
) (int64, error) {
	tColleague := table.FivenetUser.AS("colleague")
	tUserJobs := table.FivenetUserJobs
	tUserProps := table.FivenetUserProps

	condition := mysql.AND(
		tTimeClock.Job.EQ(mysql.String(q.Job)),
		tUserJobs.Job.EQ(mysql.String(q.Job)),
		mysql.OR(
			mysql.AND(tColleagueProps.AbsenceBegin.IS_NULL(), tColleagueProps.AbsenceEnd.IS_NULL()),
			tColleagueProps.AbsenceBegin.LT_EQ(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(q.Days, mysql.DAY))),
			),
			tColleagueProps.AbsenceEnd.LT_EQ(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(q.Days, mysql.DAY))),
			),
		),
		tTimeClock.UserID.NOT_IN(
			tTimeClock.SELECT(tTimeClock.UserID).FROM(tTimeClock).WHERE(mysql.AND(
				tTimeClock.Job.EQ(mysql.String(q.Job)),
				tTimeClock.Date.GT_EQ(
					mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(q.Days, mysql.DAY))),
				),
			)).GROUP_BY(tTimeClock.UserID),
		),
	)

	countStmt := tTimeClock.
		SELECT(mysql.COUNT(mysql.DISTINCT(tTimeClock.UserID)).AS("data_count.total")).
		FROM(tTimeClock.
			INNER_JOIN(tColleague, tColleague.ID.EQ(tTimeClock.UserID)).
			LEFT_JOIN(tUserProps, tUserProps.UserID.EQ(tTimeClock.UserID)).
			LEFT_JOIN(tColleagueProps, mysql.AND(tColleagueProps.UserID.EQ(tTimeClock.UserID), tColleagueProps.Job.EQ(mysql.String(q.Job)))),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) CountTimeclock(ctx context.Context, db qrm.DB, q TimeclockQuery) (int64, error) {
	condition := mysql.AND(tTimeClock.Job.EQ(mysql.String(q.Job)))

	if q.UserMode <= jobstimeclock.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF {
		condition = condition.AND(tTimeClock.UserID.EQ(mysql.Int32(q.UserID)))
	} else if len(q.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			ids[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tTimeClock.UserID.IN(ids...))
	}

	if q.Date != nil {
		if q.Mode <= jobstimeclock.TimeclockMode_TIMECLOCK_MODE_DAILY {
			condition = condition.AND(tTimeClock.Date.EQ(mysql.DateT(q.Date.GetEnd().AsTime())))
		} else if q.Mode == jobstimeclock.TimeclockMode_TIMECLOCK_MODE_WEEKLY {
			if q.Date.GetEnd() != nil {
				condition = condition.AND(
					mysql.RawBool(
						"YEARWEEK(`timeclock_entry`.`date`, 1) = YEARWEEK($date, 1)",
						mysql.RawArgs{"$date": q.Date.GetEnd().AsTime()},
					),
				)
			}
		} else {
			if q.Date.GetStart() != nil {
				condition = condition.AND(
					tTimeClock.Date.GT_EQ(mysql.DateT(q.Date.GetStart().AsTime())),
				)
			}
			if q.Date.GetEnd() != nil {
				condition = condition.AND(
					tTimeClock.Date.LT_EQ(mysql.DateT(q.Date.GetEnd().AsTime())),
				)
			}
		}
		if q.Date.GetStart() != nil && time.Since(q.Date.GetStart().AsTime()) >= TimeclockMaxDays {
			return 0, nil
		}
		if q.Date.GetEnd() != nil && time.Since(q.Date.GetEnd().AsTime()) >= TimeclockMaxDays {
			return 0, nil
		}
	}

	countExpr := mysql.COUNT(mysql.DISTINCT(tTimeClock.UserID))
	if q.PerDay {
		countExpr = mysql.COUNT(
			mysql.RawString("DISTINCT `timeclock_entry`.`date`, `timeclock_entry`.`user_id`"),
		)
	}

	countStmt := tTimeClock.
		SELECT(countExpr.AS("data_count.total")).
		FROM(tTimeClock).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListTimeclock(
	ctx context.Context,
	db qrm.DB,
	q TimeclockQuery,
) ([]*jobstimeclock.TimeclockEntry, error) {
	userInfoJob := q.Job
	tColleague := table.FivenetUser.AS("colleague")

	condition := mysql.AND(tTimeClock.Job.EQ(mysql.String(userInfoJob)))
	if q.UserMode <= jobstimeclock.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF {
		condition = condition.AND(tTimeClock.UserID.EQ(mysql.Int32(q.UserID)))
	} else if len(q.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			ids[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tTimeClock.UserID.IN(ids...))
	}

	if q.Date != nil {
		if q.Mode <= jobstimeclock.TimeclockMode_TIMECLOCK_MODE_DAILY {
			if q.Date.GetEnd() == nil {
				q.Date.End = timestamp.Now()
			}
			condition = condition.AND(tTimeClock.Date.EQ(mysql.DateT(q.Date.GetEnd().AsTime())))
		} else if q.Mode == jobstimeclock.TimeclockMode_TIMECLOCK_MODE_WEEKLY {
			if q.Date.GetEnd() != nil {
				condition = condition.AND(
					mysql.RawBool(
						"YEARWEEK(`timeclock_entry`.`date`, 1) = YEARWEEK($date, 1)",
						mysql.RawArgs{"$date": q.Date.GetEnd().AsTime()},
					),
				)
			}
		} else {
			if q.Date.GetStart() != nil {
				condition = condition.AND(
					tTimeClock.Date.GT_EQ(mysql.DateT(q.Date.GetStart().AsTime())),
				)
			}
			if q.Date.GetEnd() != nil {
				condition = condition.AND(
					tTimeClock.Date.LT_EQ(mysql.DateT(q.Date.GetEnd().AsTime())),
				)
			}
		}
	}

	groupBys := []mysql.GroupByClause{}
	countExpr := mysql.COUNT(mysql.DISTINCT(tTimeClock.UserID))
	dateSelectExpr := mysql.Expression(tTimeClock.Date)
	if q.PerDay {
		groupBys = append(groupBys, tTimeClock.Date, tTimeClock.UserID)
		countExpr = mysql.COUNT(
			mysql.RawString("DISTINCT `timeclock_entry`.`date`, `timeclock_entry`.`user_id`"),
		)
	} else {
		groupBys = append(groupBys, tTimeClock.UserID)
		dateSelectExpr = mysql.MAX(tTimeClock.Date)
	}

	countStmt := tTimeClock.
		SELECT(countExpr.AS("data_count.total")).
		FROM(tTimeClock.
			INNER_JOIN(tColleague,
				tColleague.ID.EQ(tTimeClock.UserID),
			),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if count.Total <= 0 {
		return []*jobstimeclock.TimeclockEntry{}, nil
	}

	spentTimeColumn := mysql.StringColumn("timeclock_entry.spent_time")
	orderBys := []mysql.OrderByClause{tTimeClock.Date.DESC(), spentTimeColumn.DESC()}
	if q.Sort != nil && len(q.Sort.GetColumns()) > 0 {
		orderBys = []mysql.OrderByClause{}
		for _, sc := range q.Sort.GetColumns() {
			switch sc.GetId() {
			case "date":
				if sc.GetDesc() {
					orderBys = append(orderBys, mysql.DateColumn("agg.date").DESC())
				} else {
					orderBys = append(orderBys, mysql.DateColumn("agg.date").ASC())
				}
			case rankColumn:
				if sc.GetDesc() {
					orderBys = append(orderBys, tColleagueProps.NamePrefix.DESC())
				} else {
					orderBys = append(orderBys, tColleagueProps.NamePrefix.ASC())
				}
			case nameColumn:
				if sc.GetDesc() {
					orderBys = append(orderBys, tColleague.Firstname.DESC())
				} else {
					orderBys = append(orderBys, tColleague.Firstname.ASC())
				}
			default:
				if sc.GetDesc() {
					orderBys = append(orderBys, spentTimeColumn.DESC())
				} else {
					orderBys = append(orderBys, spentTimeColumn.ASC())
				}
			}
		}
	}

	agg := tTimeClock.
		SELECT(
			tTimeClock.UserID.AS("agg.user_id"),
			dateSelectExpr.AS("agg.date"),
			mysql.MIN(tTimeClock.StartTime).AS("agg.start_time"),
			mysql.MAX(tTimeClock.EndTime).AS("agg.end_time"),
			mysql.SUM(tTimeClock.SpentTime).AS("agg.spent_time"),
		).
		FROM(tTimeClock).
		WHERE(condition).
		GROUP_BY(groupBys...).
		AsTable("agg")

	stmt := agg.
		SELECT(
			mysql.IntegerColumn("agg.user_id").AS("timeclock_entry.user_id"),
			mysql.DateColumn("agg.date").AS("timeclock_entry.date"),
			mysql.DateTimeColumn("agg.start_time").AS("timeclock_entry.start_time"),
			mysql.DateTimeColumn("agg.end_time").AS("timeclock_entry.end_time"),
			mysql.FloatColumn("agg.spent_time").AS("timeclock_entry.spent_time"),
			tColleague.ID,
			tUserJobs.Job.AS("colleague.job"),
			tUserJobs.Grade.AS("colleague.job_grade"),
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		FROM(
			agg.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(mysql.IntegerColumn("agg.user_id")),
				).
				INNER_JOIN(tUserJobs,
					mysql.AND(
						tUserJobs.UserID.EQ(tColleague.ID),
						tUserJobs.Job.EQ(mysql.String(userInfoJob)),
					),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tColleague.ID),
				).
				LEFT_JOIN(tColleagueProps,
					mysql.AND(
						tColleagueProps.UserID.EQ(tColleague.ID),
						tColleague.Job.EQ(mysql.String(userInfoJob)),
					),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		ORDER_BY(orderBys...).
		OFFSET(q.Offset).
		LIMIT(q.Limit)

	entries := []*jobstimeclock.TimeclockEntry{}
	if err := stmt.QueryContext(ctx, db, &entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return entries, nil
}

func (s *Store) ListTimeclockTimeline(
	ctx context.Context,
	db qrm.DB,
	q TimeclockQuery,
) ([]*jobstimeclock.TimeclockEntry, error) {
	tColleague := table.FivenetUser.AS("colleague")
	condition := mysql.AND(tTimeClock.Job.EQ(mysql.String(q.Job)))
	if q.UserMode <= jobstimeclock.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF {
		condition = condition.AND(tTimeClock.UserID.EQ(mysql.Int32(q.UserID)))
	} else if len(q.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			ids[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tTimeClock.UserID.IN(ids...))
	}

	if q.Date != nil {
		if q.Date.GetStart() != nil {
			condition = condition.AND(
				tTimeClock.Date.GT_EQ(mysql.DateT(q.Date.GetStart().AsTime())),
			)
		}
		if q.Date.GetEnd() != nil {
			condition = condition.AND(tTimeClock.Date.LT_EQ(mysql.DateT(q.Date.GetEnd().AsTime())))
		}
	}

	spentTimeColumn := mysql.StringColumn("timeclock_entry.spent_time")
	orderBys := []mysql.OrderByClause{tTimeClock.Date.DESC(), spentTimeColumn.DESC()}
	if q.Sort != nil && len(q.Sort.GetColumns()) > 0 {
		orderBys = []mysql.OrderByClause{}
		for _, sc := range q.Sort.GetColumns() {
			switch sc.GetId() {
			case "date":
				if sc.GetDesc() {
					orderBys = append(orderBys, mysql.DateColumn("timeclock_entry.date").DESC())
				} else {
					orderBys = append(orderBys, mysql.DateColumn("timeclock_entry.date").ASC())
				}
			case rankColumn:
				if sc.GetDesc() {
					orderBys = append(orderBys, tColleagueProps.NamePrefix.DESC())
				} else {
					orderBys = append(orderBys, tColleagueProps.NamePrefix.ASC())
				}
			case nameColumn:
				if sc.GetDesc() {
					orderBys = append(orderBys, tColleague.Firstname.DESC())
				} else {
					orderBys = append(orderBys, tColleague.Firstname.ASC())
				}
			default:
				if sc.GetDesc() {
					orderBys = append(orderBys, spentTimeColumn.DESC())
				} else {
					orderBys = append(orderBys, spentTimeColumn.ASC())
				}
			}
		}
	}

	stmt := tTimeClock.
		SELECT(
			tTimeClock.UserID,
			tTimeClock.Date,
			tTimeClock.StartTime,
			tTimeClock.EndTime,
			tTimeClock.SpentTime,
			tColleague.ID,
			tUserJobs.Job.AS("colleague.job"),
			tUserJobs.Grade.AS("colleague.job_grade"),
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		FROM(
			tTimeClock.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tTimeClock.UserID),
				).
				INNER_JOIN(tUserJobs,
					mysql.AND(
						tUserJobs.UserID.EQ(tColleague.ID),
						tUserJobs.Job.EQ(mysql.String(q.Job)),
					),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tColleague.ID),
				).
				LEFT_JOIN(tColleagueProps,
					mysql.AND(
						tColleagueProps.UserID.EQ(tColleague.ID),
						tColleagueProps.Job.EQ(mysql.String(q.Job)),
					),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(condition).
		ORDER_BY(orderBys...)

	if q.Offset > 0 {
		stmt = stmt.OFFSET(q.Offset)
	}
	if q.Limit > 0 {
		stmt = stmt.LIMIT(q.Limit)
	}

	entries := []*jobstimeclock.TimeclockEntry{}
	if err := stmt.QueryContext(ctx, db, &entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return entries, nil
}

func (s *Store) ListInactiveEmployees(
	ctx context.Context,
	db qrm.DB,
	q InactiveEmployeesQuery,
) ([]*jobscolleagues.Colleague, error) {
	tColleague := table.FivenetUser.AS("colleague")
	tUserJobs := table.FivenetUserJobs
	tUserProps := table.FivenetUserProps
	tAvatar := table.FivenetFiles.AS("profile_picture")

	condition := mysql.AND(
		tTimeClock.Job.EQ(mysql.String(q.Job)),
		tUserJobs.Job.EQ(mysql.String(q.Job)),
		mysql.OR(
			mysql.AND(tColleagueProps.AbsenceBegin.IS_NULL(), tColleagueProps.AbsenceEnd.IS_NULL()),
			tColleagueProps.AbsenceBegin.LT_EQ(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(q.Days, mysql.DAY))),
			),
			tColleagueProps.AbsenceEnd.LT_EQ(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(q.Days, mysql.DAY))),
			),
		),
		tTimeClock.UserID.NOT_IN(
			tTimeClock.SELECT(tTimeClock.UserID).FROM(tTimeClock).WHERE(mysql.AND(
				tTimeClock.Job.EQ(mysql.String(q.Job)),
				tTimeClock.Date.GT_EQ(
					mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(q.Days, mysql.DAY))),
				),
			)).GROUP_BY(tTimeClock.UserID),
		),
	)

	orderBys := []mysql.OrderByClause{}
	if q.Sort != nil && len(q.Sort.GetColumns()) > 0 {
		for _, sc := range q.Sort.GetColumns() {
			var columns []mysql.Column
			switch sc.GetId() {
			case "name":
				columns = append(columns, tColleague.Firstname, tColleague.Lastname)
			case "rank":
				fallthrough
			default:
				columns = append(columns, tColleague.JobGrade)
			}

			for _, column := range columns {
				if sc.GetDesc() {
					orderBys = append(orderBys, column.DESC())
				} else {
					orderBys = append(orderBys, column.ASC())
				}
			}
		}
	} else {
		orderBys = append(orderBys, tColleague.JobGrade.ASC())
	}

	stmt := tTimeClock.
		SELECT(
			tTimeClock.UserID,
			tColleague.ID,
			tUserJobs.Job,
			tUserJobs.Grade,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		FROM(tTimeClock.
			INNER_JOIN(tColleague, tColleague.ID.EQ(tTimeClock.UserID)).
			INNER_JOIN(tUserJobs, mysql.AND(tUserJobs.UserID.EQ(tColleague.ID), tUserJobs.Job.EQ(mysql.String(q.Job)))).
			LEFT_JOIN(tUserProps, tUserProps.UserID.EQ(tTimeClock.UserID)).
			LEFT_JOIN(tColleagueProps, mysql.AND(tColleagueProps.UserID.EQ(tTimeClock.UserID), tColleagueProps.Job.EQ(mysql.String(q.Job)))).
			LEFT_JOIN(tAvatar, tAvatar.ID.EQ(tUserProps.AvatarFileID)),
		).
		WHERE(condition).
		ORDER_BY(orderBys...).
		GROUP_BY(
			tTimeClock.UserID,
			tColleague.ID,
			tUserJobs.Job,
			tUserJobs.Grade,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID,
			tAvatar.FilePath,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		OFFSET(q.Offset).
		LIMIT(q.Limit)

	colleagues := []*jobscolleagues.Colleague{}
	if err := stmt.QueryContext(ctx, db, &colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return colleagues, nil
}
