//nolint:contextcheck // Test setup and teardown intentionally mix test contexts with background teardown contexts.
package servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
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

const (
	mysqlRootPassword = "secret"
	mysqlUserPassword = "changeme"
	mysqlSeedDBName   = "fivenettest"
	mysqlCharset      = "utf8mb4"
	mysqlCollation    = "utf8mb4_unicode_ci"
	mysqlTimezone     = "Europe/Berlin"
)

type dbServer struct {
	t *testing.T

	db      *sql.DB
	dbName  string
	release func() error
	stopped bool
}

type setupUnavailableError struct {
	err error
}

func (e *setupUnavailableError) Error() string {
	return e.err.Error()
}

func (e *setupUnavailableError) Unwrap() error {
	return e.err
}

type mysqlTestDBManager struct {
	mu sync.Mutex

	pool     dockertest.ClosablePool
	resource dockertest.ClosableResource
	hasSeed  bool

	cloneSeq  atomic.Uint64
	cloneRefs int
}

var sharedMySQLTestDBManager = &mysqlTestDBManager{}

func NewDBServer(ctx context.Context, t *testing.T, setup bool) *dbServer {
	t.Helper()

	s := &dbServer{
		t: t,
	}

	t.Cleanup(s.Stop)

	if setup {
		s.Setup(ctx)
	}

	return s
}

// Setup ensures the shared MySQL test container exists, initializes the seed
// database once, and then allocates a fresh cloned database for this server.
// Callers receive an isolated schema that already contains the migrations and
// test seed data, and Stop() releases only that clone.
func (m *dbServer) Setup(ctx context.Context) {
	m.t.Helper()

	if m.db != nil || m.release != nil {
		m.Stop()
	}

	db, dbName, release, err := sharedMySQLTestDBManager.acquire(ctx, m.t)
	if err != nil {
		var unavailable *setupUnavailableError
		if errors.As(err, &unavailable) {
			m.t.Skipf("skipping docker-backed DB tests: %v", unavailable.err)
			return
		}

		require.NoError(m.t, err, "failed to set up test database")
		return
	}

	m.db = db
	m.dbName = dbName
	m.release = release
	m.stopped = false
}

func (m *dbServer) DB() (*sql.DB, error) {
	if m.db == nil {
		return nil, errors.New(
			"test DB connection has not been established! You are accessing DB() method too early",
		)
	}

	return m.db, nil
}

func (m *dbServer) DBName() string {
	if m == nil {
		return ""
	}

	return m.dbName
}

func (m *dbServer) Stop() {
	if m == nil || m.stopped {
		return
	}
	m.stopped = true

	if m.release != nil {
		release := m.release
		m.release = nil
		defer func() {
			require.NoError(m.t, release(), "could not release cloned test database")
		}()
	}

	if m.db != nil {
		require.NoError(m.t, m.db.Close(), "could not close test database connection")
		m.db = nil
	}
}

func (m *dbServer) FxProvide() fx.Option {
	return fx.Provide(func() *sql.DB {
		return m.db
	})
}

func (m *mysqlTestDBManager) acquire(
	ctx context.Context,
	t *testing.T,
) (*sql.DB, string, func() error, error) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.ensureSeedLocked(ctx, t); err != nil {
		return nil, "", nil, err
	}

	cloneName := fmt.Sprintf("%s_clone_%06d", mysqlSeedDBName, m.cloneSeq.Add(1))

	if err := m.cloneSeedLocked(ctx, t, cloneName); err != nil {
		return nil, "", nil, err
	}

	db, err := m.openDB(cloneName, false)
	if err != nil {
		_ = m.dropDatabaseLocked(ctx, cloneName)
		return nil, "", nil, fmt.Errorf("failed to open cloned test database: %w", err)
	}

	// Make sure the clone is ready before handing it to the caller.
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		_ = m.dropDatabaseLocked(ctx, cloneName)
		return nil, "", nil, fmt.Errorf("failed to ping cloned test database: %w", err)
	}

	m.cloneRefs++

	release := func() error {
		return m.releaseClone(t, cloneName)
	}

	return db, cloneName, release, nil
}

