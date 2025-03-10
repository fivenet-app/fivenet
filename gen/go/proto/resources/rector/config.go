package rector

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	DefaultUserTrackerRefreshTime   = 3*time.Second + 250*time.Millisecond
	DefaultUserTrackerDbRefreshTime = 1 * time.Second

	DefaultDiscordSyncInterval = 15 * time.Minute
)

func (x *AppConfig) Default() {
	if x.DefaultLocale == "" {
		x.DefaultLocale = "en" // Default to English locale
	}

	if x.Auth == nil {
		x.Auth = &Auth{
			SignupEnabled: true,
			LastCharLock:  false,
		}
	}

	if x.Perms == nil {
		x.Perms = &Perms{
			Default: []*Perm{
				{
					Category: "AuthService",
					Name:     "ChooseCharacter",
				},
				{
					Category: "DocStoreService",
					Name:     "ListDocuments",
				},
				{
					Category: "DocStoreService",
					Name:     "PostComment",
				},
				{
					Category: "QualificationsService",
					Name:     "ListQualifications",
				},
				{
					Category: "WikiService",
					Name:     "ListPages",
				},
			},
		}
	}

	if x.Website == nil {
		x.Website = &Website{
			Links:     &Links{},
			StatsPage: false,
		}
	}
	if x.Website.Links == nil {
		x.Website.Links = &Links{}
	}

	if x.JobInfo == nil {
		x.JobInfo = &JobInfo{
			PublicJobs: []string{},
			HiddenJobs: []string{},
		}
	}
	if x.JobInfo.UnemployedJob == nil {
		x.JobInfo.UnemployedJob = &UnemployedJob{
			Name:  "unemployed",
			Grade: 1,
		}
	}

	if x.UserTracker == nil {
		x.UserTracker = &UserTracker{
			LivemapJobs: []string{},
		}
	}
	if x.UserTracker.RefreshTime == nil {
		x.UserTracker.RefreshTime = durationpb.New(DefaultUserTrackerRefreshTime)
	}
	if x.UserTracker.DbRefreshTime == nil {
		x.UserTracker.DbRefreshTime = durationpb.New(DefaultUserTrackerDbRefreshTime)
	}

	if x.Discord == nil {
		x.Discord = &Discord{
			Enabled: false,
		}
	}
	if x.Discord.SyncInterval == nil {
		x.Discord.SyncInterval = durationpb.New(DefaultDiscordSyncInterval)
	}
	if x.Discord.BotPresence == nil {
		x.Discord.BotPresence = &DiscordBotPresence{
			Type: DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_UNSPECIFIED,
		}
	}
}
