package config

import (
	"github.com/creasty/defaults"
)

var (
	C = &Config{}
)

type Config struct {
	LogLevel string `default:"DEBUG" yaml:"logLevel"`
	Mode     string `default:"debug" yaml:"mode"`

	HTTP     HTTP     `yaml:"http"`
	GRPC     GRPC     `yaml:"grpc"`
	Database Database `yaml:"database"`
	JWT      JWT      `yaml:"jwt"`

	FiveM FiveM `yaml:"fivem"`
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
	DSN string `yaml:"dsn"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type FiveM struct {
	PermissionRoleJobs []string `yaml:"permissionRoleJobs"`
}

func init() {
	// Set defaults on start
	defaults.Set(C)
}
