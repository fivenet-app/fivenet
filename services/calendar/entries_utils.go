package calendar

import (
	"context"
	"errors"
	"slices"
	"time"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

type calendarEntryOccurrenceKey struct {
	entryID       int64
	occurrenceKey string
}

func calendarEntryRSVPVisible(
	userInfo *userinfo.UserInfo,
	rsvpResponse calendarentries.RsvpResponses,
) mysql.BoolExpression {
	tCalendarEntryRsvp := tCalendarRSVP.AS("calendar_entry_rsvp_visibility")
	tCalendarEntryRsvpOccurrence := tCalendarRSVPOccurrence.AS(
		"calendar_entry_rsvp_occurrence_visibility",
	)

	return mysql.OR(
		mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tCalendarEntryRsvp).
				WHERE(mysql.AND(
					tCalendarEntryRsvp.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarEntryRsvp.Response.GT(mysql.Int32(int32(rsvpResponse))),
					tCalendarEntryRsvp.EntryID.EQ(tCalendarEntry.ID),
				)),
		),
		mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tCalendarEntryRsvpOccurrence).
				WHERE(mysql.AND(
					tCalendarEntryRsvpOccurrence.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarEntryRsvpOccurrence.Response.GT(mysql.Int32(int32(rsvpResponse))),
					tCalendarEntryRsvpOccurrence.EntryID.EQ(tCalendarEntry.ID),
				)),
		),
	)
}

func calendarEntryVisibility(
	userInfo *userinfo.UserInfo,
	access calendaraccess.AccessLevel,
	rsvpResponse calendarentries.RsvpResponses,
) mysql.BoolExpression {
	return mysql.OR(
		mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tCAccess).
				WHERE(mysql.AND(
					tCAccess.TargetID.EQ(tCalendarEntry.CalendarID),
					tCAccess.Access.GT_EQ(mysql.Int32(int32(access))),
					mysql.OR(
						tCAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						mysql.AND(
							tCAccess.Job.EQ(mysql.String(userInfo.GetJob())),
							tCAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
						),
					),
				)),
		),
		calendarEntryRSVPVisible(userInfo, rsvpResponse),
		tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
	)
}

func (s *Server) birthdayCalendarVisible(
	calendarID mysql.IntegerExpression,
	access calendaraccess.AccessLevel,
	userInfo *userinfo.UserInfo,
) mysql.BoolExpression {
	return mysql.AND(
		tCalendar.SystemKind.EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
		tCalendar.Job.EQ(mysql.String(userInfo.GetJob())),
		mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tCAccess).
				WHERE(mysql.AND(
					tCAccess.TargetID.EQ(calendarID),
					tCAccess.Access.GT_EQ(
						mysql.Int32(int32(access)),
					),
					mysql.OR(
						tCAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						mysql.AND(
							tCAccess.Job.EQ(mysql.String(userInfo.GetJob())),
							tCAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
						),
					),
				)),
		),
	)
}

