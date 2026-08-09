package jobs

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	docstats "github.com/fivenet-app/fivenet/v2026/pkg/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
	"github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var HousekeeperModule = fx.Module(
	"jobs.housekeeper",
	fx.Provide(
		NewHousekeeper,
	),
)

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	db    *sql.DB
	stats *docstats.Service
}

const (
	timeclockCleanupRowsAttr = "rows_updated"
	employeeCountDayAttr     = "day"
	employeeCountRowsAttr    = "employee_count_rows"
	onVacationCountRowsAttr  = "on_vacation_count_rows"
)

type HousekeeperParams struct {
	fx.In

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider

	Stats *docstats.Service
}

type HousekeeperResult struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewHousekeeper(p HousekeeperParams) HousekeeperResult {
	s := &Housekeeper{
		logger: p.Logger.Named("jobs.housekeeper"),
		tracer: p.TP.Tracer("jobs.housekeeper"),
		db:     p.DB,
		stats:  p.Stats,
	}

	return HousekeeperResult{
		Housekeeper:  s,
		CronRegister: s,
	}
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "jobs.timeclock_cleanup",
		Schedule: "@hourly",
	}); err != nil {
		return err
	}
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "jobs.stats.employee_count.recent",
		Schedule: "5 */1 * * *",
	}); err != nil {
		return err
	}

	if err := registry.UnregisterCronjob(ctx, "jobs.timeclock_handling"); err != nil {
		s.logger.Error("failed to unregister jobs.timeclock_handling", zap.Error(err))
	}
	if err := registry.UnregisterCronjob(ctx, "jobs-timeclock-handling"); err != nil {
		s.logger.Error("failed to unregister jobs.timeclock_handling", zap.Error(err))
	}

	return nil
}

func (s *Housekeeper) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add("jobs.timeclock_cleanup", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := s.tracer.Start(ctx, "jobs.timeclock_cleanup")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if err := data.Unmarshal(dest); err != nil {
			s.logger.Warn("failed to unmarshal timeclock cleanup cron data", zap.Error(err))
		}

		rowsUpdated, err := s.timeclockCleanup(ctx)
		if err != nil {
			s.logger.Error("error during timeclock cleanup", zap.Error(err))
			return err
		}
		dest.SetAttribute(timeclockCleanupRowsAttr, strconv.FormatInt(rowsUpdated, 10))

		if err := data.MarshalFrom(dest); err != nil {
			return err
		}

		return nil
	})
	h.Add(
		"jobs.stats.employee_count.recent",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "jobs.stats.employee_count.recent")
			defer span.End()

			dest := &cron.GenericCronData{
				Attributes: map[string]string{},
			}
			if err := data.Unmarshal(dest); err != nil {
				s.logger.Warn("failed to unmarshal employee count cron data", zap.Error(err))
			}

			employeeCountRows, onVacationRows, err := s.stats.BuildEmployeeCountMetrics(ctx)
			if err != nil {
				s.logger.Error("error during employee count metrics build", zap.Error(err))
				return err
			}

			dest.SetAttribute(
				employeeCountDayAttr,
				timeutils.StartOfDay(time.Now().UTC()).Format(time.DateOnly),
			)
			dest.SetAttribute(employeeCountRowsAttr, strconv.FormatInt(employeeCountRows, 10))
			dest.SetAttribute(onVacationCountRowsAttr, strconv.FormatInt(onVacationRows, 10))

			if err := data.MarshalFrom(dest); err != nil {
				return err
			}

			return nil
		},
	)

	return nil
}

func (s *Housekeeper) timeclockCleanup(ctx context.Context) (int64, error) {
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

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsUpdated, nil
}
