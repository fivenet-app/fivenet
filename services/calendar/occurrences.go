package calendar

import (
	"context"
	"fmt"
	"sort"
	"time"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"google.golang.org/protobuf/proto"
)

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
				"%s:%d:%d",
				occurrenceKeyPrefix,
				entry.GetId(),
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
