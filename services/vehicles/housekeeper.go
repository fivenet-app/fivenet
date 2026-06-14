package vehicles

import (
	context "context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var HousekeeperModule = fx.Module(
	"vehicles.housekeeper",
	fx.Provide(
		NewHousekeeper,
	),
)

const changedRowsAttributeKey = "changed_rows"

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	store  *vehiclesstore.Store
	appCfg appconfig.IConfig
}

type HousekeeperParams struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	TP        *tracesdk.TracerProvider
	AppConfig appconfig.IConfig
}

type HousekeeperResult struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewHousekeeper(p HousekeeperParams) HousekeeperResult {
	s := &Housekeeper{
		logger: p.Logger.Named("vehicles.housekeeper"),
		tracer: p.TP.Tracer("vehicles.housekeeper"),

		store:  vehiclesstore.New(p.DB, &config.CustomDB{}),
		appCfg: p.AppConfig,
	}

	return HousekeeperResult{
		Housekeeper:  s,
		CronRegister: s,
	}
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "vehicles.props.max_wanted_duration",
		Schedule: "*/2 * * * *", // Every two minutes
	}); err != nil {
		return err
	}

	return nil
}

func (s *Housekeeper) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add(
		"vehicles.props.max_wanted_duration",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "vehicles.props.max_wanted_duration")
			defer span.End()

			dest := &cron.GenericCronData{
				Attributes: map[string]string{},
			}
			if err := data.Unmarshal(dest); err != nil {
				s.logger.Warn(
					"failed to unmarshal vehicles props cleanup cron data",
					zap.Error(err),
				)
			}

			changedRows, err := s.maxWantedDurationHandling(ctx)
			if err != nil {
				s.logger.Error("error during vehicles props cleanup", zap.Error(err))
				return err
			}

			dest.SetAttribute(changedRowsAttributeKey, strconv.Itoa(changedRows))

			// Marshal the updated cron data
			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf(
					"failed to marshal updated vehicles props cleanup cron data. %w",
					err,
				)
			}

			return nil
		},
	)

	return nil
}

func (s *Housekeeper) maxWantedDurationHandling(ctx context.Context) (int, error) {
	game := s.appCfg.Get().GetGame()

	// Skip if the max wanted duration feature is disabled or not properly configured
	if !game.GetMaxWantedDurationVehicleEnabled() || game.GetMaxWantedDurationVehicle() == nil {
		return 0, nil
	}

	maxDays := game.GetMaxWantedDurationVehicle().GetSeconds() / 24 / 3600

	dest, err := s.store.ListExpiredWanted(ctx, maxDays, 100)
	if err != nil {
		return 0, err
	}

	for _, plate := range dest {
		if err := s.store.ClearWanted(ctx, plate); err != nil {
			s.logger.Error(
				"error handling vehicle props changes during cleanup",
				zap.String("plate", plate),
				zap.Error(err),
			)
			continue
		}
	}

	return len(dest), nil
}
