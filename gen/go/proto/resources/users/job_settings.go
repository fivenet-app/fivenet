package users

import "slices"

const (
	DefaultJobAbsencePastDays   = 7
	DefaultJobAbsenceFutureDays = 93 // ~3 months
)

func (x *JobSettings) Default() {
	if x.AbsencePastDays <= 0 {
		x.AbsencePastDays = DefaultJobAbsencePastDays
	}
	if x.AbsenceFutureDays <= 0 {
		x.AbsenceFutureDays = DefaultJobAbsenceFutureDays
	}
}

func (x *DiscordSyncSettings) IsStatusLogEnabled() bool {
	return x.StatusLog && x.StatusLogSettings != nil && x.StatusLogSettings.ChannelId != ""
}

func (x *DiscordSyncChanges) Add(change *DiscordSyncChange) {
	if x.Changes == nil {
		x.Changes = []*DiscordSyncChange{}
	}

	if len(x.Changes) > 0 {
		lastChange := x.Changes[len(x.Changes)-1]

		if lastChange.Plan == change.Plan {
			return
		}
	}

	x.Changes = append(x.Changes, change)

	if len(x.Changes) > 12 {
		x.Changes = slices.Delete(x.Changes, 0, 1)
	}
}
