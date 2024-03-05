package config

import (
	"sync/atomic"
	"time"

	"github.com/galexrt/fivenet/pkg/utils/dbutils"
)

type Config struct {
	cfg atomic.Pointer[BaseConfig]
}

func (c *Config) Get() *BaseConfig {
	return c.cfg.Load()
}

func (c *Config) Set(val *BaseConfig) {
	c.cfg.Store(val)
}

type BaseConfig struct {
	LogLevel string `default:"DEBUG" yaml:"logLevel"`
	Mode     string `default:"debug" yaml:"mode"`

	Tracing Tracing `yaml:"tracing"`

	HTTP       HTTP       `yaml:"http"`
	GRPC       GRPC       `yaml:"grpc"`
	Database   Database   `yaml:"database"`
	NATS       NATS       `yaml:"nats"`
	JWT        JWT        `yaml:"jwt"`
	Storage    Storage    `yaml:"storage"`
	ImageProxy ImageProxy `yaml:"imageProxy"`
	Cache      Cache      `yaml:"cache"`
	Audit      Audit      `yaml:"audit"`

	AppConfig
}

type Tracing struct {
	Enabled     bool    `default:"false" yaml:"enabled"`
	URL         string  `yaml:"url"`
	Environment string  `default:"dev" yaml:"environment"`
	Ratio       float64 `default:"0.1" yaml:"ratio"`
}

type HTTP struct {
	Listen      string   `default:":8080" yaml:"listen"`
	AdminListen string   `default:":7070" yaml:"adminListen"`
	Sessions    Sessions `yaml:"sessions"`
	Links       Links    `yaml:"links"`
	PublicURL   string   `yaml:"publicURL"`
}

type Sessions struct {
	CookieSecret string `yaml:"cookieSecret"`
	Domain       string `default:"localhost" yaml:"domain"`
}

type Links struct {
	PrivacyPolicy *string `json:"privacyPolicy"`
	Imprint       *string `json:"imprint"`
}

type GRPC struct {
	Listen string `default:":9090" yaml:"listen"`
}

type Database struct {
	// refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN string `yaml:"dsn"`

	// Connection options
	MaxOpenConns    int           `default:"32" yaml:"maxOpenConns"`
	MaxIdleConns    int           `default:"5" yaml:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `default:"15m" yaml:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `default:"60m" yaml:"connMaxLifetime"`

	Custom CustomDB `yaml:"custom"`
}

type CustomDB struct {
	Columns    dbutils.CustomColumns    `yaml:"columns"`
	Conditions dbutils.CustomConditions `yaml:"conditions"`
}

type NATS struct {
	URL string `default:"nats://localhost:4222" yaml:"url"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type Storage struct {
	Type       string            `default:"filesystem" yaml:"type"`
	Filesystem FilesystemStorage `yaml:"filesystem"`
	S3         S3Storage         `yaml:"s3"`
}

type FilesystemStorage struct {
	Path string `yaml:"path"`
}

type S3Storage struct {
	Endpoint        string `yaml:"endpoint"`
	Region          string `default:"us-east-1" yaml:"region"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `default:"true" yaml:"useSSL"`
	BucketName      string `yaml:"bucketName"`
	Prefix          string `yaml:"prefix"`
}

type ImageProxy struct {
	Enabled     bool              `default:"true" yaml:"enabled"`
	URL         string            `default:"/api/image_proxy/" yaml:"url"`
	CachePrefix string            `default:"images/" yaml:"cachePrefix"`
	Options     ImageProxyOptions `yaml:"options"`
}

type ImageProxyOptions struct {
	AllowHosts []string `yaml:"allowHosts"`
	DenyHosts  []string `yaml:"denyHosts"`
}

type Cache struct {
	RefreshTime time.Duration `default:"2m" yaml:"refreshTime"`
}

type Audit struct {
	RetentionDays *int `default:"90" yaml:"auditRetentionDays"`
}
