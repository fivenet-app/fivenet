package audit

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jsoniter "github.com/json-iterator/go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var tAudit = table.FivenetAuditLog

var json = jsoniter.ConfigFastest

var Module = fx.Module("audit",
	fx.Provide(
		New,
	),
)

type FilterFn func(in *model.FivenetAuditLog, data any)

type IAuditer interface {
	Log(in *model.FivenetAuditLog, data any, callbacks ...FilterFn)
}

type AuditStorer struct {
	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB
	wg     sync.WaitGroup
	input  chan *model.FivenetAuditLog
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	DB     *sql.DB
	Config *config.Config

	Cron croner.ICron
}

func New(p Params) IAuditer {
	ctxCancel, cancel := context.WithCancel(context.Background())

	a := &AuditStorer{
		logger: p.Logger.Named("audit"),
		tracer: p.TP.Tracer("audit"),
		db:     p.DB,
		wg:     sync.WaitGroup{},
		input:  make(chan *model.FivenetAuditLog),
	}

	// Register audit log table in housekeeper
	housekeeper.AddTable(&housekeeper.Table{
		Table:           tAudit,
		TimestampColumn: tAudit.CreatedAt,
		MinDays:         p.Config.Audit.RetentionDays,
	})

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		for range 4 {
			a.wg.Add(1)
			go a.worker(ctxCancel)
		}

		// Remove legacy audit log retention cron job
		if err := p.Cron.UnregisterCronjob(ctxStartup, "auditlog-retention"); err != nil {
			return err
		}

		return nil
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		close(a.input)
		cancel()
		a.wg.Wait()

		return nil
	}))

	return a
}

func (a *AuditStorer) worker(ctx context.Context) {
	defer a.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case in := <-a.input:
			if err := a.store(ctx, in); err != nil {
				a.logger.Error("failed to store audit log", zap.Error(err))
				continue
			}
		}
	}
}

func (a *AuditStorer) Log(in *model.FivenetAuditLog, data any, callbacks ...FilterFn) {
	if in == nil {
		return
	}

	for _, fn := range callbacks {
		fn(in, data)
	}
	in.Data = a.toJson(data)

	a.input <- in
}

func (a *AuditStorer) store(ctx context.Context, in *model.FivenetAuditLog) error {
	if in == nil {
		return nil
	}

	ctx, span := a.tracer.Start(ctx, "audit-store")
	defer span.End()

	// Remove everything but the last part of the GRPC service name
	// E.g., `service.centrum.CentrumService` becomes `CentrumService`
	service := strings.Split(in.Service, ".")
	in.Service = service[len(service)-1]

	stmt := tAudit.
		INSERT(
			tAudit.UserID,
			tAudit.UserJob,
			tAudit.TargetUserID,
			tAudit.Service,
			tAudit.Method,
			tAudit.State,
			tAudit.Data,
		).
		MODEL(in)

	if _, err := stmt.ExecContext(ctx, a.db); err != nil {
		return err
	}

	return nil
}

func (a *AuditStorer) toJson(data any) *string {
	if data == nil {
		noData := "No Data"
		return &noData
	}

	out, err := json.MarshalToString(data)
	if err != nil {
		data = "Failed to marshal data"
	}
	return &out
}
