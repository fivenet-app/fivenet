package dbsyncconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/creasty/defaults"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var watchSetupOnce sync.Once

var Module = fx.Module("dbsync.config",
	fx.Provide(
		New,
	),
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

func New(p ParamsConfig) (ResultConfig, error) {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	s := &Config{
		shutdowner: p.Shutdowner,
		v:          v,
	}

	// Viper config reading setup
	v.SetEnvPrefix("FIVENET")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")

	if configFile := os.Getenv("FIVENET_DBSYNC_FILE"); configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("dbsync")
		v.AddConfigPath(".")
		v.AddConfigPath("/config")
	}

	v.AutomaticEnv()

	if err := s.LoadConfig(); err != nil {
		return ResultConfig{}, err
	}

	cc := s.cfg.Load()
	r := ResultConfig{
		Config: s,
		Cfg: &config.Config{
			Mode:     cc.Mode,
			LogLevel: cc.LogLevel,
			Log:      cc.Log,

			HTTP: config.HTTP{
				AdminListen: cc.AdminListen,
			},

			// Ignore db requirements, dbsync doesn't need them
			IgnoreRequirements: true,
			UpdateCheck:        cc.UpdateCheck,

			Database: config.Database{
				DatabaseConnection: cc.Source.DatabaseConnection,
			},
		},
	}
	// Do not run DB migrations for dbsync
	r.Cfg.Database.SkipMigrations = true

	if cc.Tables.Accounts.Enabled {
		if cc.Tables.Accounts.Query == nil {
			return r, errors.New(
				"accounts table is enabled but no query is set. please set a query for the accounts table or disable the feature",
			)
		}
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

	// If the destination URL or token is set, populate the API config for backward
	// compatibility with older config versions
	if c.Destination.URL != "" || c.Destination.Token != "" {
		c.Destination.API = DBSyncDestinationAPI{
			URL:      c.Destination.URL,
			Token:    c.Destination.Token,
			Insecure: c.Destination.Insecure,
		}

		c.Destination.URL = ""
		c.Destination.Token = ""
	}

	if err := c.Init(); err != nil {
		return fmt.Errorf("failed to initialize config. %w", err)
	}

	// Validate config
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]
		// Skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	if err := validate.Struct(c); err != nil {
		// Build detailed validation error message
		var msg strings.Builder
		msg.WriteString("Invalid FiveNet DBSync config detected:\n")
		for _, validationErr := range func() validator.ValidationErrors {
			var target validator.ValidationErrors
			_ = errors.As(err, &target)
			return target
		}() {
			fmt.Fprintf(&msg, "- Field %q violated %s validation.\n",
				validationErr.StructNamespace(),
				validationErr.Tag())
		}
		return fmt.Errorf("%s", msg.String())
	}

	s.cfg.Store(c)

	return nil
}

func (s *Config) SetupWatch(logger *zap.Logger, restartFn func() error) {
	watchSetupOnce.Do(func() {
		s.setupWatch(logger, restartFn)
	})
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

	AdminListen string `default:"" yaml:"adminListen"`

	Source      DBSyncSource      `yaml:"source"`
	Destination DBSyncDestination `yaml:"destination"`

	Tables DBSyncSourceTables `yaml:"tables"`

	UpdateCheck config.UpdateCheck `yaml:"updateCheck"`

	TableManager TableManagerConfig `yaml:"tableManager"`

	Limits SyncLimits `yaml:"limits"`
}

type DBSyncSource struct {
	config.DatabaseConnection `yaml:",inline" mapstructure:",squash"`

	BaseDataResyncInterval time.Duration `default:"5m" yaml:"baseDataResyncInterval" validate:"gte=1"`
}

type DBSyncMethod string

const (
	DBSyncModeAPI DBSyncMethod = "api"
	DBSyncModeDB  DBSyncMethod = "db"
)

type DBSyncDestination struct {
	DBSyncDestinationAPI `yaml:",inline" mapstructure:",squash"`

	Method DBSyncMethod `default:"api" yaml:"method"`

	API DBSyncDestinationAPI `yaml:"api"`

	SyncInterval time.Duration `default:"5s"    yaml:"syncInterval" validate:"gte=1"`
	DryRun       bool          `default:"false" yaml:"dryRun"`
}