func (m *mysqlTestDBManager) ensureSeedLocked(ctx context.Context, t *testing.T) error {
	t.Helper()
	if m.hasSeed {
		return nil
	}

	if m.pool == nil {
		pool, poolErr := dockertest.NewPool(ctx, "")
		if poolErr != nil {
			return &setupUnavailableError{
				err: poolErr,
			}
		}

		m.pool = pool
	}

	if m.resource == nil {
		image, tag := loadDockerComposeServiceImage(t, "mysql")

		resource, runErr := m.pool.Run(
			ctx,
			image,
			dockertest.WithTag(tag),
			dockertest.WithEnv([]string{
				"MYSQL_ROOT_PASSWORD=" + mysqlRootPassword,
				"MYSQL_USER=" + "fivenet",
				"MYSQL_PASSWORD=" + mysqlUserPassword,
				"MYSQL_DATABASE=" + mysqlSeedDBName,
				"TZ=" + mysqlTimezone,
			}),
			dockertest.WithCmd([]string{
				"mysqld",
				"--innodb-ft-min-token-size=2",
				"--innodb-ft-max-token-size=50",
				"--default-time-zone=" + mysqlTimezone,
			}),
			dockertest.WithoutReuse(),
		)
		if runErr != nil {
			return &setupUnavailableError{
				err: runErr,
			}
		}

		m.resource = resource
	}

	if err := m.waitForSeedLocked(ctx); err != nil {
		_ = m.resetLocked(t)
		return fmt.Errorf("failed to wait for seed database to become ready: %w", err)
	}

	if err := m.prepareSeedLocked(ctx); err != nil {
		_ = m.resetLocked(t)
		return err
	}

	m.hasSeed = true
	return nil
}

func (m *mysqlTestDBManager) waitForSeedLocked(ctx context.Context) error {
	return m.pool.Retry(ctx, 0, func() error {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		db, err := m.openDB(mysqlSeedDBName, false)
		if err != nil {
			return fmt.Errorf("failed to open seed database connection: %w", err)
		}
		defer db.Close()

		if err := db.PingContext(ctx); err != nil {
			return fmt.Errorf("failed to ping seed database: %w", err)
		}

		rows, err := db.QueryContext(ctx, "SELECT 1;")
		if err != nil {
			return fmt.Errorf("failed to execute test query on seed database: %w", err)
		}
		defer rows.Close()

		return rows.Err()
	})
}

func (m *mysqlTestDBManager) prepareSeedLocked(ctx context.Context) error {
	if _, err := query.MigrateDB(
		ctx,
		zap.NewNop(),
		m.dsnForDB(mysqlSeedDBName, false),
		false,
		false,
	); err != nil {
		return fmt.Errorf("failed to migrate seed test database: %w", err)
	}

	if err := m.loadBaseDataLocked(ctx); err != nil {
		return fmt.Errorf("failed to load base data into seed database: %w", err)
	}

	return nil
}

func (m *mysqlTestDBManager) loadBaseDataLocked(ctx context.Context) error {
	path := filepath.Join(tests.TestDataSQLPath, "base_*.sql")
	files, err := filepath.Glob(path)
	if err != nil {
		return fmt.Errorf("failed to find base data sql files (%s): %w", path, err)
	}
	// Sort the found files as they might not be in lexical order which we
	// need for this case https://github.com/golang/go/issues/17153
	slices.Sort(files)

	initDB, err := m.openDB(mysqlSeedDBName, true)
	if err != nil {
		return fmt.Errorf("failed to open seed database for multi statement exec: %w", err)
	}
	defer initDB.Close()

	for _, file := range files {
		if err := m.loadSQLFileLocked(ctx, initDB, file); err != nil {
			return err
		}
	}

	return nil
}

func (m *mysqlTestDBManager) loadSQLFileLocked(
	ctx context.Context,
	initDB *sql.DB,
	file string,
) error {
	c, ioErr := os.ReadFile(file)
	if ioErr != nil {
		return fmt.Errorf("failed to read %s for tests: %w", file, ioErr)
	}

	if _, err := initDB.ExecContext(ctx, string(c)); err != nil {
		return fmt.Errorf("failed to apply %s for tests: %w", file, err)
	}

	return nil
}

