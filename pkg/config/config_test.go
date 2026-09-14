package config

import (
	"strings"
	"testing"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDemoDefaults(t *testing.T) {
	t.Parallel()
	cfg := &Config{}
	require.NoError(t, defaults.Set(cfg), "failed to set defaults")

	assert.True(
		t,
		cfg.Demo.Features.Dispatches,
		"expected demo.features.dispatches default to be true",
	)
	assert.True(
		t,
		cfg.Demo.Features.Locations,
		"expected demo.features.locations default to be true",
	)
	assert.True(
		t,
		cfg.Demo.Features.Timeclock,
		"expected demo.features.timeclock default to be true",
	)
	assert.False(t, cfg.Demo.Features.Users, "expected demo.features.users default to be false")

	assert.Equal(
		t,
		50,
		cfg.Demo.FakeUsers.Count,
		"expected demo.fakeUsers.count default 50, got %d",
		cfg.Demo.FakeUsers.Count,
	)
}

func TestAppConfigLinksFromYAML(t *testing.T) {
	t.Parallel()

	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	v.SetConfigType("yaml")
	require.NoError(t, v.ReadConfig(strings.NewReader(`
appConfig:
  initial:
    website:
      links:
        privacyPolicy: "https://example.com/privacy"
        imprint: "https://example.com/imprint"
  override:
    enabled: true
    website:
      links:
        privacyPolicy: "https://override.example.com/privacy"
`)))

	cfg := &Config{}
	require.NoError(t, defaults.Set(cfg))
	require.NoError(t, v.Unmarshal(cfg))

	require.NotNil(t, cfg.AppConfig.Initial.Website.Links)
	assert.Equal(t, "https://example.com/privacy", *cfg.AppConfig.Initial.Website.Links.PrivacyPolicy)
	assert.Equal(t, "https://example.com/imprint", *cfg.AppConfig.Initial.Website.Links.Imprint)
	assert.True(t, cfg.AppConfig.Override.Enabled)
	require.NotNil(t, cfg.AppConfig.Override.Website.Links)
	assert.Equal(
		t,
		"https://override.example.com/privacy",
		*cfg.AppConfig.Override.Website.Links.PrivacyPolicy,
	)
	assert.Nil(t, cfg.AppConfig.Override.Website.Links.Imprint)
}
