package main

import (
	"embed"

	"github.com/galexrt/arpanet/cmd"
)

//go:embed assets/*
var assets embed.FS

func main() {
	cmd.SetAssets(assets)
	cmd.Execute()
}