func (m *mysqlTestDBManager) cloneSeedLocked(
	ctx context.Context,
	t *testing.T,
	cloneName string,
) error {
	t.Helper()
	seedDB, err := m.openDB(mysqlSeedDBName, false)
	if err != nil {
		return fmt.Errorf("failed to open seed database for clone setup: %w", err)
	}
	defer seedDB.Close()

	if err := m.createDatabaseLocked(ctx, seedDB, cloneName); err != nil {
		return err
	}

	tables, err := m.listSeedTablesLocked(ctx, seedDB)
	if err != nil {
		_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
		return err
	}

	db, err := m.openDB(cloneName, false)
	if err != nil {
		_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
		return fmt.Errorf("failed to open cloned database %s: %w", cloneName, err)
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0;"); err != nil {
		_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
		return fmt.Errorf("failed to disable foreign key checks for %s: %w", cloneName, err)
	}
	defer func(ctx context.Context) {
		_, _ = db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1;")
	}(ctx)

	for _, table := range tables {
		createStmt, err := m.showCreateTableLocked(ctx, seedDB, table)
		if err != nil {
			_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
			return err
		}

		createStmt, err = rewriteCreateTableStatement(createStmt, cloneName, table)
		if err != nil {
			_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
			return err
		}

		if _, err := db.ExecContext(ctx, createStmt); err != nil {
			_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
			return fmt.Errorf(
				"failed to create cloned table %s: %w",
				qualifiedMySQLName(cloneName, table),
				err,
			)
		}

		cloneTable := qualifiedMySQLName(cloneName, table)
		seedTable := qualifiedMySQLName(mysqlSeedDBName, table)
		columns, err := m.listCopyColumnsLocked(ctx, seedDB, cloneName, table)
		if err != nil {
			_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
			return err
		}

		if _, err := db.ExecContext(
			ctx,
			fmt.Sprintf(
				"INSERT INTO %s (%s) SELECT %s FROM %s;",
				cloneTable,
				strings.Join(columns, ", "),
				strings.Join(columns, ", "),
				seedTable,
			),
		); err != nil {
			_ = m.dropDatabaseUsingLocked(ctx, seedDB, cloneName)
			return fmt.Errorf("failed to copy rows into cloned table %s: %w", cloneTable, err)
		}
	}

	return nil
}

func (m *mysqlTestDBManager) createDatabaseLocked(
	ctx context.Context,
	db *sql.DB,
	dbName string,
) error {
	stmt := fmt.Sprintf(
		"CREATE DATABASE %s CHARACTER SET %s COLLATE %s;",
		quoteMySQLIdent(dbName),
		mysqlCharset,
		mysqlCollation,
	)
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("failed to create cloned database %s: %w", dbName, err)
	}

	return nil
}

func (m *mysqlTestDBManager) listSeedTablesLocked(
	ctx context.Context,
	db *sql.DB,
) ([]string, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT table_name
		 FROM information_schema.tables
		 WHERE table_schema = ? AND table_type = 'BASE TABLE'
		 ORDER BY table_name`,
		mysqlSeedDBName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables from seed database: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan seed table name: %w", err)
		}

		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading seed table names: %w", err)
	}

	return tables, nil
}

func (m *mysqlTestDBManager) showCreateTableLocked(
	ctx context.Context,
	db *sql.DB,
	table string,
) (string, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s;", quoteMySQLIdent(table)))
	if err != nil {
		return "", fmt.Errorf("failed to show create table for %s: %w", table, err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return "", fmt.Errorf("failed while reading create statement for %s: %w", table, err)
		}
		return "", fmt.Errorf("no create statement returned for %s", table)
	}

	var tableName string
	var createStmt string
	if err := rows.Scan(&tableName, &createStmt); err != nil {
		return "", fmt.Errorf("failed to scan create statement for %s: %w", table, err)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("failed while reading create statement for %s: %w", table, err)
	}

	return createStmt, nil
}

func (m *mysqlTestDBManager) listCopyColumnsLocked(
	ctx context.Context,
	db *sql.DB,
	cloneSchema string,
	table string,
) ([]string, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT CONCAT('`+"`"+`', REPLACE(dst.column_name, '`+"`"+`', '`+"`"+"`"+`'), '`+"`"+`')
		 FROM information_schema.columns dst
		 JOIN information_schema.columns src
		   ON src.table_schema = ?
		  AND src.table_name = ?
		  AND src.column_name = dst.column_name
		 WHERE dst.table_schema = ?
		   AND dst.table_name = ?
		   AND dst.extra NOT LIKE '%GENERATED%'
		 ORDER BY dst.ordinal_position`,
		mysqlSeedDBName,
		table,
		cloneSchema,
		table,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list copy columns for %s: %w", table, err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, fmt.Errorf("failed to scan copy column for %s: %w", table, err)
		}
		columns = append(columns, column)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading copy columns for %s: %w", table, err)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("no copyable columns found for %s", table)
	}

	return columns, nil
}

