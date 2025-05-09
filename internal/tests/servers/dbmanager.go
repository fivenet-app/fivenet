package servers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/fivenet-app/fivenet/v2025/internal/tests"
	"github.com/fivenet-app/fivenet/v2025/query"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type dbServer struct {
	t *testing.T

	db       *sql.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource

	stopped bool
}

func NewDBServer(t *testing.T, setup bool) *dbServer {
	s := &dbServer{
		t: t,
	}

	if setup {
		s.Setup()
	}

	return s
}

func (m *dbServer) Setup() {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	var err error
	m.pool, err = dockertest.NewPool("")
	if err != nil {
		m.t.Fatalf("could not construct pool: %q", err)
	}

	// uses pool to try to connect to Docker
	err = m.pool.Client.Ping()
	if err != nil {
		m.t.Fatalf("could not connect to Docker: %q", err)
	}

	// pulls an image, creates a container based on it and runs it
	m.resource, err = m.pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "docker.io/library/mysql",
			Tag:        "9.3.0",
			Env: []string{
				"MYSQL_ROOT_PASSWORD=secret",
				"MYSQL_USER=fivenet",
				"MYSQL_PASSWORD=changeme",
				"MYSQL_DATABASE=fivenettest",
				"TZ=Europe/Berlin",
			},
			Cmd: []string{
				"mysqld",
				"--innodb-ft-min-token-size=2",
				"--innodb-ft-max-token-size=50",
				"--default-time-zone=Europe/Berlin",
			},
		},
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		m.t.Fatalf("could not start resource: %q", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := m.pool.Retry(func() error {
		db, err := sql.Open("mysql", m.getDSN())
		if err != nil {
			return fmt.Errorf("failed to open database connection. %w", err)
		}
		if err := db.Ping(); err != nil {
			return fmt.Errorf("failed to ping database. %w", err)
		}

		if _, err := db.Query("SELECT 1;"); err != nil {
			return fmt.Errorf("failed to execute test query on database. %w", err)
		}

		return nil
	}); err != nil {
		m.t.Fatalf("could not connect to database: %q", err)
	}

	if err := m.prepareDBForFirstUse(); err != nil {
		m.t.Fatalf("failed to prepare database for first use: %q", err)
	}

	if err := m.LoadBaseData(); err != nil {
		m.t.Fatalf("failed to load base data into database: %q", err)
	}

	m.db, err = sql.Open("mysql", m.getDSN())
	if err != nil {
		m.t.Fatalf("could not connect to database after setup: %q", err)
	}

	// Auto stop server when test is done
	m.t.Cleanup(m.Stop)
}

func (m *dbServer) DB() (*sql.DB, error) {
	if m.db == nil {
		return nil, fmt.Errorf("test DB connection has not been established! You are accessing DB() method too early")
	}

	return m.db, nil
}

func (m *dbServer) getDSN() string {
	// Using `root` isn't cool, but a workaround for now to create triggers in the database
	return fmt.Sprintf("root:secret@(127.0.0.1:%s)/fivenettest?collation=utf8mb4_unicode_ci&loc=Local&parseTime=true", m.resource.GetPort("3306/tcp"))
}

func (m *dbServer) prepareDBForFirstUse() error {
	// Load and apply premigrate.sql file
	if err := m.loadSQLFile(filepath.Join(tests.TestDataSQLPath, "initial_esx.sql")); err != nil {
		return err
	}

	// Use DB migrations to handle the rest (esx compat mode is true)
	if err := query.MigrateDB(zap.NewNop(), m.getDSN(), true); err != nil {
		m.t.Fatalf("failed to migrate test database: %v", err)
	}

	return nil
}

func (m *dbServer) getMultiStatementDB() (*sql.DB, error) {
	// Open db connection with multiStatements param so we can apply sql files
	initDB, err := sql.Open("mysql", m.getDSN()+"&parseTime=true&multiStatements=true")
	if err != nil {
		return nil, fmt.Errorf("failed to open test database connection for multi statement exec: %v", err)
	}

	return initDB, nil
}

func (m *dbServer) loadSQLFile(file string) error {
	initDB, err := m.getMultiStatementDB()
	if err != nil {
		m.t.Fatalf("failed to get mult istatement db: %v", err)
	}

	c, ioErr := os.ReadFile(file)
	if ioErr != nil {
		m.t.Fatalf("failed to read %s for tests: %v", file, ioErr)
	}
	sqlBase := string(c)
	if _, err := initDB.Exec(sqlBase); err != nil {
		m.t.Fatalf("failed to apply %s for tests: %v", file, err)
	}

	return nil
}

func (m *dbServer) LoadBaseData() error {
	path := filepath.Join(tests.TestDataSQLPath, "base_*.sql")
	files, err := filepath.Glob(path)
	if err != nil {
		m.t.Fatalf("failed to find base data sql files (%s): %v", path, err)
	}
	// Sort the found files as they might not be in lexical order which we
	// need for this case https://github.com/golang/go/issues/17153
	sort.Strings(files)

	for _, file := range files {
		if err := m.loadSQLFile(file); err != nil {
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

	// You can't defer this because os.Exit doesn't care for defer
	if err := m.pool.Purge(m.resource); err != nil {
		m.t.Fatalf("could not purge container resource: %v", err)
	}
}

func (m *dbServer) FxProvide() fx.Option {
	return fx.Provide(func() *sql.DB {
		return m.db
	})
}
