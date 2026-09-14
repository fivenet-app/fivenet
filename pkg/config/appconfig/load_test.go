package appconfig

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	serverconfig "github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyInitialConfigAndStartupOverride(t *testing.T) {
	t.Parallel()

	initialPrivacyPolicy := "https://example.com/privacy"
	overrideImprint := "https://override.example.com/imprint"
	c := &Config{
		serverConfig: &serverconfig.Config{
			AppConfig: serverconfig.AppConfig{
				Initial: serverconfig.AppConfigInitial{
					Website: serverconfig.AppConfigWebsite{
						Links: &serverconfig.AppConfigLinks{PrivacyPolicy: &initialPrivacyPolicy},
					},
				},
				Override: serverconfig.AppConfigOverride{
					Enabled: true,
					Website: serverconfig.AppConfigWebsite{
						Links: &serverconfig.AppConfigLinks{Imprint: &overrideImprint},
					},
				},
			},
		},
	}
	cfg := &settings.AppConfig{}
	cfg.Default()

	c.applyInitialConfig(cfg)
	require.True(t, cfg.GetWebsite().GetLinks().HasPrivacyPolicy())
	assert.Equal(t, initialPrivacyPolicy, cfg.GetWebsite().GetLinks().GetPrivacyPolicy())
	assert.False(t, cfg.GetWebsite().GetLinks().HasImprint())

	c.applyStartupOverride(cfg)
	assert.Equal(t, initialPrivacyPolicy, cfg.GetWebsite().GetLinks().GetPrivacyPolicy())
	require.True(t, cfg.GetWebsite().GetLinks().HasImprint())
	assert.Equal(t, overrideImprint, cfg.GetWebsite().GetLinks().GetImprint())
}
