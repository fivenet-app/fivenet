package calendar

import (
	"context"
	"slices"

	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
)

func (s *Server) finalizeCalendarEntries(
	ctx context.Context,
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

	targets := make([]citizenshydrator.ShortTarget, 0, len(entries))
	for i, entry := range entries {
		if entry.GetCreatorId() > 0 {
			targets = append(targets, citizenshydrator.ShortTarget{
				UserID: entry.GetCreatorId(),
				Set: func(user *usershort.UserShort) {
					entries[i].Creator = user
				},
			})
		}
	}
	if len(targets) == 0 {
		return entries
	}

	hydrateShort := s.hydrator.HydrateShortTargetsSafeFunc(userInfo)
	if err := hydrateShort(ctx, nil, targets); err != nil {
		return entries
	}

	return entries
}
