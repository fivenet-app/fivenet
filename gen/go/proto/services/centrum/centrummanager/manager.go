package centrummanager

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrumstate"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
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
		postals:  p.Postals,
		tracker:  p.Tracker,

		appCfg: p.AppConfig,

		unitAccess: access.NewGrouped[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitUserAccess](
			p.DB,
			table.FivenetCentrumUnits,
			&access.TargetTableColumns{
				ID:        table.FivenetCentrumUnits.ID,
				DeletedAt: table.FivenetCentrumUnits.DeletedAt,
			},
			access.NewJobs[centrum.UnitJobAccess, *centrum.UnitJobAccess, centrum.UnitAccessLevel](
				table.FivenetCentrumUnitsJobAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCentrumUnitsJobAccess.ID,
						CreatedAt: table.FivenetCentrumUnitsJobAccess.CreatedAt,
						TargetID:  table.FivenetCentrumUnitsJobAccess.UnitID,
						Access:    table.FivenetCentrumUnitsJobAccess.Access,
					},
					Job:          table.FivenetCentrumUnitsJobAccess.Job,
					MinimumGrade: table.FivenetCentrumUnitsJobAccess.MinimumGrade,
				},
				table.FivenetCentrumUnitsJobAccess.AS("unit_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCentrumUnitsJobAccess.AS("unit_job_access").ID,
						CreatedAt: table.FivenetCentrumUnitsJobAccess.AS("unit_job_access").CreatedAt,
						TargetID:  table.FivenetCentrumUnitsJobAccess.AS("unit_job_access").UnitID,
						Access:    table.FivenetCentrumUnitsJobAccess.AS("unit_job_access").Access,
					},
					Job:          table.FivenetCentrumUnitsJobAccess.AS("unit_job_access").Job,
					MinimumGrade: table.FivenetCentrumUnitsJobAccess.AS("unit_job_access").MinimumGrade,
				},
			),
			nil,
			access.NewQualifications[centrum.UnitQualificationAccess, *centrum.UnitQualificationAccess, centrum.UnitAccessLevel](
				table.FivenetCentrumUnitsQualificationsAccess,
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCentrumUnitsQualificationsAccess.ID,
						CreatedAt: table.FivenetCentrumUnitsQualificationsAccess.CreatedAt,
						TargetID:  table.FivenetCentrumUnitsQualificationsAccess.UnitID,
						Access:    table.FivenetCentrumUnitsQualificationsAccess.Access,
					},
					QualificationId: table.FivenetCentrumUnitsQualificationsAccess.QualificationID,
				},
				table.FivenetCentrumUnitsQualificationsAccess.AS("unit_qualification_access"),
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCentrumUnitsQualificationsAccess.AS("unit_qualification_access").ID,
						CreatedAt: table.FivenetCentrumUnitsQualificationsAccess.AS("unit_qualification_access").CreatedAt,
						TargetID:  table.FivenetCentrumUnitsQualificationsAccess.AS("unit_qualification_access").UnitID,
						Access:    table.FivenetCentrumUnitsQualificationsAccess.AS("unit_qualification_access").Access,
					},
					QualificationId: table.FivenetCentrumUnitsQualificationsAccess.AS("unit_qualification_access").QualificationID,
				},
			),
		),

		State: p.State,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if _, err := s.registerStream(ctxStartup); err != nil {
			return err
		}

		if err := s.loadData(ctxStartup); err != nil {
			return err
		}

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
