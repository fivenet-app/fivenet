package clientconfig

import (
	settings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
	"google.golang.org/protobuf/proto"
)

func BuildClientConfig(
	cfg *config.Config,
	providers []*ProviderConfig,
	appCfg *appconfig.Cfg,
) *ClientConfig {
	var quickButtons *settings.QuickButtons
	if appCfg.GetQuickButtons() != nil {
		quickButtons = proto.Clone(appCfg.GetQuickButtons()).(*settings.QuickButtons)
	}

	clientCfg := &ClientConfig{
		Version: version.Version,

		DefaultLocale: appCfg.DefaultLocale,

		Login: &LoginConfig{
			SignupEnabled: appCfg.Auth.GetSignupEnabled(),
			LastCharLock:  appCfg.Auth.GetLastCharLock(),
			Providers:     providers,
		},
		Discord: &Discord{
			BotEnabled: appCfg.Discord.BotId != nil && appCfg.Discord.GetBotId() != "",
		},
		Website: &Website{
			Links: &Links{
				Imprint:       appCfg.Website.GetLinks().Imprint,
				PrivacyPolicy: appCfg.Website.GetLinks().PrivacyPolicy,
			},
			StatsPage: appCfg.Website.GetStatsPage(),
		},
		FeatureGates: &FeatureGates{
			ImageProxy: cfg.ImageProxy.Enabled,
		},
		Game: &Game{
			UnemployedJobName: appCfg.JobInfo.GetUnemployedJob().GetName(),
			StartJobGrade:     cfg.Game.StartJobGrade,
		},
		System: &System{
			BannerMessageEnabled: appCfg.System.GetBannerMessageEnabled(),
			Otlp: &OTLPFrontend{
				Enabled: cfg.OTLP.Enabled,
				Url:     cfg.OTLP.Frontend.URL,
				Headers: cfg.OTLP.Frontend.Headers,
			},
		},
		Display: &Display{
			IntlLocale:   appCfg.Display.IntlLocale,
			CurrencyName: appCfg.Display.CurrencyName,
		},
		QuickButtons: quickButtons,
	}

	if appCfg.System.GetBannerMessage() != nil {
		clientCfg.System.BannerMessage = &settings.BannerMessage{
			Id:        appCfg.System.GetBannerMessage().GetId(),
			Title:     appCfg.System.GetBannerMessage().GetTitle(),
			Icon:      appCfg.System.GetBannerMessage().Icon,
			Color:     appCfg.System.GetBannerMessage().Color,
			CreatedAt: appCfg.System.GetBannerMessage().GetCreatedAt(),
			ExpiresAt: appCfg.System.GetBannerMessage().GetExpiresAt(),
		}
	}

	return clientCfg
}

func BuildProviderList(cfg *config.Config) []*ProviderConfig {
	providers := make([]*ProviderConfig, len(cfg.OAuth2.Providers))

	for i, p := range cfg.OAuth2.Providers {
		providers[i] = &ProviderConfig{
			Name:     p.Name,
			Label:    p.Label,
			Icon:     p.Icon,
			Homepage: p.Homepage,
		}
	}

	return providers
}
