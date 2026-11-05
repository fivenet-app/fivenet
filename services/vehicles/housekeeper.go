package vehicles

import (
	context "context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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

	db     *sql.DB
	store  vehiclesstore.IStore
	appCfg appconfig.IConfig
}

type HousekeeperParams struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	TP        *tracesdk.TracerProvider
	AppConfig appconfig.IConfig
	Hydrator  citizenshydrator.IHydrator
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
		store:  vehiclesstore.New(p.DB, &config.CustomDB{}, p.Hydrator),
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

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "vehicles.props.expire_custom_properties",
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

	h.Add(
		"vehicles.props.expire_custom_properties",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "vehicles.props.expire_custom_properties")
			defer span.End()

			dest := &cron.GenericCronData{Attributes: map[string]string{}}
			if err := data.Unmarshal(dest); err != nil {
				s.logger.Warn(
					"failed to unmarshal expired custom vehicle properties cron data",
					zap.Error(err),
				)
			}

			startedAt := time.Now()
			changedRows, err := s.expireCustomPropertyHandling(ctx)
			propertyregistry.ObserveExpiryCleanup(
				"vehicle",
				time.Since(startedAt),
				changedRows,
				err,
			)
			if err != nil {
				s.logger.Error(
					"error during expired custom vehicle properties cleanup",
					zap.Error(err),
				)
				return err
			}
			dest.SetAttribute(changedRowsAttributeKey, strconv.Itoa(changedRows))
			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf(
					"failed to marshal expired custom vehicle properties cron data. %w",
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
	wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, "vehicle.wanted")
	if err != nil {
		return 0, err
	}
	if !wantedEnabled {
		// Keep the stored state intact while the system property is disabled.
		return 0, nil
	}

	maxSeconds := game.GetMaxWantedDurationVehicle().GetSeconds()

	dest, err := s.store.ListExpiredWanted(ctx, maxSeconds, 100)
	if err != nil {
		return 0, err
	}

	changed := 0
	for _, plate := range dest {
		cleared, err := s.clearExpiredWanted(ctx, plate, time.Duration(maxSeconds)*time.Second)
		if err != nil {
			s.logger.Error(
				"error handling vehicle props changes during cleanup",
				zap.String("plate", plate),
				zap.Error(err),
			)
			continue
		}
		if cleared {
			changed++
		}
	}

	return changed, nil
}

// clearExpiredWanted locks the canonical wanted value and lifecycle state,
// then rechecks expiry inside the transaction. The initial candidate query is
// deliberately only a cheap work queue and must not decide whether to clear.
func (s *Housekeeper) clearExpiredWanted(
	ctx context.Context,
	plate string,
	maxDuration time.Duration,
) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	tValues := table.FivenetVehiclePropertyValues
	value := struct {
		Wanted *bool `alias:"wanted"`
	}{}
	err = tValues.SELECT(tValues.ValueBool.AS("wanted")).FROM(tValues).WHERE(mysql.AND(
		tValues.Plate.EQ(mysql.String(plate)),
		tValues.PropertyKey.EQ(mysql.String("vehicle.wanted")),
	)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &value)
	if errors.Is(err, qrm.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}

	tState := table.FivenetVehiclePropertyState
	state := struct {
		ActiveAt  *time.Time `alias:"active_at"`
		ExpiresAt *time.Time `alias:"expires_at"`
	}{}
	err = tState.SELECT(tState.ActiveAt, tState.ExpiresAt).FROM(tState).WHERE(mysql.AND(
		tState.Plate.EQ(mysql.String(plate)),
		tState.PropertyKey.EQ(mysql.String("vehicle.wanted")),
	)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &state)
	if errors.Is(err, qrm.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}

	now := time.Now()
	if value.Wanted == nil || !*value.Wanted ||
		!wantedExpired(state.ActiveAt, state.ExpiresAt, maxDuration, now) {
		return false, tx.Commit()
	}
	if err := s.store.ClearWantedTx(ctx, tx, plate); err != nil {
		return false, fmt.Errorf("clearing vehicle %q wanted state: %w", plate, err)
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func wantedExpired(activeAt, expiresAt *time.Time, maxDuration time.Duration, now time.Time) bool {
	return (activeAt != nil && !activeAt.Add(maxDuration).After(now)) ||
		(expiresAt != nil && !expiresAt.After(now))
}

// expireCustomPropertyHandling iterates lifecycle-enabled definitions first,
// keeping custom expiry isolated from the managed wanted cleanup above.
func (s *Housekeeper) expireCustomPropertyHandling(ctx context.Context) (int, error) {
	keys, err := s.store.ListExpiringCustomPropertyKeys(ctx)
	if err != nil {
		return 0, err
	}

	changed := 0
	for _, key := range keys {
		plates, err := s.store.ListExpiredCustomPropertyPlates(ctx, key, 100)
		if err != nil {
			return changed, err
		}
		for _, plate := range plates {
			cleared, err := s.store.ClearExpiredCustomProperty(ctx, plate, key)
			if err != nil {
				propertyregistry.RecordExpiryCleanupItemFailure("vehicle")
				s.logger.Error(
					"error clearing expired custom vehicle property",
					zap.String("plate", plate),
					zap.String("property_key", key),
					zap.Error(err),
				)
				continue
			}
			if cleared {
				changed++
			}
		}
	}
	return changed, nil
}
