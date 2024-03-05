package config

import "time"

type AppConfig struct {
	Game    Game    `yaml:"game"`
	Discord Discord `yaml:"discord"`
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
