package vehicles

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

type mockAppConfig struct {
	cfg *appconfig.Cfg
}

func expectExpiredWantedVehicles(
	mock sqlmock.Sqlmock,
	limit int64,
	plates ...string,
) {
	rows := sqlmock.NewRows([]string{"plate"})
	for _, plate := range plates {
		rows.AddRow(plate)
	}
	mock.ExpectQuery(`(?s)FROM fivenet_vehicle_property_values AS value.*INNER JOIN fivenet_vehicle_property_state AS state.*value.property_key.*value.value_type.*value.value_bool.*state.active_at.*state.expires_at.*LIMIT`).
		WithArgs("vehicle.wanted", 4, limit).
		WillReturnRows(rows)
}

func expectClearVehicleWanted(mock sqlmock.Sqlmock, plate string) {
	mock.ExpectQuery(`(?s)FROM fivenet_vehicle_property_values value.*LEFT JOIN fivenet_vehicle_property_state state.*value.property_key.*value.plate IN`).
		WithArgs("vehicle.wanted", plate).
		WillReturnRows(sqlmock.NewRows([]string{"plate", "updated_at", "value_bool", "reason", "active_at", "expires_at"}).
			AddRow(plate, time.Now(), true, nil, nil, nil))
	mock.ExpectBegin()
	mock.ExpectExec(`(?s)INSERT INTO fivenet_vehicle_property_values`).
		WithArgs(plate, "vehicle.wanted", 4, false).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)DELETE FROM fivenet_vehicle_property_state`).
		WithArgs(plate, "vehicle.wanted").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT INTO fivenet_vehicle_activity`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func (m *mockAppConfig) Get() *appconfig.Cfg {
	return m.cfg
}

func (m *mockAppConfig) Set(val *appconfig.Cfg) {
	m.cfg = val
}

func (m *mockAppConfig) Update(_ context.Context, val *appconfig.Cfg) error {
	m.cfg = val
	return nil
}

func (m *mockAppConfig) Reload(_ context.Context) (*appconfig.Cfg, error) {
	return m.cfg, nil
}

func (m *mockAppConfig) Subscribe() chan *appconfig.Cfg {
	return make(chan *appconfig.Cfg)
}

func (m *mockAppConfig) Unsubscribe(_ chan *appconfig.Cfg) {}

func TestHousekeeperMaxWantedDurationHandling_Disabled(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	cfg := &settings.AppConfig{}
	cfg.Default()
	cfg.Game = &settings.Game{
		MaxWantedDurationVehicleEnabled: false,
		MaxWantedDurationVehicle:        durationpb.New(24 * 3600 * 1e9),
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		store:  vehiclesstore.New(db, &config.CustomDB{}, nil),
		appCfg: &mockAppConfig{cfg: cfg},
	}

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 0, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_NoDuration(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	cfg := &settings.AppConfig{}
	cfg.Default()
	cfg.Game = &settings.Game{
		MaxWantedDurationVehicleEnabled: true,
		MaxWantedDurationVehicle:        nil,
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		store:  vehiclesstore.New(db, &config.CustomDB{}, nil),
		appCfg: &mockAppConfig{cfg: cfg},
	}

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 0, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_SystemPropertyDisabled(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	cfg := &settings.AppConfig{}
	cfg.Default()
	cfg.Game = &settings.Game{
		MaxWantedDurationVehicleEnabled: true,
		MaxWantedDurationVehicle:        durationpb.New(24 * time.Hour),
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		store:  vehiclesstore.New(db, &config.CustomDB{}, nil),
		appCfg: &mockAppConfig{cfg: cfg},
	}
	mock.ExpectQuery("SELECT .*enabled").
		WithArgs("vehicle", string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted)).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_property_system_overrides.enabled"}).AddRow(false))

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 0, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_QueryCondition(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	cfg := &settings.AppConfig{}
	cfg.Default()
	cfg.Game = &settings.Game{
		MaxWantedDurationVehicleEnabled: true,
		MaxWantedDurationVehicle:        durationpb.New(72 * time.Hour), // 3 days
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		store:  vehiclesstore.New(db, &config.CustomDB{}, nil),
		appCfg: &mockAppConfig{cfg: cfg},
	}
	mock.ExpectQuery("SELECT .*enabled").
		WithArgs("vehicle", string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted)).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_property_system_overrides.enabled"}))

	// Canonical wanted values are eligible when their managed state is old or expired.
	expectExpiredWantedVehicles(mock, 100, "ABC DEF1")
	expectClearVehicleWanted(mock, "ABC DEF1")

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_ResetMultipleVehicles(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	cfg := &settings.AppConfig{}
	cfg.Default()
	cfg.Game = &settings.Game{
		MaxWantedDurationVehicleEnabled: true,
		MaxWantedDurationVehicle:        durationpb.New(72 * time.Hour),
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		store:  vehiclesstore.New(db, &config.CustomDB{}, nil),
		appCfg: &mockAppConfig{cfg: cfg},
	}
	mock.ExpectQuery("SELECT .*enabled").
		WithArgs("vehicle", string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted)).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_property_system_overrides.enabled"}))

	// Both vehicles take the same canonical value/state cleanup path.
	expectExpiredWantedVehicles(mock, 100, "ABC DEF1", "ABC XYZ2")
	expectClearVehicleWanted(mock, "ABC DEF1")
	expectClearVehicleWanted(mock, "ABC XYZ2")

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}
