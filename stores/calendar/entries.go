package calendarstore

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jinzhu/now"
	"google.golang.org/protobuf/proto"
)

const maxCalendarEntriesLimit = int64(125)

type calendarEntryOccurrenceKey struct {
	entryID       int64
	occurrenceKey string
}

type calendarAccessEntry struct {
	ID     int64 `alias:"calendar.id"`
	Public bool  `alias:"calendar.public"`
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
	acl *access.SubjectObjectAccess,
	userInfo *userinfo.UserInfo,
	accessLevel calendaraccess.AccessLevel,
	rsvpResponse calendarentries.RsvpResponses,
) mysql.BoolExpression {
	_ = acl
	_ = accessLevel

	return mysql.OR(
		mysql.AND(
			tCalendar.Job.IS_NULL(),
			tCalendar.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
		),
		calendarEntryRSVPVisible(userInfo, rsvpResponse),
	)
}

func calendarEntriesQuery(
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
	visibility mysql.BoolExpression,
	ctes []mysql.CommonTableExpression,
	includeDeleted bool,
	limit *int64,
) mysql.Statement {
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
					mysql.OR(
						mysql.Bool(includeDeleted),
						tCalendar.DeletedAt.IS_NULL(),
					),
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
			mysql.OR(
				mysql.Bool(includeDeleted),
				tCalendarEntry.DeletedAt.IS_NULL(),
			),
			condition,
		)).
		ORDER_BY(
			tCalendarEntry.StartTime.ASC(),
			tCalendarEntry.ID.ASC(),
		)

	if limit != nil && *limit > 0 {
		stmt = stmt.
			LIMIT(*limit)
	}

	if len(ctes) > 0 {
		return mysql.WITH(ctes...)(stmt)
	}

	return stmt
}

func visibilityIDsSelect(query access.VisibilityQuery) mysql.SelectStatement {
	id := mysql.IntegerColumn("id").From(query.Table)
	return mysql.SELECT(id).FROM(query.Table)
}

func (s *Store) listCalendarEntriesVisibility(
	userInfo *userinfo.UserInfo,
	includeDeleted bool,
	accessLevel calendaraccess.AccessLevel,
	rsvpResponse calendarentries.RsvpResponses,
) (mysql.BoolExpression, mysql.BoolExpression, []mysql.CommonTableExpression) {
	visibleIDs := s.access.VisibleIDsByConditionQuery(
		userInfo,
		int32(accessLevel),
		includeDeleted,
		mysql.OR(
			tCalendar.SystemKind.IS_NULL(),
			tCalendar.SystemKind.NOT_EQ(
				mysql.Int32(
					int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
				),
			),
			mysql.AND(
				tCalendar.SystemKind.EQ(
					mysql.Int32(
						int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
					),
				),
				mysql.OR(
					mysql.Bool(userInfo.GetSuperuser()),
					tCalendar.Job.EQ(mysql.String(userInfo.GetJob())),
				),
			),
		),
	)
	visibleIDStmt := visibilityIDsSelect(visibleIDs)
	visibleCalendarCondition := tCalendar.ID.IN(visibleIDStmt)

	visibility := mysql.OR(
		tCalendar.ID.IN(
			tCalendarSubs.
				SELECT(tCalendarSubs.CalendarID).
				FROM(tCalendarSubs).
				WHERE(mysql.AND(tCalendarSubs.UserID.EQ(mysql.Int32(userInfo.GetUserId())))),
		),
		mysql.AND(
			tCalendar.Job.IS_NULL(),
			tCalendar.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
		),
		visibleCalendarCondition,
		calendarEntryRSVPVisible(userInfo, rsvpResponse),
	)

	return visibility, visibleCalendarCondition, visibleIDs.CTEs
}

