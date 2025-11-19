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
		v:          viper.New(),
	}

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

	cc := s.cfg.Load()
	r := ResultConfig{
		Config: s,
		Cfg: &config.Config{
			Mode:     cc.Mode,
			LogLevel: cc.LogLevel,
			Log:      cc.Log,
			// Ignore db requirements, dbsync doesn't need them
			IgnoreRequirements: true,
			UpdateCheck:        cc.UpdateCheck,
		},
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

	LogLevel string     `default:"INFO" yaml:"logLevel" enum:"DEBUG,INFO,WARN,ERROR,PANIC,FATAL"`
	Log      config.Log `               yaml:"log"`

	WatchConfig bool `default:"true" yaml:"watchConfig"`

	StateFile string `default:"dbsync.state.yaml" yaml:"stateFile"`

	Destination DBSyncDestination `yaml:"destination"`
	Source      DBSyncSource      `yaml:"source"`

	Tables DBSyncSourceTables `yaml:"tables"`

	UpdateCheck config.UpdateCheck `yaml:"updateCheck"`

	TableManager TableManagerConfig `yaml:"tableManager"`
}

type DBSyncSource struct {
	config.DatabaseConnection `yaml:",inline" mapstructure:",squash"`
}

type DBSyncDestination struct {
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`

	SyncInterval time.Duration `default:"5s" yaml:"syncInterval"`
}

type DBSyncSourceTables struct {
	Jobs      JobsTable      `yaml:"jobs"`
	JobGrades JobGradesTable `yaml:"jobGrades"`
	Licenses  LicensesTable  `yaml:"licenses"`

	Users            UsersTable        `yaml:"users"`
	CitizensLicenses UserLicensesTable `yaml:"userLicenses"`
	Vehicles         VehiclesTable     `yaml:"vehicles"`
}

func (c *DBSyncSourceTables) GetAllTables() []DBSyncTable {
	tables := []DBSyncTable{}

	if c.Jobs.Enabled {
		tables = append(tables, c.Jobs.DBSyncTable)
	}
	if c.JobGrades.Enabled {
		tables = append(tables, c.JobGrades.DBSyncTable)
	}
	if c.Licenses.Enabled {
		tables = append(tables, c.Licenses.DBSyncTable)
	}
	if c.Users.Enabled {
		tables = append(tables, c.Users.DBSyncTable)
	}
	if c.CitizensLicenses.Enabled {
		tables = append(tables, c.CitizensLicenses.DBSyncTable)
	}
	if c.Vehicles.Enabled {
		tables = append(tables, c.Vehicles.DBSyncTable)
	}

	return tables
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
	TableName         string         `yaml:"tableName"`
	UpdatedTimeColumn *string        `yaml:"updatedTimeColumn"`
	Query             *string        `yaml:"query"`
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

type JobsTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns JobsColumns `yaml:"columns"`

	Filters []Filter `yaml:"filters"`
}

func (c *JobsTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil {
		return prepareStringQuery(*c.Query, c.DBSyncTable, state, offset, limit)
	}

	where = append(where, getWhereCondition(c.DBSyncTable, state))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"job.name":  c.Columns.Name,
		"job.label": c.Columns.Label,
	}, where, offset, limit)
}

type JobsColumns struct {
	Name  string `yaml:"name"  default:"name"`
	Label string `yaml:"label" default:"label"`
}

type JobGradesTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns JobGradesColumns `yaml:"columns"`
}

func (c *JobGradesTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil {
		return prepareStringQuery(*c.Query, c.DBSyncTable, state, offset, limit)
	}

	where = append(where, getWhereCondition(c.DBSyncTable, state))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"job_grade.job_name": c.Columns.JobName,
		"job_grade.grade":    c.Columns.Grade,
		"job_grade.name":     c.Columns.Name,
		"job_grade.label":    c.Columns.Label,
	}, where, offset, limit)
}

type JobGradesColumns struct {
	JobName string `yaml:"jobName" default:"job_name"`
	Grade   string `yaml:"grade"   default:"grade"`
	Name    string `yaml:"name"    default:"name"`
	Label   string `yaml:"label"   default:"label"`
}

type LicensesTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns LicensesColumns `yaml:"columns"`
}

func (c *LicensesTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil {
		return prepareStringQuery(*c.Query, c.DBSyncTable, state, offset, limit)
	}
	where = append(where, getWhereCondition(c.DBSyncTable, state))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"license.type":  c.Columns.Type,
		"license.label": c.Columns.Label,
	}, where, offset, limit)
}

type LicensesColumns struct {
	Type  string `yaml:"type"  default:"type"`
	Label string `yaml:"label" default:"label"`
}

type UsersTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns UsersColumns `yaml:"columns"`

	SplitName    bool                  `default:"false" yaml:"splitName"`
	DateOfBirth  DateOfBirthNormalizer `                yaml:"dateOfBirth"`
	ValueMapping *UsersValueMappings   `                yaml:"valueMapping"`

	IgnoreEmptyName bool `default:"true" yaml:"ignoreEmptyName"`

	Filters UsersFilters `yaml:"filters"`
}

func (c *UsersTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil {
		return prepareStringQuery(*c.Query, c.DBSyncTable, state, offset, limit)
	}

	where = append(where, getWhereCondition(c.DBSyncTable, state))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"user.id":           c.Columns.ID,
		"user.identifier":   c.Columns.Identifier,
		"user.group":        c.Columns.Group,
		"user.firstname":    c.Columns.FirstName,
		"user.lastname":     c.Columns.Lastname,
		"user.dateofbirth":  c.Columns.DateOfBirth,
		"user.job":          c.Columns.Job,
		"user.job_grade":    c.Columns.JobGrade,
		"user.sex":          c.Columns.Sex,
		"user.phone_number": c.Columns.PhoneNumber,
		"user.height":       c.Columns.Height,
		"user.visum":        c.Columns.Visum,
		"user.playtime":     c.Columns.Playtime,
	}, where, offset, limit)
}

type UsersColumns struct {
	ID          string `yaml:"id"          default:"id"`
	Identifier  string `yaml:"identifier"  default:"identifier"`
	Group       string `yaml:"group"       default:"group"`
	FirstName   string `yaml:"firstName"   default:"firstname"`
	Lastname    string `yaml:"lastname"    default:"lastname"`
	DateOfBirth string `yaml:"dateOfBirth" default:"dateofbirth"`
	Job         string `yaml:"job"         default:"job"`
	JobGrade    string `yaml:"jobGrade"    default:"job_grade"`
	Sex         string `yaml:"sex"         default:"sex"`
	PhoneNumber string `yaml:"phoneNumber" default:"phone_number"`
	Height      string `yaml:"height"      default:"height"`
	Visum       string `yaml:"visum"       default:"visum"`
	Playtime    string `yaml:"playtime"    default:"playtime"`
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

type UsersFilters struct {
	Jobs []Filter `yaml:"jobs"`
}

type UserLicensesTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns UserLicensesColumns `yaml:"columns"`
}

func (c *UserLicensesTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil {
		return prepareStringQuery(*c.Query, c.DBSyncTable, state, offset, limit)
	}

	where = append(where, getWhereCondition(c.DBSyncTable, state))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"license.type":  c.Columns.Type,
		"license.owner": c.Columns.OwnerIdentifier,
	}, where, offset, limit)
}

type UserLicensesColumns struct {
	Type            string `yaml:"type"            default:"type"`
	OwnerIdentifier string `yaml:"ownerIdentifier" default:"owner"`
}

type VehiclesTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns VehiclesColumns `yaml:"columns"`
}

func (c *VehiclesTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil {
		return prepareStringQuery(*c.Query, c.DBSyncTable, state, offset, limit)
	}

	where = append(where, getWhereCondition(c.DBSyncTable, state))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"vehicle.owner": c.Columns.OwnerIdentifier,
		"vehicle.plate": c.Columns.Plate,
		"vehicle.type":  c.Columns.Type,
		"vehicle.model": c.Columns.Model,
	}, where, offset, limit)
}

type VehiclesColumns struct {
	OwnerIdentifier string  `yaml:"ownerIdentifier" default:"owner"`
	OwnerID         *string `yaml:"ownerId"`
	Plate           string  `yaml:"plate"           default:"plate"`
	Type            string  `yaml:"type"            default:"type"`
	Model           string  `yaml:"model"           default:"model"`
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

type TableManagerConfig struct {
	Enabled bool `default:"true" yaml:"enabled"`
}
