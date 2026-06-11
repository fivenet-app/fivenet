package calendar

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListCalendarEntryRSVP(
	ctx context.Context,
	req *pbcalendar.ListCalendarEntryRSVPRequest,
) (*pbcalendar.ListCalendarEntryRSVPResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntryId())))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrFailedQuery
	}
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() !=
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.checkIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		entry.GetId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	tUser := table.FivenetUser.AS("user_short")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	condition := mysql.AND(
		tCalendarRSVP.EntryID.EQ(mysql.Int64(entry.GetId())),
		tCalendarRSVP.Response.GT(
			mysql.Int32(int32(calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN)),
		),
	)

	countStmt := tCalendarRSVP.
		SELECT(
			mysql.COUNT(tCalendarRSVP.UserID).AS("data_count.total"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponse(count.Total)
	resp := &pbcalendar.ListCalendarEntryRSVPResponse{
		Pagination: pag,
	}

	if count.Total <= 0 {
		return resp, nil
	}

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tUserProps.AvatarFileID.AS("user_short.profile_picture_file_id"),
			tAvatar.FilePath.AS("user_short.profile_picture"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(condition).
		ORDER_BY(tCalendarRSVP.Response.DESC()).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetEntries() {
		if resp.GetEntries()[i].GetUser() != nil {
			jobInfoFn(resp.GetEntries()[i].GetUser())
		}
	}

	return resp, nil
}

func (s *Server) RSVPCalendarEntry(
	ctx context.Context,
	req *pbcalendar.RSVPCalendarEntryRequest,
) (*pbcalendar.RSVPCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.GetEntry() == nil {
		return nil, errorscalendar.ErrNoPerms
	}

	entry, err := s.getEntry(
		ctx,
		userInfo,
		tCalendarEntry.ID.EQ(mysql.Int64(req.GetEntry().GetEntryId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorscalendar.ErrFailedQuery
	}
	if entry.GetCalendar() != nil &&
		entry.GetCalendar().
			GetSystemKind() !=
			calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		return nil, errorscalendar.ErrNoPerms
	}

	check, err := s.checkIfUserHasAccessToCalendarEntry(
		ctx,
		entry.GetCalendarId(),
		entry.GetId(),
		userInfo,
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errorscalendar.ErrNoPerms
	}

	if entry.GetClosed() {
		return nil, errorscalendar.ErrEntryClosed
	}

	occurrenceKey := strings.TrimSpace(req.GetOccurrenceKey())
	if occurrenceKey == "" && req.GetEntry() != nil {
		occurrenceKey = strings.TrimSpace(req.GetEntry().GetOccurrenceKey())
	}

	if occurrenceKey != "" {
		if entry.GetRecurring() == nil {
			return nil, errorscalendar.ErrNoPerms
		}
		if err := validateRecurringOccurrenceKey(entry, occurrenceKey); err != nil {
			return nil, err
		}
	}

	if req.Remove != nil && req.GetRemove() && occurrenceKey == "" {
		req.Entry.Response = calendarentries.RsvpResponses_RSVP_RESPONSES_HIDDEN
	}

	if occurrenceKey != "" {
		tCalendarRsvpOccurrence := table.FivenetCalendarRsvpOccurrence

		identity, err := parseRecurringOccurrenceKey(occurrenceKey)
		if err != nil {
			return nil, err
		}

		if identity.EntryID != req.GetEntry().GetEntryId() {
			return nil, errorscalendar.ErrNoPerms
		}

		if req.Remove != nil && req.GetRemove() {
			stmt := tCalendarRsvpOccurrence.
				DELETE().
				WHERE(mysql.AND(
					tCalendarRsvpOccurrence.EntryID.EQ(mysql.Int64(identity.EntryID)),
					tCalendarRsvpOccurrence.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tCalendarRsvpOccurrence.RecurrenceVersion.EQ(
						mysql.Int32(identity.RecurrenceVersion),
					),
					tCalendarRsvpOccurrence.RecurrenceID.EQ(
						mysql.TimestampT(identity.RecurrenceID),
					),
				)).
				LIMIT(1)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
			}
		} else {
			stmt := tCalendarRsvpOccurrence.
				INSERT(
					tCalendarRsvpOccurrence.EntryID,
					tCalendarRsvpOccurrence.RecurrenceVersion,
					tCalendarRsvpOccurrence.RecurrenceID,
					tCalendarRsvpOccurrence.OccurrenceKey,
					tCalendarRsvpOccurrence.UserID,
					tCalendarRsvpOccurrence.Response,
				).
				VALUES(
					identity.EntryID,
					identity.RecurrenceVersion,
					identity.RecurrenceID,
					occurrenceKey,
					userInfo.GetUserId(),
					req.GetEntry().GetResponse(),
				).
				ON_DUPLICATE_KEY_UPDATE(
					tCalendarRsvpOccurrence.Response.SET(
						mysql.Int32(int32(req.GetEntry().GetResponse())),
					),
					tCalendarRsvpOccurrence.OccurrenceKey.SET(mysql.String(occurrenceKey)),
				)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
				}
			}
		}
	} else {
		tCalendarRSVP := table.FivenetCalendarRsvp
		stmt := tCalendarRSVP.
			INSERT(
				tCalendarRSVP.EntryID,
				tCalendarRSVP.UserID,
				tCalendarRSVP.Response,
			).
			VALUES(
				req.GetEntry().GetEntryId(),
				userInfo.GetUserId(),
				req.GetEntry().GetResponse(),
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCalendarRSVP.Response.SET(mysql.Int32(int32(req.GetEntry().GetResponse()))),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
			}
		}
	}

	rsvpEntry, err := s.getRSVPCalendarEntry(
		ctx,
		req.GetEntry().GetEntryId(),
		userInfo.GetUserId(),
		occurrenceKey,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	if rsvpEntry.GetUser() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, rsvpEntry.GetUser())
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbcalendar.RSVPCalendarEntryResponse{
		Entry: rsvpEntry,
	}, nil
}