func calendarEntriesQuery(
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
	visibility mysql.BoolExpression,
	limit *int64,
) mysql.SelectStatement {
	tCreator := table.FivenetUser.AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarEntry.CreatedAt,
			tCalendarEntry.UpdatedAt,
			tCalendarEntry.DeletedAt,
			tCalendarEntry.CalendarID,
			tCalendar.ID,
			tCalendar.Job,
			tCalendar.SystemKind,
			tCalendar.Name,
			tCalendar.Color,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendarEntry.Job,
			tCalendarEntry.StartTime,
			tCalendarEntry.EndTime,
			tCalendarEntry.Title,
			tCalendarEntry.Content,
			tCalendarEntry.Closed,
			tCalendarEntry.RsvpOpen,
			tCalendarEntry.Recurring,
			tCalendarEntry.RecurringUntil,
			tCalendarEntry.RecurrenceVersion,
			tCalendarEntry.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
			tAvatar.FilePath.AS("creator.profile_picture"),
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
		).
		FROM(tCalendarEntry.
			INNER_JOIN(tCalendar,
				mysql.AND(
					tCalendar.ID.EQ(tCalendarEntry.CalendarID),
					tCalendar.DeletedAt.IS_NULL(),
				),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tCalendarRSVP,
				mysql.AND(
					tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarRSVP.EntryID.EQ(tCalendarEntry.ID),
				),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(mysql.AND(
			visibility,
			condition,
		)).
		ORDER_BY(
			tCalendarEntry.StartTime.ASC(),
			tCalendarEntry.ID.ASC(),
		)

	if limit != nil && *limit > 0 {
		stmt = stmt.LIMIT(*limit)
	}

	return stmt
}

func (s *Server) loadExpandedCalendarEntries(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
	visibility mysql.BoolExpression,
	rangeStart, rangeEnd time.Time,
	limit *int64,
) ([]*calendarentries.CalendarEntry, error) {
	stmt := calendarEntriesQuery(userInfo, condition, visibility, limit)

	entries := []*calendarentries.CalendarEntry{}
	if err := stmt.QueryContext(ctx, s.db, &entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	expanded, err := s.expandCalendarEntries(ctx, userInfo, entries, rangeStart, rangeEnd)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if err := s.overlayCalendarEntryRSVPOverrides(ctx, userInfo, expanded); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return expanded, nil
}

func (s *Server) overlayCalendarEntryRSVPOverrides(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	entries []*calendarentries.CalendarEntry,
) error {
	occurrenceKeys := map[calendarEntryOccurrenceKey]struct{}{}
	entryIDs := map[int64]struct{}{}
	for i := range entries {
		occurrence := entries[i].GetOccurrence()
		if occurrence == nil ||
			occurrence.GetKind() != calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING ||
			occurrence.GetKey() == "" ||
			occurrence.GetSourceEntryId() <= 0 {
			continue
		}

		key := calendarEntryOccurrenceKey{
			entryID:       occurrence.GetSourceEntryId(),
			occurrenceKey: occurrence.GetKey(),
		}
		occurrenceKeys[key] = struct{}{}
		entryIDs[key.entryID] = struct{}{}
	}

	if len(occurrenceKeys) == 0 {
		return nil
	}

	entryIDExprs := make([]mysql.Expression, 0, len(entryIDs))
	for entryID := range entryIDs {
		entryIDExprs = append(entryIDExprs, mysql.Int64(entryID))
	}

	stmt := tCalendarRSVPOccurrence.
		SELECT(
			tCalendarRSVPOccurrence.EntryID,
			tCalendarRSVPOccurrence.OccurrenceKey,
			tCalendarRSVPOccurrence.CreatedAt,
			tCalendarRSVPOccurrence.UserID,
			tCalendarRSVPOccurrence.Response,
		).
		FROM(tCalendarRSVPOccurrence).
		WHERE(mysql.AND(
			tCalendarRSVPOccurrence.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			tCalendarRSVPOccurrence.EntryID.IN(entryIDExprs...),
		))

	overrides := []*calendarentries.CalendarEntryRSVP{}
	if err := stmt.QueryContext(ctx, s.db, &overrides); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	overrideByKey := make(
		map[calendarEntryOccurrenceKey]*calendarentries.CalendarEntryRSVP,
		len(overrides),
	)
	for i := range overrides {
		if overrides[i] == nil || overrides[i].GetOccurrenceKey() == "" {
			continue
		}
		overrideByKey[calendarEntryOccurrenceKey{
			entryID:       overrides[i].GetEntryId(),
			occurrenceKey: overrides[i].GetOccurrenceKey(),
		}] = overrides[i]
	}

	for i := range entries {
		occurrence := entries[i].GetOccurrence()
		if occurrence == nil ||
			occurrence.GetKind() != calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING ||
			occurrence.GetKey() == "" ||
			occurrence.GetSourceEntryId() <= 0 {
			continue
		}

		override := overrideByKey[calendarEntryOccurrenceKey{
			entryID:       occurrence.GetSourceEntryId(),
			occurrenceKey: occurrence.GetKey(),
		}]
		if override == nil {
			continue
		}

		entries[i].Rsvp = override
	}

	return nil
}

func (s *Server) finalizeCalendarEntries(
	entries []*calendarentries.CalendarEntry,
	userInfo *userinfo.UserInfo,
) []*calendarentries.CalendarEntry {
	slices.SortFunc(entries, func(left, right *calendarentries.CalendarEntry) int {
		l := left.GetStartTime().AsTime()
		r := right.GetStartTime().AsTime()
		if l.Before(r) {
			return -1
		}
		if l.After(r) {
			return 1
		}
		if left.GetCalendarId() < right.GetCalendarId() {
			return -1
		}
		if left.GetCalendarId() > right.GetCalendarId() {
			return 1
		}
		if left.GetId() < right.GetId() {
			return -1
		}
		if left.GetId() > right.GetId() {
			return 1
		}
		return 0
	})

	if s.enricher == nil {
		return entries
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range entries {
		if entries[i].GetCreator() != nil {
			jobInfoFn(entries[i].GetCreator())
		}
	}

	return entries
}

func filterUpcomingCalendarEntries(
	entries []*calendarentries.CalendarEntry,
	userInfo *userinfo.UserInfo,
) []*calendarentries.CalendarEntry {
	if len(entries) == 0 {
		return entries
	}

	filtered := make([]*calendarentries.CalendarEntry, 0, len(entries))
	for i := range entries {
		if entries[i] == nil {
			continue
		}

		if entries[i].GetOccurrence() != nil &&
			entries[i].GetOccurrence().
				GetKind() ==
				calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_BIRTHDAY {
			filtered = append(filtered, entries[i])
			continue
		}

		if entries[i].GetCreatorId() == userInfo.GetUserId() {
			filtered = append(filtered, entries[i])
			continue
		}

		if entries[i].GetRsvp() != nil &&
			entries[i].GetRsvp().
				GetResponse() >
				calendarentries.RsvpResponses_RSVP_RESPONSES_NO {
			filtered = append(filtered, entries[i])
		}
	}

	return filtered
}

func recurrenceShapeChanged(
	oldEntry *calendarentries.CalendarEntry,
	newEntry *calendarentries.CalendarEntry,
) bool {
	if !timestampEqual(oldEntry.GetStartTime(), newEntry.GetStartTime()) {
		return true
	}

	if !timestampEqual(oldEntry.GetEndTime(), newEntry.GetEndTime()) {
		return true
	}

	if !proto.Equal(oldEntry.GetRecurring(), newEntry.GetRecurring()) {
		return true
	}

	// If recurring_until is represented separately from Recurring.Until,
	// compare it here too.
	if !timestampEqual(oldEntry.GetRecurringUntil(), newEntry.GetRecurringUntil()) {
		return true
	}

	return false
}

func timestampEqual(a, b *timestamp.Timestamp) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.AsTime().Equal(b.AsTime())
}