func (m *mysqlTestDBManager) dropDatabaseUsingLocked(
	ctx context.Context,
	db *sql.DB,
	dbName string,
) error {
	_, err := db.ExecContext(
		ctx,
		fmt.Sprintf("DROP DATABASE IF EXISTS %s;", quoteMySQLIdent(dbName)),
	)
	if err != nil {
		return fmt.Errorf("failed to drop cloned database %s: %w", dbName, err)
	}

	return nil
}

func (m *mysqlTestDBManager) releaseClone(t *testing.T, cloneName string) error {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cloneRefs > 0 {
		m.cloneRefs--
	}

	err := m.dropDatabaseLocked(t.Context(), cloneName)

	if m.cloneRefs == 0 {
		if resetErr := m.resetLocked(t); resetErr != nil {
			if err != nil {
				return errors.Join(err, resetErr)
			}
			return resetErr
		}
	}

	return err
}

func (m *mysqlTestDBManager) resetLocked(t *testing.T) error {
	t.Helper()
	if m.resource == nil {
		m.hasSeed = false
		return nil
	}

	resource := m.resource
	m.resource = nil
	m.hasSeed = false
	m.cloneRefs = 0

	var errs []error
	if err := resource.Close(t.Context()); err != nil {
		errs = append(errs, fmt.Errorf("failed to close shared mysql test container: %w", err))
	}

	if m.pool != nil {
		pool := m.pool
		m.pool = nil
		if err := pool.Close(t.Context()); err != nil {
			errs = append(errs, fmt.Errorf("failed to close shared mysql test pool: %w", err))
		}
	}

	return errors.Join(errs...)
}

func (m *mysqlTestDBManager) dropDatabaseLocked(ctx context.Context, dbName string) error {
	db, err := m.openDB(mysqlSeedDBName, false)
	if err != nil {
		return fmt.Errorf("failed to open database manager connection for drop: %w", err)
	}
	defer db.Close()

	_, err = db.ExecContext(
		ctx,
		fmt.Sprintf("DROP DATABASE IF EXISTS %s;", quoteMySQLIdent(dbName)),
	)
	if err != nil {
		return fmt.Errorf("failed to drop cloned database %s: %w", dbName, err)
	}

	return nil
}

func (m *mysqlTestDBManager) openDB(dbName string, multiStatements bool) (*sql.DB, error) {
	dsn := m.dsnForDB(dbName, multiStatements)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *mysqlTestDBManager) dsnForDB(dbName string, multiStatements bool) string {
	dsn := fmt.Sprintf(
		"root:%s@(127.0.0.1:%s)/%s?collation=%s&loc=Local&parseTime=true",
		mysqlRootPassword,
		m.resource.GetPort("3306/tcp"),
		dbName,
		mysqlCollation,
	)
	if multiStatements {
		dsn += "&multiStatements=true"
	}

	return dsn
}

func quoteMySQLIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func qualifiedMySQLName(schema, name string) string {
	return quoteMySQLIdent(schema) + "." + quoteMySQLIdent(name)
}

func rewriteCreateTableStatement(createStmt, schema, table string) (string, error) {
	oldPrefix := fmt.Sprintf("CREATE TABLE %s", quoteMySQLIdent(table))
	newPrefix := fmt.Sprintf("CREATE TABLE %s", qualifiedMySQLName(schema, table))
	if !strings.HasPrefix(createStmt, oldPrefix) {
		return "", fmt.Errorf("unexpected CREATE TABLE statement for %s", table)
	}

	return strings.Replace(createStmt, oldPrefix, newPrefix, 1), nil
}
