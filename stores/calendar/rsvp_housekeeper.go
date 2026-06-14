package calendar

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) cleanupStaleCalendarRSVPOccurrences(
	ctx context.Context,
	tx qrm.DB,
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
		))

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *Store) cleanupOrphanedCalendarRSVPOccurrences(
	ctx context.Context,
	tx qrm.DB,
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
		))

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *Store) CleanupCalendarRSVPOccurrences(ctx context.Context) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -30)

	rowsAffected := int64(0)

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
