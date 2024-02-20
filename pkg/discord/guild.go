package discord

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/discord/embeds"
	"github.com/galexrt/fivenet/pkg/discord/modules"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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
	g.mutex.Lock()
	defer g.mutex.Unlock()

	start := time.Now()

	if !g.ready.Load() {
		g.logger.Debug("discord guild is not ready yet, running setup")

		if err := g.setup(); err != nil {
			return err
		}
	}

	if g.guild.Unavailable {
		g.logger.Warn("discord guild is unavailable, skipping sync run")
		return nil
	}

	settings, err := g.getSyncSettings(g.bot.ctx, g.Job)
	if err != nil {
		return err
	}

	errs := multierr.Combine()
	if settings.IsStatusLogEnabled() {
		if err := g.sendStartStatusLog(*settings.StatusLogSettings.ChannelId); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	logs := []*discordgo.MessageEmbed{}
	for key, module := range g.modules {
		g.logger.Debug("running discord guild module", zap.String("dc_module", key))

		moduleLogs, err := module.Run(settings)
		if err != nil {
			errs = multierr.Append(errs, err)
		}

		logs = append(logs, moduleLogs...)
	}

	if settings.IsStatusLogEnabled() {
		if err := g.sendStatusLog(*settings.StatusLogSettings.ChannelId, logs, time.Since(start), errs); err != nil {
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

func (g *Guild) sendStartStatusLog(channelId string) error {
	channel, err := g.bot.discord.Channel(channelId)
	if err != nil {
		return err
	}

	if _, err := g.bot.discord.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  "Starting sync...",
		Author: embeds.EmbedAuthor,
		Color:  embeds.ColorInfo,
	}); err != nil {
		return err
	}

	return nil
}

func (g *Guild) sendStatusLog(channelId string, logs []*discordgo.MessageEmbed, duration time.Duration, errs error) error {
	channel, err := g.bot.discord.Channel(channelId)
	if err != nil {
		return err
	}

	if errs != nil {
		logs = append(logs, &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Errors during sync",
			Description: fmt.Sprintf("Following errors occured during sync:\n```\n%v\n```", errs),
			Author:      embeds.EmbedAuthor,
			Color:       embeds.ColorFailure,
		})
	}

	logs = append(logs, &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Sync completed!",
		Description: fmt.Sprintf("Completed in %s.", duration),
		Author:      embeds.EmbedAuthor,
		Color:       embeds.ColorSuccess,
	})

	if _, err := g.bot.discord.ChannelMessageSendEmbeds(channel.ID, logs); err != nil {
		return err
	}

	return nil
}

func (g *Guild) getSyncSettings(ctx context.Context, job string) (*users.DiscordSyncSettings, error) {
	stmt := tJobProps.
		SELECT(
			tJobProps.DiscordSyncSettings,
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		).
		LIMIT(1)

	var dest users.JobProps
	if err := stmt.QueryContext(ctx, g.bot.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	// Make sure the defaults are set
	dest.Default(job)

	return dest.DiscordSyncSettings, nil
}
