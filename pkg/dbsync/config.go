package dbsync

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
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

	c := &DBSync{}
	if err := defaults.Set(c); err != nil {
		return fmt.Errorf("failed to set config defaults: %w", err)
	}

	if err := v.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	s.cfg = c

	return nil
}

type DBSync struct {
	Enabled bool `default:"false" yaml:"enabled"`

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
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
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
	Enabled           bool    `yaml:"enabled"`
	UpdatedTimeColumn *string `yaml:"updatedTimeColumn"`
	Query             string  `yaml:"query"`
}

type UsersDBSyncTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	SplitName bool `default:"false" yaml:"splitName"`
}
