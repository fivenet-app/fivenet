package modules

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/zap"
)

var (
	tOauth2Accs    = table.FivenetOauth2Accounts
	tAccs          = table.FivenetAccounts
	tUsers         = table.Users.AS("users")
	tJobsUserProps = table.FivenetJobsUserProps
)

var Modules = map[string]NewModuleFunc{}

type NewModuleFunc func(*BaseModule) (Module, error)

type Module interface {
	Plan(ctx context.Context) (*types.Plan, []*discordgo.MessageEmbed, error)
}

func GetModule(name string, base *BaseModule) (Module, error) {
	fn, ok := Modules[name]
	if !ok {
		return nil, fmt.Errorf("no module found by name %s", name)
	}

	// "Wrap" logger with module name
	base.logger = base.logger.Named(name)

	return fn(base)
}

type BaseModule struct {
	logger   *zap.Logger
	db       *sql.DB
	discord  *discordgo.Session
	guild    *discordgo.Guild
	job      string
	cfg      *config.Discord
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.Enricher

	settings *users.DiscordSyncSettings
}

func NewBaseModule(logger *zap.Logger, db *sql.DB, discord *discordgo.Session, guild *discordgo.Guild, job string, cfg *config.Discord, appCfg appconfig.IConfig, enricher *mstlystcdata.Enricher, settings *users.DiscordSyncSettings) *BaseModule {
	return &BaseModule{
		logger:   logger,
		db:       db,
		discord:  discord,
		guild:    guild,
		job:      job,
		cfg:      cfg,
		appCfg:   appCfg,
		enricher: enricher,

		settings: settings,
	}
}
