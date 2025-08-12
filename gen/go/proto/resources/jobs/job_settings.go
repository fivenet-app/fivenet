package jobs

import "slices"

const (
	DefaultJobAbsencePastDays   = 7
	DefaultJobAbsenceFutureDays = 93 // ~3 months
)

func (x *JobSettings) Default() {
	if x.GetAbsencePastDays() <= 0 {
		x.AbsencePastDays = DefaultJobAbsencePastDays
	}
	if x.GetAbsenceFutureDays() <= 0 {
		x.AbsenceFutureDays = DefaultJobAbsenceFutureDays
	}
}

func (x *DiscordSyncSettings) IsStatusLogEnabled() bool {
	return x.GetStatusLog() && x.GetStatusLogSettings() != nil &&
		x.GetStatusLogSettings().GetChannelId() != ""
}

func (x *DiscordSyncChanges) Add(change *DiscordSyncChange) {
	if x.Changes == nil {
		x.Changes = []*DiscordSyncChange{}
	}

	if len(x.GetChanges()) > 0 {
		lastChange := x.GetChanges()[len(x.GetChanges())-1]

		if lastChange.GetPlan() == change.GetPlan() {
			return
		}
	}

	x.Changes = append(x.Changes, change)

	if len(x.GetChanges()) > 12 {
		x.Changes = slices.Delete(x.GetChanges(), 0, 1)
	}
}
