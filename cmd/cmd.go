package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fivenet-app/fivenet/v2026/cmd/envs"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

type CLI struct {
	VersionFlag kong.VersionFlag `name:"version" short:"v" aliases:"V" help:"Print version information and quit"`

	Config             string        `help:"Config file path"                                  env:"FIVENET_CONFIG_FILE"`
	StartTimeout       time.Duration `help:"App start timeout duration"                        env:"FIVENET_START_TIMEOUT"       default:"180s"`
	SkipMigrations     *bool         `help:"Disable the automatic DB migrations on startup."   env:"FIVENET_SKIP_DB_MIGRATIONS"`
	IgnoreRequirements *bool         `help:"Ignore database and Nats requirements on startup." env:"FIVENET_IGNORE_REQUIREMENTS"`

	Version VersionCmd `help:"Print version information and quit." cmd:""`

	Server  ServerCmd  `cmd:"" help:"Run FiveNet server."`
	Worker  WorkerCmd  `cmd:"" help:"Run FiveNet worker."`
	Discord DiscordCmd `cmd:"" help:"Run FiveNet Discord bot."`
	DBSync  DBSyncCmd  `cmd:"" help:"Run FiveNet database sync." name:"dbsync" alias:"db-sync"`

	Update UpdateCmd `cmd:"" help:"Check for updates and update the FiveNet binary." alias:"upd"`

	Tools      ToolsCmd      `cmd:"" help:"Run FiveNet tools/helpers."`
	Migrations MigrationsCmd `cmd:"" help:"Run FiveNet migration helpers."`
}

type VersionCmd struct{}

func (c *VersionCmd) Run() error {
	fmt.Println(version.Version)
	return nil
}

type ToolsCmd struct {
	DB ToolsDBCmd `cmd:""`

	Sync ToolsSyncCmd `cmd:""`
}

type MigrationsCmd struct {
	HTMLToJSON MigrationsHTMLToJSONCmd `cmd:"" help:"Migrate documents, comments, etc., from (raw) HTML format to the legacy custom JSON format." name:"htmltojson"`
	Filestore  MigrationsFilestoreCmd  `cmd:"" help:"Migrate files from the old database format to the new filestore format."                     name:"filestore"`

	StatsBackfill MigrationsStatsBackfillCmd `cmd:"" help:"Backfill stats for documents." name:"statsbackfill"`
}

func (c *CLI) AfterApply(cli *CLI) error {
	// Cli flag overrides env var
	if cli.Config != "" {
		if err := os.Setenv(envs.ConfigFileEnvVar, cli.Config); err != nil {
			panic(err)
		}
	}
	if cli.SkipMigrations != nil {
		if err := os.Setenv(
			envs.SkipDBMigrationsEnv,
			strconv.FormatBool(*cli.SkipMigrations),
		); err != nil {
			panic(err)
		}
	}
	if cli.IgnoreRequirements != nil {
		if err := os.Setenv(
			envs.IgnoreRequirementsEnv,
			strconv.FormatBool(*cli.IgnoreRequirements),
		); err != nil {
			panic(err)
		}
	}

	return nil
}