type DBSyncDestinationAPI struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`

	Timeout time.Duration `yaml:"timeout" default:"15s" validate:"gte=1"`

	Insecure          bool `yaml:"insecure"`
	TransportSecurity bool `yaml:"transportSecurity" default:"true"`
}

type DBSyncSourceTables struct {
	Jobs      JobsTable      `yaml:"jobs"`
	JobGrades JobGradesTable `yaml:"jobGrades"`

	Licenses LicensesTable `yaml:"licenses"`

	Accounts AccountsTable `yaml:"accounts"`

	Users            UsersTable            `yaml:"users"`
	UserLicenses     UserLicensesTable     `yaml:"userLicenses"`
	UserJobs         UserJobsTable         `yaml:"userJobs"`
	UserPhoneNumbers UserPhoneNumbersTable `yaml:"userPhoneNumbers"`

	Vehicles VehiclesTable `yaml:"vehicles"`
}

func (c *DBSyncSourceTables) GetAllTables() []DBSyncTable {
	tables := []DBSyncTable{}

	if c == nil {
		return tables
	}

	// Base data
	if c.Jobs.Enabled {
		tables = append(tables, c.Jobs.DBSyncTable)
	}
	if c.Licenses.Enabled {
		tables = append(tables, c.Licenses.DBSyncTable)
	}

	// Main data
	if c.Accounts.Enabled {
		tables = append(tables, c.Accounts.DBSyncTable)
	}
	if c.Users.Enabled {
		tables = append(tables, c.Users.DBSyncTable)
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

	// Compiled regex pattern for internal use (compiled during config load)
	CompiledPattern *regexp.Regexp `yaml:"-"`
}

type DBSyncTable struct {
	Enabled           bool           `yaml:"enabled"                     default:"false"`
	TableName         string         `yaml:"tableName"`
	UpdatedTimeColumn *string        `yaml:"updatedTimeColumn,omitempty"`
	Query             *string        `yaml:"query,omitempty"`
	SyncInterval      *time.Duration `yaml:"syncInterval,omitempty"                      validate:"omitempty,gte=1"`
}

func (c *DBSyncTable) GetSyncInterval() *time.Duration {
	// Minimum sync interval is 1 second
	if c.SyncInterval != nil && *c.SyncInterval >= 1*time.Second {
		return c.SyncInterval
	}

	return nil
}

type JobsTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns JobsColumns `yaml:"columns"`

	Filters []Filter `yaml:"filters"`
}

func (c *JobsTable) GetQuery(
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		return prepareStringQuery(*c.Query, c.DBSyncTable, nil, limit, "")
	}

	where = append(where, getWhereCondition(c.DBSyncTable, nil, ""))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"job.name":  c.Columns.Name,
		"job.label": c.Columns.Label,
	}, where, limit, []string{c.Columns.Name})
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
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		q := prepareStringQuery(*c.Query, c.DBSyncTable, state, limit, "")
		q = strings.ReplaceAll(q, "$jobName", "?")
		return q
	}

	where = append(where, fmt.Sprintf("%#q = ?", c.Columns.JobName))
	where = append(where, getWhereCondition(c.DBSyncTable, state, ""))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"job_grade.job_name": c.Columns.JobName,
		"job_grade.grade":    c.Columns.Grade,
		"job_grade.name":     c.Columns.Name,
		"job_grade.label":    c.Columns.Label,
	}, where, limit, []string{c.Columns.JobName, c.Columns.Grade})
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
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		return prepareStringQuery(*c.Query, c.DBSyncTable, nil, limit, "")
	}
	where = append(where, getWhereCondition(c.DBSyncTable, nil, ""))
	return buildQueryFromColumns(c.TableName, map[string]string{
		"license.type":  c.Columns.Type,
		"license.label": c.Columns.Label,
	}, where, limit, []string{c.Columns.Type})
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

	ResyncInterval *time.Duration `yaml:"resyncInterval,omitempty" validate:"omitempty,gte=1"`
}

func (c *UsersTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		where = slices.DeleteFunc(where, func(cond string) bool {
			return strings.TrimSpace(cond) == ""
		})
		if len(where) > 0 {
			return prepareStringQueryWithWhereCondition(*c.Query, limit, where)
		}

		return prepareStringQuery(*c.Query, c.DBSyncTable, state, limit, c.Columns.ID)
	}

	orderBy := []string{}
	if c.UpdatedTimeColumn != nil && *c.UpdatedTimeColumn != "" {
		// Ensure updated time column is always the first in the order by for consistent pagination.
		orderBy = append([]string{*c.UpdatedTimeColumn + " ASC"}, orderBy...)
	}
	orderBy = append(orderBy, c.Columns.ID+" ASC")

	columns := map[string]string{
		"user.id":           c.Columns.ID,
		"user.identifier":   c.Columns.Identifier,
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
	}
	if c.UpdatedTimeColumn != nil && *c.UpdatedTimeColumn != "" {
		columns["user.updated_at"] = *c.UpdatedTimeColumn
	}

	where = append(where, getWhereCondition(c.DBSyncTable, state, c.Columns.ID))
	return buildQueryFromColumns(c.TableName, columns, where, limit, orderBy)
}

// GetSyncQuery returns the users query while always combining state cursor conditions with extra where conditions.
func (c *UsersTable) GetSyncQuery(
	state *TableSyncState,
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		return prepareStringQueryWithStateAndWhere(
			*c.Query,
			c.DBSyncTable,
			state,
			limit,
			c.Columns.ID,
			where,
		)
	}

	return c.GetQuery(state, 0, limit, where...)
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

func (c *ValueMapping) Process(input string) string {
	val, ok := c.Values[input]
	if ok {
		return val
	} else if c.Fallback != nil && *c.Fallback != "" {
		// Only return fallback if it's not an empty string
		return *c.Fallback
	}

	return input
}

type UsersFilters struct {
	Jobs []Filter `yaml:"jobs"`
}

type UserLicensesTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns UserLicensesColumns `yaml:"columns"`
}

func (c *UserLicensesTable) GetQuery(
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		return prepareStringQuery(*c.Query, c.DBSyncTable, nil, limit, "")
	}

	where = append(where, getWhereCondition(c.DBSyncTable, nil, ""))
	where = append(where, "`"+c.Columns.OwnerIdentifier+"` = $identifier")
	return buildQueryFromColumns(c.TableName, map[string]string{
		"license.type":  c.Columns.Type,
		"license.owner": c.Columns.OwnerIdentifier,
	}, where, limit, []string{c.Columns.Type, c.Columns.OwnerIdentifier})
}

type UserLicensesColumns struct {
	Type            string `yaml:"type"            default:"type"`
	OwnerIdentifier string `yaml:"ownerIdentifier" default:"owner"`
}

type UserJobsTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`
}

