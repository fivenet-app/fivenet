package stats

import (
	"context"
	"time"

	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Service) BuildEmployeeCountMetrics(ctx context.Context) error {
	day := timeutils.StartOfDay(time.Now().UTC())
	tRollup := table.FivenetStatsDailyRollup
	tUserJobs := table.FivenetUserJobs

	var unemployedJobName string
	if s.appCfg != nil {
		appCfg := s.appCfg.Get()
		unemployedJob := appCfg.GetJobInfo().GetUnemployedJob()
		if unemployedJob != nil {
			unemployedJobName = unemployedJob.GetName()
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tRollup.
		DELETE().
		WHERE(mysql.AND(
			tRollup.Day.EQ(mysql.DateT(day)),
			tRollup.SourceKind.EQ(mysql.String(SourceKindEmployeeCount)),
			tRollup.SourceKey.EQ(mysql.String("fivenet_user_jobs")),
			tRollup.MetricKey.IN(
				mysql.String("employee_count"),
				mysql.String("on_vacation_count"),
			),
		)).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	stmt := tRollup.
		INSERT(
			tRollup.Day,
			tRollup.Job,
			tRollup.SourceKind,
			tRollup.SourceKey,
			tRollup.MetricKey,
			tRollup.Dimension1,
			tRollup.Dimension2,
			tRollup.Dimension3,
			tRollup.Value,
		).
		QUERY(
			tUserJobs.
				SELECT(
					mysql.DateT(day),
					tUserJobs.Job,
					mysql.String(SourceKindEmployeeCount),
					mysql.String("fivenet_user_jobs"),
					mysql.String("employee_count"),
					mysql.String(""),
					mysql.String(""),
					mysql.String(""),
					mysql.COUNT(mysql.STAR),
				).
				FROM(tUserJobs).
				WHERE(tUserJobs.Job.NOT_EQ(mysql.String(unemployedJobName))).
				GROUP_BY(tUserJobs.Job),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	tColleagueProps := table.FivenetJobColleagueProps
	vacationStmt := tRollup.
		INSERT(
			tRollup.Day,
			tRollup.Job,
			tRollup.SourceKind,
			tRollup.SourceKey,
			tRollup.MetricKey,
			tRollup.Dimension1,
			tRollup.Dimension2,
			tRollup.Dimension3,
			tRollup.Value,
		).
		QUERY(
			tUserJobs.
				SELECT(
					mysql.DateT(day),
					tUserJobs.Job,
					mysql.String(SourceKindEmployeeCount),
					mysql.String("fivenet_user_jobs"),
					mysql.String("on_vacation_count"),
					mysql.String(""),
					mysql.String(""),
					mysql.String(""),
					mysql.COUNT(mysql.DISTINCT(tUserJobs.UserID)),
				).
				FROM(
					tUserJobs.
						INNER_JOIN(
							tColleagueProps,
							mysql.AND(
								tColleagueProps.UserID.EQ(tUserJobs.UserID),
								tColleagueProps.Job.EQ(tUserJobs.Job),
								tColleagueProps.DeletedAt.IS_NULL(),
							),
						),
				).
				WHERE(mysql.AND(
					tColleagueProps.AbsenceBegin.IS_NOT_NULL(),
					tColleagueProps.AbsenceEnd.IS_NOT_NULL(),
					tColleagueProps.AbsenceBegin.LT_EQ(mysql.CURRENT_DATE()),
					tColleagueProps.AbsenceEnd.GT_EQ(mysql.CURRENT_DATE()),
				)).
				GROUP_BY(tUserJobs.Job),
		)

	if _, err := vacationStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) QueryEmployeeCountOverTime(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
	period pbstats.StatsPeriod,
) ([]*DailyValue, error) {
	return s.QueryPeriodValues(
		ctx,
		startDay,
		endDay,
		job,
		SourceKindEmployeeCount,
		"fivenet_user_jobs",
		"employee_count",
		period,
	)
}

func (s *Service) QueryEmployeeSeriesOverTime(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
	period pbstats.StatsPeriod,
) ([]*PeriodSeriesValue, error) {
	periodExpr := periodStartExpr(period)
	query := `
SELECT
  ` + periodExpr + ` AS day,
  metric_key AS ` + "`key`" + `,
  SUM(value) AS value
FROM fivenet_stats_daily_rollup
WHERE day >= ?
  AND day <= ?
  AND job = ?
  AND source_kind = ?
  AND source_key = ?
  AND metric_key IN ('employee_count', 'on_vacation_count')
GROUP BY ` + periodExpr + `, metric_key
ORDER BY ` + periodExpr + ` ASC, metric_key ASC
`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		timeutils.StartOfDay(startDay).Format(time.DateOnly),
		timeutils.StartOfDay(endDay).Format(time.DateOnly),
		job,
		SourceKindEmployeeCount,
		"fivenet_user_jobs",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*PeriodSeriesValue{}
	for rows.Next() {
		item := &PeriodSeriesValue{}
		if err := rows.Scan(&item.Day, &item.Key, &item.Value); err != nil {
			return nil, err
		}

		item.Label = item.Key
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
