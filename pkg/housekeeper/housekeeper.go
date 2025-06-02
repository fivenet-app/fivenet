package housekeeper

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

var Module = fx.Module("db_housekeeper",
	fx.Provide(
		New,
	),
)

const DefaultDeleteLimit = 500
const (
	lastJobName       = "last_job_name"
	lastTableMapIndex = "last_key"
)

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	db *sql.DB

	getTablesListFn func() map[string]*Table

	dryRun bool
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider
}

type Result struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func New(p Params) Result {
	h := &Housekeeper{
		logger: p.Logger.Named("housekeeper"),
		tracer: p.TP.Tracer("housekeeper"),

		db: p.DB,

		getTablesListFn: func() map[string]*Table {
			tableListsMu.Lock()
			defer tableListsMu.Unlock()

			return tablesList
		},

		dryRun: false,
	}

	return Result{
		Housekeeper:  h,
		CronRegister: h,
	}
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	registry.UnregisterCronjob(ctx, "housekeeper.run")
	registry.UnregisterCronjob(ctx, "housekeeper.job_delete")

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "housekeeper.hard_delete",
		Schedule: "*/2 * * * *", // Every 2 minutes
		Timeout:  durationpb.New(50 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "housekeeper.soft_delete_job",
		Schedule: "*/2 * * * *", // Every 2 minutes
		Timeout:  durationpb.New(50 * time.Second),
	}); err != nil {
		return err
	}

	return nil
}

func (h *Housekeeper) RegisterCronjobHandlers(hand *croner.Handlers) error {
	hand.Add("housekeeper.soft_delete_job", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "housekeeper.soft_delete_job")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := data.Data.UnmarshalTo(dest); err != nil {
			h.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		if err := h.runJobSoftDelete(ctx, dest); err != nil {
			return fmt.Errorf("error during housekeeper (soft delete). %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper (soft delete) cron data. %w", err)
		}

		return nil
	})

	hand.Add("housekeeper.hard_delete", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "housekeeper.hard_delete")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := data.Data.UnmarshalTo(dest); err != nil {
			h.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		if err := h.runHardDelete(ctx, dest); err != nil {
			return fmt.Errorf("error during housekeeper (hard delete). %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper (hard delete) cron data. %w", err)
		}

		return nil
	})

	return nil
}
