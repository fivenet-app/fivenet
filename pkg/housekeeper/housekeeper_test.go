package housekeeper

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRunHousekeeper(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock logger
	logger := zap.NewNop()

	tablesList := map[string]*Table{
		"test_table": {
			Table:      jet.NewTable("", "test_table", ""),
			DateColumn: jet.DateColumn("date_column"),
			MinDays:    30,
		},
	}

	// Create a Housekeeper instance and add mock table to list
	h := &Housekeeper{
		logger: logger,
		db:     db,

		getTablesListFn: func() map[string]*Table {
			return tablesList
		},
	}

	// Define test data
	data := &cron.GenericCronData{
		Attributes: map[string]string{},
	}

	// Mock DELETE query
	mock.ExpectExec(`DELETE FROM test_table .+date_column IS NOT NULL AND .date_column <= .+CURRENT_DATE - INTERVAL 30 DAY. AS DATE.+ LIMIT \?;`).
		WithArgs(2000).
		WillReturnResult(sqlmock.NewResult(0, 10)) // Simulate 10 rows affected

	// Run the method
	err = h.runHousekeeper(context.Background(), data)
	assert.NoError(t, err)

	// Verify the last table key is set accordingly
	assert.Equal(t, "test_table", data.Attributes[lastTableMapIndex])

	// Mock DELETE query of test_table
	mock.ExpectExec(`DELETE FROM test_table .+date_column IS NOT NULL AND .date_column <= .+CURRENT_DATE - INTERVAL 30 DAY. AS DATE.+ LIMIT \?;`).
		WithArgs(2000).
		WillReturnResult(sqlmock.NewResult(0, 10)) // Simulate 10 rows affected

	// Run the method again with the same table list
	err = h.runHousekeeper(context.Background(), data)
	assert.NoError(t, err)

	// Verify the last table key is set accordingly
	assert.Equal(t, "test_table", data.Attributes[lastTableMapIndex])

	// Test with a second table in the list
	tablesList["zsecond_table"] = &Table{
		Table:           jet.NewTable("", "zsecond_table", ""),
		TimestampColumn: jet.TimestampColumn("timestamp_column"),
		MinDays:         7,
	}

	// Mock DELETE query for second table
	mock.ExpectExec(`DELETE FROM zsecond_table .+timestamp_column IS NOT NULL AND .timestamp_column <= .+CURRENT_DATE - INTERVAL 7 DAY.+ LIMIT \?;`).
		WithArgs(2000).
		WillReturnResult(sqlmock.NewResult(0, 10)) // Simulate 10 rows affected

	// Run the method
	err = h.runHousekeeper(context.Background(), data)
	assert.NoError(t, err)

	// Verify the last table key is set to the second table
	assert.Equal(t, "zsecond_table", data.Attributes[lastTableMapIndex])

	// Mock DELETE query of test_table
	mock.ExpectExec(`DELETE FROM test_table .+date_column IS NOT NULL AND .date_column <= .+CURRENT_DATE - INTERVAL 30 DAY. AS DATE.+ LIMIT \?;`).
		WithArgs(2000).
		WillReturnResult(sqlmock.NewResult(0, 10)) // Simulate 10 rows affected

	// Run the method again it should "rollover" to the first table in the list
	err = h.runHousekeeper(context.Background(), data)
	assert.NoError(t, err)

	// Verify the last table key is set to the first test table
	assert.Equal(t, "test_table", data.Attributes[lastTableMapIndex])

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
