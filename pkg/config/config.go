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

	OAuth2  OAuth2  `yaml:"oauth2"`
	Game    Game    `yaml:"game"`
	Discord Discord `yaml:"discord"`
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

type Audit struct {
	RetentionDays *int `default:"90" yaml:"auditRetentionDays"`
}

type Game struct {
	Auth           Auth           `yaml:"auth"`
	UnemployedJob  UnemployedJob  `yaml:"unemployedJob"`
	PublicJobs     []string       `yaml:"publicJobs"`
	HiddenJobs     []string       `yaml:"hiddenJobs"`
	Livemap        Livemap        `yaml:"livemap"`
	DispatchCenter DispatchCenter `yaml:"dispatchCenter"`
}

type Auth struct {
	SignupEnabled      bool     `default:"true" yaml:"signupEnabled"`
	SuperuserGroups    []string `yaml:"superuserGroups"`
	DefaultPermissions []Perm   `yaml:"defaultPermissions"`
}

type UnemployedJob struct {
	Name  string `default:"unemployed" yaml:"job"`
	Grade int32  `default:"1" yaml:"grade"`
}

type Livemap struct {
	RefreshTime   time.Duration `default:"3s350ms" yaml:"refreshTime"`
	DBRefreshTime time.Duration `default:"1s" yaml:"dbRefreshTime"`
	Jobs          []string      `yaml:"jobs"`
	PostalsFile   string        `default:".output/public/data/postals.json" yaml:"postalsFile"`
}

type DispatchCenter struct {
	ConvertJobs []string `yaml:"convertJobs"`
}

type Perm struct {
	Category string `yaml:"category"`
	Name     string `yaml:"name"`
}

type Discord struct {
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
	Enabled             bool     `default:"false" yaml:"enabled"`
	GradeRoleFormat     string   `default:"[%grade%] %grade_label%" yaml:"gradeRoleFormat"`
	EmployeeRoleFormat  string   `default:"%s Personal" yaml:"employeeRoleFormat"`
	NicknameRegex       string   `yaml:"nicknameRegex"`
	IgnoreJobs          []string `yaml:"ignoreJobs"`
	UnemployedRoleName  string   `default:"Citizen" yaml:"unemployedRoleName"`
	JobsAbsceneRoleName string   `default:"Absent" yaml:"jobsAbsceneRoleName"`
}

type DiscordGroupSync struct {
	Enabled bool                        `default:"false" yaml:"enabled"`
	Mapping map[string]DiscordGroupRole `yaml:"omitempty,mapping"`
}

type DiscordGroupRole struct {
	RoleName    string `yaml:"roleName"`
	Permissions *int64 `yaml:"omitempty,permissions"`
	Color       string `yaml:"color"`
	NotSameJob  bool   `yaml:"notSameJob"`
}

type DiscordCommands struct {
	Enabled bool `default:"false" yaml:"enabled"`
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
