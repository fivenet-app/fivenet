package calendar

import (
	"slices"

	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
)

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
