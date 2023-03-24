package main

import (
	"embed"
	"math/rand"
	"time"

	"github.com/galexrt/arpanet/cmd"
)

//go:embed dist/*
var assets embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())

	cmd.SetAssets(assets)
	cmd.Execute()
}
