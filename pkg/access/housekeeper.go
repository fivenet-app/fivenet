package access

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	subjectCleanupOrphanSubjectsRemovedAttr = "orphan_subjects_removed"
	subjectCleanupStaleJobGradeRemovedAttr  = "stale_job_grade_subjects_removed"
)

var Module = fx.Module("access.housekeeper",
	fx.Provide(NewHousekeeper),
)

type Housekeeper struct {
	logger   *zap.Logger
	tracer   trace.Tracer
	resolver *SubjectResolver
}

type HousekeeperParams struct {
	fx.In

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider
}

type HousekeeperResult struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewHousekeeper(p HousekeeperParams) HousekeeperResult {
	h := &Housekeeper{
		logger:   p.Logger.Named("access.housekeeper"),
		tracer:   p.TP.Tracer("access.housekeeper"),
		resolver: NewSubjectResolver(p.DB),
	}

	return HousekeeperResult{
		Housekeeper:  h,
		CronRegister: h,
	}
}

func (h *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "access.subjects.cleanup",
		Schedule: "*/10 * * * *",
		Timeout:  durationpb.New(1 * time.Minute),
	}); err != nil {
		return err
	}

	return nil
}

func (h *Housekeeper) RegisterCronjobHandlers(hand *croner.Handlers) error {
	hand.Add("access.subjects.cleanup", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "access.subjects.cleanup")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if err := data.Unmarshal(dest); err != nil {
			h.logger.Warn("failed to unmarshal access subjects cleanup cron data", zap.Error(err))
		}

		var errs error
		orphanRemoved, err := h.resolver.CleanupOrphanSubjects(ctx, h.resolver.db)
		if err != nil {
			h.logger.Error("error during orphan subject cleanup", zap.Error(err))
			errs = multierr.Append(errs, fmt.Errorf("error during orphan subject cleanup. %w", err))
		}
		dest.SetAttribute(
			subjectCleanupOrphanSubjectsRemovedAttr,
			strconv.FormatInt(orphanRemoved, 10),
		)

		staleRemoved, err := h.resolver.CleanupStaleJobGradeSubjects(ctx, h.resolver.db)
		if err != nil {
			h.logger.Error("error during stale job grade subject cleanup", zap.Error(err))
			if errs != nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf("error during stale job grade subject cleanup. %w", err),
				)
			}
		}
		dest.SetAttribute(
			subjectCleanupStaleJobGradeRemovedAttr,
			strconv.FormatInt(staleRemoved, 10),
		)

		if err := data.MarshalFrom(dest); err != nil {
			return fmt.Errorf(
				"failed to marshal updated access subjects cleanup cron data. %w",
				err,
			)
		}

		return errs
	})

	return nil
}