func (s *Store) ListCalendarEntries(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	opts ListCalendarEntriesOptions,
) ([]*calendarentries.CalendarEntry, error) {
	rsvpResponse := calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN
	if opts.ShowHidden {
		rsvpResponse = calendarentries.RsvpResponses_RSVP_RESPONSES_UNSPECIFIED
	}
	includeDeleted := userInfo.GetSuperuser()
	visibility, visibleCalendarCondition, ctes := s.listCalendarEntriesVisibility(
		userInfo,
		includeDeleted,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		rsvpResponse,
	)

	condition := mysql.AND(
		mysql.OR(
			tCalendar.ID.IN(
				tCalendarSubs.
					SELECT(tCalendarSubs.CalendarID).
					FROM(tCalendarSubs).
					WHERE(mysql.AND(tCalendarSubs.UserID.EQ(mysql.Int32(userInfo.GetUserId())))),
			),
			tCalendarEntry.ID.IN(
				tCalendarRSVP.
					SELECT(tCalendarRSVP.EntryID).
					FROM(tCalendarRSVP).
					WHERE(mysql.AND(
						tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tCalendarRSVP.Response.GT(mysql.Int32(int32(rsvpResponse))),
					)),
			),
			tCalendarEntry.ID.IN(
				tCalendarRSVPOccurrence.
					SELECT(tCalendarRSVPOccurrence.EntryID).
					FROM(tCalendarRSVPOccurrence).
					WHERE(mysql.AND(
						tCalendarRSVPOccurrence.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tCalendarRSVPOccurrence.Response.GT(mysql.Int32(int32(rsvpResponse))),
						tCalendarRSVPOccurrence.RecurrenceVersion.EQ(
							tCalendarEntry.RecurrenceVersion,
						),
					)),
			),
			mysql.AND(
				tCalendar.Job.IS_NULL(),
				tCalendar.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
			),
			visibleCalendarCondition,
			calendarEntryRSVPVisible(userInfo, rsvpResponse),
		),
	)

	if opts.After != nil {
		condition = condition.AND(
			tCalendar.UpdatedAt.GT_EQ(mysql.TimestampT(opts.After.AsTime())),
		)
	}

	baseDate := now.New(
		time.Date(int(opts.Year), time.Month(opts.Month), 1, 0, 0, 0, 0, time.UTC),
	)
	startDate := baseDate.BeginningOfMonth()
	endDate := baseDate.EndOfMonth()

	condition = condition.AND(tCalendarEntry.StartTime.LT_EQ(mysql.DateTimeT(endDate)))

	dateWindowCondition := mysql.OR(
		mysql.AND(
			tCalendarEntry.Recurring.IS_NULL(),
			mysql.OR(
				tCalendarEntry.EndTime.IS_NULL().
					AND(tCalendarEntry.StartTime.GT_EQ(mysql.DateTimeT(startDate))),
				tCalendarEntry.EndTime.GT_EQ(mysql.DateTimeT(startDate)),
			),
		),
		mysql.AND(
			tCalendarEntry.Recurring.IS_NOT_NULL(),
			mysql.OR(
				tCalendarEntry.RecurringUntil.IS_NULL(),
				tCalendarEntry.RecurringUntil.GT_EQ(mysql.DateTimeT(startDate)),
			),
		),
	)

	regularCondition := condition.AND(dateWindowCondition).AND(mysql.OR(
		tCalendar.SystemKind.IS_NULL(),
		tCalendar.SystemKind.NOT_EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
	))

	birthdayCondition := mysql.AND(
		condition,
		tCalendar.SystemKind.EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
		mysql.OR(
			mysql.Bool(userInfo.GetSuperuser()),
			tCalendar.Job.EQ(mysql.String(userInfo.GetJob())),
		),
	)

	if len(opts.CalendarIDs) > 0 {
		ids := []mysql.Expression{}
		for i := range opts.CalendarIDs {
			if opts.CalendarIDs[i] == 0 {
				continue
			}
			ids = append(ids, mysql.Int64(opts.CalendarIDs[i]))
		}

		regularCondition = regularCondition.AND(tCalendarEntry.CalendarID.IN(ids...))
		birthdayCondition = birthdayCondition.AND(tCalendarEntry.CalendarID.IN(ids...))
	}

	regularEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		regularCondition,
		visibility,
		ctes,
		startDate,
		endDate,
		new(maxCalendarEntriesLimit),
		includeDeleted,
	)
	if err != nil {
		return nil, err
	}

	birthdayEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		birthdayCondition,
		visibility,
		ctes,
		startDate,
		endDate,
		nil,
		includeDeleted,
	)
	if err != nil {
		return nil, err
	}

	return append(regularEntries, birthdayEntries...), nil
}

