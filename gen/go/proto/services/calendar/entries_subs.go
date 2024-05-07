package calendar

import (
	"context"

	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) createOrDeleteEntrySubscription(ctx context.Context, calendarId uint64, entryId uint64, userId int32, subscribe bool, confirmed bool, muted bool) error {
	if subscribe {
		tCalendarEntriesSubs := table.FivenetCalendarEntriesSubs
		stmt := tCalendarEntriesSubs.
			INSERT(
				tCalendarEntriesSubs.CalendarID,
				tCalendarEntriesSubs.EntryID,
				tCalendarEntriesSubs.UserID,
				tCalendarEntriesSubs.Confirmed,
				tCalendarEntriesSubs.Muted,
			).
			VALUES(
				calendarId,
				entryId,
				userId,
				confirmed,
				muted,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCalendarEntriesSubs.Confirmed.SET(jet.Bool(confirmed)),
				tCalendarEntriesSubs.Muted.SET(jet.Bool(muted)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	} else {
		stmt := tCalendarEntriesSubs.
			DELETE().
			WHERE(jet.AND(
				tCalendarEntriesSubs.CalendarID.EQ(jet.Uint64(calendarId)),
				tCalendarEntriesSubs.CalendarID.EQ(jet.Uint64(entryId)),
				tCalendarEntriesSubs.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}
