package config

import (
	"fmt"
	"time"

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
	viper.AddConfigPath("/config")
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
	OAuth2   OAuth2   `yaml:"oauth2"`

	Game Game `yaml:"game"`
}

type Sentry struct {
	ServerDSN   string `yaml:"serverDSN"`
	Environment string `default:"dev" yaml:"environment"`
	ClientDSN   string `default:"" yaml:"clientDSN"`
}

type HTTP struct {
	Listen   string   `default:":8080" yaml:"listen"`
	Sessions Sessions `yaml:"sessions"`
}

type GRPC struct {
	Listen    string `default:":9090" yaml:"listen"`
	ClientURL string `default:"/grpc" yaml:"clientURL"`
}

type Sessions struct {
	CookieSecret string `yaml:"cookieSecret"`
	Domain       string `default:"localhost" yaml:"domain"`
}

type Database struct {
	// refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN    string `yaml:"dsn"`
	DBName string `yaml:"dbName"`

	MaxOpenConns    int           `default:"32" yaml:"maxOpenConns"`
	MaxIdleConns    int           `default:"5" yaml:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `default:"15m" yaml:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `default:"60m" yaml:"connMaxLifetime"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type Game struct {
	SuperuserGroups    []string      `yaml:"superuserGroups"`
	UnemployedJob      UnemployedJob `yaml:"unemployedJob"`
	PublicJobs         []string      `yaml:"publicJobs"`
	PermissionRoleJobs []string      `yaml:"permissionRoleJobs"`
	LivemapJobs        []string      `yaml:"livemapJobs"`
}

type UnemployedJob struct {
	Name  string `default:"unemployed" yaml:"job"`
	Grade int32  `default:"1" yaml:"grade"`
}

type OAuth2 struct {
	Providers []*OAuth2Provider
}

type OAuth2ProviderType string

const (
	OAuth2ProviderGeneric OAuth2ProviderType = "generic"
	OAuth2ProviderDiscord OAuth2ProviderType = "discord"
)

type OAuth2Provider struct {
	Name          string             `yaml:"name"`
	Label         string             `yaml:"label"`
	Homepage      string             `yaml:"homepage"`
	Type          OAuth2ProviderType `yaml:"type"`
	DefaultAvatar string             `yaml:"defaultAvatar"`
	RedirectURL   string             `yaml:"redirectURL"`
	ClientID      string             `yaml:"clientID"`
	ClientSecret  string             `yaml:"clientSecret"`
	Scopes        []string           `yaml:"scopes"`
	Endpoints     OAuth2Endpoints    `yaml:"endpoints"`
	Mapping       *OAuth2Mapping     `yaml:"omitempty,mapping"`
}

type OAuth2Endpoints struct {
	AuthURL     string `yaml:"authURL"`
	TokenURL    string `yaml:"tokenURL"`
	UserInfoURL string `yaml:"userInfoURL"`
}

type OAuth2Mapping struct {
	ID       string `yaml:"id"`
	Username string `yaml:"username"`
	Avatar   string `yaml:"avatar"`
}
