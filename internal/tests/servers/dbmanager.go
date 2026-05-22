package servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/internal/tests"
	"github.com/fivenet-app/fivenet/v2026/query"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type dbServer struct {
	t *testing.T

	db       *sql.DB
	pool     dockertest.Pool
	resource dockertest.Resource

	stopped bool
}

func NewDBServer(t *testing.T, setup bool) *dbServer {
	t.Helper()

	s := &dbServer{
		t: t,
	}

	if setup {
		s.Setup(t.Context())
	}

	return s
}

func (m *dbServer) Setup(ctx context.Context) {
	image, tag := loadDockerComposeServiceImage(m.t, "mysql")

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	m.pool = dockertest.NewPoolT(m.t, "")

	// Pulls image, creates a container based on it and runs it.
	// We disable reuse explicitly to guarantee test isolation.
	m.resource = m.pool.RunT(
		m.t,
		image,
		dockertest.WithTag(tag),
		dockertest.WithEnv([]string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_USER=fivenet",
			"MYSQL_PASSWORD=changeme",
			"MYSQL_DATABASE=fivenettest",
			"TZ=Europe/Berlin",
		}),
		dockertest.WithCmd([]string{
			"mysqld",
			"--innodb-ft-min-token-size=2",
			"--innodb-ft-max-token-size=50",
			"--default-time-zone=Europe/Berlin",
		}),
		dockertest.WithoutReuse(),
	)

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := m.pool.Retry(ctx, 0, func() error {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		db, err := sql.Open("mysql", m.getDSN())
		if err != nil {
			return fmt.Errorf("failed to open database connection. %w", err)
		}
		if err := db.PingContext(ctx); err != nil {
			return fmt.Errorf("failed to ping database. %w", err)
		}

		rows, err := db.QueryContext(ctx, "SELECT 1;")
		if err != nil {
			return fmt.Errorf("failed to execute test query on database. %w", err)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("error in rows. %w", err)
		}
		defer rows.Close()

		return nil
	}); err != nil {
		require.NoError(m.t, err, "could not connect to database")
	}

	require.NoError(m.t, m.prepareDBForFirstUse(ctx), "failed to prepare database for first use")

	require.NoError(m.t, m.LoadBaseData(ctx), "failed to load base data into database")

	var err error
	m.db, err = sql.Open("mysql", m.getDSN())
	require.NoError(m.t, err, "could not connect to database after setup")
}

func (m *dbServer) DB() (*sql.DB, error) {
	if m.db == nil {
		return nil, errors.New(
			"test DB connection has not been established! You are accessing DB() method too early",
		)
	}

	return m.db, nil
}

func (m *dbServer) getDSN() string {
	// Using `root` isn't cool, but a workaround for now to create triggers in the database
	return fmt.Sprintf(
		"root:secret@(127.0.0.1:%s)/fivenettest?collation=utf8mb4_unicode_ci&loc=Local&parseTime=true",
		m.resource.GetPort("3306/tcp"),
	)
}

func (m *dbServer) prepareDBForFirstUse(ctx context.Context) error {
	// Use DB migrations to handle the rest (esx compat mode is true)
	_, err := query.MigrateDB(ctx, zap.NewNop(), m.getDSN(), false, false)
	require.NoError(m.t, err, "failed to migrate test database")

	return nil
}

func (m *dbServer) getMultiStatementDB() (*sql.DB, error) {
	// Open db connection with multiStatements param so we can apply sql files
	initDB, err := sql.Open("mysql", m.getDSN()+"&parseTime=true&multiStatements=true")
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open test database connection for multi statement exec. %w",
			err,
		)
	}

	return initDB, nil
}

func (m *dbServer) loadSQLFile(ctx context.Context, file string) error {
	initDB, err := m.getMultiStatementDB()
	require.NoError(m.t, err, "failed to get multi statement db")

	c, ioErr := os.ReadFile(file)
	require.NoError(m.t, ioErr, "failed to read %s for tests", file)
	sqlBase := string(c)
	_, err = initDB.ExecContext(ctx, sqlBase)
	require.NoError(m.t, err, "failed to apply %s for tests", file)

	return nil
}

func (m *dbServer) LoadBaseData(ctx context.Context) error {
	path := filepath.Join(tests.TestDataSQLPath, "base_*.sql")
	files, err := filepath.Glob(path)
	require.NoError(m.t, err, "failed to find base data sql files (%s)", path)
	// Sort the found files as they might not be in lexical order which we
	// need for this case https://github.com/golang/go/issues/17153
	slices.Sort(files)

	for _, file := range files {
		if err := m.loadSQLFile(ctx, file); err != nil {
			return err
		}
	}

	return nil
}

func (m *dbServer) Stop() {
	if m.stopped {
		return
	}
	m.stopped = true

	if m.db != nil {
		require.NoError(m.t, m.db.Close(), "could not close test database connection")
	}
}

func (m *dbServer) FxProvide() fx.Option {
	return fx.Provide(func() *sql.DB {
		return m.db
	})
}
