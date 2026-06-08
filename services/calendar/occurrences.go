package calendar

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

const (
	birthdayCalendarColor = "rose"
)

type birthdayColleague struct {
	UserID      int32  `alias:"user_id"`
	Firstname   string `alias:"firstname"`
	Lastname    string `alias:"lastname"`
	Dateofbirth string `alias:"dateofbirth"`
}

func (s *Server) ensureJobBirthdayCalendar(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
) (*calendarresource.Calendar, error) {
	job := strings.TrimSpace(userInfo.GetJob())
	if job == "" {
		return nil, nil
	}

	jobInfo := s.enricher.GetJobByName(job)
	title := birthdayCalendarTitle(s.i18n, s.appCfg, job, jobInfo)

	existing, err := s.getCalendar(
		ctx,
		userInfo,
		mysql.AND(
			tCalendar.DeletedAt.IS_NULL(),
			tCalendar.Job.EQ(mysql.String(job)),
			tCalendar.SystemKind.EQ(
				mysql.Int32(
					int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		tCalendar := table.FivenetCalendar
		if _, err := tCalendar.
			UPDATE(
				tCalendar.Name,
				tCalendar.Description,
				tCalendar.Public,
				tCalendar.Closed,
				tCalendar.CreatorID,
				tCalendar.CreatorJob,
				tCalendar.SystemKind,
			).
			SET(
				title,
				mysql.String("System-managed birthday calendar"),
				false,
				true,
				mysql.NULL,
				mysql.String(job),
				mysql.Int32(int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS)),
			).
			WHERE(mysql.AND(
				tCalendar.ID.EQ(mysql.Int64(existing.GetId())),
			)).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		if err := ensureBirthdayCalendarAccess(
			ctx,
			tx,
			existing.GetId(),
			job,
			jobInfo,
		); err != nil {
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		return existing, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	tCalendar := table.FivenetCalendar
	stmt := tCalendar.
		INSERT(
			tCalendar.Job,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.CreatorID,
			tCalendar.CreatorJob,
			tCalendar.SystemKind,
		).
		VALUES(
			mysql.String(job),
			title,
			mysql.String("System-managed birthday calendar"),
			mysql.Bool(false),
			mysql.Bool(true),
			birthdayCalendarColor,
			mysql.NULL,
			mysql.String(job),
			mysql.Int32(int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS)),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCalendar.Name.SET(mysql.String(title)),
			tCalendar.Description.SET(mysql.String("System-managed birthday calendar")),
			tCalendar.Public.SET(mysql.Bool(false)),
			tCalendar.Closed.SET(mysql.Bool(true)),
			tCalendar.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tCalendar.CreatorID.SET(mysql.IntExp(mysql.NULL)),
			tCalendar.CreatorJob.SET(mysql.String(job)),
			tCalendar.SystemKind.SET(mysql.Int32(int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS))),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	if lastID <= 0 {
		tCalendar := table.FivenetCalendar.AS("calendar")
		tCreator := table.FivenetUser.AS("creator")
		tAvatar := table.FivenetFiles.AS("profile_picture")
		row := &calendarresource.Calendar{}
		selectStmt := tCalendar.
			SELECT(
				tCalendar.ID,
				tCalendar.CreatedAt,
				tCalendar.UpdatedAt,
				tCalendar.DeletedAt,
				tCalendar.Job,
				tCalendar.Name,
				tCalendar.Description,
				tCalendar.Public,
				tCalendar.Closed,
				tCalendar.Color,
				tCalendar.CreatorID,
				tCreator.ID,
				tCreator.Job,
				tCreator.JobGrade,
				tCreator.Firstname,
				tCreator.Lastname,
				tCreator.Dateofbirth,
				tCreator.PhoneNumber,
				tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
				tAvatar.FilePath.AS("creator.profile_picture"),
				tCalendarSubs.CalendarID,
				tCalendarSubs.UserID,
				tCalendarSubs.CreatedAt,
				tCalendarSubs.Confirmed,
				tCalendarSubs.Muted,
				tCalendar.SystemKind,
			).
			FROM(tCalendar.
				LEFT_JOIN(tCreator,
					tCalendar.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCalendar.CreatorID),
				).
				LEFT_JOIN(tCalendarSubs,
					mysql.AND(
						tCalendarSubs.CalendarID.EQ(tCalendar.ID),
						tCalendarSubs.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
			).
			WHERE(mysql.AND(
				tCalendar.DeletedAt.IS_NULL(),
				tCalendar.Job.EQ(mysql.String(job)),
				tCalendar.SystemKind.EQ(
					mysql.Int32(
						int32(
							calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS,
						),
					),
				),
			)).
			LIMIT(1)
		if err := selectStmt.QueryContext(ctx, tx, row); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, err
			}
		}
		if row.GetId() > 0 {
			lastID = row.GetId()
		}
	}

	if lastID <= 0 {
		return nil, fmt.Errorf("unable to create birthday calendar for job %s", job)
	}

	if err := ensureBirthdayCalendarAccess(ctx, tx, lastID, job, jobInfo); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.getCalendar(
		ctx,
		userInfo,
		tCalendar.ID.EQ(mysql.Int64(lastID)),
	)
}

func (s *Server) expandCalendarEntries(
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

func (s *Server) expandCalendarEntryOccurrences(
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

	return expandRecurringEntry(entry, rangeStart, rangeEnd), nil
}

func expandRecurringEntry(
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
	occurrenceIndex := 0

	for {
		if occurrenceStart.After(rangeEnd) {
			break
		}

		if until := entry.GetRecurring().
			GetUntil(); until != nil &&
			occurrenceStart.After(until.AsTime()) {
			break
		}

		if occurrenceIndex > 0 ||
			entryOverlapsRange(occurrenceStart, entry.GetEndTime(), rangeStart, rangeEnd) {
			clone := proto.Clone(entry).(*calendarentries.CalendarEntry)
			clone.StartTime = timestamp.New(occurrenceStart)
			if entry.GetEndTime() != nil {
				end := occurrenceStart.Add(duration)
				clone.EndTime = timestamp.New(end)
			}
			key := fmt.Sprintf(
				"%s:%d:%d",
				occurrenceKeyPrefix,
				entry.GetId(),
				occurrenceStart.Unix(),
			)
			if sourceUserID != nil {
				key = fmt.Sprintf(
					"%s:%d:%04d:%02d:%02d",
					occurrenceKeyPrefix,
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

		occurrenceIndex++
		occurrenceStart = nextRecurringOccurrence(
			occurrenceStart,
			interval,
			entry.GetRecurring().GetEvery(),
		)
	}

	return out
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
