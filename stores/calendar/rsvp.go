package calendarstore

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

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

func (s *Store) GetRSVPCalendarEntry(
	ctx context.Context,
	entryID int64,
	userID int32,
	occurrenceKey string,
) (*calendarentries.CalendarEntryRSVP, error) {
	if occurrenceKey != "" {
		if dest, err := s.getOccurrenceRSVPCalendarEntry(
			ctx,
			entryID,
			userID,
			occurrenceKey,
		); err != nil {
			return nil, err
		} else if dest != nil {
			return dest, nil
		}
	}

	return s.getSeriesRSVPCalendarEntry(ctx, entryID, userID)
}

func (s *Store) ListCalendarEntryRSVP(
	ctx context.Context,
	opts ListCalendarEntryRSVPOptions,
	userInfo *userinfo.UserInfo,
) (*pbcalendar.ListCalendarEntryRSVPResponse, error) {
	tUser := table.FivenetUser.AS("user_short")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	condition := mysql.AND(
		tCalendarRSVP.EntryID.EQ(mysql.Int64(opts.EntryID)),
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
			return nil, err
		}
	}

	pag, limit := opts.Pagination.GetResponse(count.Total)
	resp := &pbcalendar.ListCalendarEntryRSVPResponse{Pagination: pag}
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
		OFFSET(opts.Pagination.GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Store) SetCalendarEntryRSVP(
	ctx context.Context,
	entry *calendarentries.CalendarEntryRSVP,
	userInfo *userinfo.UserInfo,
	occurrenceKey string,
	remove bool,
) error {
	if occurrenceKey != "" {
		tCalendarRsvpOccurrence := table.FivenetCalendarRsvpOccurrence
		identity, err := parseRecurringOccurrenceKey(occurrenceKey)
		if err != nil {
			return err
		}

		if identity.EntryID != entry.GetEntryId() {
			return errorscalendar.ErrNoPerms
		}

		if remove {
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
				return err
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
					entry.GetResponse(),
				).
				ON_DUPLICATE_KEY_UPDATE(
					tCalendarRsvpOccurrence.Response.SET(mysql.Int32(int32(entry.GetResponse()))),
					tCalendarRsvpOccurrence.OccurrenceKey.SET(mysql.String(occurrenceKey)),
				)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}
		}
		return nil
	}

	tCalendarRsvp := table.FivenetCalendarRsvp
	stmt := tCalendarRsvp.
		INSERT(
			tCalendarRsvp.EntryID,
			tCalendarRsvp.UserID,
			tCalendarRsvp.Response,
		).
		VALUES(
			entry.GetEntryId(),
			userInfo.GetUserId(),
			entry.GetResponse(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCalendarRsvp.Response.SET(mysql.Int32(int32(entry.GetResponse()))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *Store) getSeriesRSVPCalendarEntry(
	ctx context.Context,
	entryID int64,
	userID int32,
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
			tCalendarRSVP.EntryID.EQ(mysql.Int64(entryID)),
			tCalendarRSVP.UserID.EQ(mysql.Int32(userID)),
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

func (s *Store) getOccurrenceRSVPCalendarEntry(
	ctx context.Context,
	entryID int64,
	userID int32,
	occurrenceKey string,
) (*calendarentries.CalendarEntryRSVP, error) {
	identity, err := parseRecurringOccurrenceKey(occurrenceKey)
	if err != nil {
		return nil, err
	}

	if identity.EntryID != entryID {
		return nil, errorscalendar.ErrNoPerms
	}

	tUser := table.FivenetUser.AS("user_short")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	tOccurrence := table.FivenetCalendarRsvpOccurrence.AS("calendar_entry_rsvp")

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
			tOccurrence.EntryID.EQ(mysql.Int64(entryID)),
			tOccurrence.UserID.EQ(mysql.Int32(userID)),
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

func (s *Store) ValidateRecurringOccurrenceKey(
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
