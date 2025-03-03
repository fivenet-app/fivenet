package servers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/fivenet-app/fivenet/internal/tests"
	"github.com/fivenet-app/fivenet/query"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.uber.org/zap"
)

var TestDBServer = &dbServer{}

type dbServer struct {
	db       *sql.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource

	stopped bool
}

func (m *dbServer) Setup() error {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	var err error
	m.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not construct pool: %q", err)
	}

	// uses pool to try to connect to Docker
	err = m.pool.Client.Ping()
	if err != nil {
		return fmt.Errorf("could not connect to Docker: %q", err)
	}

	// pulls an image, creates a container based on it and runs it
	m.resource, err = m.pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "docker.io/library/mysql",
			Tag:        "9.2.0",
			Env: []string{
				"MYSQL_ROOT_PASSWORD=secret",
				"MYSQL_USER=fivenet",
				"MYSQL_PASSWORD=changeme",
				"MYSQL_DATABASE=fivenettest",
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
		return fmt.Errorf("could not start resource: %q", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := m.pool.Retry(func() error {
		db, err := sql.Open("mysql", m.getDSN())
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to database: %q", err)
	}

	if err := m.prepareDBForFirstUse(); err != nil {
		return err
	}

	if err := m.LoadBaseData(); err != nil {
		return err
	}

	m.db, err = sql.Open("mysql", m.getDSN())
	if err != nil {
		return fmt.Errorf("could not connect to database after setup: %q", err)
	}

	return nil
}

func (m *dbServer) DB() (*sql.DB, error) {
	if m.db == nil {
		return nil, fmt.Errorf("test DB connection has not been established! You are accessing DB() method too early")
	}

	return m.db, nil
}

func (m *dbServer) getDSN() string {
	// Using `root` isn't cool, but a workaround for now to create triggers in the database
	return fmt.Sprintf("root:secret@(127.0.0.1:%s)/fivenettest?collation=utf8mb4_unicode_ci&loc=Local", m.resource.GetPort("3306/tcp"))
}

func (m *dbServer) prepareDBForFirstUse() error {
	// Load and apply premigrate.sql file
	if err := m.loadSQLFile(filepath.Join(tests.TestDataSQLPath, "initial_esx.sql")); err != nil {
		return err
	}

	// Use DB migrations to handle the rest (esx compat mode is true)
	if err := query.MigrateDB(zap.NewNop(), m.getDSN(), true); err != nil {
		return fmt.Errorf("failed to migrate test database: %w", err)
	}

	return nil
}

func (m *dbServer) getMultiStatementDB() (*sql.DB, error) {
	// Open db connection with multiStatements param so we can apply sql files
	initDB, err := sql.Open("mysql", m.getDSN()+"&multiStatements=true")
	if err != nil {
		return nil, fmt.Errorf("failed to open test database connection for multi statement exec: %w", err)
	}

	return initDB, nil
}

func (m *dbServer) loadSQLFile(file string) error {
	initDB, err := m.getMultiStatementDB()
	if err != nil {
		return fmt.Errorf("failed to get mult istatement db: %w", err)
	}

	c, ioErr := os.ReadFile(file)
	if ioErr != nil {
		return fmt.Errorf("failed to read %s for tests: %w", file, ioErr)
	}
	sqlBase := string(c)
	if _, err := initDB.Exec(sqlBase); err != nil {
		return fmt.Errorf("failed to apply %s for tests: %w", file, err)
	}

	return nil
}

func (m *dbServer) LoadBaseData() error {
	path := filepath.Join(tests.TestDataSQLPath, "base_*.sql")
	files, err := filepath.Glob(path)
	if err != nil {
		return fmt.Errorf("failed to find base data sql files (%s): %w", path, err)
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

func (m *dbServer) Stop() error {
	if m.stopped {
		return nil
	}
	m.stopped = true

	// You can't defer this because os.Exit doesn't care for defer
	if err := m.pool.Purge(m.resource); err != nil {
		return fmt.Errorf("could not purge container resource: %w", err)
	}

	return nil
}

// Reset truncates all `fivenet_*` tables and loads the base test data
func (m *dbServer) Reset() error {
	initDB, err := m.getMultiStatementDB()
	if err != nil {
		return fmt.Errorf("failed to get mutli statement db: %w", err)
	}

	rows, err := initDB.Query("SHOW TABLES LIKE 'fivenet_%';")
	if err != nil {
		return fmt.Errorf("failed to list fivenet tables in test database: %q", err)
	}

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name to string: %q", err)
		}

		// Placeholders aren't supported for table names, see
		// https://github.com/go-sql-driver/mysql/issues/848#issuecomment-414910152`
		if _, err := initDB.Exec("SET FOREIGN_KEY_CHECKS = 0; TRUNCATE TABLE `" + tableName + "`; SET FOREIGN_KEY_CHECKS = 1;"); err != nil {
			log.Printf("failed to truncate %s table: %s", tableName, err)
		}
	}

	// Load base test data after every reset
	return m.LoadBaseData()
}
