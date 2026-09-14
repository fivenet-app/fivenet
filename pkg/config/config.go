//nolint:tagliatelle // To prevent linting issues with the tagliatelle linter (due to, e.g., `URL`, `TTL`).
package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/zaputils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
)

const (
	CwdConfigDir       = "."
	ContainerConfigDir = "/config"
	SystemConfigDir    = "/etc/fivenet"
)

type Config struct {
	// Mode `debug` or `release` mode. `debug` mode affects performance and logging; `release` mode is the default and should be used in production.
	Mode string `default:"release" yaml:"mode"`

	// LogLevel controls the minimum severity emitted by the application logger.
	LogLevel string `default:"INFO" yaml:"logLevel" enum:"DEBUG,INFO,WARN,ERROR,PANIC,FATAL"`
	Log      Log    `               yaml:"log"`

	// Secret is the app-wide secret used to encrypt and decrypt stored data. It must be exactly 32 bytes and must differ from JWT.Secret.
	Secret string `yaml:"secret"`

	// IgnoreRequirements skips startup checks for database and NATS connectivity. Use only when those services are intentionally unavailable at boot.
	IgnoreRequirements bool `default:"false" yaml:"ignoreRequirements"`

	// Demo configures development-only demo data generation. It is deliberately omitted from config.example.yaml.
	Demo Demo `yaml:"demo"`

	// JWT configures JWT signing.
	JWT JWT `yaml:"jwt"`
	// HTTP configures the public and administration HTTP servers.
	HTTP HTTP `yaml:"http"`
	// Database configures the MySQL connection.
	Database Database `yaml:"database"`
	// NATS configures the NATS and JetStream connection.
	NATS NATS `yaml:"nats"`
	// Storage configures file and object storage.
	Storage Storage `yaml:"storage"`
	// ImageProxy configures the remote image proxy.
	ImageProxy ImageProxy `yaml:"imageProxy"`
	// Icons configures serving Iconify icon sets.
	Icons Icons `yaml:"icons"`
	// Audit configures audit-log retention.
	Audit Audit `yaml:"audit"`
	// OAuth2 configures external OAuth2 identity providers.
	OAuth2 OAuth2 `yaml:"oauth2"`
	// PostalsFile is the path to the postals data file.
	PostalsFile string `yaml:"postalsFile" default:".output/public/data/postals.json" validate:"filepath"`
	// Auth configures administrative access and permission caching.
	Auth Auth `yaml:"auth"`
	// DispatchCenter configures dispatch-center integration.
	DispatchCenter DispatchCenter `yaml:"dispatchCenter"`
	// Discord configures the Discord bot and role synchronization.
	Discord Discord `yaml:"discord"`
	// Game configures game-server conventions.
	Game Game `yaml:"game"`
	// Sync configures the external synchronization API.
	Sync Sync `yaml:"sync"`
	// OTLP configures OpenTelemetry tracing.
	OTLP OTLPConfig `yaml:"otlp"`
	// UpdateCheck configures periodic update checks.
	UpdateCheck UpdateCheck `yaml:"updateCheck"`
}

type Log struct {
	// LogToStderr If enabled logs are sent to stderr instead of stdout.
	LogToStderr bool `default:"false" yaml:"logToStderr"`
	// LogToFile enables logging to a file instead of stdout.
	LogToFile bool `default:"false" yaml:"logToFile"`
	// File contains the configuration for file logging (if LogToFile is true, make sure to configure this).
	File LogFile `yaml:"file"`

	// LevelOverrides overrides LogLevel for individual components. Empty values use LogLevel.
	LevelOverrides LogLevelOverrides `yaml:"levelOverrides"`
}

type LogFile struct {
	// Path to the log file.
	Path string `yaml:"path" default:"./fivenet.log"`
	// Rotation contains the configuration for log rotation.
	Rotation LogRotation `yaml:"rotation"`
}

