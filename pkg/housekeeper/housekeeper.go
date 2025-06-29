// Package housekeeper provides background database cleanup and maintenance jobs.
package housekeeper

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Module provides the fx module for the housekeeper package.
var Module = fx.Module("db_housekeeper",
	fx.Provide(
		New,
	),
)

// DefaultDeleteLimit is the default limit for delete operations in a single run.
const DefaultDeleteLimit = 500

const (
	// lastJobName is the key for the last job name attribute.
	lastJobName = "last_job_name"
	// lastTableMapIndex is the key for the last table map index attribute.
	lastTableMapIndex = "last_key"
)

// Housekeeper is responsible for running background cleanup and maintenance jobs on the database.
type Housekeeper struct {
	// logger is the logger instance for housekeeper operations.
	logger *zap.Logger
	// tracer is the OpenTelemetry tracer for housekeeper spans.
	tracer trace.Tracer

	// db is the SQL database connection used for cleanup jobs.
	db *sql.DB

	// getTablesListFn returns the list of tables to be processed by the housekeeper.
	getTablesListFn func() map[string]*Table

	// dryRun, if true, prevents actual deletion and only simulates the job.
	dryRun bool
}

// Params defines the dependencies required to construct a Housekeeper.
type Params struct {
	fx.In

	// LC is the fx lifecycle for managing start/stop hooks.
	LC fx.Lifecycle

	// Logger is the application's logger.
	Logger *zap.Logger
	// DB is the SQL database connection.
	DB *sql.DB
	// TP is the OpenTelemetry tracer provider.
	TP *tracesdk.TracerProvider
	// Config for log level overrides.
	Cfg *config.Config
}

// Result defines the outputs provided by the Housekeeper module.
type Result struct {
	fx.Out

	// Housekeeper is the main housekeeper instance.
	Housekeeper *Housekeeper
	// CronRegister allows the housekeeper to register cronjobs with the croner system.
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

// New constructs a new Housekeeper and returns it as an fx result.
func New(p Params) Result {
	logger := p.Logger.WithOptions(zap.IncreaseLevel(p.Cfg.LogLevelOverrides.Get(config.LoggingComponentHousekeeper, p.Cfg.LogLevel))).
		Named("housekeeper")

	h := &Housekeeper{
		logger: logger,
		tracer: p.TP.Tracer("housekeeper"),

		// Assign the database connection.
		db: p.DB,

		// Provide a function to get the list of tables for housekeeping.
		getTablesListFn: func() map[string]*Table {
			tableListsMu.Lock()
			defer tableListsMu.Unlock()

			return tablesList
		},

		// By default, do not run in dry-run mode.
		dryRun: false,
	}

	return Result{
		Housekeeper:  h,
		CronRegister: h,
	}
}

// RegisterCronjobs unregisters old housekeeper jobs and registers new cronjobs for hard and soft delete operations.
func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	// Unregister legacy or obsolete cronjobs.
	registry.UnregisterCronjob(ctx, "housekeeper.run")
	registry.UnregisterCronjob(ctx, "housekeeper.job_delete")

	// Register the hard delete cronjob to run every 2 minutes.
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "housekeeper.hard_delete",
		Schedule: "*/2 * * * *", // Every 2 minutes
		Timeout:  durationpb.New(50 * time.Second),
	}); err != nil {
		return err
	}

	// Register the soft delete cronjob to run every 2 minutes.
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "housekeeper.soft_delete_job",
		Schedule: "*/2 * * * *", // Every 2 minutes
		Timeout:  durationpb.New(50 * time.Second),
	}); err != nil {
		return err
	}

	return nil
}

// RegisterCronjobHandlers registers the handlers for housekeeper cronjobs with the croner system.
func (h *Housekeeper) RegisterCronjobHandlers(hand *croner.Handlers) error {
	// Handler for the soft delete job.
	hand.Add("housekeeper.soft_delete_job", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "housekeeper.soft_delete_job")
		defer span.End()

		// Prepare the cron data structure.
		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data, _ = anypb.New(&cron.GenericCronData{})
		}
		if err := data.Data.UnmarshalTo(dest); err != nil {
			h.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		// Unmarshal the cron data
		if err := data.Data.UnmarshalTo(dest); err != nil {
			h.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		// Run the soft delete job logic
		if err := h.runJobSoftDelete(ctx, dest); err != nil {
			return fmt.Errorf("error during housekeeper (soft delete). %w", err)
		}

		// Marshal the updated cron data
		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper (soft delete) cron data. %w", err)
		}

		return nil
	})

	// Handler for the hard delete job.
	hand.Add("housekeeper.hard_delete", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "housekeeper.hard_delete")
		defer span.End()

		// Prepare the cron data structure
		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data, _ = anypb.New(&cron.GenericCronData{})
		}
		// Unmarshal the cron data
		if err := data.Data.UnmarshalTo(dest); err != nil {
			h.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		// Run the hard delete job logic
		if err := h.runHardDelete(ctx, dest); err != nil {
			return fmt.Errorf("error during housekeeper (hard delete). %w", err)
		}

		// Marshal the updated cron data
		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper (hard delete) cron data. %w", err)
		}

		return nil
	})

	return nil
}
