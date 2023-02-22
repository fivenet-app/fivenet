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

	HTTP     HTTP             `yaml:"http"`
	Database PostgresDatabase `yaml:"database"`
	FiveM    FiveM            `yaml:"fiveM"`
}

type HTTP struct {
	Listen   string   `default:":12345" yaml:"listen"`
	Sessions Sessions `yaml:"sessions"`
}

type Sessions struct {
	CookieSecret string `yaml:"cookieSecret"`
}

type PostgresDatabase struct {
	DSN string `yaml:"dsn"`
}

type FiveM struct {
	Database MySQLDatabase `yaml:"database"`
}

type MySQLDatabase struct {
	Host     string `default:"localhost" yaml:"host"`
	Port     int    `default:"3306" yaml:"3306"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
}

func init() {
	// Set defaults on start
	defaults.Set(C)
}