type LogRotation struct {
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	MaxSize int `default:"10" yaml:"maxSize"`
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `default:"7" yaml:"maxBackups"`
	// MaxAge is the maximum number of days to retain old log files based on the timestamp.
	MaxAge int `default:"14" yaml:"maxAge"`
	// Deprecated: Use `Compression` instead. This field is kept for backward compatibility.
	Compress bool `default:"true" yaml:"compress"`
	// Compression determines the compression algorithm to use for rotated log files. It can be "none", "gzip", or "zstd".
	Compression string `default:"zstd" yaml:"compress"`
	// Rotation interval for log rotation (e.g., daily rotation).
	RotationInterval time.Duration `default:"24h" yaml:"rotationInterval"`
}

type LoggingComponent string

const (
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
	Enabled   bool         `default:"false"  yaml:"enabled"`
	Seed      uint64       `default:"42"     yaml:"seed"`
	TargetJob string       `default:"police" yaml:"targetJob"`
	Features  DemoFeatures `                 yaml:"features"`
	FakeUsers DemoFakeUser `                 yaml:"fakeUsers"`
}

type DemoFeatures struct {
	Dispatches      bool `default:"true"  yaml:"dispatches"`
	Locations       bool `default:"true"  yaml:"locations"`
	Markers         bool `default:"true"  yaml:"markers"`
	Timeclock       bool `default:"true"  yaml:"timeclock"`
	CalendarEntries bool `default:"false" yaml:"calendarEntries"`
	Users           bool `default:"false" yaml:"users"`
	Vehicles        bool `default:"false" yaml:"vehicles"`
}

type DemoFakeUser struct {
	Count int `default:"50" yaml:"count"`
}

type HTTP struct {
	// Listen is the address for the public HTTP server.
	Listen string `default:":8080" yaml:"listen"`
	// AdminListen is the address for metrics and Go debug endpoints. Leave it empty to disable the admin server.
	AdminListen string `default:":7070" yaml:"adminListen"`
	// Sessions configures session cookies.
	Sessions Sessions `yaml:"sessions"`
	// Links configures optional legal links displayed by the frontend.
	Links Links `yaml:"links"`
	// PublicURL is the canonical public base URL, including scheme and host.
	PublicURL string `yaml:"publicURL"`
	// Origins lists browser origins allowed to access the frontend and API.
	Origins []string `default:"" yaml:"origins"`
	// TrustedProxies lists reverse-proxy IPs or CIDRs trusted to forward client IP and protocol headers.
	TrustedProxies []string `yaml:"trustedProxies"`
}

type Sessions struct {
	// CookieSecret signs and encrypts session cookies.
	CookieSecret string `yaml:"cookieSecret"`
	// Domain is the cookie domain.
	Domain string `yaml:"domain" default:"localhost"`
}

type Links struct {
	// PrivacyPolicy is an optional URL to the privacy policy.
	PrivacyPolicy *string `yaml:"privacyPolicy"`
	// Imprint is an optional URL to the legal notice.
	Imprint *string `yaml:"imprint"`
}

// Database represents the configuration for connecting to a MySQL database.
// It includes credentials, connection settings, and additional options.
type Database struct {
	DatabaseConnection `yaml:",inline" mapstructure:",squash"`

	// Custom contains additional custom database configuration options.
	Custom CustomDB `yaml:"custom"`
}

type DatabaseConnection struct {
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
	// DisableLocking disables the use of table locking in the database (mainly for migrations).
	DisableLocking bool `default:"false" yaml:"disableLocking"`

	// ESXCompat enables compatibility mode for ESX-specific database configurations.
	ESXCompat bool `default:"false" yaml:"esxCompat"`
	// SkipMigrations indicates whether to skip database migrations on startup.
	SkipMigrations bool `default:"false" yaml:"skipMigrations"`
}

type CustomDB struct {
	// Columns maps FiveNet fields to custom database column names. Set a column to "-" to disable it.
	Columns dbutils.CustomColumns `yaml:"columns"`
	// Conditions configures custom database query conditions.
	Conditions dbutils.CustomConditions `yaml:"conditions"`
}

