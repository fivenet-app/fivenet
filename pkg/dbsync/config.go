package dbsync

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/spf13/viper"
)

func (s *Sync) loadConfig() error {
	v := viper.New()
	// Viper config reading setup
	v.SetEnvPrefix("FIVENET")
	v.SetConfigType("yaml")

	if configFile := os.Getenv("FIVENET_DBSYNC_FILE"); configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("dbsync")
		v.AddConfigPath(".")
		v.AddConfigPath("/config")
	}

	// Find and read the dbsync config file
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	c := &config.DBSync{}
	if err := defaults.Set(c); err != nil {
		return fmt.Errorf("failed to set config defaults: %w", err)
	}

	if err := v.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	s.cfg = c

	return nil
}
