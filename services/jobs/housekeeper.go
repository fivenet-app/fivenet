package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	docstats "github.com/fivenet-app/fivenet/v2026/pkg/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
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
	store jobsstore.IStore
}

const (
	timeclockCleanupRowsAttr = "rows_updated"

	employeeCountDayAttr    = "day"
	employeeCountRowsAttr   = "employee_count_rows"
	onVacationCountRowsAttr = "on_vacation_count_rows"

	recountGroupCountsLastGroupIDAttr = "last_group_id"
	recountGroupCountsProcessedAttr   = "processed_groups"
	recountGroupCountsBatchSize       = 20
)

type HousekeeperParams struct {
	fx.In

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider

	Stats *docstats.Service
	Store jobsstore.IStore
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
		store:  p.Store,
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
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "jobs.groups.recount_cached_counts",
		Schedule: "*/5 * * * *",
		Timeout:  durationpb.New(2 * time.Minute),
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
	h.Add(
		"jobs.groups.recount_cached_counts",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "jobs.groups.recount_cached_counts")
			defer span.End()

			dest := &cron.GenericCronData{
				Attributes: map[string]string{},
			}
			if err := data.Unmarshal(dest); err != nil {
				s.logger.Warn("failed to unmarshal group recount cron data", zap.Error(err))
			}

			lastGroupID := int64(0)
			if raw := dest.GetAttribute(recountGroupCountsLastGroupIDAttr); raw != "" {
				if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
					lastGroupID = parsed
				}
			}

			lastGroupID, processed, err := s.recountGroupCounts(ctx, lastGroupID)
			if err != nil {
				s.logger.Error("error during group recount", zap.Error(err))
				return err
			}

			dest.SetAttribute(recountGroupCountsProcessedAttr, strconv.FormatInt(processed, 10))
			if lastGroupID > 0 {
				dest.SetAttribute(
					recountGroupCountsLastGroupIDAttr,
					strconv.FormatInt(lastGroupID, 10),
				)
			}

			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("failed to marshal group recount cron data. %w", err)
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

func (s *Housekeeper) recountGroupCounts(
	ctx context.Context,
	lastGroupID int64,
) (int64, int64, error) {
	tJobGroups := table.FivenetJobGroups
	groupIDs := make([]int64, 0, recountGroupCountsBatchSize)

	if lastGroupID > 0 {
		// Walk forward from the last processed ID, then wrap to the start if needed.
		var nextGroupIDs []int64
		nextStmt := tJobGroups.
			SELECT(tJobGroups.ID).
			FROM(tJobGroups).
			WHERE(tJobGroups.ID.GT(mysql.Int64(lastGroupID))).
			ORDER_BY(tJobGroups.ID.ASC()).
			LIMIT(recountGroupCountsBatchSize)

		if err := nextStmt.QueryContext(
			ctx,
			s.db,
			&nextGroupIDs,
		); err != nil &&
			!errors.Is(err, qrm.ErrNoRows) {
			return 0, 0, fmt.Errorf("failed to load next job group ids for recount. %w", err)
		}
		groupIDs = append(groupIDs, nextGroupIDs...)

		if remaining := recountGroupCountsBatchSize - len(groupIDs); remaining > 0 {
			var wrappedGroupIDs []int64
			wrapStmt := tJobGroups.
				SELECT(tJobGroups.ID).
				FROM(tJobGroups).
				WHERE(tJobGroups.ID.LT_EQ(mysql.Int64(lastGroupID))).
				ORDER_BY(tJobGroups.ID.ASC()).
				LIMIT(int64(remaining))

			if err := wrapStmt.QueryContext(
				ctx,
				s.db,
				&wrappedGroupIDs,
			); err != nil &&
				!errors.Is(err, qrm.ErrNoRows) {
				return 0, 0, fmt.Errorf("failed to wrap job group ids for recount. %w", err)
			}
			groupIDs = append(groupIDs, wrappedGroupIDs...)
		}
	} else {
		stmt := tJobGroups.
			SELECT(tJobGroups.ID).
			FROM(tJobGroups).
			ORDER_BY(tJobGroups.ID.ASC()).
			LIMIT(recountGroupCountsBatchSize)

		if err := stmt.QueryContext(
			ctx,
			s.db,
			&groupIDs,
		); err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return 0, 0, fmt.Errorf("failed to load initial job group ids for recount. %w", err)
		}
	}

	processed := 0
	for _, groupID := range groupIDs {
		if err := s.store.RecountGroupStats(ctx, s.db, groupID); err != nil {
			return 0, 0, fmt.Errorf(
				"failed to recount cached job group %d stats. %w",
				groupID,
				err,
			)
		}
		processed++
		lastGroupID = groupID
	}

	return lastGroupID, int64(processed), nil
}
