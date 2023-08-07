package audit

import (
	"context"
	"database/sql"
	"sync"

	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jsoniter "github.com/json-iterator/go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	audit = table.FivenetAuditLog
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type IAuditer interface {
	Log(in *model.FivenetAuditLog, data any)
}

type AuditStorer struct {
	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB
	ctx    context.Context
	wg     sync.WaitGroup
	input  chan *model.FivenetAuditLog
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	DB     *sql.DB
}

func New(p Params) IAuditer {
	ctx, cancel := context.WithCancel(context.Background())

	a := &AuditStorer{
		logger: p.Logger,
		tracer: p.TP.Tracer("audit-storer"),
		db:     p.DB,
		ctx:    ctx,
		wg:     sync.WaitGroup{},
		input:  make(chan *model.FivenetAuditLog),
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		for i := 0; i < 4; i++ {
			a.wg.Add(1)
			go a.worker()
		}
		return nil
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) {
		close(a.input)
		cancel()
		a.wg.Wait()
	}))

	return a
}

func (a *AuditStorer) worker() {
	defer a.wg.Done()
	for {
		select {
		case in := <-a.input:
			if err := a.store(in); err != nil {
				a.logger.Error("failed to store audit log", zap.Error(err))
				continue
			}
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *AuditStorer) Log(in *model.FivenetAuditLog, data any) {
	if in == nil {
		return
	}

	in.Data = a.toJson(data)
	a.input <- in
}

func (a *AuditStorer) store(in *model.FivenetAuditLog) error {
	ctx, span := a.tracer.Start(a.ctx, "audit-log-store")
	defer span.End()

	stmt := audit.
		INSERT(
			audit.UserID,
			audit.UserJob,
			audit.TargetUserID,
			audit.Service,
			audit.Method,
			audit.State,
			audit.Data,
		).
		MODEL(in)

	if _, err := stmt.ExecContext(ctx, a.db); err != nil {
		return err
	}

	return nil
}

func (a *AuditStorer) toJson(data any) *string {
	if data != nil {
		data, err := json.MarshalToString(data)
		if err != nil {
			data = "Failed to marshal data"
		}
		return &data
	}

	noData := "No Data"
	return &noData
}