func (s *Store) GetUpcomingEntries(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	opts GetUpcomingEntriesOptions,
) ([]*calendarentries.CalendarEntry, error) {
	rangeStart := time.Now().Add(-1 * time.Minute)
	rangeEnd := time.Now().Add(time.Duration(opts.Seconds) * time.Second)
	includeDeleted := userInfo.GetSuperuser()
	visibility, visibleCalendarCondition, ctes := s.listCalendarEntriesVisibility(
		userInfo,
		includeDeleted,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN,
	)

	condition := mysql.AND(
		mysql.OR(
			tCalendarEntry.ID.IN(
				tCalendarRSVP.
					SELECT(tCalendarRSVP.EntryID).
					FROM(tCalendarRSVP).
					WHERE(mysql.AND(
						tCalendarRSVP.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tCalendarRSVP.Response.GT(
							mysql.Int32(int32(calendarentries.RsvpResponses_RSVP_RESPONSES_NO)),
						),
					)),
			),
			tCalendarEntry.ID.IN(
				tCalendarRSVPOccurrence.
					SELECT(tCalendarRSVPOccurrence.EntryID).
					FROM(tCalendarRSVPOccurrence).
					WHERE(mysql.AND(
						tCalendarRSVPOccurrence.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tCalendarRSVPOccurrence.Response.GT(
							mysql.Int32(int32(calendarentries.RsvpResponses_RSVP_RESPONSES_NO)),
						),
						tCalendarRSVPOccurrence.RecurrenceVersion.EQ(
							tCalendarEntry.RecurrenceVersion,
						),
					)),
			),
			visibleCalendarCondition,
			tCalendarEntry.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
		),
		tCalendarEntry.StartTime.LT_EQ(mysql.TimestampT(rangeEnd)),
	)

	regularCondition := condition.AND(mysql.OR(
		tCalendar.SystemKind.IS_NULL(),
		tCalendar.SystemKind.NOT_EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
	))
	birthdayCondition := mysql.AND(
		condition,
		tCalendar.SystemKind.EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
		mysql.OR(
			mysql.Bool(userInfo.GetSuperuser()),
			tCalendar.Job.EQ(mysql.String(userInfo.GetJob())),
		),
	)

	regularEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		regularCondition,
		visibility,
		ctes,
		rangeStart,
		rangeEnd,
		new(maxCalendarEntriesLimit),
		includeDeleted,
	)
	if err != nil {
		return nil, err
	}

	birthdayEntries, err := s.loadExpandedCalendarEntries(
		ctx,
		userInfo,
		birthdayCondition,
		visibility,
		ctes,
		rangeStart,
		rangeEnd,
		nil,
		includeDeleted,
	)
	if err != nil {
		return nil, err
	}

	return append(regularEntries, birthdayEntries...), nil
}

