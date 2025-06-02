package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumbrokers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrummanager"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCentrumSettings,
		JobColumn:       table.FivenetCentrumSettings.Job,
		DeletedAtColumn: table.FivenetCentrumSettings.DeletedAt,

		MinDays: 30,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCentrumUnits,
		IDColumn:        table.FivenetCentrumUnits.ID,
		JobColumn:       table.FivenetCentrumUnits.Job,
		DeletedAtColumn: table.FivenetCentrumUnits.DeletedAt,

		MinDays: 30,

		DependantTables: []*housekeeper.Table{
			{
				Table:      table.FivenetCentrumUnitsStatus,
				IDColumn:   table.FivenetCentrumUnitsStatus.ID,
				ForeignKey: table.FivenetCentrumUnitsStatus.UnitID,
			},
		},
	})
}

type Server struct {
	pbcentrum.CentrumServiceServer

	logger   *zap.Logger
	wg       sync.WaitGroup
	ctx      context.Context
	jsCons   jetstream.ConsumeContext
	db       *sql.DB
	ps       perms.Permissions
	aud      audit.IAuditer
	js       *events.JSWrapper
	tracker  tracker.ITracker
	postals  postals.Postals
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.UserAwareEnricher

	brokers *centrumbrokers.Brokers
	state   *centrummanager.Manager
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
	Brokers   *centrumbrokers.Brokers
}

func NewServer(p Params) (*Server, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Server{
		logger: p.Logger.Named("centrum"),
		wg:     sync.WaitGroup{},
		ctx:    ctxCancel,

		db:       p.DB,
		ps:       p.Perms,
		aud:      p.Audit,
		js:       p.JS,
		tracker:  p.Tracker,
		postals:  p.Postals,
		appCfg:   p.AppConfig,
		enricher: p.Enricher,

		brokers: p.Brokers,
		state:   p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := s.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to subscribe to events: %w", err)
		}

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
