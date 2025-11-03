package dbsync

import (
	"fmt"
	"os"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/creasty/defaults"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	shutdowner fx.Shutdowner

	v *viper.Viper

	cfg atomic.Pointer[DBSyncConfig]
}

type ParamsConfig struct {
	fx.In

	LC         fx.Lifecycle
	Shutdowner fx.Shutdowner
}

type ResultConfig struct {
	fx.Out

	Config *Config
	Cfg    *config.Config
}

func NewConfig(p ParamsConfig) (ResultConfig, error) {
	s := &Config{
		shutdowner: p.Shutdowner,
	}

	r := ResultConfig{
		Config: s,
		Cfg:    &config.Config{},
	}

	s.v = viper.New()
	// Viper config reading setup
	s.v.SetEnvPrefix("FIVENET")
	s.v.SetConfigType("yaml")

	if configFile := os.Getenv("FIVENET_DBSYNC_FILE"); configFile != "" {
		s.v.SetConfigFile(configFile)
	} else {
		s.v.SetConfigName("dbsync")
		s.v.AddConfigPath(".")
		s.v.AddConfigPath("/config")
	}

	if err := s.LoadConfig(); err != nil {
		return r, err
	}

	return r, nil
}

func (s *Config) Load() *DBSyncConfig {
	return s.cfg.Load()
}

func (s *Config) LoadConfig() error {
	// Find and read the dbsync config file
	if err := s.v.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error in config file format/syntax. %w", err)
	}

	c := &DBSyncConfig{}
	if err := defaults.Set(c); err != nil {
		return fmt.Errorf("failed to set defaults in config file. %w", err)
	}

	if err := s.v.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to read config file data into the system. %w", err)
	}

	if err := c.Init(); err != nil {
		return fmt.Errorf("failed to initialize config. %w", err)
	}

	s.cfg.Store(c)

	return nil
}

func (s *Config) setupWatch(logger *zap.Logger, restartFn func() error) {
	s.v.WatchConfig()
	s.v.OnConfigChange(func(_ fsnotify.Event) {
		logger.Info("config change detected, reloading config")
		if err := s.LoadConfig(); err != nil {
			logger.Error("failed to hot reload config", zap.Error(err))
			return
		}

		if err := restartFn(); err != nil {
			if err := s.shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
				logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
			}
			logger.Error("failed to restart dbsync", zap.Error(err))
			return
		}
	})
}

type DBSyncConfig struct {
	Mode string `default:"release" yaml:"mode"`

	Log config.Log `yaml:"log"`

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
	Jobs      JobsDBSyncTable `yaml:"jobs"`
	JobGrades DBSyncTable     `yaml:"jobGrades"`
	Licenses  DBSyncTable     `yaml:"licenses"`

	Users            UsersDBSyncTable `yaml:"users"`
	CitizensLicenses DBSyncTable      `yaml:"userLicenses"`
	Vehicles         DBSyncTable      `yaml:"vehicles"`
}

type FilterAction string

const (
	// Replace the matching pattern with the replacement string.
	FilterActionReplace FilterAction = "replace"
	// Drop the whole record if the pattern matches.
	FilterActionDrop FilterAction = "drop"
)

type Filter struct {
	Pattern     string       `yaml:"pattern"`
	Action      FilterAction `yaml:"action"      default:"replace"`
	Replacement string       `yaml:"replacement"`

	compiledPattern *regexp.Regexp
}

type DBSyncTable struct {
	Enabled           bool           `yaml:"enabled"`
	UpdatedTimeColumn *string        `yaml:"updatedTimeColumn"`
	Query             string         `yaml:"query"`
	SyncInterval      *time.Duration `yaml:"syncInterval"`
}

func (c *DBSyncTable) GetSyncInterval() *time.Duration {
	// Minimum sync interval is 1 second
	if c.SyncInterval != nil && *c.SyncInterval <= 1*time.Second {
		interval := 1 * time.Second
		return &interval
	}

	return c.SyncInterval
}

type JobsDBSyncTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Filters []Filter `yaml:"filters"`
}

type UsersDBSyncTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	SplitName    bool                  `default:"false" yaml:"splitName"`
	DateOfBirth  DateOfBirthNormalizer `                yaml:"dateOfBirth"`
	ValueMapping *UsersValueMappings   `                yaml:"valueMapping"`

	IgnoreEmptyName bool `default:"true" yaml:"ignoreEmptyName"`

	Filters UsersFilters `yaml:"filters"`
}

type UsersFilters struct {
	Jobs []Filter `yaml:"jobs"`
}

type DateOfBirthNormalizer struct {
	Formats      []string `yaml:"formats"`
	OutputFormat string   `yaml:"output"  default:""`
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
	if input == nil {
		return
	}

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

func (c *DBSyncConfig) Init() error {
	// Compile filter regex patterns
	for k := range c.Tables.Jobs.Filters {
		filter := &c.Tables.Jobs.Filters[k]
		var err error
		filter.compiledPattern, err = regexp.Compile(filter.Pattern)
		if err != nil {
			return fmt.Errorf("failed to compile regex for filter %d: %w", k, err)
		}
	}

	for k := range c.Tables.Users.Filters.Jobs {
		filter := &c.Tables.Users.Filters.Jobs[k]
		var err error
		filter.compiledPattern, err = regexp.Compile(filter.Pattern)
		if err != nil {
			return fmt.Errorf("failed to compile regex for user job filter %d: %w", k, err)
		}
	}

	return nil
}
