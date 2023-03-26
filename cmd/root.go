package cmd

import (
	"fmt"
	"os"

	"github.com/galexrt/arpanet/pkg/config"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

var rootCmd = &cobra.Command{
	Use:     "arpanet",
	Version: version.Print("arpanet"),
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
