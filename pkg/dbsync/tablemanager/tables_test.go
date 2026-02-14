package dbsynctablemanager

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTableManager_CheckTables(t *testing.T) {
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tableManager := &TableManager{
		logger: logger,
		dryRun: false,
	}

	ctx := t.Context()

	col := "updated_at"
	tables := []dbsyncconfig.DBSyncTable{
		{
			TableName:         "vehicles",
			UpdatedTimeColumn: &col,
		},
		{
			TableName:         "jobs",
			UpdatedTimeColumn: nil,
		},
		{
			TableName:         "",
			UpdatedTimeColumn: &col,
		},
	}

	// Mock database responses
	mock.ExpectQuery("SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = DATABASE\\(\\) AND table_name = \\? LIMIT 1").
		WithArgs("vehicles").
		WillReturnRows(sqlmock.NewRows([]string{"TABLE_NAME"}).AddRow("vehicles"))

	mock.ExpectQuery("SELECT c.COLUMN_NAME FROM information_schema.COLUMNS c WHERE c.TABLE_SCHEMA = DATABASE\\(\\) AND c.TABLE_NAME = \\? AND LOWER\\(c.EXTRA\\) LIKE '%on update current_timestamp%' LIMIT 1").
		WithArgs("vehicles").
		WillReturnRows(sqlmock.NewRows([]string{"COLUMN_NAME"}))

	mock.ExpectExec("ALTER TABLE `vehicles` ADD `updated_at` datetime\\(3\\) on update CURRENT_TIMESTAMP\\(3\\) NULL").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = tableManager.CheckTables(ctx, db, tables)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTableManager_checkIfTableExists(t *testing.T) {
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tableManager := &TableManager{
		logger: logger,
	}

	ctx := t.Context()
	tableName := "test_table"

	// Mock database responses
	mock.ExpectQuery("SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = DATABASE\\(\\) AND table_name = \\? LIMIT 1").
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"TABLE_NAME"}).AddRow(tableName))

	err = tableManager.checkIfTableExists(ctx, db, tableName)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTableManager_checkIfTableHasUpdatedAtColumn(t *testing.T) {
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tableManager := &TableManager{
		logger: logger,
	}

	ctx := t.Context()
	tableName := "test_table"

	// Mock database responses
	mock.ExpectQuery("SELECT c.COLUMN_NAME FROM information_schema.COLUMNS c WHERE c.TABLE_SCHEMA = DATABASE\\(\\) AND c.TABLE_NAME = \\? AND LOWER\\(c.EXTRA\\) LIKE '%on update current_timestamp%' LIMIT 1").
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"COLUMN_NAME"}).AddRow("updated_at"))

	hasUpdatedAt, err := tableManager.checkIfTableHasUpdatedAtColumn(ctx, db, tableName)
	assert.NoError(t, err)
	assert.True(t, hasUpdatedAt)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTableManager_addUpdatedAtColumnToTable(t *testing.T) {
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tableManager := &TableManager{
		logger: logger,
	}

	ctx := t.Context()
	tableName := "test_table"
	columnName := "updated_at"

	// Mock database responses
	mock.ExpectExec("ALTER TABLE `test_table` ADD `updated_at` datetime\\(3\\) on update CURRENT_TIMESTAMP\\(3\\) NULL").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("ALTER TABLE `test_table` ADD INDEX `idx_test_table_updated_at` \\(`updated_at`\\)").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = tableManager.addUpdatedAtColumnToTable(ctx, db, tableName, columnName)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
