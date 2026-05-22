package housekeeper

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSoftDeleteJobData(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	// Mock dependencies
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to create mock db")
	defer db.Close()

	housekeeper := &Housekeeper{
		logger: logger,
		db:     db,
	}

	// Define table and job details
	table := &Table{
		Table:           mysql.NewTable("", "calendars", ""),
		IDColumn:        mysql.IntegerColumn("id"),
		JobColumn:       mysql.StringColumn("job"),
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),

		MinDays: 30,

		DependantTables: []*Table{
			{
				Table:           mysql.NewTable("", "calendar_entries", ""),
				IDColumn:        mysql.IntegerColumn("id"),
				DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
				ForeignKey:      mysql.IntegerColumn("calendar_id"),

				MinDays: 30,

				DependantTables: []*Table{
					{
						Table:      mysql.NewTable("", "calendar_rsvp", ""),
						ForeignKey: mysql.IntegerColumn("entry_id"),
					},
				},
			},
			{
				Table:      mysql.NewTable("", "calendar_subscriptions", ""),
				ForeignKey: mysql.IntegerColumn("calendar_id"),
			},
			{
				Table:      mysql.NewTable("", "calendar_subscriptions", ""),
				ForeignKey: mysql.IntegerColumn("calendar_id"),
			},
		},
	}
	jobName := "test_job"

	// Mock queries for main table
	mock.ExpectExec("UPDATE calendars SET deleted_at = CURRENT_TIMESTAMP WHERE \\( \\(job = \\?\\) AND \\(deleted_at IS NULL\\) \\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 10))

	// Mock queries for dependant table `calendar_entries`
	mock.ExpectExec("(?s)UPDATE calendar_entries SET deleted_at = CURRENT_TIMESTAMP WHERE .*\\(job = \\?\\).*deleted_at IS NULL.*calendar_id IN.*SELECT id AS \"id\" FROM calendars.*\\(job = \\?\\).*deleted_at IS NULL.*LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Execute the function
	var r int64
	r, err = housekeeper.SoftDeleteJobData(ctx, table, jobName)
	assert.NoError(t, err, "SoftDeleteJobData failed (%d)", r)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet(), "unmet expectations")
}

func TestHardDelete(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	// Mock dependencies
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to create mock db")
	defer db.Close()

	housekeeper := &Housekeeper{
		logger: logger,
		db:     db,
	}

	// Define table details
	table := &Table{
		Table:           mysql.NewTable("", "calendars", ""),
		IDColumn:        mysql.IntegerColumn("id"),
		JobColumn:       mysql.StringColumn("job"),
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),

		MinDays: 30,

		DependantTables: []*Table{
			{
				Table:           mysql.NewTable("", "calendar_entries", ""),
				IDColumn:        mysql.IntegerColumn("id"),
				ForeignKey:      mysql.IntegerColumn("calendar_id"),
				DeletedAtColumn: mysql.TimestampColumn("deleted_at"),

				MinDays: 30,

				DependantTables: []*Table{
					{
						Table:      mysql.NewTable("", "calendar_rsvp", ""),
						ForeignKey: mysql.IntegerColumn("entry_id"),
					},
				},
			},
			{
				Table:      mysql.NewTable("", "calendar_subscriptions", ""),
				ForeignKey: mysql.IntegerColumn("calendar_id"),
			},
		},
	}

	// Mock queries for dependant table `calendar_rsvp`
	mock.ExpectExec("(?s)DELETE FROM calendar_rsvp WHERE .*entry_id IN.*SELECT id AS \"id\" FROM calendar_entries.*deleted_at IS NOT NULL.*deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\).*LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Mock queries for dependant table `calendar_entries`
	mock.ExpectExec("(?s)DELETE FROM calendar_entries WHERE .*calendar_id IN.*SELECT id AS \"id\" FROM calendars.*deleted_at IS NOT NULL.*deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\).*LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 10))

	// Mock queries for dependant table `calendar_subscriptions`
	mock.ExpectExec("(?s)DELETE FROM calendar_subscriptions WHERE .*calendar_id IN.*SELECT id AS \"id\" FROM calendars.*deleted_at IS NOT NULL.*deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\).*LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Mock queries for main table `calendars`
	mock.ExpectExec("DELETE FROM calendars WHERE \\( \\(deleted_at IS NOT NULL\\) AND \\(deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\)\\) \\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Execute the function
	var r int64
	r, err = housekeeper.HardDelete(ctx, table)
	assert.NoError(t, err, "HardDelete failed (%d)", r)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet(), "unmet expectations")
}