func (c *UserJobsTable) GetQuery(
	offset int64,
	limit int64,
	where ...string,
) string {
	return prepareStringQuery(*c.Query, c.DBSyncTable, nil, limit, "")
}

type UserPhoneNumbersTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`
}

func (c *UserPhoneNumbersTable) GetQuery(
	offset int64,
	limit int64,
	where ...string,
) string {
	return prepareStringQuery(*c.Query, c.DBSyncTable, nil, limit, "")
}

type VehiclesTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	Columns VehiclesColumns `yaml:"columns"`

	ResyncInterval *time.Duration `yaml:"resyncInterval,omitempty" validate:"omitempty,gte=1"`
}

func (c *VehiclesTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		return prepareStringQuery(*c.Query, c.DBSyncTable, state, limit, c.Columns.Plate)
	}

	orderBy := []string{}
	if c.UpdatedTimeColumn != nil && *c.UpdatedTimeColumn != "" {
		// Ensure updated time column is always the first in the order by for consistent pagination.
		orderBy = append([]string{*c.UpdatedTimeColumn + " ASC"}, orderBy...)
	}
	orderBy = append(orderBy, c.Columns.Plate+" ASC")

	columns := map[string]string{
		"vehicle.ownerIdentifier": c.Columns.OwnerIdentifier,
		"vehicle.plate":           c.Columns.Plate,
		"vehicle.type":            c.Columns.Type,
		"vehicle.model":           c.Columns.Model,
	}
	if c.UpdatedTimeColumn != nil && *c.UpdatedTimeColumn != "" {
		columns["vehicle.updated_at"] = *c.UpdatedTimeColumn
	}

	where = append(where, getWhereCondition(c.DBSyncTable, state, c.Columns.Plate))
	return buildQueryFromColumns(c.TableName, columns, where, limit, orderBy)
}

// GetSyncQuery returns the vehicles query while always combining state cursor conditions with extra where conditions.
func (c *VehiclesTable) GetSyncQuery(
	state *TableSyncState,
	limit int64,
	where ...string,
) string {
	if c.Query != nil && *c.Query != "" {
		return prepareStringQueryWithStateAndWhere(
			*c.Query,
			c.DBSyncTable,
			state,
			limit,
			c.Columns.Plate,
			where,
		)
	}

	return c.GetQuery(state, 0, limit, where...)
}

type VehiclesColumns struct {
	OwnerIdentifier string  `yaml:"ownerIdentifier" default:"owner"`
	OwnerID         *string `yaml:"ownerId"`
	Plate           string  `yaml:"plate"           default:"plate"`
	Type            string  `yaml:"type"            default:"type"`
	Model           string  `yaml:"model"           default:"model"`
}

type AccountsTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`
}

