package calendar

import (
	"context"

	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) createOrDeleteSubscription(ctx context.Context, calendarId uint64, entryId *uint64, userId int32, subscribe bool, muted bool) error {
	if subscribe {
		tCalendarSubs := table.FivenetCalendarSubs
		stmt := tCalendarSubs.
			INSERT(
				tCalendarSubs.CalendarID,
				tCalendarSubs.EntryID,
				tCalendarSubs.UserID,
				tCalendarSubs.Muted,
			).
			VALUES(
				calendarId,
				entryId,
				userId,
				muted,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCalendarSubs.Muted.SET(jet.Bool(muted)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	} else {
		entryIdColumn := jet.IntExp(jet.NULL)
		if entryId != nil {
			entryIdColumn = jet.Uint64(*entryId)
		}

		stmt := tCalendarSubs.
			DELETE().
			WHERE(jet.AND(
				tCalendarSubs.CalendarID.EQ(jet.Uint64(calendarId)),
				tCalendarSubs.EntryID.EQ(entryIdColumn),
				tCalendarSubs.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}
