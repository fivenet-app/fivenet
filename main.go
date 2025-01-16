package main

import (
	"os"
	"runtime"
	"time"

	"github.com/alecthomas/kong"

	// GRPC Services
	"github.com/fivenet-app/fivenet/cmd"
	// Modules
)

func main() {
	// https://github.com/DataDog/go-profiler-notes/blob/main/block.md#overhead
	// Thanks, to the authors of this document!
	runtime.SetBlockProfileRate(20000)
	runtime.SetMutexProfileFraction(100)

	ctx := kong.Parse(&cmd.Cli)

	// Cli flag overrides env var
	if cmd.Cli.Config != "" {
		if err := os.Setenv("FIVENET_CONFIG_FILE", cmd.Cli.Config); err != nil {
			panic(err)
		}
	}
	if cmd.Cli.StartTimeout <= 0 {
		cmd.Cli.StartTimeout = 180 * time.Second
	}

	err := ctx.Run(&cmd.Context{})
	ctx.FatalIfErrorf(err)
}
