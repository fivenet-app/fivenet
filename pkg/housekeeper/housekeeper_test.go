package housekeeper

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
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
	require.NoError(t, err, "SoftDeleteJobData failed (%d)", r)

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet(), "unmet expectations")
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
	mock.ExpectExec("DELETE FROM calendars WHERE \\( \\(deleted_at IS NOT NULL\\) AND \\(deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\)\\) \\) ORDER BY deleted_at ASC, id ASC LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Execute the function
	var r int64
	r, err = housekeeper.HardDelete(ctx, table)
	require.NoError(t, err, "HardDelete failed (%d)", r)

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet(), "unmet expectations")
}

func TestHardDeleteStringKeyDependencyIsBatched(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	h := &Housekeeper{logger: zap.NewNop(), db: db}
	table := &Table{
		Table:           mysql.NewTable("", "property_definitions", ""),
		KeyColumnString: mysql.StringColumn("property_key"),
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
		MinDays:         30,
		DependantTables: []*Table{{
			Table:            mysql.NewTable("", "property_values", ""),
			ForeignKeyString: mysql.StringColumn("property_key"),
		}},
	}

	mock.ExpectExec("(?s)DELETE FROM property_values WHERE .*property_key IN.*SELECT property_key AS \"property_key\" FROM property_definitions.*deleted_at IS NOT NULL.*LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, DefaultDeleteLimit))
	mock.ExpectExec("DELETE FROM property_definitions WHERE .*deleted_at IS NOT NULL.*LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = h.HardDelete(t.Context(), table)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRunHardDeleteRecordsProgressAttributes(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	table := &Table{
		Table:           mysql.NewTable("", "fivenet_audit_log", ""),
		IDColumn:        mysql.IntegerColumn("id"),
		DeletedAtColumn: mysql.TimestampColumn("created_at"),
		MinDays:         90,
	}
	h := &Housekeeper{
		logger: zap.NewNop(),
		db:     db,
		getTablesListFn: func() map[string]*Table {
			return map[string]*Table{"fivenet_audit_log": table}
		},
	}

	mock.ExpectExec("DELETE FROM fivenet_audit_log WHERE .*created_at <= .* LIMIT \\?;").
		WillReturnResult(sqlmock.NewResult(0, DefaultDeleteLimit))

	data := &cron.GenericCronData{Attributes: map[string]string{}}
	rows, err := h.runHardDelete(t.Context(), data)
	require.NoError(t, err)
	assert.EqualValues(t, DefaultDeleteLimit, rows)
	assert.Equal(t, "fivenet_audit_log", data.GetAttribute(currentTable))
	assert.Equal(t, "true", data.GetAttribute(backlogRemaining))
	assert.Empty(t, data.GetAttribute(lastTableMapIndex))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRetentionDeleteOrderByUsesConfiguredAgeColumn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		table *Table
		want  string
	}{
		{
			name: "date",
			table: &Table{
				DateColumn: mysql.DateColumn("entry_date"),
				IDColumn:   mysql.IntegerColumn("id"),
			},
			want: "entry_date ASC, id ASC",
		},
		{
			name: "timestamp",
			table: &Table{
				TimestampColumn: mysql.TimestampColumn("created_at"),
				IDColumn:        mysql.IntegerColumn("id"),
			},
			want: "created_at ASC, id ASC",
		},
		{
			name: "deleted at",
			table: &Table{
				DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
				IDColumn:        mysql.IntegerColumn("id"),
			},
			want: "deleted_at ASC, id ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			orderBy := retentionDeleteOrderBy(tt.table)
			require.Len(t, orderBy, 2)

			stmt := mysql.NewTable("", "cleanup", "").
				DELETE().
				WHERE(mysql.Bool(true)).
				ORDER_BY(orderBy...).
				LIMIT(DefaultDeleteLimit)
			assert.Contains(t, stmt.DebugSql(), "ORDER BY "+tt.want)
		})
	}
}

func TestTableRegistrationKeyUsesActiveRetentionColumn(t *testing.T) {
	t.Parallel()

	table := mysql.NewTable("", "fivenet_centrum_markers", "")
	assert.Equal(t, "fivenet_centrum_markers:deleted_at", tableRegistrationKey(&Table{
		Table:           table,
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
	}))
	assert.Equal(t, "fivenet_centrum_markers:expires_at", tableRegistrationKey(&Table{
		Table:           table,
		TimestampColumn: mysql.TimestampColumn("expires_at"),
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
	}))
}

func TestAddTableKeepsSeparatePoliciesForSameTable(t *testing.T) {
	t.Parallel()

	tableListsMu.Lock()
	originalTablesList := tablesList
	tablesList = map[string]*Table{}
	tableListsMu.Unlock()
	t.Cleanup(func() {
		tableListsMu.Lock()
		tablesList = originalTablesList
		tableListsMu.Unlock()
	})

	table := mysql.NewTable("", "fivenet_centrum_markers", "")
	AddTable(&Table{
		Table:           table,
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
	})
	AddTable(&Table{
		Table:           table,
		TimestampColumn: mysql.TimestampColumn("expires_at"),
		DeletedAtColumn: mysql.TimestampColumn("deleted_at"),
	})

	tableListsMu.Lock()
	defer tableListsMu.Unlock()
	assert.Contains(t, tablesList, "fivenet_centrum_markers:deleted_at")
	assert.Contains(t, tablesList, "fivenet_centrum_markers:expires_at")
}
