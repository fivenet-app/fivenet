package dbsync

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var TableManagerModule = fx.Module(
	"dbsync.table_manager",
	fx.Provide(NewTableManager),
)

type TableManager struct {
	logger *zap.Logger

	dryRun bool
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
		logger: p.Logger.Named("dbsync.table_manager"),
		dryRun: !p.Config.Load().TableManager.Enabled,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		return t.CheckTables(ctxStartup, p.DB, p.Config.Load().Tables.GetAllTables())
	}))

	return t
}

func (t *TableManager) CheckTables(
	ctx context.Context,
	db qrm.DB,
	tables []DBSyncTable,
) error {
	if len(tables) == 0 {
		return nil
	}

	for _, table := range tables {
		if table.TableName == "" || table.UpdatedTimeColumn == nil {
			continue
		}

		if err := t.checkIfTableExists(ctx, db, table.TableName); err != nil {
			return fmt.Errorf("table %q check failed: %w", table.TableName, err)
		}

		hasUpdatedAt, err := t.checkIfTableHasUpdatedAtColumn(ctx, db, table.TableName)
		if err != nil {
			return fmt.Errorf(
				"table %q updated_at column check failed: %w",
				table.TableName,
				err,
			)
		}

		if !hasUpdatedAt {
			if t.dryRun {
				t.logger.Info(
					"dry run enabled, skipping adding non-existent updated_at column to table",
					zap.String("table", table.TableName),
					zap.String("column", *table.UpdatedTimeColumn),
				)

				table.UpdatedTimeColumn = nil
				continue
			}

			columnName := *table.UpdatedTimeColumn
			if err := t.addUpdatedAtColumnToTable(ctx, db, table.TableName, columnName); err != nil {
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
}

func (t *TableManager) checkIfTableExists(
	ctx context.Context,
	db qrm.Queryable,
	tableName string,
) error {
	rows, err := db.QueryContext(ctx, `
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
	db qrm.Queryable,
	tableName string,
) (bool, error) {
	rows, err := db.QueryContext(ctx, `
        SELECT c.COLUMN_NAME
        FROM information_schema.COLUMNS c
        WHERE c.TABLE_SCHEMA = DATABASE() AND c.TABLE_NAME = ? AND LOWER(c.EXTRA) LIKE '%on update current_timestamp%'
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

// addUpdatedAtColumnToTable adds a new column to the specified table that automatically updates
// its value to the current timestamp whenever the row is updated using MySQLs `ON UPDATE` system.
func (t *TableManager) addUpdatedAtColumnToTable(
	ctx context.Context,
	db qrm.Executable,
	tableName string,
	columnName string,
) error {
	query := `ALTER TABLE ` + "`" + tableName + "` " +
		`ADD ` + "`" + columnName + "`" + ` datetime(3) on update CURRENT_TIMESTAMP(3) NULL`
	if _, err := db.ExecContext(ctx, query); err != nil {
		t.logger.Debug(
			"alter table on updated time column to table",
			zap.String("column_name", columnName),
			zap.String("table_name", tableName),
			zap.String("query", query),
		)
		return fmt.Errorf(
			"failed to add on update timestamp column (%q) to table %q. %w",
			columnName,
			tableName,
			err,
		)
	}

	return nil
}
