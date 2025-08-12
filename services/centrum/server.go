package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	pkggrpc "github.com/fivenet-app/fivenet/v2025/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/units"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
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

	// Remove unit statuses after 14 days
	housekeeper.AddTable(&housekeeper.Table{
		Table:    table.FivenetCentrumUnitsStatus,
		IDColumn: table.FivenetCentrumUnitsStatus.ID,

		TimestampColumn: table.FivenetCentrumUnitsStatus.CreatedAt,

		MinDays: 14,
	})

	// Remove heatmaps of deleted jobs
	housekeeper.AddTable(&housekeeper.Table{
		Table:     table.FivenetCentrumDispatchesHeatmaps,
		JobColumn: table.FivenetCentrumDispatchesHeatmaps.Job,

		TimestampColumn: table.FivenetCentrumDispatchesHeatmaps.GeneratedAt,

		MinDays: 14,
	})
}

type Server struct {
	pbcentrum.CentrumServiceServer

	logger *zap.Logger
	tracer trace.Tracer
	wg     sync.WaitGroup
	ctx    context.Context
	jsCons jetstream.ConsumeContext

	db       *sql.DB
	ps       perms.Permissions
	aud      audit.IAuditer
	js       *events.JSWrapper
	tracker  tracker.ITracker
	postals  postals.Postals
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.UserAwareEnricher
	jobs     *mstlystcdata.Jobs

	helpers     *helpers.Helpers
	settings    *settings.SettingsDB
	dispatchers *dispatchers.DispatchersDB
	units       *units.UnitDB
	dispatches  *dispatches.DispatchDB
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
	Enricher  *mstlystcdata.UserAwareEnricher
	Jobs      *mstlystcdata.Jobs

	Helpers     *helpers.Helpers
	Settings    *settings.SettingsDB
	Dispatchers *dispatchers.DispatchersDB
	Units       *units.UnitDB
	Dispatches  *dispatches.DispatchDB
}

type Result struct {
	fx.Out

	Server       *Server
	Service      pkggrpc.Service     `group:"grpcservices"`
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewServer(p Params) (Result, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Server{
		logger: p.Logger.Named("centrum"),
		tracer: p.TP.Tracer("mstlystcdata-cache"),
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
		jobs:     p.Jobs,

		helpers:     p.Helpers,
		settings:    p.Settings,
		dispatchers: p.Dispatchers,
		units:       p.Units,
		dispatches:  p.Dispatches,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if _, err := eventscentrum.RegisterStream(ctxStartup, s.js); err != nil {
			return fmt.Errorf("failed to register stream. %w", err)
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			if err := s.loadData(ctxCancel); err != nil {
				s.logger.Error("failed to load initial centrum data", zap.Error(err))
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

	return Result{
		Server:       s,
		Service:      s,
		CronRegister: s,
	}, nil
}

func (s *Server) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.dispatch.heatmap",
		Schedule: "*/5 * * * *", // Every 5 minutes
		Timeout:  durationpb.New(3 * time.Minute),
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) RegisterCronjobHandlers(hand *croner.Handlers) error {
	hand.Add("centrum.dispatch.heatmap", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := s.tracer.Start(ctx, "centrum.dispatch.heatmap")
		defer span.End()

		if err := s.generateDispatchHeatmaps(ctx); err != nil {
			s.logger.Error("failed to generate centrum dispatch heatmaps for jobs", zap.Error(err))
			return err
		}

		return nil
	})

	return nil
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcentrum.RegisterCentrumServiceServer(srv, s)
}

// GetPermsRemap returns the permissions re-mapping for the services.
func (s *Server) GetPermsRemap() map[string]string {
	return pbcentrum.PermsRemap
}

func (s *Server) loadData(ctx context.Context) error {
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := s.settings.LoadFromDB(gctx, ""); err != nil {
			return fmt.Errorf("failed to load settings from DB. %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := s.dispatchers.LoadFromDB(ctx, ""); err != nil {
			return fmt.Errorf("failed to load dispatchers from DB. %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := s.units.LoadFromDB(ctx, 0); err != nil {
			return fmt.Errorf("failed to load units from DB. %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to load initial data. %w", err)
	}

	if _, err := s.dispatches.LoadFromDB(ctx, nil); err != nil {
		return fmt.Errorf("failed to load dispatches from DB. %w", err)
	}

	return nil
}
