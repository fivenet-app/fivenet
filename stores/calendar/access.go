package calendarstore

import (
	"context"
	"errors"

	calendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CheckIfUserHasAccessToCalendar(
	ctx context.Context,
	calendarID int64,
	userInfo *userinfo.UserInfo,
	accessLevel calendaraccess.AccessLevel,
	publicOk bool,
) (bool, error) {
	out, err := s.CheckIfUserHasAccessToCalendarIDs(
		ctx,
		userInfo,
		accessLevel,
		publicOk,
		calendarID,
	)
	return len(out) > 0, err
}

func (s *Store) CheckIfUserHasAccessToCalendarIDs(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	accessLevel calendaraccess.AccessLevel,
	publicOk bool,
	calendarIDs ...int64,
) ([]*calendarAccessEntry, error) {
	var dest []*calendarAccessEntry
	if len(calendarIDs) == 0 {
		return dest, nil
	}

	tCreator := tCalendar.AS("creator")

	ids := make([]mysql.Expression, len(calendarIDs))
	for i := range calendarIDs {
		ids[i] = mysql.Int64(calendarIDs[i])
	}

	condition := mysql.AND(
		tCalendar.Job.IS_NULL(),
		tCalendar.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
	)
	if publicOk {
		condition = condition.OR(tCalendar.Public.IS_TRUE())
	}

	var accessExists mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		accessExists = s.access.ACLAccessExistsCondition(tCalendar.ID, userInfo, int32(accessLevel))
	}

	visibleCondition := mysql.OR(accessExists, condition)
	if userInfo.GetSuperuser() {
		visibleCondition = mysql.OR(
			tCalendar.SystemKind.IS_NULL(),
			tCalendar.SystemKind.NOT_EQ(
				mysql.Int32(
					int32(calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
				),
			),
			s.birthdayCalendarVisible(tCalendar.ID, accessLevel, userInfo),
		)
	}

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
		).
		FROM(tCalendar.
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		WHERE(mysql.AND(
			tCalendar.ID.IN(ids...),
			tCalendar.DeletedAt.IS_NULL(),
			visibleCondition,
		)).
		ORDER_BY(tCalendar.ID.DESC())

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) CheckIfUserHasAccessToCalendarEntry(
	ctx context.Context,
	calendarID int64,
	entryID int64,
	userInfo *userinfo.UserInfo,
	accessLevel calendaraccess.AccessLevel,
	publicOk bool,
) (bool, error) {
	out, err := s.CheckIfUserHasAccessToCalendarEntryIDs(ctx, userInfo, publicOk, entryID)
	if err != nil {
		return false, err
	}

	if len(out) > 0 {
		return true, nil
	}

	check, err := s.CheckIfUserHasAccessToCalendar(ctx, calendarID, userInfo, accessLevel, publicOk)
	if err != nil {
		return false, err
	}

	return check, nil
}

func (s *Store) CheckIfUserHasAccessToCalendarEntryIDs(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	publicOk bool,
	entryIDs ...int64,
) ([]int64, error) {
	var dest []int64
	if len(entryIDs) == 0 {
		return dest, nil
	}

	if userInfo.GetSuperuser() {
		dest = append(dest, entryIDs...)
		return dest, nil
	}

	ids := make([]mysql.Expression, len(entryIDs))
	for i := range entryIDs {
		ids[i] = mysql.Int64(entryIDs[i])
	}

	condition := mysql.Bool(false)
	if publicOk {
		condition = tCalendar.Public.IS_TRUE()
	}

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarRSVP.EntryID,
		).
		FROM(tCalendarEntry.
			LEFT_JOIN(tCalendarRSVP,
				tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			).
			INNER_JOIN(tCalendar,
				mysql.AND(
					tCalendar.ID.EQ(tCalendarEntry.CalendarID),
					tCalendar.DeletedAt.IS_NULL(),
				),
			),
		).
		WHERE(mysql.AND(
			tCalendarEntry.DeletedAt.IS_NULL(),
			tCalendarRSVP.EntryID.IN(ids...),
			mysql.OR(
				mysql.AND(
					tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarEntry.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
				),
				tCalendarRSVP.EntryID.IS_NOT_NULL(),
				condition,
			),
		)).
		ORDER_BY(tCalendarEntry.ID.DESC())

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
