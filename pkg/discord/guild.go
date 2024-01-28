package discord

import (
	"fmt"

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

		logger:   b.logger.Named("guild").With(zap.String("job", job), zap.String("discord_guild_id", guild.ID)),
		bot:      b,
		guild:    guild,
		modules:  ms,
		commands: cmds,
	}, nil
}

func (g *Guild) Setup() error {
	g.logger.Info("setting up guild")

	if _, err := g.bot.discord.Guild(g.guild.ID); err != nil {
		return fmt.Errorf("failed to retrieve guild. %w", err)
	}

	// Make sure that the guild roles are cached in state
	if len(g.guild.Roles) == 0 {
		if _, err := g.bot.discord.GuildRoles(g.guild.ID); err != nil {
			return fmt.Errorf("failed to retrieve roles for guild. %w", err)
		}
	}

	if g.commands != nil {
		if err := g.commands.Register(g.bot.discord); err != nil {
			return fmt.Errorf("failed to register bot commands. %w", err)
		}
	}

	return nil
}

func (g *Guild) Run() error {
	if g.guild.Unavailable {
		g.logger.Warn("discord guild is unavailable, skipping sync run")
		return nil
	}

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
	return nil
}
