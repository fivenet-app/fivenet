// Package main implements the entrypoint for the application.
package main

import (
	"runtime"

	"github.com/alecthomas/kong"
	"github.com/fivenet-app/fivenet/v2026/cmd"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

func main() {
	// https://github.com/DataDog/go-profiler-notes/blob/main/block.md#overhead
	// Thanks, to the authors of this document!
	runtime.SetBlockProfileRate(20000)
	runtime.SetMutexProfileFraction(100)

	cli := &cmd.CLI{}
	ctx := kong.Parse(cli,
		kong.Vars{
			"version": version.Version,
		},
		kong.Bind(cli),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
