package jobs

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/pkg/croner"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	db *sql.DB
}

type HousekeeperParams struct {
	fx.In

	Logger       *zap.Logger
	DB           *sql.DB
	TP           *tracesdk.TracerProvider
	CronHandlers *croner.Handlers
}

func NewHousekeeper(p HousekeeperParams) *Housekeeper {
	s := &Housekeeper{
		logger: p.Logger.Named("jobs_housekeeper"),
		tracer: p.TP.Tracer("jobs_housekeeper"),
		db:     p.DB,
	}

	p.CronHandlers.Add("jobs.timeclock_cleanup", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := s.tracer.Start(ctx, "jobs.timeclock_cleanup")
		defer span.End()

		if err := s.timeclockCleanup(ctx); err != nil {
			s.logger.Error("error during timeclock cleanup", zap.Error(err))
			return err
		}

		return nil
	})

	p.CronHandlers.Add("jobs.timeclock_handling", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := s.tracer.Start(ctx, "jobs-timeclock-handling")
		defer span.End()

		if err := s.timeclockHandler(ctx); err != nil {
			s.logger.Error("error during timeclock handling", zap.Error(err))
			return err
		}

		return nil
	})

	return s
}

func (s *Housekeeper) timeclockCleanup(ctx context.Context) error {
	stmt := tTimeClock.
		UPDATE().
		SET(
			tTimeClock.StartTime.SET(jet.TimestampExp(jet.NULL)),
		).
		WHERE(jet.AND(
			tTimeClock.Date.BETWEEN(
				jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(14, jet.DAY))),
				jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(2, jet.DAY))),
			),
			tTimeClock.StartTime.IS_NOT_NULL(),
			tTimeClock.EndTime.IS_NULL(),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Housekeeper) timeclockHandler(ctx context.Context) error {
	stmt := tTimeClock.
		UPDATE().
		SET(
			tTimeClock.SpentTime.SET(jet.FloatExp(jet.Raw("`timeclock_entry`.`spent_time` + CAST((TIMESTAMPDIFF(SECOND, `timeclock_entry`.`start_time`, `timeclock_entry`.`end_time`) / 3600) AS DECIMAL(10,2))"))),
			tTimeClock.StartTime.SET(jet.TimestampExp(jet.NULL)),
			tTimeClock.EndTime.SET(jet.TimestampExp(jet.NULL)),
		).
		WHERE(jet.AND(
			tTimeClock.StartTime.IS_NOT_NULL(),
			tTimeClock.EndTime.IS_NOT_NULL(),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
