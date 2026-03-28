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
	tRollup := table.FivenetStatsDailyRollup
	periodExpr := periodStartDateExpr(period)

	stmt := tRollup.
		SELECT(
			periodExpr.AS("day"),
			tRollup.MetricKey.AS("key"),
			mysql.SUM(tRollup.Value).AS("value"),
		).
		FROM(tRollup).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.EQ(mysql.String(job)),
			tRollup.SourceKind.EQ(mysql.String(SourceKindEmployeeCount)),
			tRollup.SourceKey.EQ(mysql.String("fivenet_user_jobs")),
			tRollup.MetricKey.IN(
				mysql.String("employee_count"),
				mysql.String("on_vacation_count"),
			),
		)).
		GROUP_BY(periodExpr, tRollup.MetricKey).
		ORDER_BY(periodExpr.ASC(), tRollup.MetricKey.ASC())

	rowsDest := []*PeriodSeriesValue{}
	if err := stmt.QueryContext(ctx, s.db, &rowsDest); err != nil {
		return nil, err
	}

	items := make([]*PeriodSeriesValue, 0, len(rowsDest))
	for i := range rowsDest {
		rowsDest[i].Label = rowsDest[i].Key
		items = append(items, rowsDest[i])
	}

	return items, nil
}
