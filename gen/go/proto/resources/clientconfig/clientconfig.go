package clientconfig

import (
	settings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
)

func BuildClientConfig(cfg *config.Config, providers []*ProviderConfig, appCfg *appconfig.Cfg) *ClientConfig {
	clientCfg := &ClientConfig{
		Version: version.Version,

		DefaultLocale: appCfg.DefaultLocale,

		Login: &LoginConfig{
			SignupEnabled: appCfg.Auth.SignupEnabled,
			LastCharLock:  appCfg.Auth.LastCharLock,
			Providers:     providers,
		},
		Discord: &Discord{
			BotEnabled: appCfg.Discord.BotId != nil && *appCfg.Discord.BotId != "",
		},
		Website: &Website{
			Links: &Links{
				Imprint:       appCfg.Website.Links.Imprint,
				PrivacyPolicy: appCfg.Website.Links.PrivacyPolicy,
			},
			StatsPage: appCfg.Website.StatsPage,
		},
		FeatureGates: &FeatureGates{
			ImageProxy: cfg.ImageProxy.Enabled,
		},
		Game: &Game{
			UnemployedJobName: appCfg.JobInfo.UnemployedJob.Name,
			StartJobGrade:     cfg.Game.StartJobGrade,
		},
		System: &System{
			BannerMessageEnabled: appCfg.System.BannerMessageEnabled,
			BannerMessage:        appCfg.System.BannerMessage,
			Otlp: &OTLPFrontend{
				Enabled: cfg.OTLP.Enabled,
				Url:     cfg.OTLP.Frontend.URL,
				Headers: cfg.OTLP.Frontend.Headers,
			},
		},
	}

	if appCfg.System.BannerMessage != nil {
		clientCfg.System.BannerMessage = &settings.BannerMessage{
			Id:        appCfg.System.BannerMessage.Id,
			Title:     appCfg.System.BannerMessage.Title,
			Icon:      appCfg.System.BannerMessage.Icon,
			Color:     appCfg.System.BannerMessage.Color,
			CreatedAt: appCfg.System.BannerMessage.CreatedAt,
			ExpiresAt: appCfg.System.BannerMessage.ExpiresAt,
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
