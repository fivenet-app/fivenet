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
	"github.com/fivenet-app/fivenet/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCentrumMarkers,
		TimestampColumn: table.FivenetCentrumMarkers.ExpiresAt,
		MinDays:         3,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCentrumMarkers,
		TimestampColumn: table.FivenetCentrumMarkers.DeletedAt,
		MinDays:         7,
	})
}

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

	markersCache        *xsync.MapOf[string, []*livemap.MarkerMarker]
	markersDeletedCache *xsync.MapOf[string, []uint64]

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
	Users        *map[string][]*livemap.UserMarker
	MarkerUpdate *livemap.MarkerMarker
	MarkerDelete *uint64
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

		markersCache:        xsync.NewMapOf[string, []*livemap.MarkerMarker](),
		markersDeletedCache: xsync.NewMapOf[string, []uint64](),

		broker: broker.New[*brokerEvent](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		go s.broker.Start(ctxCancel)

		if err := s.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return err
		}

		go func() {
			updateCh := p.Tracker.Subscribe()
			defer p.Tracker.Unsubscribe(updateCh)

			for {
				select {
				case <-ctxCancel.Done():
					return

				case event := <-updateCh:
					if len(event.Removed) == 0 {
						continue
					}

					// Group removed user markers by job
					grouped := map[string][]*livemap.UserMarker{}
					for _, um := range event.Removed {
						if _, ok := grouped[um.Job]; !ok {
							grouped[um.Job] = []*livemap.UserMarker{}
						}

						grouped[um.Job] = append(grouped[um.Job], um)
					}

					s.broker.Publish(&brokerEvent{
						Users: &grouped,
					})
				}
			}
		}()

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
