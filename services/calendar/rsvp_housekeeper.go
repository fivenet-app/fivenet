package calendar

import (
	"context"
	"database/sql"
	"time"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

const cleanupLimit = 1000

func (s *Server) cleanupStaleCalendarRSVPOccurrences(
	ctx context.Context,
	tx *sql.Tx,
	cutoff time.Time,
) (int64, error) {
	tCalendarEntry := table.FivenetCalendarEntries
	tCalendarRSVPOccurrence := table.FivenetCalendarRsvpOccurrence

	stmt := tCalendarRSVPOccurrence.
		DELETE().
		USING(
			tCalendarRSVPOccurrence.
				INNER_JOIN(
					tCalendarEntry,
					tCalendarRSVPOccurrence.EntryID.EQ(tCalendarEntry.ID),
				),
		).
		WHERE(mysql.AND(
			tCalendarRSVPOccurrence.CreatedAt.LT(mysql.TimestampT(cutoff)),
			tCalendarRSVPOccurrence.RecurrenceVersion.NOT_EQ(
				tCalendarEntry.RecurrenceVersion,
			),
		)).
		LIMIT(1000)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *Server) cleanupOrphanedCalendarRSVPOccurrences(
	ctx context.Context,
	tx *sql.Tx,
	cutoff time.Time,
) (int64, error) {
	tCalendarEntry := table.FivenetCalendarEntries
	tCalendarRSVPOccurrence := table.FivenetCalendarRsvpOccurrence

	stmt := tCalendarRSVPOccurrence.
		DELETE().
		USING(
			tCalendarRSVPOccurrence.
				LEFT_JOIN(
					tCalendarEntry,
					tCalendarRSVPOccurrence.EntryID.EQ(tCalendarEntry.ID),
				),
		).
		WHERE(mysql.AND(
			tCalendarRSVPOccurrence.CreatedAt.LT(mysql.TimestampT(cutoff)),
			tCalendarEntry.ID.IS_NULL(),
		)).
		LIMIT(1000)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *Server) cleanupCalendarRSVPOccurrences(ctx context.Context) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -30)

	rowsAffected := int64(0)

	// DB transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return rowsAffected, err
	}
	defer tx.Rollback()

	rows, err := s.cleanupStaleCalendarRSVPOccurrences(ctx, tx, cutoff)
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected += rows

	rows, err = s.cleanupOrphanedCalendarRSVPOccurrences(ctx, tx, cutoff)
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected += rows

	if err := tx.Commit(); err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}
