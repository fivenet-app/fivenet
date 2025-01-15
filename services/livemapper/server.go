package livemapper

import (
	"context"
	"database/sql"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	pblivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pblivemapper.LivemapperServiceServer

	logger *zap.Logger
	jsCons jetstream.ConsumeContext

	tracer   trace.Tracer
	db       *sql.DB
	js       *events.JSWrapper
	ps       perms.Permissions
	enricher *mstlystcdata.Enricher
	tracker  tracker.ITracker
	aud      audit.IAuditer
	appCfg   appconfig.IConfig

	markersCache *xsync.MapOf[string, []*livemap.MarkerMarker]

	broker *broker.Broker[*brokerEvent]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	JS        *events.JSWrapper
	Perms     perms.Permissions
	Enricher  *mstlystcdata.Enricher
	Config    *config.Config
	Tracker   tracker.ITracker
	Audit     audit.IAuditer
	AppConfig appconfig.IConfig
}

type brokerEvent struct {
	Send events.Type
}

func NewServer(p Params) *Server {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Server{
		logger: p.Logger,

		tracer:   p.TP.Tracer("livemapper-cache"),
		db:       p.DB,
		js:       p.JS,
		ps:       p.Perms,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		aud:      p.Audit,
		appCfg:   p.AppConfig,

		markersCache: xsync.NewMapOf[string, []*livemap.MarkerMarker](),

		broker: broker.New[*brokerEvent](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		go s.broker.Start(ctxCancel)

		if err := s.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return err
		}

		go func() {
			for {
				select {
				case <-ctxCancel.Done():
					return

				case <-time.After(30 * time.Second):
					if err := s.refreshData(ctxCancel); err != nil {
						s.logger.Error("failed periodic livemap marker refresh", zap.Error(err))
					}
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		if s.jsCons != nil {
			s.jsCons.Stop()
			s.jsCons = nil
		}

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pblivemapper.RegisterLivemapperServiceServer(srv, s)
}

func (s *Server) refreshData(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "livemap-refresh-cache")
	defer span.End()

	if err := s.refreshMarkers(ctx); err != nil {
		return err
	}

	return nil
}