type NATS struct {
	// URL is the NATS server URL.
	URL string `default:"nats://localhost:4222" yaml:"url"`
	// Replicas is the number of JetStream replicas.
	Replicas int `default:"1" yaml:"replicas"`
	// NKey is an optional NATS NKey seed used for authentication.
	NKey *string `yaml:"nKey"`
	// Creds is an optional path to a NATS credentials file.
	Creds *string `yaml:"creds"`
	// PublishAsyncMaxPending NATS Jetstream publish async max pending limit. NATS' default is `4000`.
	PublishAsyncMaxPending int `default:"2000" yaml:"publishAsyncMaxPending"`
}

type JWT struct {
	// Secret is the JWT signing secret. It must differ from the app-wide Secret.
	Secret string `yaml:"secret"`
}

type StorageType string

const (
	StorageTypeS3         StorageType = "s3"
	StorageTypeFilesystem StorageType = "filesystem"
	StorageTypeNoop       StorageType = "noop"
)

type Storage struct {
	// Type selects the storage backend: "filesystem", "s3", or "noop". Only the matching backend configuration is used.
	Type StorageType `default:"filesystem" yaml:"type"`
	// Filesystem configures local filesystem storage.
	Filesystem FilesystemStorage `yaml:"filesystem"`
	// S3 configures S3-compatible object storage.
	S3 S3Storage `yaml:"s3"`

	// MetricsEnabled enables periodic collection of storage metrics.
	MetricsEnabled bool `default:"true" yaml:"metricsEnabled"`
	// MetricsInterval is the period between storage metric collections.
	MetricsInterval time.Duration `default:"15m" yaml:"metricsInterval"`
}

func ValidateStorage(fl validator.StructLevel) {
	storage, ok := fl.Current().Interface().(Storage)
	if !ok {
		fl.ReportError(storage, "Storage", "storage", "storage", "")
		return
	}

	var prefix string
	switch storage.Type {
	case StorageTypeFilesystem:
		prefix = storage.Filesystem.Prefix
	case StorageTypeS3:
		prefix = storage.S3.Prefix

	case StorageTypeNoop:
		return

	default:
		fl.ReportError(storage, "Type", "type", "storagetype", "")
		return
	}

	if prefix != "" {
		p := filepath.Clean(filepath.FromSlash(prefix))
		if p == "." || p == ".." || strings.HasPrefix(p, ".."+string(os.PathSeparator)) {
			fl.ReportError(storage, "Prefix", "prefix", "storageprefix", "")
			return
		}
		if filepath.IsAbs(p) || filepath.VolumeName(p) != "" {
			fl.ReportError(storage, "Prefix", "prefix", "storageprefix", "")
			return
		}
	}
}

type FilesystemStorage struct {
	// Path is the local filesystem root.
	Path string `yaml:"path"`
	// Prefix is an optional relative key prefix.
	Prefix string `yaml:"prefix"`
}

type S3Storage struct {
	// Endpoint is the S3-compatible API endpoint.
	Endpoint string `yaml:"endpoint"`
	// Region is the S3 region.
	Region string `yaml:"region" default:"us-east-1"`
	// AccessKeyID is the S3 access key ID.
	AccessKeyID string `yaml:"accessKeyID"`
	// SecretAccessKey is the S3 secret access key.
	SecretAccessKey string `yaml:"secretAccessKey"`
	// UseSSL controls whether the S3 endpoint is contacted over TLS.
	UseSSL bool `yaml:"useSSL" default:"true"`
	// BucketName is the target S3 bucket.
	BucketName string `yaml:"bucketName"`
	// Prefix is an optional relative key prefix.
	Prefix string `yaml:"prefix"`
	// Retries is the number of storage operation retries.
	Retries int `yaml:"retries" default:"3"`
	// CheckOnStartup verifies S3 access during startup.
	CheckOnStartup bool `yaml:"checkOnStartup" default:"false"`
}

type ImageProxy struct {
	// AllowHosts is a list of hostnames or wildcard patterns that are allowed to be proxied. If empty, all hostnames are allowed.
	// Deny rules take precedence over allow rules.
	AllowHosts []string `yaml:"allowHosts"`
	// DenyHosts is a list of hostnames or wildcard patterns that are denied from being proxied. If empty, no hostnames are denied.
	DenyHosts []string `yaml:"denyHosts"`
	// AllowIPs allow public IP addresses to be proxied. If false, only hostnames are allowed.
	AllowIPs bool `yaml:"allowIPs" default:"false"`
	// MinimumCacheDuration is the minimum duration for which proxied images should be cached. If not set, defaults to 60 minutes.
	MinimumCacheDuration time.Duration `yaml:"minimumCacheDuration" default:"60m"`
}

