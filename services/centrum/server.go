package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	pbcentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/services/centrum/centrummanager"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pbcentrum.CentrumServiceServer

	logger *zap.Logger
	wg     sync.WaitGroup
	jsCons jetstream.ConsumeContext

	tracer   trace.Tracer
	db       *sql.DB
	ps       perms.Permissions
	aud      audit.IAuditer
	js       *events.JSWrapper
	tracker  tracker.ITracker
	postals  postals.Postals
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.UserAwareEnricher

	brokersMutex sync.RWMutex
	brokers      map[string]*broker.Broker[*pbcentrum.StreamResponse]

	state *centrummanager.Manager
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	Perms     perms.Permissions
	Audit     audit.IAuditer
	JS        *events.JSWrapper
	Config    *config.Config
	AppConfig appconfig.IConfig
	Tracker   tracker.ITracker
	Postals   postals.Postals
	Manager   *centrummanager.Manager
	Enricher  *mstlystcdata.UserAwareEnricher
}

func NewServer(p Params) (*Server, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	brokers := map[string]*broker.Broker[*pbcentrum.StreamResponse]{}

	s := &Server{
		logger: p.Logger.Named("centrum"),
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		db:       p.DB,
		ps:       p.Perms,
		aud:      p.Audit,
		js:       p.JS,
		tracker:  p.Tracker,
		postals:  p.Postals,
		appCfg:   p.AppConfig,
		enricher: p.Enricher,

		brokersMutex: sync.RWMutex{},
		brokers:      brokers,

		state: p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		s.handleAppConfigUpdate(ctxCancel, p.AppConfig.Get())

		if err := s.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to subscribe to events: %w", err)
		}

		// Handle app config updates
		go func() {
			configUpdateCh := p.AppConfig.Subscribe()
			for {
				select {
				case <-ctxCancel.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					if cfg == nil {
						continue
					}
					s.handleAppConfigUpdate(ctxCancel, cfg)
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		if s.jsCons != nil {
			s.jsCons.Stop()
			s.jsCons = nil
		}

		return nil
	}))

	return s, nil
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcentrum.RegisterCentrumServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbcentrum.PermsRemap
}
