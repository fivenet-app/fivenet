package housekeeper

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
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
	mock.ExpectExec("(?s)UPDATE calendar_entries SET deleted_at = CURRENT_TIMESTAMP WHERE .*deleted_at IS NULL.*calendar_id IN.*SELECT id AS \"id\" FROM calendars.*\\(job = \\?\\).*deleted_at IS NULL.*LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Execute the function
	var r int64
	r, err = housekeeper.SoftDeleteJobData(ctx, table, jobName)
	assert.NoError(t, err, "SoftDeleteJobData failed (%d)", r)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet(), "unmet expectations")
}

func TestMarkRowsAsDeleted_NoParentJobColumnLeakOnChildWithoutJob(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	core, observed := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	h := &Housekeeper{
		logger: logger,
		dryRun: true,
	}

	parent := &Table{
		Table:           mysql.NewTable("", "fivenet_mailer_emails", ""),
		IDColumn:        mysql.IntegerColumn("id"),
		JobColumn:       mysql.StringColumn("job"),
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
	}

	child := &Table{
		Table:           mysql.NewTable("", "fivenet_mailer_threads", ""),
		ForeignKey:      mysql.IntegerColumn("creator_email_id"),
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
	}

	_, err := h.markRowsAsDeleted(ctx, parent, child, "unemployed")
	require.NoError(t, err)

	var query string
	for _, e := range observed.All() {
		if e.Message == "dry run markRowsAsDeleted statement" {
			query = e.ContextMap()["query"].(string)
			break
		}
	}

	require.NotEmpty(t, query, "expected dry run query log")
	assert.NotContains(t, query, "fivenet_mailer_emails.job")
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