func (s *Server) getRSVPCalendarEntry(
	ctx context.Context,
	entryId int64,
	userId int32,
	occurrenceKey string,
) (*calendarentries.CalendarEntryRSVP, error) {
	if occurrenceKey != "" {
		if dest, err := s.getOccurrenceRSVPCalendarEntry(
			ctx,
			entryId,
			userId,
			occurrenceKey,
		); err != nil {
			return nil, err
		} else if dest != nil {
			return dest, nil
		}
	}

	return s.getSeriesRSVPCalendarEntry(ctx, entryId, userId)
}

func (s *Server) getSeriesRSVPCalendarEntry(
	ctx context.Context,
	entryId int64,
	userId int32,
) (*calendarentries.CalendarEntryRSVP, error) {
	tUser := table.FivenetUser.AS("user_short")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tCalendarRSVP.
		SELECT(
			tCalendarRSVP.EntryID,
			tCalendarRSVP.CreatedAt,
			tCalendarRSVP.UserID,
			tCalendarRSVP.Response,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tUserProps.AvatarFileID.AS("user_short.profile_picture_file_id"),
			tAvatar.FilePath.AS("user_short.profile_picture"),
		).
		FROM(tCalendarRSVP.
			LEFT_JOIN(tUser,
				tCalendarRSVP.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(mysql.AND(
			tCalendarRSVP.EntryID.EQ(mysql.Int64(entryId)),
			tCalendarRSVP.UserID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	var dest calendarentries.CalendarEntryRSVP
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetEntryId() == 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Server) getOccurrenceRSVPCalendarEntry(
	ctx context.Context,
	entryId int64,
	userId int32,
	occurrenceKey string,
) (*calendarentries.CalendarEntryRSVP, error) {
	identity, err := parseRecurringOccurrenceKey(occurrenceKey)
	if err != nil {
		return nil, err
	}

	if identity.EntryID != entryId {
		return nil, errorscalendar.ErrNoPerms
	}

	tUser := table.FivenetUser.AS("user_short")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	tOccurrence := table.FivenetCalendarRsvpOccurrence

	stmt := tOccurrence.
		SELECT(
			tOccurrence.EntryID,
			tOccurrence.CreatedAt,
			tOccurrence.UserID,
			tOccurrence.Response,
			tOccurrence.OccurrenceKey,

			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,

			tUserProps.AvatarFileID.AS("user_short.profile_picture_file_id"),
			tAvatar.FilePath.AS("user_short.profile_picture"),
		).
		FROM(tOccurrence.
			LEFT_JOIN(tUser, tOccurrence.UserID.EQ(tUser.ID)).
			LEFT_JOIN(tUserProps, tUserProps.UserID.EQ(tUser.ID)).
			LEFT_JOIN(tAvatar, tAvatar.ID.EQ(tUserProps.AvatarFileID)),
		).
		WHERE(mysql.AND(
			tOccurrence.EntryID.EQ(mysql.Int64(entryId)),
			tOccurrence.UserID.EQ(mysql.Int32(userId)),
			tOccurrence.RecurrenceVersion.EQ(mysql.Int32(identity.RecurrenceVersion)),
			tOccurrence.RecurrenceID.EQ(mysql.TimestampT(identity.RecurrenceID)),
		)).
		LIMIT(1)

	var dest calendarentries.CalendarEntryRSVP
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetEntryId() == 0 {
		return nil, nil
	}

	return &dest, nil
}

type recurringOccurrenceIdentity struct {
	EntryID           int64
	RecurrenceVersion int32
	RecurrenceID      time.Time
	OccurrenceUnix    int64
}

func parseRecurringOccurrenceKey(
	occurrenceKey string,
) (*recurringOccurrenceIdentity, error) {
	parts := strings.Split(occurrenceKey, ":")
	if len(parts) != 4 || parts[0] != "recurring" {
		return nil, errorscalendar.ErrNoPerms
	}

	entryID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, errorscalendar.ErrNoPerms
	}

	version, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return nil, errorscalendar.ErrNoPerms
	}

	occurrenceUnix, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return nil, errorscalendar.ErrNoPerms
	}

	return &recurringOccurrenceIdentity{
		EntryID:           entryID,
		RecurrenceVersion: int32(version),
		OccurrenceUnix:    occurrenceUnix,
		RecurrenceID:      time.Unix(occurrenceUnix, 0).UTC(),
	}, nil
}

func validateRecurringOccurrenceKey(
	entry *calendarentries.CalendarEntry,
	occurrenceKey string,
) error {
	identity, err := parseRecurringOccurrenceKey(occurrenceKey)
	if err != nil {
		return err
	}

	if identity.EntryID != entry.GetId() {
		return errorscalendar.ErrNoPerms
	}

	if identity.RecurrenceVersion != entry.GetRecurrenceVersion() {
		return errorscalendar.ErrNoPerms
	}

	if entry.GetStartTime() == nil || entry.GetRecurring() == nil {
		return errorscalendar.ErrNoPerms
	}

	targetTime := time.Unix(identity.OccurrenceUnix, 0)
	occurrenceStart := entry.GetStartTime().AsTime()

	interval := entry.GetRecurring().GetCount()
	if interval <= 0 {
		interval = 1
	}

	for !occurrenceStart.After(targetTime) {
		if until := entry.GetRecurring().GetUntil(); until != nil &&
			occurrenceStart.After(until.AsTime()) {
			break
		}

		if occurrenceStart.Unix() == identity.OccurrenceUnix {
			return nil
		}

		occurrenceStart = nextRecurringOccurrence(
			occurrenceStart,
			interval,
			entry.GetRecurring().GetEvery(),
		)
	}

	return errorscalendar.ErrNoPerms
}
