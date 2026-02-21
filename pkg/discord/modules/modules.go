package modules

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	jobssettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"go.uber.org/zap"
)

type IState interface {
	Member(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error)

	AddRole(
		guildID discord.GuildID,
		userID discord.UserID,
		roleID discord.RoleID,
		data api.AddRoleData,
	) error
}

var (
	tAccsOauth2     = table.FivenetAccountsOauth2
	tColleagueProps = table.FivenetJobColleagueProps
)

var Modules = map[string]NewModuleFunc{}

type NewModuleFunc func(*BaseModule, *broker.Broker[any]) (Module, error)

type Module interface {
	GetName() string
	Plan(ctx context.Context) (*discordtypes.State, []discord.Embed, error)
}

func GetModule(name string, base *BaseModule, events *broker.Broker[any]) (Module, error) {
	fn, ok := Modules[name]
	if !ok {
		return nil, fmt.Errorf("no module found by name %s", name)
	}

	// "Wrap" logger with module name
	base.logger = base.logger.Named(name)

	return fn(base, events)
}

type BaseModule struct {
	ctx      context.Context
	logger   *zap.Logger
	db       *sql.DB
	discord  IState
	guild    discord.Guild
	job      string
	cfg      *config.Discord
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.Enricher

	oauth2ProviderName string

	settings *atomic.Pointer[jobssettings.DiscordSyncSettings]
}

func NewBaseModule(
	ctx context.Context,
	logger *zap.Logger,
	db *sql.DB,
	discord *state.State,
	guild discord.Guild,
	job string,
	cfg *config.Discord,
	appCfg appconfig.IConfig,
	enricher *mstlystcdata.Enricher,
	oauth2ProviderName string,
	settings *atomic.Pointer[jobssettings.DiscordSyncSettings],
) *BaseModule {
	bm := &BaseModule{
		ctx:      ctx,
		logger:   logger,
		db:       db,
		discord:  discord,
		guild:    guild,
		job:      job,
		cfg:      cfg,
		appCfg:   appCfg,
		enricher: enricher,

		oauth2ProviderName: oauth2ProviderName,

		settings: settings,
	}

	return bm
}

func (m *BaseModule) GetOAuth2ProviderName() string {
	if m.oauth2ProviderName != "" {
		return m.oauth2ProviderName
	}
	// Sane default
	return "discord"
}
