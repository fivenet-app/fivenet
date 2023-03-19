package config

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

var (
	C = &Config{}
)

func init() {
	// Set defaults on start
	defaults.Set(C)
}

func InitConfigWithViper() {
	// Viper Config reading setup
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.Unmarshal(C)
}

type Config struct {
	LogLevel string `default:"DEBUG" yaml:"logLevel"`
	Mode     string `default:"debug" yaml:"mode"`

	Sentry   Sentry   `yaml:"sentry"`
	HTTP     HTTP     `yaml:"http"`
	GRPC     GRPC     `yaml:"grpc"`
	Database Database `yaml:"database"`
	JWT      JWT      `yaml:"jwt"`

	FiveM FiveM `yaml:"fivem"`
}

type Sentry struct {
	DSN         string `yaml:"dsn"`
	Environment string `default:"dev" yaml:"environment"`
}

type HTTP struct {
	Listen   string   `default:":8080" yaml:"listen"`
	Sessions Sessions `yaml:"sessions"`
}

type GRPC struct {
	Listen string `default:":9090" yaml:"listen"`
}

type Sessions struct {
	CookieSecret string `yaml:"cookieSecret"`
}

type Database struct {
	// refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN    string `yaml:"dsn"`
	DBName string `yaml:"dbName"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type FiveM struct {
	PermissionRoleJobs []string `yaml:"permissionRoleJobs"`
}
