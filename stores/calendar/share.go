package calendar

import (
	"context"
	"errors"
	"slices"

	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ShareCalendarEntry(
	ctx context.Context,
	tx qrm.DB,
	entryID int64,
	inUserIds []int32,
) ([]int32, error) {
	userIds := make([]mysql.Expression, len(inUserIds))
	for i := range inUserIds {
		userIds[i] = mysql.Int32(inUserIds[i])
	}

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.UserID,
		).
		FROM(tCalendarRSVP).
		WHERE(mysql.AND(
			tCalendarRSVP.EntryID.EQ(mysql.Int64(entryID)),
			tCalendarRSVP.UserID.IN(userIds...),
		))

	var currentRSVPs []*calendarentries.CalendarEntryRSVP
	if err := stmt.QueryContext(ctx, tx, &currentRSVPs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	newUsers := []int32{}
	if len(currentRSVPs) == 0 {
		newUsers = append(newUsers, inUserIds...)
	} else {
		for _, rsvp := range currentRSVPs {
			if !slices.Contains(inUserIds, rsvp.GetUserId()) {
				newUsers = append(newUsers, rsvp.GetUserId())
			}
		}
	}

	tCalendarRSVP := tCalendarRSVP
	insertStmt := tCalendarRSVP.
		INSERT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		)
	for i := range userIds {
		insertStmt = insertStmt.VALUES(
			entryID,
			userIds[i],
			calendarentries.RsvpResponses_RSVP_RESPONSES_INVITED,
		)
	}

	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, err
		}
	}

	return newUsers, nil
}

func (s *Store) GetUserShortByID(ctx context.Context, userID int32) (*usershort.UserShort, error) {
	tUsers := table.FivenetUser.AS("user_short")

	stmt := tUsers.
		SELECT(
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.PhoneNumber,
		).
		FROM(tUsers).
		WHERE(tUsers.ID.EQ(mysql.Int32(userID))).
		LIMIT(1)

	sourceUser := &usershort.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, sourceUser); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return sourceUser, nil
}
