package settings

import (
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	DefaultUserTrackerRefreshTime   = 3*time.Second + 250*time.Millisecond
	DefaultUserTrackerDbRefreshTime = 1 * time.Second

	DefaultDiscordSyncInterval = 15 * time.Minute
)

func (x *AppConfig) Default() {
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
					Category: "auth.AuthService",
					Name:     "ChooseCharacter",
				},
				{
					Category: "completor.CompletorService",
					Name:     "CompleteJobs",
				},
				{
					Category: "documents.DocumentsService",
					Name:     "ListDocuments",
				},
				{
					Category: "qualifications.QualificationsService",
					Name:     "ListQualifications",
				},
				{
					Category: "wiki.WikiService",
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
		x.UserTracker = &UserTracker{}
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

	if x.System == nil {
		x.System = &System{
			BannerMessageEnabled: false,
		}
	}
	if x.System.BannerMessage != nil {
		if x.System.BannerMessage.CreatedAt == nil {
			x.System.BannerMessage.CreatedAt = timestamp.Now()
		}
	}
}
