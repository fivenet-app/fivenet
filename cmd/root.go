package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/galexrt/arpanet/pkg/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	assets embed.FS
)

var rootCmd = &cobra.Command{
	Use: "arpanet",
}

func init() {
	cobra.OnInitialize(config.InitConfigWithViper)
}

func Execute() {
	// Logger Setup
	loggerConfig := zap.NewProductionConfig()
	level, err := zapcore.ParseLevel(config.C.LogLevel)
	if err != nil {
		panic("failed to parse log level from config")
	}
	loggerConfig.Level.SetLevel(level)

	logger, err = loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Run Command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SetAssets(fs embed.FS) {
	assets = fs
}
