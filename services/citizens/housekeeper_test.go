package citizens

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
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
		store:  citizensstore.New(db, &config.CustomDB{}),
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
		store:  citizensstore.New(db, &config.CustomDB{}),
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
		MaxWantedDurationUser:        durationpb.New(72 * time.Hour), // 3 days
	}

	s := &Housekeeper{
		logger: zap.NewNop(),
		db:     db,
		appCfg: &mockAppConfig{cfg: cfg},
		store:  citizensstore.New(db, &config.CustomDB{}),
	}

	// Assert the key eligibility condition:
	// wanted IS TRUE AND (wanted_at < CURRENT_TIMESTAMP - INTERVAL maxDays DAY OR wanted_till < CURRENT_TIMESTAMP).
	expectedQuery := regexp.QuoteMeta(
		`SELECT fivenet_user_props.user_id AS "fivenet_user_props.user_id" FROM fivenet_user_props ` +
			`WHERE ( (fivenet_user_props.wanted IS TRUE) AND ` +
			`( (fivenet_user_props.wanted_at < (CURRENT_TIMESTAMP - INTERVAL 3 DAY)) ` +
			`OR (fivenet_user_props.wanted_till < CURRENT_TIMESTAMP) ) ) LIMIT ?;`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(int32(42)))

	// Retrieve user props for matched user
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_props AS user_props`)+`(?s).*`+regexp.QuoteMeta(`LEFT JOIN fivenet_files AS mugshot ON`)).
		WithArgs(int32(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_props.user_id",
			"user_props.updated_at",
			"user_props.wanted",
			"user_props.job",
			"user_props.job_grade",
			"user_props.traffic_infraction_points",
			"user_props.traffic_infraction_points_updated_at",
			"user_props.open_fines",
			"user_props.mugshot_file_id",
			"mugshot.mugshot_file_id",
			"file_path",
		}).AddRow(
			int32(42),
			nil,
			true,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_labels INNER JOIN fivenet_user_labels_job AS label ON`)+`(?s).*`).
		WithArgs(int32(42), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_user_labels_job.id", "fivenet_user_labels_job.job", "fivenet_user_labels_job.name", "fivenet_user_labels_job.color"}))

	// Make sure the query flips wanted to false for matched user.
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_props.*ON DUPLICATE KEY UPDATE.*wanted = \\?.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// User activity insert for wanted change.
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_activity`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, changedRows)
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
		store:  citizensstore.New(db, &config.CustomDB{}),
	}

	// Two eligible users are returned by the selection query
	mock.ExpectQuery(`(?s)FROM .*fivenet_user_props.*WHERE .*wanted.*wanted_at.*wanted_till.*LIMIT \?`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id"}).
				AddRow(int32(42)).
				AddRow(int32(43)),
		)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_props AS user_props`)+`(?s).*`+regexp.QuoteMeta(`LEFT JOIN fivenet_files AS mugshot ON`)).
		WithArgs(int32(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_props.user_id",
			"user_props.updated_at",
			"user_props.wanted",
			"user_props.job",
			"user_props.job_grade",
			"user_props.traffic_infraction_points",
			"user_props.traffic_infraction_points_updated_at",
			"user_props.open_fines",
			"user_props.mugshot_file_id",
			"mugshot.mugshot_file_id",
			"file_path",
		}).AddRow(
			int32(42),
			nil,
			true,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_labels INNER JOIN fivenet_user_labels_job AS label ON`)+`(?s).*`).
		WithArgs(int32(42), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_user_labels_job.id", "fivenet_user_labels_job.job", "fivenet_user_labels_job.name", "fivenet_user_labels_job.color"}))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_props.*ON DUPLICATE KEY UPDATE.*wanted = \\?.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_activity`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_props AS user_props`)+`(?s).*`+regexp.QuoteMeta(`LEFT JOIN fivenet_files AS mugshot ON`)).
		WithArgs(int32(43), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_props.user_id",
			"user_props.updated_at",
			"user_props.wanted",
			"user_props.job",
			"user_props.job_grade",
			"user_props.traffic_infraction_points",
			"user_props.traffic_infraction_points_updated_at",
			"user_props.open_fines",
			"user_props.mugshot_file_id",
			"mugshot.mugshot_file_id",
			"file_path",
		}).AddRow(
			int32(43),
			nil,
			true,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_labels INNER JOIN fivenet_user_labels_job AS label ON`)+`(?s).*`).
		WithArgs(int32(43), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_user_labels_job.id", "fivenet_user_labels_job.job", "fivenet_user_labels_job.name", "fivenet_user_labels_job.color"}))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_props.*ON DUPLICATE KEY UPDATE.*wanted = \\?.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_activity`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	changedRows, err := s.maxWantedDurationHandling(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, changedRows)
	require.NoError(t, mock.ExpectationsWereMet())
}
