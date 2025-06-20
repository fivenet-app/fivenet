package livemap

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	pblivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCentrumMarkers,
		IDColumn:        table.FivenetCentrumMarkers.ID,
		DeletedAtColumn: table.FivenetCentrumMarkers.DeletedAt,
		JobColumn:       table.FivenetCentrumMarkers.Job,

		MinDays: 7,
	})
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCentrumMarkers,
		IDColumn:        table.FivenetCentrumMarkers.ID,
		TimestampColumn: table.FivenetCentrumMarkers.ExpiresAt,
		DeletedAtColumn: table.FivenetCentrumMarkers.DeletedAt,
		JobColumn:       table.FivenetCentrumMarkers.Job,

		MinDays: 5,
	})

	// User locations - Make sure to delete them after 1 day when not updated (buggy events from server)
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCentrumUserLocations,
		DeletedAtColumn: table.FivenetCentrumUserLocations.UpdatedAt,
		JobColumn:       table.FivenetCentrumUserLocations.Job,

		MinDays: 1,
	})
}

type Server struct {
	pblivemap.LivemapServiceServer

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

	markersCache        *xsync.Map[string, []*livemap.MarkerMarker]
	markersDeletedCache *xsync.Map[string, []uint64]

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
	MarkerUpdate *livemap.MarkerMarker
	MarkerDelete *uint64
}

func NewServer(p Params) *Server {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Server{
		logger: p.Logger,

		tracer:   p.TP.Tracer("livemap"),
		db:       p.DB,
		js:       p.JS,
		ps:       p.Perms,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		aud:      p.Audit,
		appCfg:   p.AppConfig,

		markersCache:        xsync.NewMap[string, []*livemap.MarkerMarker](),
		markersDeletedCache: xsync.NewMap[string, []uint64](),

		broker: broker.New[*brokerEvent](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		go s.broker.Start(ctxCancel)

		if err := s.registerStreamAndConsumer(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to register subscriptions. %w", err)
		}

		go func() {
			for {
				if err := s.refreshData(ctxCancel); err != nil {
					s.logger.Error("failed periodic livemap marker refresh", zap.Error(err))
				}

				select {
				case <-ctxCancel.Done():
					return

				case <-time.After(30 * time.Second):
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
	pblivemap.RegisterLivemapServiceServer(srv, s)
}

func (s *Server) refreshData(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "livemap.refresh-cache")
	defer span.End()

	if err := s.refreshMarkers(ctx); err != nil {
		return err
	}

	return nil
}