func (s *Store) GetEntry(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
) (*calendarentries.CalendarEntry, error) {
	tCreator := table.FivenetUser.AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	includeDeleted := userInfo.GetSuperuser()

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarEntry.CreatedAt,
			tCalendarEntry.UpdatedAt,
			tCalendarEntry.DeletedAt,
			tCalendarEntry.CalendarID,
			tCalendar.ID,
			tCalendar.Name,
			tCalendar.Color,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
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
			tCalendarEntry.CreatorJob,
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
					mysql.OR(
						mysql.Bool(includeDeleted),
						tCalendar.DeletedAt.IS_NULL(),
					),
				),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCalendarEntry.CreatorID),
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
		WHERE(condition).
		LIMIT(1)

	dest := &calendarentries.CalendarEntry{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Store) UpsertCalendarEntry(
	ctx context.Context,
	tx qrm.DB,
	entry *calendarentries.CalendarEntry,
	oldEntry *calendarentries.CalendarEntry,
	userInfo *userinfo.UserInfo,
) (int64, error) {
	tCalendarEntry := table.FivenetCalendarEntries

	if entry.GetId() > 0 {
		values := []interface{}{
			mysql.String(entry.GetTitle()),
			entry.GetContent(),
			dbutils.TimestampToMySQL(entry.GetStartTime()),
			dbutils.TimestampToMySQL(entry.GetEndTime()),
			mysql.Bool(entry.GetClosed()),
			mysql.Bool(entry.GetRsvpOpen()),
			entry.GetRecurring(),
			dbutils.TimestampToMySQL(entry.GetRecurring().GetUntil()),
		}

		if recurrenceShapeChanged(oldEntry, entry) {
			values = append(
				values,
				tCalendarEntry.RecurrenceVersion.SET(
					tCalendarEntry.RecurrenceVersion.ADD(mysql.Int32(1)),
				),
			)
		} else {
			values = append(values, oldEntry.RecurrenceVersion)
		}

		stmt := tCalendarEntry.
			UPDATE(
				tCalendarEntry.Title,
				tCalendarEntry.Content,
				tCalendarEntry.StartTime,
				tCalendarEntry.EndTime,
				tCalendarEntry.Closed,
				tCalendarEntry.RsvpOpen,
				tCalendarEntry.Recurring,
				tCalendarEntry.RecurringUntil,
				tCalendarEntry.RecurrenceVersion,
			).
			SET(values[0], values[1:]...).
			WHERE(mysql.AND(
				tCalendarEntry.ID.EQ(mysql.Int64(entry.GetId())),
				tCalendarEntry.CalendarID.EQ(mysql.Int64(entry.GetCalendarId())),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return 0, err
		}

		return entry.GetId(), nil
	}

	stmt := tCalendarEntry.
		INSERT(
			tCalendarEntry.CalendarID,
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
			tCalendarEntry.CreatorJob,
		).
		VALUES(
			entry.GetCalendarId(),
			userInfo.GetJob(),
			entry.GetStartTime(),
			entry.GetEndTime(),
			entry.GetTitle(),
			entry.GetContent(),
			entry.GetClosed(),
			entry.GetRsvpOpen(),
			entry.GetRecurring(),
			entry.GetRecurring().GetUntil(),
			1,
			userInfo.GetUserId(),
			userInfo.GetJob(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (s *Store) DeleteCalendarEntry(
	ctx context.Context,
	tx qrm.DB,
	entryID int64,
	calendarID int64,
	deletedAt *timestamp.Timestamp,
) error {
	tCalendarEntry := table.FivenetCalendarEntries

	stmt := tCalendarEntry.
		UPDATE(
			tCalendarEntry.DeletedAt,
		).
		SET(
			tCalendarEntry.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt)),
		).
		WHERE(mysql.AND(
			tCalendarEntry.CalendarID.EQ(mysql.Int64(calendarID)),
			tCalendarEntry.ID.EQ(mysql.Int64(entryID)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) loadExpandedCalendarEntries(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
	visibility mysql.BoolExpression,
	ctes []mysql.CommonTableExpression,
	rangeStart, rangeEnd time.Time,
	limit *int64,
	includeDeleted bool,
) ([]*calendarentries.CalendarEntry, error) {
	stmt := calendarEntriesQuery(userInfo, condition, visibility, ctes, includeDeleted, limit)

	entries := []*calendarentries.CalendarEntry{}
	if err := stmt.QueryContext(ctx, s.db, &entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	expanded, err := s.expandCalendarEntries(ctx, userInfo, entries, rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}

	if err := s.overlayCalendarEntryRSVPOverrides(ctx, userInfo, expanded); err != nil {
		return nil, err
	}

	return expanded, nil
}

func (s *Store) expandCalendarEntries(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	entries []*calendarentries.CalendarEntry,
	rangeStart, rangeEnd time.Time,
) ([]*calendarentries.CalendarEntry, error) {
	out := make([]*calendarentries.CalendarEntry, 0, len(entries))
	for i := range entries {
		if entries[i] == nil {
			continue
		}

		expanded, err := s.expandCalendarEntryOccurrences(
			ctx,
			userInfo,
			entries[i],
			rangeStart,
			rangeEnd,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, expanded...)
	}

	sort.SliceStable(out, func(i, j int) bool {
		left := out[i].GetStartTime().AsTime()
		right := out[j].GetStartTime().AsTime()
		if left.Equal(right) {
			if out[i].GetCalendarId() == out[j].GetCalendarId() {
				return out[i].GetId() < out[j].GetId()
			}
			return out[i].GetCalendarId() < out[j].GetCalendarId()
		}
		return left.Before(right)
	})

	return out, nil
}

func (s *Store) expandCalendarEntryOccurrences(
	_ context.Context,
	_ *userinfo.UserInfo,
	entry *calendarentries.CalendarEntry,
	rangeStart, rangeEnd time.Time,
) ([]*calendarentries.CalendarEntry, error) {
	if entry == nil || entry.GetStartTime() == nil {
		return nil, nil
	}

	if entry.GetRecurring() == nil {
		if !entryOverlapsRange(
			entry.GetStartTime().AsTime(),
			entry.GetEndTime(),
			rangeStart,
			rangeEnd,
		) {
			return nil, nil
		}

		clone := proto.Clone(entry).(*calendarentries.CalendarEntry)
		clone.Occurrence = &calendarentries.CalendarEntryOccurrence{
			Key: fmt.Sprintf(
				"manual:%d:%d",
				clone.GetId(),
				clone.GetStartTime().AsTime().Unix(),
			),
			Kind:          calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_MANUAL,
			SourceEntryId: &clone.Id,
			AllDay:        clone.GetEndTime() == nil,
		}

		return []*calendarentries.CalendarEntry{clone}, nil
	}

	return s.expandRecurringEntry(entry, rangeStart, rangeEnd), nil
}

func (s *Store) expandRecurringEntry(
	entry *calendarentries.CalendarEntry,
	rangeStart, rangeEnd time.Time,
) []*calendarentries.CalendarEntry {
	if entry == nil || entry.GetStartTime() == nil || entry.GetRecurring() == nil {
		return nil
	}

	interval := entry.GetRecurring().GetCount()
	if interval <= 0 {
		interval = 1
	}

	duration := time.Duration(0)
	if entry.GetEndTime() != nil {
		duration = entry.GetEndTime().AsTime().Sub(entry.GetStartTime().AsTime())
	}

	occurrenceKind := calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING
	occurrenceKeyPrefix := "recurring"
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() ==
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS {
		occurrenceKind = calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_BIRTHDAY
		occurrenceKeyPrefix = "birthday"
	}

	var sourceUserID *int32
	if occurrenceKind == calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_BIRTHDAY {
		if creatorID := entry.GetCreatorId(); creatorID > 0 {
			sourceUserID = &creatorID
		}
	}

	out := []*calendarentries.CalendarEntry{}
	occurrenceStart := entry.GetStartTime().AsTime()
	for !occurrenceStart.After(rangeEnd) {
		if until := entry.GetRecurring().
			GetUntil(); until != nil &&
			occurrenceStart.After(until.AsTime()) {
			break
		}

		occurrenceEnd := (*timestamp.Timestamp)(nil)
		if entry.GetEndTime() != nil {
			end := occurrenceStart.Add(duration)
			occurrenceEnd = timestamp.New(end)
		}

		if entryOverlapsRange(occurrenceStart, occurrenceEnd, rangeStart, rangeEnd) {
			clone := proto.Clone(entry).(*calendarentries.CalendarEntry)
			clone.StartTime = timestamp.New(occurrenceStart)
			clone.EndTime = occurrenceEnd

			if entry.GetEndTime() != nil {
				end := occurrenceStart.Add(duration)
				clone.EndTime = timestamp.New(end)
			}

			key := fmt.Sprintf(
				"%s:%d:%d:%d",
				occurrenceKeyPrefix,
				entry.GetId(),
				entry.GetRecurrenceVersion(),
				occurrenceStart.Unix(),
			)
			if sourceUserID != nil {
				key = fmt.Sprintf(
					"%s:%d:%d:%04d:%02d:%02d",
					occurrenceKeyPrefix,
					entry.GetCalendarId(),
					*sourceUserID,
					occurrenceStart.Year(),
					occurrenceStart.Month(),
					occurrenceStart.Day(),
				)
			}
			clone.Occurrence = &calendarentries.CalendarEntryOccurrence{
				Key:           key,
				Kind:          occurrenceKind,
				SourceEntryId: &clone.Id,
				SourceUserId:  sourceUserID,
				AllDay:        clone.GetEndTime() == nil,
			}
			out = append(out, clone)
		}

		occurrenceStart = nextRecurringOccurrence(
			occurrenceStart,
			interval,
			entry.GetRecurring().GetEvery(),
		)
	}

	return out
}

func (s *Store) overlayCalendarEntryRSVPOverrides(
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
		overrideByKey[calendarEntryOccurrenceKey{entryID: overrides[i].GetEntryId(), occurrenceKey: overrides[i].GetOccurrenceKey()}] = overrides[i]
	}

	for i := range entries {
		occurrence := entries[i].GetOccurrence()
		if occurrence == nil ||
			occurrence.GetKind() != calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING ||
			occurrence.GetKey() == "" ||
			occurrence.GetSourceEntryId() <= 0 {
			continue
		}

		override := overrideByKey[calendarEntryOccurrenceKey{entryID: occurrence.GetSourceEntryId(), occurrenceKey: occurrence.GetKey()}]
		if override == nil {
			continue
		}

		entries[i].Rsvp = override
	}

	return nil
}

func nextRecurringOccurrence(
	start time.Time,
	interval int32,
	every calendarentries.CalendarEntryRecurringEvery,
) time.Time {
	switch every {
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_DAY:
		return start.AddDate(0, 0, int(interval))
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_WEEK:
		return start.AddDate(0, 0, 7*int(interval))
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_MONTH:
		return start.AddDate(0, int(interval), 0)
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_YEAR:
		return start.AddDate(int(interval), 0, 0)
	default:
		return start.AddDate(0, 0, int(interval))
	}
}

func entryOverlapsRange(
	start time.Time,
	end *timestamp.Timestamp,
	rangeStart, rangeEnd time.Time,
) bool {
	if start.After(rangeEnd) {
		return false
	}

	if end == nil {
		return !start.Before(rangeStart) && !start.After(rangeEnd)
	}

	endTime := end.AsTime()
	return !endTime.Before(rangeStart) && !start.After(rangeEnd)
}

func recurrenceShapeChanged(
	oldEntry *calendarentries.CalendarEntry,
	newEntry *calendarentries.CalendarEntry,
) bool {
	if oldEntry == nil || newEntry == nil {
		return false
	}

	if !timestampEqual(oldEntry.GetStartTime(), newEntry.GetStartTime()) {
		return true
	}

	if !timestampEqual(oldEntry.GetEndTime(), newEntry.GetEndTime()) {
		return true
	}

	if !proto.Equal(oldEntry.GetRecurring(), newEntry.GetRecurring()) {
		return true
	}

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

func (s *Store) FilterUpcomingCalendarEntries(
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
