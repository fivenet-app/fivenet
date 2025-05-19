package config

import (
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
)

type Config struct {
	LogLevel string `default:"DEBUG" yaml:"logLevel"`
	Mode     string `default:"debug" yaml:"mode"`

	LogLevels map[string][]string `yaml:"logLevels"`

	Tracing Tracing `yaml:"tracing"`

	HTTP        HTTP       `yaml:"http"`
	Database    Database   `yaml:"database"`
	NATS        NATS       `yaml:"nats"`
	JWT         JWT        `yaml:"jwt"`
	Storage     Storage    `yaml:"storage"`
	ImageProxy  ImageProxy `yaml:"imageProxy"`
	Audit       Audit      `yaml:"audit"`
	OAuth2      OAuth2     `yaml:"oauth2"`
	PostalsFile string     `default:".output/public/data/postals.json" yaml:"postalsFile"`

	Auth           Auth           `yaml:"auth"`
	DispatchCenter DispatchCenter `yaml:"dispatchCenter"`

	Discord Discord `yaml:"discord"`

	Game Game `yaml:"game"`

	Sync Sync `yaml:"sync"`
}

type TracingExporter string

const (
	TracingExporter_StdoutTrace   TracingExporter = "stdout"
	TracingExporter_OTLPTraceGRPC TracingExporter = "otlptracegrpc"
	TracingExporter_OTLPTraceHTTP TracingExporter = "otlptracehttp"
)

type Tracing struct {
	Enabled     bool            `default:"false" yaml:"enabled"`
	Type        TracingExporter `default:"stdout" yaml:"type"`
	URL         string          `yaml:"url"`
	Insecure    bool            `yaml:"insecure"`
	Timeout     time.Duration   `default:"10s" yaml:"timeout"`
	Environment string          `default:"dev" yaml:"environment"`
	Ratio       float64         `default:"0.1" yaml:"ratio"`
	Attributes  []string        `yaml:"attributes"`
}

type HTTP struct {
	Listen         string   `default:":8080" yaml:"listen"`
	AdminListen    string   `default:":7070" yaml:"adminListen"`
	Sessions       Sessions `yaml:"sessions"`
	Links          Links    `yaml:"links"`
	PublicURL      string   `yaml:"publicURL"`
	Origins        []string `default:"" yaml:"origins"`
	TrustedProxies []string `yaml:"trustedProxies"`
}

type Sessions struct {
	CookieSecret string `yaml:"cookieSecret"`
	Domain       string `default:"localhost" yaml:"domain"`
}

type Links struct {
	PrivacyPolicy *string `json:"privacyPolicy"`
	Imprint       *string `json:"imprint"`
}

type Database struct {
	// Refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN string `yaml:"dsn"`

	Net       string `default:"tcp" yaml:"net"`
	Host      string `yaml:"host"`
	Port      int32  `default:"3306" yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	Collation string `default:"utf8mb4_unicode_ci" yaml:"collation"`

	// Connection options
	MaxOpenConns    int           `default:"32" yaml:"maxOpenConns"`
	MaxIdleConns    int           `default:"5" yaml:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `default:"15m" yaml:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `default:"60m" yaml:"connMaxLifetime"`

	ESXCompat bool `default:"false" yaml:"esxCompat"`

	Custom CustomDB `yaml:"custom"`
}

type CustomDB struct {
	Columns    dbutils.CustomColumns    `yaml:"columns"`
	Conditions dbutils.CustomConditions `yaml:"conditions"`
}

type NATS struct {
	URL      string  `default:"nats://localhost:4222" yaml:"url"`
	Replicas int     `default:"1" yaml:"replicas"`
	NKey     *string `yaml:"nKey"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

type StorageType string

const (
	StorageTypeS3         StorageType = "s3"
	StorageTypeFilesystem StorageType = "filesystem"
	StorageTypeNoop       StorageType = "noop"
)

type Storage struct {
	Type       StorageType       `default:"filesystem" yaml:"type"`
	Filesystem FilesystemStorage `yaml:"filesystem"`
	S3         S3Storage         `yaml:"s3"`
}

type FilesystemStorage struct {
	Path   string `yaml:"path"`
	Prefix string `yaml:"prefix"`
}

type S3Storage struct {
	Endpoint        string `yaml:"endpoint"`
	Region          string `default:"us-east-1" yaml:"region"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `default:"true" yaml:"useSSL"`
	BucketName      string `yaml:"bucketName"`
	Prefix          string `yaml:"prefix"`
	Retries         int    `default:"10" yaml:"retries"`
	UsePreSigned    bool   `default:"true" yaml:"usePresigned"`
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
	RetentionDays int `default:"180" yaml:"auditRetentionDays"`
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
	Icon          *string            `yaml:"icon"`
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

type Auth struct {
	SuperuserGroups []string `yaml:"superuserGroups"`
	SuperuserUsers  []string `yaml:"superuserUsers"`
}

type DispatchCenter struct {
	Type        string   `default:"gksphone" yaml:"type"`
	ConvertJobs []string `yaml:"convertJobs"`
}

type Discord struct {
	Enabled      bool                `default:"false" yaml:"enabled"`
	DryRun       bool                `default:"false" yaml:"dryRun"`
	Token        string              `yaml:"token"`
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
	Enabled             bool   `default:"false" yaml:"enabled"`
	GradeRoleFormat     string `default:"[%grade%] %grade_label%" yaml:"gradeRoleFormat"`
	EmployeeRoleFormat  string `default:"%s Personal" yaml:"employeeRoleFormat"`
	UnemployedRoleName  string `default:"Citizen" yaml:"unemployedRoleName"`
	JobsAbsceneRoleName string `default:"Absent" yaml:"jobsAbsceneRoleName"`
}

type DiscordGroupSync struct {
	Enabled bool                        `default:"false" yaml:"enabled"`
	Mapping map[string]DiscordGroupRole `yaml:"omitempty,mapping"`
}

type DiscordGroupRole struct {
	RoleName    string  `yaml:"roleName"`
	Permissions *uint64 `yaml:"omitempty,permissions"`
	Color       string  `yaml:"color"`
	NotSameJob  bool    `yaml:"notSameJob"`
}

type DiscordCommands struct {
	Enabled bool `default:"false" yaml:"enabled"`
}

type Game struct {
	StartJobGrade              int32 `default:"0" yaml:"startJobGrade"`
	CleanupRolesForMissingJobs bool  `default:"false" yaml:"cleanupRolesForMissingJobs"`
}

type Sync struct {
	Enabled   bool     `yaml:"enabled"`
	APITokens []string `yaml:"apiTokens"`
}
