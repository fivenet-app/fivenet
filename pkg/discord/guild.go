package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/discord/commands"
	"github.com/galexrt/fivenet/pkg/discord/modules"
	"go.uber.org/zap"
)

type Guild struct {
	Job string `alias:"job"`
	ID  string `alias:"id"`

	logger   *zap.Logger
	bot      *Bot
	guild    *discordgo.Guild
	modules  map[string]modules.Module
	commands *commands.Cmds
}

func NewGuild(b *Bot, guild *discordgo.Guild, job string) (*Guild, error) {
	ms := map[string]modules.Module{}

	base := modules.NewBaseModule(b.ctx, b.logger.Named("module").With(zap.String("job", job), zap.String("discord_guild_id", guild.ID)),
		b.db, b.discord, guild, job, b.cfg, b.enricher)

	gModules := []string{}
	if b.cfg.UserInfoSync.Enabled {
		gModules = append(gModules, "userinfo")
	}
	if b.cfg.GroupSync.Enabled {
		gModules = append(gModules, "groupsync")
	}

	var err error
	for _, name := range gModules {
		ms[name], err = modules.GetModule(name, base)
		if err != nil {
			return nil, err
		}
	}

	var cmds *commands.Cmds
	if b.cfg.Commands.Enabled {
		cmds = commands.New(guild.ID)
	}

	return &Guild{
		Job: job,
		ID:  guild.ID,

		logger:   b.logger.Named("discord_guild").With(zap.String("job", job), zap.String("discord_guild_id", guild.ID)),
		bot:      b,
		guild:    guild,
		modules:  ms,
		commands: cmds,
	}, nil
}

func (g *Guild) Setup() error {
	g.logger.Info("setting up guild")
	if g.commands != nil {
		return g.commands.Register(g.bot.discord)
	}

	return nil
}

func (g *Guild) Run() error {
	for key, module := range g.modules {
		g.logger.Debug("running discord guild module", zap.String("dc_module", key))
		if err := module.Run(); err != nil {
			return err
		}
	}

	g.logger.Info("completed sync run")

	return nil
}

func (g *Guild) Stop() error {
	if g.commands != nil {
		return g.commands.Unregister(g.bot.discord)
	}

	return nil
}
