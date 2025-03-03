package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/fivenet-app/fivenet/cmd/envs"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		Load,
	),
)

type Result struct {
	fx.Out

	Config        *Config
	DiscordConfig *Discord
}

func Load() (Result, error) {
	v := viper.New()
	// Viper config reading setup
	v.SetEnvPrefix("FIVENET")
	v.SetConfigType("yaml")

	if configFile := os.Getenv(envs.ConfigFileEnvVar); configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("/config")
	}

	res := Result{}
	// Find and read the config file
	if err := v.ReadInConfig(); err != nil {
		return res, fmt.Errorf("fatal error config file: %w", err)
	}

	c := &Config{}
	if err := defaults.Set(c); err != nil {
		return res, fmt.Errorf("failed to set config defaults: %w", err)
	}
	res.Config = c

	res.DiscordConfig = &c.Discord

	if err := v.Unmarshal(c); err != nil {
		return res, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return res, nil
}

var TestModule = fx.Module("config_test",
	fx.Provide(
		LoadTestConfig,
	),
)

func LoadTestConfig() (*Config, error) {
	c := &Config{}
	if err := defaults.Set(c); err != nil {
		return nil, fmt.Errorf("failed to set config defaults: %w", err)
	}

	// Set audit log retention days high so they won't run in "short" tests
	c.Audit.RetentionDays = 365

	return c, nil
}
