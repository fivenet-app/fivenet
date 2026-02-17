package vehicles

import (
	context "context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

	db     *sql.DB
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

		db:     p.DB,
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

	tVehicleProps := table.FivenetVehiclesProps
	stmt := tVehicleProps.
		SELECT(
			tVehicleProps.Plate,
		).
		FROM(tVehicleProps).
		WHERE(mysql.AND(
			tVehicleProps.Wanted.IS_TRUE(),
			mysql.OR(
				tVehicleProps.WantedAt.LT(
					mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(maxDays, "DAY")),
				),
				tVehicleProps.WantedTill.LT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		LIMIT(100)

	var dest []string
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, err
	}

	for _, plate := range dest {
		props := &vehiclesprops.VehicleProps{}
		if err := props.LoadFromDB(ctx, s.db, plate); err != nil {
			s.logger.Error(
				"error loading vehicle props for cleanup",
				zap.String("plate", plate),
				zap.Error(err),
			)
			continue
		}

		// Create a copy of the original props and modify the wanted status
		in := proto.Clone(props).(*vehiclesprops.VehicleProps)

		wanted := false
		in.Wanted = &wanted
		in.WantedReason = nil
		in.WantedAt = nil

		// Update vehicle props
		if err := props.HandleChanges(ctx, s.db, in); err != nil {
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
