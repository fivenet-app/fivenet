package audit

import (
	"context"
	"database/sql"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/croner"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var RetentionModule = fx.Module("audit_retention",
	fx.Provide(
		NewRetention,
	),
)

type Retention struct {
	tracer trace.Tracer
	logger *zap.Logger
	db     *sql.DB

	auditRetentionDays *int
}

type RetentionParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	DB     *sql.DB
	Config *config.Config

	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

func NewRetention(p RetentionParams) *Retention {
	r := &Retention{
		logger: p.Logger.Named("audit_retention"),
		tracer: p.TP.Tracer("audit-retention"),
		db:     p.DB,

		auditRetentionDays: p.Config.Audit.RetentionDays,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		p.CronHandlers.Add("auditlog-retention", func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := r.tracer.Start(ctx, "audit-retention")
			defer span.End()

			if err := r.Run(ctx); err != nil {
				r.logger.Error("error during audit store cleanup", zap.Error(err))
				return err
			}

			return nil
		})

		if err := p.Cron.RegisterCronjob(ctxStartup, &cron.Cronjob{
			Name:     "auditlog-retention",
			Schedule: "@30minutes",
		}); err != nil {
			return err
		}

		return nil
	}))

	return r
}

func (r *Retention) Run(ctx context.Context) error {
	if r.auditRetentionDays != nil {
		// Now minus retention days count
		t := time.Now().AddDate(0, 0, -*r.auditRetentionDays)

		if err := r.cleanup(ctx, t); err != nil {
			return err
		}
	}

	return nil
}

func (r *Retention) cleanup(ctx context.Context, before time.Time) error {
	r.logger.Debug("starting audit store cleanup", zap.Time("before_time", before))

	stmt := tAudit.
		DELETE().
		WHERE(tAudit.CreatedAt.LT_EQ(
			jet.TimestampT(before),
		)).
		LIMIT(5000)

	res, err := stmt.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	r.logger.Info("audit store cleanup completed", zap.Int64("affected_rows", affected))

	return nil
}
