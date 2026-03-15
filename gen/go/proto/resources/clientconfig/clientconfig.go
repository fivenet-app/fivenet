package clientconfig

import (
	settings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
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

		DefaultLocale: appCfg.GetDefaultLocale(),

		Auth: &Auth{
			SignupEnabled: appCfg.GetAuth().GetSignupEnabled(),
			LastCharLock:  appCfg.GetAuth().GetLastCharLock(),
			Providers:     providers,
		},
		Discord: &Discord{
			BotEnabled: appCfg.GetDiscord().GetBotId() != "",
		},
		Website: &Website{
			Links: &Links{
				Imprint:       appCfg.GetWebsite().GetLinks().Imprint,
				PrivacyPolicy: appCfg.GetWebsite().GetLinks().PrivacyPolicy,
			},
			StatsPage: appCfg.GetWebsite().GetStatsPage(),
		},
		FeatureGates: &FeatureGates{},
		Game: &Game{
			UnemployedJobName: appCfg.GetJobInfo().GetUnemployedJob().GetName(),
			StartJobGrade:     cfg.Game.StartJobGrade,

			Livemap: &Livemap{
				EnableCayoPerico: appCfg.GetLivemap().GetEnableCayoPerico(),
			},

			MaxWantedDurationUserEnabled:    appCfg.GetGame().GetMaxWantedDurationUserEnabled(),
			MaxWantedDurationUser:           proto.Clone(appCfg.GetGame().GetMaxWantedDurationUser()).(*durationpb.Duration),
			MaxWantedDurationVehicleEnabled: appCfg.GetGame().GetMaxWantedDurationVehicleEnabled(),
			MaxWantedDurationVehicle:        proto.Clone(appCfg.GetGame().GetMaxWantedDurationVehicle()).(*durationpb.Duration),
		},
		System: &System{
			BannerMessageEnabled: appCfg.GetSystem().GetBannerMessageEnabled(),
			Otlp: &OTLPFrontend{
				Enabled: cfg.OTLP.Enabled,
				Url:     cfg.OTLP.Frontend.URL,
				Headers: cfg.OTLP.Frontend.Headers,
			},
		},
		Display: &Display{
			IntlLocale:   appCfg.GetDisplay().GetIntlLocale(),
			CurrencyName: appCfg.GetDisplay().GetCurrencyName(),
		},
		QuickButtons: quickButtons,
		Data: &settings.Data{
			Mode: appCfg.GetData().GetMode(),
		},
	}

	if appCfg.GetSystem().GetBannerMessage() != nil {
		clientCfg.System.BannerMessage = proto.Clone(appCfg.GetSystem().GetBannerMessage()).(*settings.BannerMessage)
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
