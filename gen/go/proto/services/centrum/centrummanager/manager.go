package centrummanager

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrumstate"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/tracker"
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
