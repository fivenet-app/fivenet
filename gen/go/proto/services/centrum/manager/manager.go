package manager

import (
	"context"
	"database/sql"

	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	"github.com/galexrt/fivenet/pkg/config/appconfig"
	"github.com/galexrt/fivenet/pkg/coords/postals"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/nats-io/nats.go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("centrum_manager", fx.Provide(
	New,
))

type Manager struct {
	ctx    context.Context
	logger *zap.Logger

	tracer   trace.Tracer
	db       *sql.DB
	js       nats.JetStreamContext
	enricher *mstlystcdata.Enricher
	tracker  tracker.ITracker
	postals  postals.Postals

	appCfg appconfig.IConfig

	*state.State
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	JS        nats.JetStreamContext
	Enricher  *mstlystcdata.Enricher
	Postals   postals.Postals
	Tracker   tracker.ITracker
	AppConfig appconfig.IConfig

	State *state.State
}

func New(p Params) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Manager{
		ctx:    ctx,
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

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := s.loadData(); err != nil {
			return err
		}

		if err := s.registerSubscriptions(ctx); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s
}
