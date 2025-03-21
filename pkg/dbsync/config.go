package dbsync

import (
	"fmt"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func (s *Sync) loadConfig(shutdowner fx.Shutdowner) error {
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

	loadConfig := func() error {
		// Find and read the dbsync config file
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("fatal error config file: %w", err)
		}

		c := &DBSyncConfig{}
		if err := defaults.Set(c); err != nil {
			return fmt.Errorf("failed to set config defaults: %w", err)
		}

		if err := v.Unmarshal(c); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		s.cfg.Store(c)

		return nil
	}

	if err := loadConfig(); err != nil {
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(_ fsnotify.Event) {
		s.logger.Info("config change detected, reloading config")
		if err := loadConfig(); err != nil {
			s.logger.Error("failed to hot reload config", zap.Error(err))
			return
		}

		if err := s.restart(); err != nil {
			if err := shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
				s.logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
			}
			s.logger.Error("failed to restart dbsync", zap.Error(err))
			return
		}
	})

	return nil
}

type DBSyncConfig struct {
	WatchConfig bool `default:"true" yaml:"watchConfig"`

	StateFile string `default:"dbsync.state.yaml" yaml:"stateFile"`

	Destination DBSyncDestination `yaml:"destination"`
	Source      DBSyncSource      `yaml:"source"`

	Tables DBSyncSourceTables `yaml:"tables"`
}

type DBSyncSource struct {
	// Refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN string `yaml:"dsn"`
}

type DBSyncDestination struct {
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`

	SyncInterval time.Duration `default:"5s" yaml:"syncInterval"`
}

type DBSyncSourceTables struct {
	Jobs      DBSyncTable `yaml:"jobs"`
	JobGrades DBSyncTable `yaml:"jobGrades"`
	Licenses  DBSyncTable `yaml:"licenses"`

	Users        UsersDBSyncTable `yaml:"users"`
	UserLicenses DBSyncTable      `yaml:"userLicenses"`
	Vehicles     DBSyncTable      `yaml:"vehicles"`
}

type DBSyncTable struct {
	Enabled           bool           `yaml:"enabled"`
	UpdatedTimeColumn *string        `yaml:"updatedTimeColumn"`
	Query             string         `yaml:"query"`
	SyncInterval      *time.Duration `yaml:"syncInterval"`
}

func (c *DBSyncTable) GetSyncInterval() *time.Duration {
	return c.SyncInterval
}

type UsersDBSyncTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	SplitName    bool                  `default:"false" yaml:"splitName"`
	DateOfBirth  DateOfBirthNormalizer `yaml:"dateOfBirth"`
	ValueMapping *UsersValueMappings   `yaml:"valueMapping"`
}

type DateOfBirthNormalizer struct {
	Formats      []string `yaml:"formats"`
	OutputFormat string   `default:"" yaml:"output"`
}

type UsersValueMappings struct {
	Sex *ValueMapping `yaml:"sex"`
}

type ValueMapping struct {
	Fallback *string           `yaml:"fallback"`
	Values   map[string]string `yaml:"values"`
}

func (c *ValueMapping) IsEmpty() bool {
	if c == nil {
		return true
	}

	return len(c.Values) == 0
}

func (c *ValueMapping) Process(input *string) {
	val, ok := c.Values[*input]
	if !ok {
		if c.Fallback != nil {
			*input = *c.Fallback
		}
	} else {
		*input = val
	}
}

type DBSyncTableSyncInterval interface {
	GetSyncInterval() *time.Duration
}

func (c *DBSyncConfig) GetSyncInterval(table DBSyncTableSyncInterval) time.Duration {
	if table != nil && table.GetSyncInterval() != nil {
		interval := table.GetSyncInterval()
		return *interval
	}

	return c.Destination.SyncInterval
}
