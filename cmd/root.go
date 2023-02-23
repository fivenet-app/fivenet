package cmd

import (
	"fmt"
	"os"

	"github.com/galexrt/arpanet/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

var rootCmd = &cobra.Command{
	Use: "arpanet",
}

func init() {
	cobra.OnInitialize(initConfig)
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

func initConfig() {
	// Viper Config reading setup
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.Unmarshal(config.C)
}