type Icons struct {
	// If true, the backend server will act as a proxy for the Iconify API (URL specified via `APIURL` setting).
	Proxy bool `default:"false" yaml:"proxy"`
	// If you are using the proxy mode, make sure to support the Iconify project: https://iconify.design/sponsors/
	APIURL string `default:"https://api.iconify.design" yaml:"apiUrl"`
	// Path to the directory containing icon sets (used when proxy is disabled; this path works with the official FiveNet container images).
	Path string `default:"./icons" yaml:"path"`
}

type Cache struct {
	// RefreshTime specifies the duration after which cached data should be refreshed. Must be greater than or equal to 1 second.
	RefreshTime time.Duration `default:"2m" yaml:"refreshTime" validate:"gte=1"`
}

type Audit struct {
	// RetentionDays specifies the number of days to retain audit logs before they are purged. Must be greater than or equal to 1.
	RetentionDays int `default:"180" yaml:"retentionDays" validate:"gte=1"`
}

type OAuth2 struct {
	// Providers is a list of configured OAuth2 providers. Each provider must have a unique name and type (except `type: generic`).
	Providers []*OAuth2Provider `yaml:"providers" validate:"dive"`
}

func (c *OAuth2) GetProviderByType(t OAuth2ProviderType) *OAuth2Provider {
	for _, provider := range c.Providers {
		if provider.Type == t {
			return provider
		}
	}
	return nil
}

// OAuth2ProviderNameMaxLen is the maximum allowed length for OAuth2 provider names.
// Provider name and OAuth2 requests use struct tags to validate the name length, so this constant is used in tests.
const OAuth2ProviderNameMaxLen = 64

type OAuth2ProviderType string

const (
	OAuth2ProviderGeneric OAuth2ProviderType = "generic"
	OAuth2ProviderDiscord OAuth2ProviderType = "discord"
	OAuth2ProviderSteam   OAuth2ProviderType = "steam"
)

type OAuth2Provider struct {
	// Name is the internal provider name used in config and callback routes.
	Name string `yaml:"name" validate:"required,min=1,max=64"`
	// Label is the display label shown on the login button.
	Label string `yaml:"label"`
	// Homepage is the provider homepage shown for reference in the UI.
	Homepage string `yaml:"homepage"`
	// Type is the provider implementation used by FiveNet.
	Type OAuth2ProviderType `yaml:"type"`
	// Icon is an Iconify icon name or image URL shown on the login button.
	Icon *string `yaml:"icon"`
	// DefaultAvatar is an optional fallback avatar URL for users without a provider avatar.
	DefaultAvatar string `yaml:"defaultAvatar"`
	// RedirectURL is the callback URL registered with the OAuth2 provider.
	RedirectURL string `yaml:"redirectURL"`
	// ClientID is the OAuth2 client ID issued by the provider.
	ClientID string `yaml:"clientID"`
	// ClientSecret is the OAuth2 client secret issued by the provider.
	ClientSecret string `yaml:"clientSecret"`
	// Scopes lists the OAuth2 scopes requested from the provider.
	Scopes []string `yaml:"scopes"`
	// Endpoints optionally overrides the provider's OAuth2 endpoints.
	Endpoints OAuth2Endpoints `yaml:"endpoints"`
	// Mapping tells FiveNet which fields to read from the provider response.
	Mapping *OAuth2Mapping `yaml:"mapping,omitempty"`
}

func (p OAuth2Provider) GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  p.RedirectURL,
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Scopes:       p.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   p.Endpoints.AuthURL,
			TokenURL:  p.Endpoints.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

type OAuth2Endpoints struct {
	// AuthURL is the authorization endpoint used to start the OAuth2 flow.
	AuthURL string `yaml:"authURL"`
	// TokenURL is the token endpoint used to exchange the authorization code.
	TokenURL string `yaml:"tokenURL"`
	// UserInfoURL is the endpoint used to fetch the user profile.
	UserInfoURL string `yaml:"userInfoURL"`
}

