package housekeeper

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

func TestSoftDeleteJobData(t *testing.T) {
	ctx := t.Context()

	// Mock dependencies
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
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
	mock.ExpectExec("UPDATE calendars SET deleted_at = CURRENT_TIMESTAMP WHERE \\(.+\\(job = \\?\\) AND deleted_at IS NULL.+\\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 10))

	// Mock queries for dependant table `calendar_entries`
	mock.ExpectExec("UPDATE calendar_entries SET deleted_at = CURRENT_TIMESTAMP WHERE \\(job = \\?\\) AND.+deleted_at IS NULL AND \\(calendar_id IN .+ SELECT id AS \"id\" FROM calendars WHERE .+\\(job = \\?\\) AND deleted_at IS NULL.+ LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Execute the function
	var r int64
	r, err = housekeeper.SoftDeleteJobData(ctx, table, jobName)
	if err != nil {
		t.Errorf("SoftDeleteJobData failed (%d): %v", r, err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestHardDelete(t *testing.T) {
	ctx := t.Context()

	// Mock dependencies
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
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
	mock.ExpectExec("DELETE FROM calendar_rsvp WHERE \\(entry_id IN \\(\\( SELECT id AS \"id\" FROM calendar_entries WHERE \\( deleted_at IS NOT NULL AND \\(deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\)\\) \\) \\)\\)\\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Mock queries for dependant table `calendar_entries`
	mock.ExpectExec("DELETE FROM calendar_entries WHERE \\(calendar_id IN \\(\\( SELECT id AS \"id\" FROM calendars WHERE \\( deleted_at IS NOT NULL AND \\(deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\)\\) \\) \\)\\)\\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 10))

	// Mock queries for dependant table `calendar_subscriptions`
	mock.ExpectExec("DELETE FROM calendar_subscriptions WHERE \\(calendar_id IN \\(\\( SELECT id AS \"id\" FROM calendars WHERE \\( deleted_at IS NOT NULL AND \\(deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\)\\) \\) \\)\\)\\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Mock queries for main table `calendars`
	mock.ExpectExec("DELETE FROM calendars WHERE \\( deleted_at IS NOT NULL AND \\(deleted_at <= \\(CURRENT_DATE - INTERVAL 30 DAY\\)\\) \\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Execute the function
	var r int64
	r, err = housekeeper.HardDelete(ctx, table)
	if err != nil {
		t.Errorf("HardDelete failed (%d): %v", r, err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
