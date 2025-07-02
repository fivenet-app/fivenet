package config

import (
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/zaputils"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Mode string `default:"release" yaml:"mode"`

	LogLevel string `default:"DEBUG" enum:"" yaml:"logLevel"`
	// Any empty log level will be set to the `.LogLevel`
	LogLevelOverrides LogLevelOverrides `yaml:"logLevelOverrides"`

	// Secret used to encrypt/decrypt data in, e.g., the database
	Secret string `yaml:"secret"`

	Demo Demo `yaml:"demo"`

	JWT            JWT            `yaml:"jwt"`
	HTTP           HTTP           `yaml:"http"`
	Database       Database       `yaml:"database"`
	NATS           NATS           `yaml:"nats"`
	Storage        Storage        `yaml:"storage"`
	ImageProxy     ImageProxy     `yaml:"imageProxy"`
	Audit          Audit          `yaml:"audit"`
	OAuth2         OAuth2         `yaml:"oauth2"`
	PostalsFile    string         `default:".output/public/data/postals.json" yaml:"postalsFile"`
	Auth           Auth           `yaml:"auth"`
	DispatchCenter DispatchCenter `yaml:"dispatchCenter"`
	Discord        Discord        `yaml:"discord"`
	Game           Game           `yaml:"game"`
	Sync           Sync           `yaml:"sync"`
	OTLP           OTLPConfig     `yaml:"otlp"`
}

type LoggingComponent string

var (
	LoggingComponentKVStore     LoggingComponent = "kvstore"
	LoggingComponentCron        LoggingComponent = "cron"
	LoggingComponentPerms       LoggingComponent = "perms"
	LoggingComponentHousekeeper LoggingComponent = "housekeeper"
)

type LogLevelOverrides map[string]string

func (l LogLevelOverrides) Get(component LoggingComponent, defaultLevel string) zapcore.Level {
	if level, ok := l[string(component)]; ok && level != "" {
		return zaputils.StringToLevel(level)
	}

	return zaputils.StringToLevel(defaultLevel)
}

type Demo struct {
	Enabled   bool     `default:"false" yaml:"enabled"`
	TargetJob string   `default:"police" yaml:"targetJob"`
	Users     []string `yaml:"users"`
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

// Database represents the configuration for connecting to a MySQL database.
// It includes credentials, connection settings, and additional options.
type Database struct {
	// DSN is the Data Source Name used to connect to the database.
	// Refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details.
	DSN string `yaml:"dsn"`

	// Net specifies the network type to use for the connection (e.g., "tcp").
	Net string `default:"tcp" yaml:"net"`

	// Host is the hostname or IP address of the MySQL server.
	Host string `yaml:"host"`

	// Port is the port number on which the MySQL server is listening.
	Port int32 `default:"3306" yaml:"port"`

	// Username is the username for authenticating with the MySQL server.
	Username string `yaml:"username"`

	// Password is the password for authenticating with the MySQL server.
	Password string `yaml:"password"`

	// Database is the name of the specific database to connect to.
	Database string `yaml:"database"`

	// Collation specifies the character set collation to use for the connection.
	Collation string `default:"utf8mb4_unicode_ci" yaml:"collation"`

	// MaxOpenConns defines the maximum number of open connections to the database.
	MaxOpenConns int `default:"32" yaml:"maxOpenConns"`

	// MaxIdleConns defines the maximum number of idle connections to the database.
	MaxIdleConns int `default:"5" yaml:"maxIdleConns"`

	// ConnMaxIdleTime specifies the maximum amount of time a connection can remain idle.
	ConnMaxIdleTime time.Duration `default:"15m" yaml:"connMaxIdleTime"`

	// ConnMaxLifetime specifies the maximum amount of time a connection can remain open.
	ConnMaxLifetime time.Duration `default:"60m" yaml:"connMaxLifetime"`

	// ESXCompat enables compatibility mode for ESX-specific database configurations.
	ESXCompat bool `default:"false" yaml:"esxCompat"`

	// Custom contains additional custom database configuration options.
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
	Retries         int    `default:"3" yaml:"retries"`
	CheckOnStartup  bool   `default:"false" yaml:"checkOnStartup"`
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

	PermsCacheSize int           `default:"1024" yaml:"permsCacheSize"`
	PermsCacheTTL  time.Duration `default:"30s" yaml:"permsCacheTTL"`
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

type OtelExporter string

const (
	TracingExporter_StdoutTrace OtelExporter = "stdout"
	TracingExporter_OTLPGRPC    OtelExporter = "grpc"
	TracingExporter_OTLPHTTP    OtelExporter = "http"
)

type OTLPConfig struct {
	Enabled     bool          `default:"false" yaml:"enabled"`
	Type        OtelExporter  `default:"stdout" yaml:"type"`
	URL         string        `yaml:"url"`
	Insecure    bool          `yaml:"insecure"`
	Timeout     time.Duration `default:"10s" yaml:"timeout"`
	Environment string        `default:"dev" yaml:"environment"`
	Ratio       float64       `default:"0.1" yaml:"ratio"`
	Attributes  []string      `yaml:"attributes"`
	// Headers to send with OTLP HTTP requests
	Headers map[string]string `yaml:"headers,omitempty"`
	// Compression type for OTLP HTTP requests
	Compression string             `default:"none" yaml:"compression"`
	Frontend    OTLPFrontendConfig `yaml:"frontend"`
}

type OTLPFrontendConfig struct {
	// Public URL for traces and other instrumentation (if set, only then instrumentation is enabled in the frontend)
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}