type OAuth2Mapping struct {
	// ID is the unique user identifier claim.
	ID string `yaml:"id"`
	// Username is the username or display-name claim.
	Username string `yaml:"username"`
	// Avatar is the avatar URL claim.
	Avatar string `yaml:"avatar"`
}

type Auth struct {
	// SuperuserGroups is used as a legacy fallback if the ConfigAdminGroups is empty.
	SuperuserGroups []string `yaml:"superuserGroups"`
	// SuperuserUsers is used as a legacy fallback if the ConfigAdminUsers is empty.
	SuperuserUsers []string `yaml:"superuserUsers"`

	// JobAdminGroups is used to define groups that have basic administrative privileges for managing jobs.
	JobAdminGroups []string `yaml:"jobAdminGroups"`
	// JobAdminUsers is used to define users that have basic administrative privileges for managing jobs.
	JobAdminUsers []string `yaml:"jobAdminUsers"`

	// ConfigAdminGroups is used to define groups that have full config administrative privileges for managing jobs and the configuration.
	ConfigAdminGroups []string `yaml:"configAdminGroups"`
	// ConfigAdminUsers is used to define users that have full config administrative privileges for managing jobs and the configuration.
	ConfigAdminUsers []string `yaml:"configAdminUsers"`

	// PermsCacheSize is the maximum number of permission entries to cache in memory. If set to 0, caching is disabled.
	PermsCacheSize int `default:"1024" yaml:"permsCacheSize"`
	// PermsCacheTTL is the time-to-live for cached permission entries. After this duration, cached entries will be invalidated and refreshed from the source. Must be greater than or equal to 1 second.
	PermsCacheTTL time.Duration `default:"30s" yaml:"permsCacheTTL"`
}

type DispatchCenter struct {
	// Enabled enables dispatch-center integration for phone plugin dispatch systems.
	Enabled bool `default:"false" yaml:"enabled"`
	// Type is the dispatch-center implementation to use.
	Type string `default:"gksphone" yaml:"type"`
	// ConvertJobs lists jobs whose dispatches should be converted.
	ConvertJobs []string `yaml:"convertJobs"`
}

type Discord struct {
	// Enabled enables the Discord integration.
	Enabled bool `default:"false" yaml:"enabled"`
	// DryRun prevents role changes and messages, for testing.
	DryRun bool `default:"false" yaml:"dryRun"`
	// Sync controls whether any Discord synchronization is performed.
	Sync bool `default:"true" yaml:"sync"`
	// Token is the Discord bot token.
	Token string `yaml:"token"`

	// UserInfoSync configures synchronization of job, grade, and qualification information to Discord roles.
	UserInfoSync DiscordUserInfoSync `yaml:"userInfoSync"`
	// GroupSync configures synchronization of server groups to Discord roles.
	GroupSync DiscordGroupSync `yaml:"groupSync"`
	// Qualifications configures synchronization of qualifications to Discord roles.
	Qualifications DiscordQualificationsSync `yaml:"qualifications"`
	// Commands configures Discord command registration and availability.
	Commands DiscordCommands `yaml:"commands"`
	// CalendarReminders enables calendar reminders in Discord.
	CalendarReminders bool `yaml:"calendarReminders" default:"true"`
}

func (c *Discord) IsAnyEnabled() bool {
	return c.Enabled || c.Sync || c.Commands.Enabled
}

type DiscordPresence struct {
	GameStatus         *string `yaml:"gameStatus"`
	ListeningStatus    *string `yaml:"listeningStatus"`
	StreamingStatus    *string `yaml:"streamingStatus"`
	StreamingStatusUrl *string `yaml:"streamingStatusUrl"`
	WatchStatus        *string `yaml:"watchStatus"`
}

