package docstore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/pkg/croner"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type WorkflowParams struct {
	fx.In

	LC fx.Lifecycle

	Logger       *zap.Logger
	DB           *sql.DB
	TP           *tracesdk.TracerProvider
	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

type Workflow struct {
	logger *zap.Logger
	tracer trace.Tracer

	db *sql.DB
}

func NewWorkflow(p WorkflowParams) *Workflow {
	w := &Workflow{
		logger: p.Logger,
		tracer: p.TP.Tracer("docstore_workflow"),
		db:     p.DB,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "docstore.workflow_run",
			Schedule: "@always", // Every minute
		}); err != nil {
			return err
		}

		return nil
	}))

	p.CronHandlers.Add("docstore.workflow_run", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := w.tracer.Start(ctx, "docstore.workflow_run")
		defer span.End()

		dest := &documents.CronData{}
		if err := anypb.UnmarshalTo(data.Data, dest, proto.UnmarshalOptions{}); err != nil {
			return fmt.Errorf("failed to unmarshal document workflow cron data. %w", err)
		}

		if err := w.handle(ctx, dest); err != nil {
			return fmt.Errorf("error during docstore workflow handling. %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal document workflow cron data. %w", err)
		}

		return nil
	})

	return w
}

func (w *Workflow) handle(ctx context.Context, data *documents.CronData) error {
	// TODO

	return nil
}
