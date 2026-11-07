package citizens

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
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

func expectCanonicalWantedProps(mock sqlmock.Sqlmock, userID int32) {
	mock.ExpectQuery(`FROM fivenet_user_property_values`).WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_user_property_values.user_id", "fivenet_user_property_values.property_key", "fivenet_user_property_values.value_type", "fivenet_user_property_values.value_string", "fivenet_user_property_values.value_integer", "fivenet_user_property_values.value_number", "fivenet_user_property_values.value_bool", "fivenet_user_property_values.value_timestamp", "fivenet_user_property_values.value_date", "fivenet_user_property_values.updated_at"}).
			AddRow(userID, "user.wanted", uint8(4), nil, nil, nil, true, nil, nil, time.Now()))
	mock.ExpectQuery(`FROM fivenet_user_property_state`).WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "property_key", "reason", "active_at", "expires_at", "updated_at"}))
	mock.ExpectQuery(`FROM fivenet_user_files AS mugshot_file`).
		WithArgs(userID, "mugshot", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"mugshot_file_id", "mugshot.mugshot_file_id", "file_path"}))
}

func expectLockedExpiredWanted(mock sqlmock.Sqlmock, userID int32) {
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user_property_values.*WHERE .*user_id = \?.*property_key = \?.*FOR UPDATE`).
		WithArgs(userID, "user.wanted", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"wanted"}).AddRow(true))
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user_property_state.*WHERE .*user_id = \?.*property_key = \?.*FOR UPDATE`).
		WithArgs(userID, "user.wanted", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"active_at", "expires_at"}).
			AddRow(time.Now().Add(-7*24*time.Hour), nil))
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
		MaxWantedDurationUserEnabled: false,
		MaxWantedDurationUser:        durationpb.New(24 * 3600 * 1e9),
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
		store: citizensstore.New(citizensstore.Params{
			DB:           db,
			CustomDB:     &config.CustomDB{},
			LabelsAccess: access.NewCitizenLabelsSubjectObjectAccess(db),
		}),
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
		MaxWantedDurationUserEnabled: true,
		MaxWantedDurationUser:        nil,
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
		store: citizensstore.New(citizensstore.Params{
			DB:           db,
			CustomDB:     &config.CustomDB{},
			LabelsAccess: access.NewCitizenLabelsSubjectObjectAccess(db),
		}),
	}

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
		MaxWantedDurationUserEnabled: true,
		MaxWantedDurationUser:        durationpb.New(90 * time.Minute),
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
		store: citizensstore.New(citizensstore.Params{
			DB:           db,
			CustomDB:     &config.CustomDB{},
			LabelsAccess: access.NewCitizenLabelsSubjectObjectAccess(db),
		}),
	}
	mock.ExpectQuery(`SELECT .*fivenet_property_system_overrides.*enabled.*WHERE .*scope.*property_key`).
		WithArgs("user", string(permscitizens.CitizensServiceListCitizensPropertiesPermValueUserWanted)).
		WillReturnRows(sqlmock.NewRows([]string{"enabled"}).AddRow(true))

	// Wanted's lifecycle is stored in property_state, while the boolean lives in
	// property_values.
	expectedQuery := `(?s)FROM fivenet_user_property_values AS value INNER JOIN fivenet_user_property_state AS state.*value\.property_key = \?.*value\.value_bool IS TRUE.*state\.active_at <= \(CURRENT_TIMESTAMP - INTERVAL 5400 SECOND\).*state\.expires_at <= CURRENT_TIMESTAMP.*ORDER BY.*CASE.*state\.expires_at.*state\.active_at.*LIMIT \?`
	mock.ExpectQuery(expectedQuery).
		WithArgs("user.wanted", uint8(4), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(int32(42)))

	mock.ExpectBegin()
	expectLockedExpiredWanted(mock, 42)

	// Retrieve canonical wanted value and its optional lifecycle/file relations.
	expectCanonicalWantedProps(mock, 42)

	// Make sure the query flips wanted to false for matched user.
	mock.ExpectExec(`(?s)INSERT INTO fivenet_user_property_values`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)DELETE FROM fivenet_user_property_state`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// User activity insert for wanted change.
	mock.ExpectExec(`(?s)INSERT.*fivenet_user_activity`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperClearExpiredWanted_SkipsRenewedCandidate(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	s := &Housekeeper{logger: zap.NewNop(), db: db}
	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user_property_values.*WHERE .*user_id = \?.*property_key = \?.*FOR UPDATE`).
		WithArgs(int32(42), "user.wanted", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"wanted"}).AddRow(true))
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_user_property_state.*WHERE .*user_id = \?.*property_key = \?.*FOR UPDATE`).
		WithArgs(int32(42), "user.wanted", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"active_at", "expires_at"}).
			AddRow(time.Now(), nil))
	mock.ExpectCommit()

	cleared, err := s.clearExpiredWanted(t.Context(), 42, time.Hour)
	require.NoError(t, err)
	assert.False(t, cleared)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHousekeeperMaxWantedDurationHandling_ResetMultipleUsers(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	cfg := &settings.AppConfig{}
	cfg.Default()
	cfg.Game = &settings.Game{
		MaxWantedDurationUserEnabled: true,
		MaxWantedDurationUser:        durationpb.New(72 * time.Hour),
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
		store: citizensstore.New(citizensstore.Params{
			DB:           db,
			CustomDB:     &config.CustomDB{},
			LabelsAccess: access.NewCitizenLabelsSubjectObjectAccess(db),
		}),
	}
	mock.ExpectQuery(`SELECT .*fivenet_property_system_overrides.*enabled.*WHERE .*scope.*property_key`).
		WithArgs("user", string(permscitizens.CitizensServiceListCitizensPropertiesPermValueUserWanted)).
		WillReturnRows(sqlmock.NewRows([]string{"enabled"}).AddRow(true))

	// Two eligible users are returned by the selection query
	mock.ExpectQuery(`(?s)FROM fivenet_user_property_values AS value INNER JOIN fivenet_user_property_state AS state.*value\.property_key = \?.*value\.value_bool IS TRUE.*state\.active_at <= \(CURRENT_TIMESTAMP - INTERVAL 259200 SECOND\).*state\.expires_at <= CURRENT_TIMESTAMP.*ORDER BY.*CASE.*state\.expires_at.*state\.active_at.*LIMIT \?`).
		WithArgs("user.wanted", uint8(4), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id"}).
				AddRow(int32(42)).
				AddRow(int32(43)),
		)

	mock.ExpectBegin()
	expectLockedExpiredWanted(mock, 42)
	expectCanonicalWantedProps(mock, 42)

	mock.ExpectExec(`(?s)INSERT INTO fivenet_user_property_values`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)DELETE FROM fivenet_user_property_state`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT.*fivenet_user_activity`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	expectLockedExpiredWanted(mock, 43)
	expectCanonicalWantedProps(mock, 43)

	mock.ExpectExec(`(?s)INSERT INTO fivenet_user_property_values`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)DELETE FROM fivenet_user_property_state`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT.*fivenet_user_activity`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}