type DiscordUserInfoSync struct {
	// Enabled enables user-information role synchronization.
	Enabled bool `default:"false" yaml:"enabled"`
	// GradeRoleFormat formats job-grade role names. It supports %grade% and %grade_label%.
	GradeRoleFormat string `default:"[%grade%] %grade_label%" yaml:"gradeRoleFormat"`
	// EmployeeRoleFormat formats employee-count role names. It receives the employee count via %s.
	EmployeeRoleFormat string `default:"%s Personal" yaml:"employeeRoleFormat"`
	// UnemployedRoleName is the role assigned to unemployed users.
	UnemployedRoleName string `default:"Citizen" yaml:"unemployedRoleName"`
	// JobsAbsceneRoleName is the role assigned to users absent from their job.
	JobsAbsceneRoleName string `default:"Absent" yaml:"jobsAbsceneRoleName"`
}

type DiscordGroupSync struct {
	// Enabled enables group-to-role synchronization.
	Enabled bool `default:"false" yaml:"enabled"`
	// Mapping maps server group names to Discord role settings.
	Mapping map[string]DiscordGroupRole `yaml:"mapping,omitempty"`
}

type DiscordQualificationsSync struct {
	// Enabled enables qualification-to-role synchronization.
	Enabled bool `default:"false" yaml:"enabled"`
}

type DiscordGroupRole struct {
	// RoleName is the Discord role name.
	RoleName string `yaml:"roleName"`
	// Permissions optionally sets the Discord permission bitset for the role.
	Permissions *int64 `yaml:"permissions,omitempty"`
	// Color is the Discord role color as a hexadecimal value.
	Color string `yaml:"color"`
	// NotSameJob only assigns this role when the user's job differs from the group's job.
	NotSameJob bool `yaml:"notSameJob"`
}

type DiscordCommands struct {
	// Enabled enables Discord command registration.
	Enabled bool `default:"false" yaml:"enabled"`

	// Absent enables the absent command.
	Absent bool `default:"true" yaml:"absent"`
	// Fivenet enables the fivenet command.
	Fivenet bool `default:"true" yaml:"fivenet"`
	// Help enables the help command.
	Help bool `default:"true" yaml:"help"`
	// Sync enables the sync command.
	Sync bool `default:"true" yaml:"sync"`
}

type Game struct {
	// StartJobGrade is the first job-grade number. Some servers start at 0, others at 1.
	StartJobGrade int32 `default:"0" yaml:"startJobGrade"`
	// CleanupRolesForMissingJobs deletes roles for jobs that no longer exist.
	CleanupRolesForMissingJobs bool `default:"false" yaml:"cleanupRolesForMissingJobs"`
}

type Sync struct {
	// Enabled enables the sync API endpoint required by the FiveNet plugin and DBSync function.
	Enabled bool `yaml:"enabled"`
	// APITokens lists API tokens accepted for sync requests.
	APITokens []string `yaml:"apiTokens"`
}

type OtelExporter string

const (
	TracingExporterStdoutTrace OtelExporter = "stdout"
	TracingExporterOTLPGRPC    OtelExporter = "grpc"
	TracingExporterOTLPHTTP    OtelExporter = "http"
)

type OTLPConfig struct {
	Enabled     bool          `default:"false"  yaml:"enabled"`
	Type        OtelExporter  `default:"stdout" yaml:"type"`
	URL         string        `                 yaml:"url"`
	Insecure    bool          `                 yaml:"insecure"`
	Timeout     time.Duration `default:"10s"    yaml:"timeout"`
	Environment string        `default:"dev"    yaml:"environment"`
	Ratio       float64       `default:"0.1"    yaml:"ratio"`
	Attributes  []string      `                 yaml:"attributes"`
	// Headers to send with OTLP HTTP requests
	Headers map[string]string `yaml:"headers,omitempty"`
	// Compression type for OTLP HTTP requests
	Compression string             `default:"none" yaml:"compression"`
	Frontend    OTLPFrontendConfig `               yaml:"frontend"`
}

type OTLPFrontendConfig struct {
	// Public URL for traces and other instrumentation (if set, only then instrumentation is enabled in the frontend)
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type UpdateCheck struct {
	// Enabled enables periodic update checks.
	Enabled bool `default:"true" yaml:"enabled"`
	// Interval is the period between update checks.
	Interval time.Duration `default:"6h" yaml:"interval"`
}
