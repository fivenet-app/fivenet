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
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrummanager"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

type Result struct {
	fx.Out

	Server       *Server
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewServer(p Params) (*Server, error) {
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

		state: p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if _, err := eventscentrum.RegisterStream(ctxStartup, s.js); err != nil {
			return fmt.Errorf("failed to register stream. %w", err)
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

func (s *Server) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "centrum.dispatch.heatmap",
		Schedule: "*/15 * * * *", // Every 15 minutes
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

func (s *Server) GetPermsRemap() map[string]string {
	return pbcentrum.PermsRemap
}
