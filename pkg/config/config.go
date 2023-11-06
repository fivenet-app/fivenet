package config

import (
	"fmt"
	"time"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		Load,
	),
)

type Result struct {
	fx.Out

	Config *Config
}

func Load() (*Config, error) {
	// Viper Config reading setup
	viper.SetEnvPrefix("FIVENET")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	c := &Config{}
	if err := defaults.Set(c); err != nil {
		return nil, fmt.Errorf("failed to set config defaults: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return c, nil
}

func LoadTest() (*Config, error) {
	c := &Config{}
	if err := defaults.Set(c); err != nil {
		return nil, fmt.Errorf("failed to set config defaults: %w", err)
	}

	return c, nil
}

type Config struct {
	LogLevel string `default:"DEBUG" yaml:"logLevel"`
	Mode     string `default:"debug" yaml:"mode"`

	Sentry  Sentry  `yaml:"sentry"`
	Tracing Tracing `yaml:"tracing"`

	HTTP     HTTP     `yaml:"http"`
	GRPC     GRPC     `yaml:"grpc"`
	Database Database `yaml:"database"`
	NATS     NATS     `yaml:"nats"`
	JWT      JWT      `yaml:"jwt"`
	OAuth2   OAuth2   `yaml:"oauth2"`
	Cache    Cache    `yaml:"cache"`

	Game    Game    `yaml:"game"`
	Discord Discord `yaml:"discord"`

	Storage Storage `yaml:"storage"`
}

type Sentry struct {
	ServerDSN   string `yaml:"serverDSN"`
	Environment string `default:"dev" yaml:"environment"`
	ClientDSN   string `default:"" yaml:"clientDSN"`
}

type Tracing struct {
	Enabled     bool   `default:"false" yaml:"enabled"`
	URL         string `yaml:"url"`
	Environment string `default:"dev" yaml:"environment"`
}

type HTTP struct {
	Listen      string   `default:":8080" yaml:"listen"`
	AdminListen string   `default:":7070" yaml:"adminListen"`
	Sessions    Sessions `yaml:"sessions"`
	Links       Links    `yaml:"links"`
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
}

type NATS struct {
	URL         string `default:"nats://localhost:4222" yaml:"url"`
	WorkerCount int    `default:"5" yaml:"workerCount"`
}

type JWT struct {
	Secret string `yaml:"secret"`
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

type Cache struct {
	RefreshTime time.Duration `default:"2m" yaml:"refreshTime"`
}

type Game struct {
	SignupEnabled      bool           `default:"true" yaml:"signupEnabled"`
	AuditRetentionDays *int           `default:"90" yaml:"auditRetentionDays"`
	SuperuserGroups    []string       `yaml:"superuserGroups"`
	UnemployedJob      UnemployedJob  `yaml:"unemployedJob"`
	PublicJobs         []string       `yaml:"publicJobs"`
	HiddenJobs         []string       `yaml:"hiddenJobs"`
	Livemap            Livemap        `yaml:"livemap"`
	DispatchCenter     DispatchCenter `yaml:"dispatchCenter"`
	DefaultPermissions []Perm         `yaml:"defaultPermissions"`
}
type UnemployedJob struct {
	Name  string `default:"unemployed" yaml:"job"`
	Grade int32  `default:"1" yaml:"grade"`
}

type Livemap struct {
	RefreshTime time.Duration `default:"3s350ms" yaml:"refreshTime"`
	Jobs        []string      `yaml:"jobs"`
	PostalsFile string        `default:".output/public/data/postals.json" yaml:"postalsFile"`
}

type DispatchCenter struct {
	ConvertJobs []string `yaml:"convertJobs"`
}

type Perm struct {
	Category string `yaml:"category"`
	Name     string `yaml:"name"`
}

type Discord struct {
	Bot DiscordBot `yaml:"bot"`
}

type DiscordBot struct {
	Enabled      bool                `default:"false" yaml:"enabled"`
	SyncInterval time.Duration       `default:"15m" yaml:"syncInterval"`
	InviteURL    string              `yaml:"inviteURL"`
	Token        string              `yaml:"token"`
	Presence     DiscordPresence     `yaml:"presence,omitempty"`
	UserInfoSync DiscordUserInfoSync `yaml:"userInfoSync"`
	GroupSync    DiscordGroupSync    `yaml:"groupSync"`
	Commands     DiscordCommands     `yaml:"commands"`
}

type DiscordPresence struct {
	GameStatus         *string `yaml:"gameStatus"`
	ListeningStatus    *string `yaml:"listeningStatus"`
	StreamingStatus    *string `yaml:"streamingStatus"`
	StreamingStatusUrl *string `yaml:"streamingStatusUrl"`
	WatchStatus        *string `yaml:"watchStatus"`
}

type DiscordUserInfoSync struct {
	Enabled       bool   `default:"false" yaml:"enabled"`
	RoleFormat    string `default:"[%02d] %s" yaml:"roleFormat"`
	NicknameRegex string `yaml:"nicknameRegex"`
}

type DiscordGroupSync struct {
	Enabled bool                        `default:"false" yaml:"enabled"`
	Mapping map[string]DiscordGroupRole `yaml:"omitempty,mapping"`
}

type DiscordGroupRole struct {
	RoleName    string `yaml:"roleName"`
	Permissions *int64 `yaml:"omitempty,permissions"`
	Color       string `yaml:"color"`
}

type DiscordCommands struct {
	Enabled bool `default:"false" yaml:"enabled"`
}

type Storage struct {
	Enabled         bool   `default:"false" yaml:"enabled"`
	Endpoint        string `yaml:"endpoint"`
	Region          string `default:"us-east-1" yaml:"region"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `default:"true" yaml:"useSSL"`
	BucketName      string `yaml:"bucketName"`
}
