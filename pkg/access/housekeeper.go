package access

import (
	"context"
	"database/sql"
	"fmt"
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
	subjectCleanupCronName     = "access.subjects.cleanup"
	subjectCleanupCronSchedule = "*/10 * * * *"
	subjectCleanupCronTimeout  = time.Minute
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
		Name:     subjectCleanupCronName,
		Schedule: subjectCleanupCronSchedule,
		Timeout:  durationpb.New(subjectCleanupCronTimeout),
	}); err != nil {
		return err
	}

	return nil
}

func (h *Housekeeper) RegisterCronjobHandlers(hand *croner.Handlers) error {
	hand.Add(subjectCleanupCronName, func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, subjectCleanupCronName)
		defer span.End()

		var errs error
		if err := h.resolver.CleanupOrphanSubjects(ctx, h.resolver.db); err != nil {
			h.logger.Error("error during orphan subject cleanup", zap.Error(err))
			errs = multierr.Append(errs, fmt.Errorf("error during orphan subject cleanup. %w", err))
		}

		if err := h.resolver.CleanupStaleJobGradeSubjects(ctx, h.resolver.db); err != nil {
			h.logger.Error("error during stale job grade subject cleanup", zap.Error(err))
			if errs == nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf("error during stale job grade subject cleanup. %w", err),
				)
			}
		}

		return errs
	})

	return nil
}
