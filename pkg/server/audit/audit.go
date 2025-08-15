package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// tAudit is a reference to the audit log table in the database.
var tAudit = table.FivenetAuditLog

// Default values for audit system.
const (
	bufferSize  = 128
	workerCount = 4
)

// Module provides the fx module for audit logging.
var Module = fx.Module("audit",
	fx.Provide(
		New,
	),
)

// FilterFn is a callback function type for filtering or modifying audit entries before logging.
type FilterFn func(in *audit.AuditEntry, data any)

// IAuditer defines the interface for logging audit entries.
type IAuditer interface {
	// Log records an audit entry with optional data and filter callbacks.
	Log(in *audit.AuditEntry, data any, callbacks ...FilterFn)
}

// AuditStorer implements IAuditer and manages asynchronous audit log storage.
type AuditStorer struct {
	// logger is used for logging errors and information.
	logger *zap.Logger
	// tracer is used for tracing audit log operations.
	tracer trace.Tracer
	// db is the database connection for storing audit logs.
	db *sql.DB
	// wg is a wait group for managing worker goroutines.
	wg sync.WaitGroup
	// input is the channel for incoming audit entries to be processed.
	input chan *audit.AuditEntry
}

// Params contains dependencies for constructing an AuditStorer instance.
type Params struct {
	fx.In

	// LC is the application lifecycle for registering hooks.
	LC fx.Lifecycle

	// Logger is the logger instance for logging.
	Logger *zap.Logger
	// TP is the tracer provider for distributed tracing.
	TP *tracesdk.TracerProvider
	// DB is the database connection.
	DB *sql.DB
	// Config is the application configuration.
	Config *config.Config
}

// New creates a new AuditStorer, registers the audit table with the housekeeper, and starts worker goroutines.
func New(p Params) IAuditer {
	ctxCancel, cancel := context.WithCancel(context.Background())

	a := &AuditStorer{
		logger: p.Logger.Named("audit"),
		tracer: p.TP.Tracer("audit"),
		db:     p.DB,
		wg:     sync.WaitGroup{},
		input:  make(chan *audit.AuditEntry, bufferSize),
	}

	// Register audit log table in housekeeper for retention management.
	housekeeper.AddTable(&housekeeper.Table{
		Table:           tAudit,
		IDColumn:        tAudit.ID,
		DeletedAtColumn: tAudit.CreatedAt,

		MinDays: p.Config.Audit.RetentionDays,
	})

	// Start worker goroutines for processing audit entries.
	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		for range workerCount {
			a.wg.Add(1)
			go a.worker(ctxCancel)
		}
		return nil
	}))
	// Stop workers and wait for completion on shutdown.
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		close(a.input)
		a.wg.Wait()
		return nil
	}))

	return a
}

// worker processes audit entries from the input channel and stores them in the database.
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

// Log records an audit entry, applies filter callbacks, serializes data, and queues it for storage.
func (a *AuditStorer) Log(in *audit.AuditEntry, data any, callbacks ...FilterFn) {
	if in == nil {
		return
	}

	for _, fn := range callbacks {
		fn(in, data)
	}
	in.Data = a.toJson(data)

	// Prevent panic if channel is closed
	select {
	case a.input <- in:
		// sent successfully

	default:
		// channel full, drop or log warning
		a.logger.Warn("audit log channel full, dropping entry")
	}
}

// store saves an audit entry to the database, extracting the service name and tracing the operation.
func (a *AuditStorer) store(ctx context.Context, in *audit.AuditEntry) error {
	if in == nil {
		return nil
	}

	ctx, span := a.tracer.Start(ctx, "audit-store")
	defer span.End()

	// Remove everything but the last part of the GRPC service name
	// E.g., `service.centrum.CentrumService` becomes `CentrumService`
	service := strings.Split(in.GetService(), ".")
	in.Service = service[len(service)-1]

	stmt := tAudit.
		INSERT(
			tAudit.UserID,
			tAudit.UserJob,
			tAudit.TargetUserID,
			tAudit.TargetUserJob,
			tAudit.Service,
			tAudit.Method,
			tAudit.State,
			tAudit.Data,
		).
		VALUES(
			in.UserId,
			in.GetUserJob(),
			in.TargetUserId,
			in.TargetUserJob,
			in.GetService(),
			in.GetMethod(),
			in.GetState(),
			in.GetData(),
		)

	if _, err := stmt.ExecContext(ctx, a.db); err != nil {
		return err
	}

	return nil
}

// toJson serializes the provided data to a JSON string pointer for storage in the audit log.
func (a *AuditStorer) toJson(data any) *string {
	if data == nil {
		noData := "No Data"
		return &noData
	}

	outB, err := json.Marshal(data)
	if err != nil {
		errStr := "Failed to marshal data"
		return &errStr
	}
	out := string(outB)
	return &out
}
