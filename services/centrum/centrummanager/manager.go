package centrummanager

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumbrokers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumstate"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("centrum_manager",
	fx.Provide(
		New,
	))

type Manager struct {
	logger *zap.Logger
	jsCons jetstream.ConsumeContext

	tracer   trace.Tracer
	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	tracker  tracker.ITracker
	postals  postals.Postals
	brokers  *centrumbrokers.Brokers

	appCfg appconfig.IConfig

	unitAccess *access.Grouped[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitUserAccess, *centrum.UnitUserAccess, centrum.UnitQualificationAccess, *centrum.UnitQualificationAccess, centrum.UnitAccessLevel]

	*centrumstate.State
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	JS        *events.JSWrapper
	Enricher  *mstlystcdata.Enricher
	Postals   postals.Postals
	Tracker   tracker.ITracker
	AppConfig appconfig.IConfig
	Brokers   *centrumbrokers.Brokers

	State *centrumstate.State
}

func New(p Params) *Manager {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Manager{
		logger: p.Logger.Named("centrum.state"),

		tracer:   p.TP.Tracer("centrum-manager"),
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		postals:  p.Postals,
		brokers:  p.Brokers,

		appCfg: p.AppConfig,

		unitAccess: access.NewGrouped[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitUserAccess](
			p.DB,
			table.FivenetCentrumUnits,
			&access.TargetTableColumns{
				ID:        table.FivenetCentrumUnits.ID,
				DeletedAt: table.FivenetCentrumUnits.DeletedAt,
			},
			access.NewJobs[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitAccessLevel](
				table.FivenetCentrumUnitsAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCentrumUnitsAccess.ID,
						TargetID: table.FivenetCentrumUnitsAccess.TargetID,
						Access:   table.FivenetCentrumUnitsAccess.Access,
					},
					Job:          table.FivenetCentrumUnitsAccess.Job,
					MinimumGrade: table.FivenetCentrumUnitsAccess.MinimumGrade,
				},
				table.FivenetCentrumUnitsAccess.AS("unit_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCentrumUnitsAccess.AS("unit_job_access").ID,
						TargetID: table.FivenetCentrumUnitsAccess.AS("unit_job_access").TargetID,
						Access:   table.FivenetCentrumUnitsAccess.AS("unit_job_access").Access,
					},
					Job:          table.FivenetCentrumUnitsAccess.AS("unit_job_access").Job,
					MinimumGrade: table.FivenetCentrumUnitsAccess.AS("unit_job_access").MinimumGrade,
				},
			),
			nil,
			access.NewQualifications[centrum.UnitQualificationAccess, *centrum.UnitQualificationAccess, centrum.UnitAccessLevel](
				table.FivenetCentrumUnitsAccess,
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCentrumUnitsAccess.ID,
						TargetID: table.FivenetCentrumUnitsAccess.TargetID,
						Access:   table.FivenetCentrumUnitsAccess.Access,
					},
					QualificationId: table.FivenetCentrumUnitsAccess.QualificationID,
				},
				table.FivenetCentrumUnitsAccess.AS("unit_qualification_access"),
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCentrumUnitsAccess.AS("unit_qualification_access").ID,
						TargetID: table.FivenetCentrumUnitsAccess.AS("unit_qualification_access").TargetID,
						Access:   table.FivenetCentrumUnitsAccess.AS("unit_qualification_access").Access,
					},
					QualificationId: table.FivenetCentrumUnitsAccess.AS("unit_qualification_access").QualificationID,
				},
			),
		),

		State: p.State,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if _, err := s.registerStream(ctxStartup); err != nil {
			return err
		}

		go func() {
			if err := s.loadData(ctxCancel); err != nil {
				s.logger.Error("failed to load centrum data", zap.Error(err))
				return
			}
		}()

		if err := s.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return err
		}

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

func (s *Manager) GetUnitAccess() *access.Grouped[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitUserAccess, *centrum.UnitUserAccess, centrum.UnitQualificationAccess, *centrum.UnitQualificationAccess, centrum.UnitAccessLevel] {
	return s.unitAccess
}
