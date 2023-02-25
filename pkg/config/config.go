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
}

type HTTP struct {
	Listen   string   `default:":8181" yaml:"listen"`
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

func init() {
	// Set defaults on start
	defaults.Set(C)
}
