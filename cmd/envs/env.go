// Package envs provides environment variable constants for FiveNet (CLI) configuration.
package envs

const (
	// ConfigFileEnvVar is the environment variable name for specifying the Fivenet config file path.
	ConfigFileEnvVar = "FIVENET_CONFIG_FILE"

	// SkipDBMigrationsEnv is the environment variable name to skip database migrations if set.
	SkipDBMigrationsEnv = "FIVENET_SKIP_DB_MIGRATIONS"

	// IgnoreRequirementsEnv is the environment variable name to ignore database and nats requirements if set.
	IgnoreRequirementsEnv = "FIVENET_IGNORE_REQUIREMENTS"
)
