package housekeeper

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

func TestSoftDeleteJobData(t *testing.T) {
	ctx := context.Background()

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
		Table:           jet.NewTable("", "calendars", ""),
		DeletedAtColumn: jet.TimestampColumn("deleted_at"),
		JobColumn:       jet.StringColumn("job"),
		IDColumn:        jet.IntegerColumn("id"),
		MinDays:         30,

		DependentTables: []*Table{
			{
				Table:           jet.NewTable("", "calendar_entries", ""),
				DeletedAtColumn: jet.TimestampColumn("deleted_at"),
				ForeignKey:      jet.IntegerColumn("calendar_id"),
				IDColumn:        jet.IntegerColumn("id"),

				DependentTables: []*Table{
					{
						Table:      jet.NewTable("", "calendar_rsvp", ""),
						ForeignKey: jet.IntegerColumn("entry_id"),
					},
				},
			},
			{
				Table:      jet.NewTable("", "calendar_subscriptions", ""),
				ForeignKey: jet.IntegerColumn("calendar_id"),
			},
			{
				Table:      jet.NewTable("", "calendar_subscriptions", ""),
				ForeignKey: jet.IntegerColumn("calendar_id"),
			},
		},
	}
	jobName := "test_job"

	// Mock queries for main table
	mock.ExpectExec("UPDATE calendars SET deleted_at = CURRENT_TIMESTAMP WHERE \\(.+\\(job = \\?\\) AND deleted_at IS NULL.+\\) LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 10))

	// Mock queries for dependent table `calendar_entries`
	mock.ExpectExec("UPDATE calendar_entries SET deleted_at = CURRENT_TIMESTAMP WHERE .+\\(job = \\?\\) AND deleted_at IS NULL AND \\(calendar_id IN \\(\\( SELECT id AS \"id\" FROM calendars WHERE .+\\(job = \\?\\) AND deleted_at IS NULL .+ LIMIT \\?;").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 5))

	// Execute the function
	err = housekeeper.SoftDeleteJobData(ctx, table, jobName)
	if err != nil {
		t.Errorf("SoftDeleteJobData failed: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
