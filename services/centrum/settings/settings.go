package settings

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SettingsDB struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher

	store *store.Store[centrum.Settings, *centrum.Settings]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	JS       *events.JSWrapper
	DB       *sql.DB
	Cfg      *config.Config
	Enricher *mstlystcdata.Enricher
}

func New(p Params) *SettingsDB {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.Named("centrum.settings")
	d := &SettingsDB{
		logger:   logger,
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		st, err := store.New[centrum.Settings, *centrum.Settings](
			ctxCancel,
			logger,
			p.JS,
			"centrum_settings",
		)
		if err != nil {
			return err
		}

		if err := st.Start(ctxCancel, false); err != nil {
			return err
		}
		d.store = st

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return d
}

func (s *SettingsDB) LoadFromDB(ctx context.Context, job string) error {
	tCentrumSettings := table.FivenetCentrumSettings.AS("settings")

	stmt := tCentrumSettings.
		SELECT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Type,
			tCentrumSettings.Public,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
			tCentrumSettings.PredefinedStatus,
			tCentrumSettings.Timings,
			tCentrumSettings.Access,
			tCentrumSettings.Configuration,
		).
		FROM(tCentrumSettings)

	if job != "" {
		stmt = stmt.
			WHERE(
				tCentrumSettings.Job.EQ(jet.String(job)),
			).
			LIMIT(1)
	}

	var dest []*centrum.Settings
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for _, settings := range dest {
		settings.Default(settings.GetJob())

		if err := s.updateInKV(ctx, settings.GetJob(), settings); err != nil {
			return err
		}
	}

	return nil
}

func (s *SettingsDB) updateDB(ctx context.Context, job string, settings *centrum.Settings) error {
	tCentrumSettings := table.FivenetCentrumSettings

	stmt := tCentrumSettings.
		INSERT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Type,
			tCentrumSettings.Public,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
			tCentrumSettings.PredefinedStatus,
			tCentrumSettings.Timings,
			tCentrumSettings.Access,
			tCentrumSettings.Configuration,
		).
		VALUES(
			job,
			settings.GetEnabled(),
			settings.GetType(),
			settings.GetPublic(),
			settings.GetMode(),
			settings.GetFallbackMode(),
			settings.GetPredefinedStatus(),
			settings.GetTimings(),
			settings.GetAccess(),
			settings.GetConfiguration(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCentrumSettings.Enabled.SET(jet.Bool(settings.GetEnabled())),
			tCentrumSettings.Type.SET(jet.Int32(int32(settings.GetType()))),
			tCentrumSettings.Public.SET(jet.Bool(settings.GetPublic())),
			tCentrumSettings.Mode.SET(jet.Int32(int32(settings.GetMode()))),
			tCentrumSettings.FallbackMode.SET(jet.Int32(int32(settings.GetFallbackMode()))),
			tCentrumSettings.PredefinedStatus.SET(jet.StringExp(jet.Raw("VALUES(`predefined_status`)"))),
			tCentrumSettings.Timings.SET(jet.StringExp(jet.Raw("VALUES(`timings`)"))),
			tCentrumSettings.Access.SET(jet.StringExp(jet.Raw("VALUES(`access`)"))),
			tCentrumSettings.Configuration.SET(jet.StringExp(jet.Raw("VALUES(`configuration`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	// Load settings from database so they are updated in the "cache"
	if err := s.LoadFromDB(ctx, job); err != nil {
		return err
	}

	return nil
}

func (s *SettingsDB) Update(
	ctx context.Context,
	job string,
	in *centrum.Settings,
) (*centrum.Settings, error) {
	current, err := s.Get(ctx, job)
	if err != nil {
		return nil, err
	}

	if in.Job == "" {
		in.Job = job
	}

	current.Merge(in)

	if err := s.updateDB(ctx, job, current); err != nil {
		return nil, err
	}

	return current, nil
}
