package vehicles

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

type mockAppConfig struct {
	cfg *appconfig.Cfg
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
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
	}

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 0, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_NoDuration(t *testing.T) {
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
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
	}

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 0, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_QueryCondition(t *testing.T) {
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
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
	}

	// Assert the key eligibility condition:
	// wanted IS TRUE AND (wanted_at < CURRENT_TIMESTAMP - INTERVAL maxDays DAY OR wanted_till < CURRENT_TIMESTAMP).
	expectedQuery := regexp.QuoteMeta(
		`SELECT fivenet_vehicles_props.plate AS "fivenet_vehicles_props.plate" FROM fivenet_vehicles_props ` +
			`WHERE ( fivenet_vehicles_props.wanted IS TRUE AND ` +
			`( (fivenet_vehicles_props.wanted_at < (CURRENT_TIMESTAMP - INTERVAL 3 DAY)) ` +
			`OR (fivenet_vehicles_props.wanted_till < CURRENT_TIMESTAMP) ) ) LIMIT ?;`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"plate"}).AddRow("ABC DEF1"))

	// Retrieve user props for matched user
	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*plate.*LIMIT \?`).
		WithArgs("ABC DEF1", sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"plate", "wanted"}).
				AddRow("ABC DEF1", true),
		)

	// Make sure the query flips wanted to false for matched user.
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_vehicles_props.*ON DUPLICATE KEY UPDATE.*wanted = \\?.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_ResetMultipleVehicles(t *testing.T) {
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
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
	}

	// Two eligible users are returned by the selection query
	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*wanted.*wanted_at.*wanted_till.*LIMIT \?`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"plate"}).
				AddRow("ABC DEF1").
				AddRow("ABC XYZ2"),
		)

	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*plate.*LIMIT \?`).
		WithArgs("ABC DEF1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"plate", "wanted"}).AddRow("ABC DEF1", true))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_vehicles_props.*ON DUPLICATE KEY UPDATE.*wanted = \\?.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`(?s)FROM .*fivenet_vehicles_props.*WHERE .*plate.*LIMIT \?`).
		WithArgs("ABC XYZ2", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"plate", "wanted"}).AddRow("ABC XYZ2", true))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_vehicles_props.*ON DUPLICATE KEY UPDATE.*wanted = \\?.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}
