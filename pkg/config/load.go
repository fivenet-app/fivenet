package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/fivenet-app/fivenet/v2025/cmd/envs"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module is the Fx module for providing the application config.
var Module = fx.Module("config",
	fx.Provide(
		Load,
	),
)

// Result is the output struct for the Load function, containing the main config and Discord config.
type Result struct {
	fx.Out

	// Config is the main application configuration
	Config *Config
	// DiscordConfig is a pointer to the Discord configuration
	DiscordConfig *Discord
}

// Load reads the application configuration from file and environment variables, sets defaults, and returns a Result.
// Returns an error if the config cannot be loaded, parsed, or defaults cannot be set.
func Load() (Result, error) {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	// Viper config reading setup
	v.SetEnvPrefix("FIVENET")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")

	if configFile := os.Getenv(envs.ConfigFileEnvVar); configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("/config")
	}

	v.AutomaticEnv()

	res := Result{}
	// Find and read the config file
	if err := v.ReadInConfig(); err != nil {
		return res, fmt.Errorf("fatal error config file. %w", err)
	}

	c := &Config{}
	if err := defaults.Set(c); err != nil {
		return res, fmt.Errorf("failed to set config defaults. %w", err)
	}
	res.Config = c
	res.DiscordConfig = &c.Discord

	if err := v.Unmarshal(c); err != nil {
		return res, fmt.Errorf("failed to unmarshal config. %w", err)
	}

	// Handle non-DSN database connection details
	if c.Database.DSN == "" {
		m := mysql.NewConfig()
		m.Net = c.Database.Net
		m.Addr = fmt.Sprintf("%s:%d", c.Database.Host, c.Database.Port)
		m.User = c.Database.Username
		m.Passwd = c.Database.Password
		m.DBName = c.Database.Database
		m.Collation = c.Database.Collation
		m.ParseTime = true

		c.Database.DSN = m.FormatDSN()
	}

	// Ensure origins are lower case
	for i := range c.HTTP.Origins {
		c.HTTP.Origins[i] = strings.ToLower(c.HTTP.Origins[i])
	}

	return res, nil
}

// TestModule is the Fx module for providing a test config.
var TestModule = fx.Module("config_test",
	fx.Provide(LoadTestConfig),
)

// LoadTestConfig returns a config with defaults set for testing.
// Sets audit log retention days high so they won't run in "short" tests.
func LoadTestConfig() (*Config, error) {
	c := &Config{}
	if err := defaults.Set(c); err != nil {
		return nil, fmt.Errorf("failed to set config defaults. %w", err)
	}

	// Set audit log retention days high so they won't run in "short" tests
	c.Audit.RetentionDays = 365

	return c, nil
}
