// Package main implements the entrypoint for the application.
package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fivenet-app/fivenet/v2025/cmd"
	"github.com/fivenet-app/fivenet/v2025/cmd/envs"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
)

func main() {
	// https://github.com/DataDog/go-profiler-notes/blob/main/block.md#overhead
	// Thanks, to the authors of this document!
	runtime.SetBlockProfileRate(20000)
	runtime.SetMutexProfileFraction(100)

	ctx := kong.Parse(&cmd.Cli,
		kong.Vars{
			"version": version.Version,
		},
	)

	// Cli flag overrides env var
	if cmd.Cli.Config != "" {
		if err := os.Setenv(envs.ConfigFileEnvVar, cmd.Cli.Config); err != nil {
			panic(err)
		}
	}
	if cmd.Cli.SkipMigrations != nil {
		if err := os.Setenv(envs.SkipDBMigrationsEnv, strconv.FormatBool(*cmd.Cli.SkipMigrations)); err != nil {
			panic(err)
		}
	}
	if cmd.Cli.IgnoreRequirements != nil {
		if err := os.Setenv(envs.IgnoreRequirementsEnv, strconv.FormatBool(*cmd.Cli.IgnoreRequirements)); err != nil {
			panic(err)
		}
	}

	if cmd.Cli.StartTimeout <= 0 {
		cmd.Cli.StartTimeout = 180 * time.Second
	}

	err := ctx.Run(&cmd.Context{})
	ctx.FatalIfErrorf(err)
}
