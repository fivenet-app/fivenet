package discord

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/discord/modules"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Guild struct {
	Job string `alias:"job"`
	ID  string `alias:"id"`

	mutex sync.Mutex
	ready atomic.Bool

	logger  *zap.Logger
	bot     *Bot
	guild   *discordgo.Guild
	modules map[string]modules.Module
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

	return &Guild{
		Job: job,
		ID:  guild.ID,

		mutex: sync.Mutex{},
		ready: atomic.Bool{},

		logger:  b.logger.Named("guild").With(zap.String("job", job), zap.String("discord_guild_id", guild.ID)),
		bot:     b,
		guild:   guild,
		modules: ms,
	}, nil
}

func (g *Guild) setup() error {
	g.logger.Info("setting up guild")

	if _, err := g.bot.discord.Guild(g.guild.ID); err != nil {
		return fmt.Errorf("failed to retrieve guild info from discord api. %w", err)
	}

	// Make sure that the guild roles are cached in state
	if len(g.guild.Roles) == 0 {
		if _, err := g.bot.discord.GuildRoles(g.guild.ID); err != nil {
			return fmt.Errorf("failed to retrieve roles for guild during setup. %w", err)
		}
	}

	g.ready.Store(true)

	return nil
}

func (g *Guild) Run() error {
	if !g.ready.Load() {
		g.logger.Debug("discord guild is not ready yet, running setup")

		if err := g.setup(); err != nil {
			return err
		}
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.guild.Unavailable {
		g.logger.Warn("discord guild is unavailable, skipping sync run")
		return nil
	}

	errs := multierr.Combine()
	for key, module := range g.modules {
		g.logger.Debug("running discord guild module", zap.String("dc_module", key))
		if err := module.Run(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	g.logger.Info("completed sync run")

	return errs
}

func (g *Guild) Stop() error {
	g.ready.Store(false)
	return nil
}