func (c *AccountsTable) GetQuery(
	state *TableSyncState,
	offset int64,
	limit int64,
	where ...string,
) string {
	return prepareStringQuery(*c.Query, c.DBSyncTable, state, limit, "license")
}

type DBSyncTableSyncInterval interface {
	GetSyncInterval() *time.Duration
}

func (c *DBSyncConfig) GetSyncInterval(table DBSyncTableSyncInterval) time.Duration {
	if table != nil && table.GetSyncInterval() != nil {
		// Only return the interval if it's at least 1 second
		interval := table.GetSyncInterval()
		if *interval >= 1*time.Second {
			return *interval
		}
	}

	return c.Destination.SyncInterval
}

func (c *DBSyncConfig) Init() error {
	// Compile filter regex patterns
	for k := range c.Tables.Jobs.Filters {
		filter := &c.Tables.Jobs.Filters[k]
		var err error
		filter.CompiledPattern, err = regexp.Compile(filter.Pattern)
		if err != nil {
			return fmt.Errorf("failed to compile regex for filter %d. %w", k, err)
		}
	}

	for k := range c.Tables.Users.Filters.Jobs {
		filter := &c.Tables.Users.Filters.Jobs[k]
		var err error
		filter.CompiledPattern, err = regexp.Compile(filter.Pattern)
		if err != nil {
			return fmt.Errorf("failed to compile regex for user job filter %d. %w", k, err)
		}
	}

	// Unset updatedTimeColumn for tables where it's not applicable, to avoid confusion.
	c.Tables.JobGrades.UpdatedTimeColumn = nil
	c.Tables.UserJobs.UpdatedTimeColumn = nil
	c.Tables.UserLicenses.UpdatedTimeColumn = nil
	c.Tables.UserPhoneNumbers.UpdatedTimeColumn = nil

	// Remove OFFSET from queries if present, as pagination is handled differently since v2026.3.0.
	if c.Tables.Jobs.Query != nil {
		*c.Tables.Jobs.Query = strings.ReplaceAll(*c.Tables.Jobs.Query, "OFFSET $offset", "")
	}
	if c.Tables.Licenses.Query != nil {
		*c.Tables.Licenses.Query = strings.ReplaceAll(
			*c.Tables.Licenses.Query,
			"OFFSET $offset",
			"",
		)
	}
	if c.Tables.Accounts.Query != nil {
		*c.Tables.Accounts.Query = strings.ReplaceAll(
			*c.Tables.Accounts.Query,
			"OFFSET $offset",
			"",
		)
	}
	if c.Tables.Users.Query != nil {
		*c.Tables.Users.Query = strings.ReplaceAll(*c.Tables.Users.Query, "OFFSET $offset", "")

		if strings.TrimSpace(*c.Tables.Users.Query) != "" &&
			!strings.Contains(*c.Tables.Users.Query, "$whereCondition") {
			zap.L().Error(
				"users custom query missing required $whereCondition placeholder",
			)

			return errors.New(
				"users table custom query must contain $whereCondition placeholder",
			)
		}
	}
	if c.Tables.Vehicles.Query != nil {
		*c.Tables.Vehicles.Query = strings.ReplaceAll(
			*c.Tables.Vehicles.Query,
			"OFFSET $offset",
			"",
		)
	}

	return nil
}

type TableManagerConfig struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

// SyncLimits defines limits for the sync process, such as maximum number of records to sync per table.
// Accounts, Users and Vehicles split the data into multiple requests if they exceed the API limits.
type SyncLimits struct {
	Jobs     int64 `default:"200" yaml:"jobs"     validate:"omitempty,gte=1,lte=200"`
	Licenses int64 `default:"200" yaml:"licenses" validate:"omitempty,gte=1,lte=200"`

	Accounts int64 `default:"200" yaml:"accounts" validate:"omitempty,gte=1,lte=1000"`
	Users    int64 `default:"150" yaml:"users"    validate:"omitempty,gte=1,lte=1000"`
	Vehicles int64 `default:"500" yaml:"vehicles" validate:"omitempty,gte=1,lte=1500"`
}
