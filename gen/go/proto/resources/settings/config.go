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

	// https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
	DefaultDiscordBotPermissions = -42471713

	DefaultDisplayIntlLocale = "en-US"

	DefaultPenaltyCalculatorMaxCount  = uint32(10)
	DefaultPenaltyCalculatorMaxLeeway = uint32(25)
)

func (x *AppConfig) Default() {
	if x.GetAuth() == nil {
		x.Auth = &Auth{
			SignupEnabled: true,
			LastCharLock:  false,
		}
	}

	if x.GetPerms() == nil {
		x.Perms = &Perms{}
	}
	if x.GetPerms().GetDefault() == nil {
		x.Perms.Default = []*Perm{
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
		}
	}

	if x.GetWebsite() == nil {
		x.Website = &Website{}
	}
	if x.GetWebsite().GetLinks() == nil {
		x.Website.Links = &Links{}
	}

	if x.GetJobInfo() == nil {
		x.JobInfo = &JobInfo{
			PublicJobs: []string{},
			HiddenJobs: []string{},
		}
	}
	if x.GetJobInfo().GetUnemployedJob() == nil {
		x.JobInfo.UnemployedJob = &UnemployedJob{
			Name:  "unemployed",
			Grade: 1,
		}
	}

	if x.GetUserTracker() == nil {
		x.UserTracker = &UserTracker{}
	}
	if x.GetUserTracker().GetRefreshTime() == nil {
		x.UserTracker.RefreshTime = durationpb.New(DefaultUserTrackerRefreshTime)
	}
	if x.GetUserTracker().GetDbRefreshTime() == nil {
		x.UserTracker.DbRefreshTime = durationpb.New(DefaultUserTrackerDbRefreshTime)
	}

	if x.GetDiscord() == nil {
		x.Discord = &Discord{
			BotPermissions: DefaultDiscordBotPermissions,
		}
	}
	if x.GetDiscord().GetSyncInterval() == nil {
		x.Discord.SyncInterval = durationpb.New(DefaultDiscordSyncInterval)
	}
	if x.GetDiscord().GetBotPresence() == nil {
		status := "FiveNet"
		url := "https://fivenet.app"

		x.Discord.BotPresence = &DiscordBotPresence{
			Type:   DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_GAME,
			Status: &status,
			Url:    &url,
		}
	}
	if x.GetDiscord().GetBotPermissions() == 0 {
		x.Discord.BotPermissions = DefaultDiscordBotPermissions
	}

	if x.GetSystem() == nil {
		x.System = &System{}
	}
	if x.GetSystem().GetBannerMessage() != nil {
		if x.GetSystem().GetBannerMessage().GetCreatedAt() == nil {
			x.System.BannerMessage.CreatedAt = timestamp.Now()
		}
	}

	if x.GetDisplay() == nil {
		x.Display = &Display{}
	}
	if x.GetDisplay().GetIntlLocale() == "" {
		defaultIntlLocale := DefaultDisplayIntlLocale

		x.Display.IntlLocale = &defaultIntlLocale
	}
	if x.GetDisplay().GetCurrencyName() == "" {
		x.Display.CurrencyName = "USD"
	}

	if x.GetQuickButtons() == nil {
		x.QuickButtons = &QuickButtons{}
	}
	if x.GetQuickButtons().GetPenaltyCalculator() == nil {
		x.QuickButtons.PenaltyCalculator = &PenaltyCalculator{}
	}
	if x.GetQuickButtons().GetPenaltyCalculator().GetMaxCount() == 0 {
		maxCount := DefaultPenaltyCalculatorMaxCount
		x.QuickButtons.PenaltyCalculator.MaxCount = &maxCount
	}
	if x.GetQuickButtons().GetPenaltyCalculator().GetMaxLeeway() == 0 {
		maxLeeway := DefaultPenaltyCalculatorMaxLeeway
		x.QuickButtons.PenaltyCalculator.MaxLeeway = &maxLeeway
	}
}
