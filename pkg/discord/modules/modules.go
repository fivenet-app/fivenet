package modules

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"sync/atomic"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/zap"
)

var (
	tOauth2Accs    = table.FivenetOauth2Accounts
	tAccs          = table.FivenetAccounts
	tJobsUserProps = table.FivenetJobsUserProps
)

var Modules = map[string]NewModuleFunc{}

type NewModuleFunc func(*BaseModule, *broker.Broker[any]) (Module, error)

type Module interface {
	GetName() string
	Plan(ctx context.Context) (*types.State, []discord.Embed, error)
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
	discord  *state.State
	guild    discord.Guild
	job      string
	cfg      *config.Discord
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.Enricher

	settings atomic.Pointer[users.DiscordSyncSettings]
}

func NewBaseModule(ctx context.Context, logger *zap.Logger, db *sql.DB, discord *state.State, guild discord.Guild, job string, cfg *config.Discord, appCfg appconfig.IConfig, enricher *mstlystcdata.Enricher, settings *users.DiscordSyncSettings) *BaseModule {
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

		settings: atomic.Pointer[users.DiscordSyncSettings]{},
	}
	bm.settings.Store(settings)

	return bm
}

func (m *BaseModule) checkIfJobIgnored(job string) bool {
	// Ignore certain jobs when syncing (e.g., "temporary" jobs), example:
	// "ambulance" job Discord, and an user is currently in the ignored job, e.g., "army".
	ignoredJobs := m.appCfg.Get().Discord.IgnoredJobs
	if m.job != job && slices.Contains(ignoredJobs, job) {
		return true
	}

	return false
}

func (m *BaseModule) SetSettings(settings *users.DiscordSyncSettings) {
	m.settings.Store(settings)
}
