package audit

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	jsoniter "github.com/json-iterator/go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tAudit = table.FivenetAuditLog
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
	Config *config.Config
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

		if p.Config.Game.AuditRetentionDays != nil {
			// Now minus retention days
			t := time.Now().AddDate(0, 0, -*p.Config.Game.AuditRetentionDays)
			if err := a.Cleanup(t); err != nil {
				return err
			}
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

func (a *AuditStorer) Cleanup(before time.Time) error {
	stmt := tAudit.
		DELETE().
		WHERE(tAudit.CreatedAt.LT_EQ(
			jet.TimestampT(before),
		))

	res, err := stmt.ExecContext(a.ctx, a.db)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	a.logger.Info("audit store cleanup completed", zap.Int64("affected_rows", affected))

	return nil
}

func (a *AuditStorer) Log(in *model.FivenetAuditLog, data any) {
	if in == nil {
		return
	}

	in.Data = a.toJson(data)
	a.input <- in
}

func (a *AuditStorer) store(in *model.FivenetAuditLog) error {
	if in == nil {
		return nil
	}

	ctx, span := a.tracer.Start(a.ctx, "audit-log-store")
	defer span.End()

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
