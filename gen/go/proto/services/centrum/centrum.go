package centrum

import (
	"context"
	"database/sql"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/nats-io/nats.go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	ErrFailedQuery = status.Error(codes.Internal, "errors.CentrumService.ErrFailedQuery")
)

type Server struct {
	CentrumServiceServer

	ctx    context.Context
	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB
	p      perms.Permissions
	a      audit.IAuditer

	events   *events.Eventus
	eventSub *nats.Subscription

	dispatches syncx.Map[string, *syncx.Map[uint64, *dispatch.Dispatch]]
	units      syncx.Map[string, *syncx.Map[uint64, *dispatch.Unit]]
	unitsUsers syncx.Map[int32, uint64]

	broker *utils.Broker[interface{}]
}

func NewServer(ctx context.Context, logger *zap.Logger, tp *tracesdk.TracerProvider, db *sql.DB, p perms.Permissions, aud audit.IAuditer, e *events.Eventus) *Server {
	broker := utils.NewBroker[interface{}](ctx)
	go broker.Start()

	return &Server{
		ctx:    ctx,
		logger: logger,

		tracer: tp.Tracer("centrum-cache"),

		db:     db,
		p:      p,
		a:      aud,
		events: e,

		dispatches: syncx.Map[string, *syncx.Map[uint64, *dispatch.Dispatch]]{},
		units:      syncx.Map[string, *syncx.Map[uint64, *dispatch.Unit]]{},
		unitsUsers: syncx.Map[int32, uint64]{},

		broker: broker,
	}
}

func (s *Server) Start() {
	if err := s.registerEvents(); err != nil {
		s.logger.Error("failed to register events", zap.Error(err))
	}

	go func() {
		for {
			if err := s.refresh(); err != nil {
				s.logger.Error("failed to refresh centrum data", zap.Error(err))
			}

			select {
			case <-s.ctx.Done():
				s.broker.Stop()
				return
			case <-time.After(10 * time.Second):
			}
		}
	}()
}

func (s *Server) refresh() error {
	ctx, span := s.tracer.Start(s.ctx, "livemap-refresh-cache")
	defer span.End()

	if err := s.loadDispatches(ctx); err != nil {
		s.logger.Error("failed to load dispatches", zap.Error(err))
	}
	if err := s.loadUnits(ctx); err != nil {
		s.logger.Error("failed to load units", zap.Error(err))
	}

	return nil
}
