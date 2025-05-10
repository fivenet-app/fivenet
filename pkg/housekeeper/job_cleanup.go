package housekeeper

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
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

type JobCleanup struct {
	logger *zap.Logger
	tracer trace.Tracer

	db *sql.DB

	getTablesListFn func() map[string]*JobTable
}

type JobCleanupParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider

	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

func NewJobCleanup(p JobCleanupParams) (*JobCleanup, error) {
	c := &JobCleanup{
		logger: p.Logger.Named("housekeeper"),
		tracer: p.TP.Tracer("housekeeper"),
		db:     p.DB,
		getTablesListFn: func() map[string]*JobTable {
			return fkTablesList
		},
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "housekeeper.run",
			Schedule: "*/5 * * * *", // Every 5 minutes
			Timeout:  durationpb.New(1 * time.Minute),
		}); err != nil {
			return err
		}

		return nil
	}))

	p.CronHandlers.Add("housekeeper.job_cleanup", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := c.tracer.Start(ctx, "housekeeper.run")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := data.Data.UnmarshalTo(dest); err != nil {
			c.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		if err := c.run(ctx, dest); err != nil {
			return fmt.Errorf("error during docstore workflow handling. %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper cron data. %w", err)
		}

		return nil
	})

	return c, nil
}

func (c *JobCleanup) run(ctx context.Context, data *cron.GenericCronData) error {
	fkTableListsMu.Lock()
	defer fkTableListsMu.Unlock()

	tablesList := c.getTablesListFn()

	keys := []string{}
	for key := range tablesList {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil
	}

	lastTblKey, ok := data.Attributes[lastTableMapIndex]
	if !ok {
		// Take first table
		lastTblKey = keys[0]
	} else {
		idx := slices.Index(keys, lastTblKey)
		if idx == -1 || len(keys) <= idx+1 {
			c.logger.Debug("last table key not found in keys, starting from the beginning again")
			lastTblKey = keys[0]
		} else {
			lastTblKey = keys[idx+1]
		}
	}

	tbl, ok := tablesList[lastTblKey]
	if !ok {
		return nil
	}

	_ = tbl
	// TODO create either UPDATE or DELETE query and execute it

	data.Attributes[lastTableMapIndex] = lastTblKey

	return nil
}
