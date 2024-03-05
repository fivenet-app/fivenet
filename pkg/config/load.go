package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		Load,
	),
)

func Load() (*BaseConfig, error) {
	// Viper Config reading setup
	viper.SetEnvPrefix("FIVENET")
	viper.SetConfigType("yaml")

	if configFile := os.Getenv("FIVENET_CONFIG_FILE"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config")
	}

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	c := &BaseConfig{}
	if err := defaults.Set(c); err != nil {
		return nil, fmt.Errorf("failed to set config defaults: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return c, nil
}

func LoadTest() (*BaseConfig, error) {
	c := &BaseConfig{}
	if err := defaults.Set(c); err != nil {
		return nil, fmt.Errorf("failed to set config defaults: %w", err)
	}

	return c, nil
}
