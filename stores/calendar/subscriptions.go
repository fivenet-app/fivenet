package calendar

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) SetSubscription(
	ctx context.Context,
	calendarID int64,
	userID int32,
	subscribe bool,
	confirmed bool,
	muted bool,
) error {
	if subscribe {
		stmt := tCalendarSubs.
			INSERT(
				tCalendarSubs.CalendarID,
				tCalendarSubs.UserID,
				tCalendarSubs.Confirmed,
				tCalendarSubs.Muted,
			).
			VALUES(
				calendarID,
				userID,
				confirmed,
				muted,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCalendarSubs.Confirmed.SET(mysql.Bool(confirmed)),
				tCalendarSubs.Muted.SET(mysql.Bool(muted)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
		return nil
	}

	stmt := tCalendarSubs.
		DELETE().
		WHERE(mysql.AND(
			tCalendarSubs.CalendarID.EQ(mysql.Int64(calendarID)),
			tCalendarSubs.UserID.EQ(mysql.Int32(userID)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetCalendarSub(
	ctx context.Context,
	userID int32,
	calendarID int64,
) (*calendar.CalendarSub, error) {
	stmt := tCalendarSubs.
		SELECT(
			tCalendarSubs.CalendarID,
			tCalendarSubs.UserID,
			tCalendarSubs.CreatedAt,
			tCalendarSubs.Confirmed,
			tCalendarSubs.Muted,
		).
		FROM(tCalendarSubs).
		WHERE(mysql.AND(
			tCalendarSubs.CalendarID.EQ(mysql.Int64(calendarID)),
			tCalendarSubs.UserID.EQ(mysql.Int32(userID)),
		))

	var dest calendar.CalendarSub
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetCalendarId() == 0 || dest.GetUserId() == 0 {
		return nil, nil
	}

	return &dest, nil
}
