package dbsync

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TableManager struct {
	logger *zap.Logger
	db     *sql.DB
}

type TableManagerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB

	Config *Config
}

func NewTableManager(p TableManagerParams) *TableManager {
	t := &TableManager{
		logger: p.Logger,
		db:     p.DB,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		cfg := p.Config.cfg.Load()
		tables := cfg.Tables.GetAllTables()
		if len(tables) == 0 {
			return nil
		}

		for _, table := range tables {
			if table.TableName == "" || table.UpdatedTimeColumn == nil {
				continue
			}

			if err := t.checkIfTableExists(ctxStartup, table.TableName); err != nil {
				return fmt.Errorf("table %q check failed: %w", table.TableName, err)
			}

			hasUpdatedAt, err := t.checkIfTableHasUpdatedAtColumn(ctxStartup, table.TableName)
			if err != nil {
				return fmt.Errorf(
					"table %q updated_at column check failed: %w",
					table.TableName,
					err,
				)
			}

			if !hasUpdatedAt {
				columnName := *table.UpdatedTimeColumn
				if err := t.addUpdatedAtColumnToTable(ctxStartup, table.TableName, columnName); err != nil {
					return fmt.Errorf(
						"adding updated_at column to table %q failed: %w",
						table.TableName,
						err,
					)
				}
				t.logger.Info(
					"added updated_at column to table",
					zap.String("table", table.TableName),
					zap.String("column", columnName),
				)
			}
		}

		return nil
	}))

	return t
}

func (t *TableManager) checkIfTableExists(
	ctx context.Context,
	tableName string,
) error {
	rows, err := t.db.QueryContext(ctx, `
        SELECT TABLE_NAME
        FROM information_schema.tables
        WHERE table_schema = DATABASE() AND table_name = ?
        LIMIT 1`, tableName)
	if err != nil {
		return fmt.Errorf("failed to check if table %q exists. %w", tableName, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("table %q does not exist in source database", tableName)
	}

	return nil
}

func (t *TableManager) checkIfTableHasUpdatedAtColumn(
	ctx context.Context,
	tableName string,
) (bool, error) {
	rows, err := t.db.QueryContext(ctx, `
        SELECT
            c.COLUMN_NAME
        FROM
	        information_schema.COLUMNS c
        WHERE
	        c.TABLE_SCHEMA = DATABASE() AND c.TABLE_NAME = ? AND LOWER(c.EXTRA) LIKE '%on update current_timestamp%'
        LIMIT 1`, tableName)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check if table %q has an on update timestamp column. %w",
			tableName,
			err,
		)
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	return true, nil
}

// TODO logic to add a `updated_at` (on update current_timestamp) column to tables if enabled and no such column exists (e.g., ESX's `owned_vehicles` is a prime candidate)

func (t *TableManager) addUpdatedAtColumnToTable(
	ctx context.Context,
	tableName string,
	columnName string,
) error {
	if _, err := t.db.ExecContext(ctx, `
        ALTER TABLE `+`"`+tableName+"`"+`
            ADD `+"`"+columnName+"`"+` datetime(3) on update CURRENT_TIMESTAMP(3) NULL`); err != nil {
		return fmt.Errorf(
			"failed to add on update timestamp column (%q) to table %q. %w",
			columnName,
			tableName,
			err,
		)
	}

	return nil
}
