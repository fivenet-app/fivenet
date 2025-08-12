package reqs

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type Table struct {
	Name      string
	Charset   string
	Collation string
}

type DBReqs struct {
	mu sync.Mutex

	db *sql.DB

	version string

	migrationVersion uint
	migrationDirty   bool

	dbCharset   string
	dbCollation string

	tables []Table
}

func NewDBReqs(db *sql.DB) *DBReqs {
	return &DBReqs{
		db: db,
	}
}

func (r *DBReqs) ReplaceDB(newDB *sql.DB) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db = newDB
}

func (r *DBReqs) GetVersion() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.version
}

func (r *DBReqs) SetMigrationState(version uint, dirty bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.migrationVersion = version
	r.migrationDirty = dirty
}

func (r *DBReqs) ValidateVersion(ctx context.Context) error {
	err := r.db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&r.version)
	if err != nil {
		return err
	}

	// Parse version string, e.g., "8.0.32", "10.4.17-MariaDB", etc.
	// At least MariaDB's version string has the "extra info" kinda semver compatible.
	// We will just check the major and minor version numbers.
	var major, minor int
	_, err = fmt.Sscanf(r.version, "%d.%d", &major, &minor)
	if err != nil {
		return fmt.Errorf("failed to parse DB version: %w", err)
	}

	if major < 8 {
		return fmt.Errorf(
			"database version %s is not supported, requires at least MySQL 8.0 or MariaDB 11.4",
			r.version,
		)
	}

	return nil
}

func (r *DBReqs) ValidateTables(ctx context.Context) error {
	tables, err := r.validateTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to validate database tables. %w", err)
	}

	r.tables = tables

	// Check if any tables are mismatched
	if len(tables) > 0 {
		return fmt.Errorf(
			"database (charset: %q, collation: %q) tables are mismatched: %v",
			r.dbCharset,
			r.dbCollation,
			tables,
		)
	}

	return nil
}

func (r *DBReqs) GetTables() []Table {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tables
}

func (r *DBReqs) validateTables(ctx context.Context) ([]Table, error) {
	var dbName string
	err := r.db.
		QueryRowContext(ctx, "SELECT DATABASE()").
		Scan(&dbName)
	if err != nil {
		return nil, err
	}

	// Get database charset and collation
	var dbCharset, dbCollation string
	err = r.db.
		QueryRowContext(ctx, `
        SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME
        FROM information_schema.SCHEMATA
        WHERE SCHEMA_NAME = ?
    `, dbName).
		Scan(&dbCharset, &dbCollation)
	if err != nil {
		return nil, err
	}
	r.dbCharset = dbCharset
	r.dbCollation = dbCollation

	// Get tables and their charsets/collations
	prefix := "fivenet_"
	rows, err := r.db.
		Query(`
        SELECT TABLE_NAME, TABLE_COLLATION
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = ? AND TABLE_NAME LIKE ?
    `, dbName, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mismatched []Table
	for rows.Next() {
		var tableName, tableCollation string
		if err := rows.Scan(&tableName, &tableCollation); err != nil {
			return nil, err
		}
		var tableCharset string
		err = r.db.
			QueryRowContext(ctx, `
            SELECT CHARACTER_SET_NAME
            FROM information_schema.COLLATIONS
            WHERE COLLATION_NAME = ?
        `, tableCollation).
			Scan(&tableCharset)
		if err != nil {
			return nil, err
		}
		if tableCharset != dbCharset || tableCollation != dbCollation {
			mismatched = append(mismatched, Table{
				Name:      tableName,
				Charset:   tableCharset,
				Collation: tableCollation,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mismatched, nil
}

func (r *DBReqs) GetDBCharsetAndCollation() (string, string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.dbCharset, r.dbCollation
}

func (r *DBReqs) GetMigrationState() (uint, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.migrationVersion, r.migrationDirty
}
