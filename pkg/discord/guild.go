package discord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/discord/modules"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const discordLogsEmbedChunkSize = 5

type Guild struct {
	mutex  sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	job string
	id  string

	logger *zap.Logger
	bot    *Bot
	guild  *discordgo.Guild

	lastVersion int
	modules     []string
}

func NewGuild(ctx context.Context, b *Bot, guild *discordgo.Guild, job string) (*Guild, error) {
	ctx, cancel := context.WithCancel(ctx)

	modules := []string{}
	if b.cfg.GroupSync.Enabled {
		modules = append(modules, "groupsync")
	}
	if b.cfg.UserInfoSync.Enabled {
		modules = append(modules, "userinfo")
	}

	modules = append(modules, "qualifications")

	return &Guild{
		mutex:  sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,

		job: job,
		id:  guild.ID,

		logger:  b.logger.Named("guild").With(zap.String("job", job), zap.String("discord_guild_id", guild.ID)),
		bot:     b,
		guild:   guild,
		modules: modules,
	}, nil
}

func (g *Guild) Run() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	start := time.Now()

	if g.guild.Unavailable {
		g.logger.Warn("discord guild is unavailable, skipping sync run")
		return nil
	}

	if g.lastVersion == g.bot.discord.State.Version {
		g.logger.Warn("discord state version is same", zap.Int("discord_state_last_version", g.lastVersion), zap.Int("discord_state_version", g.bot.discord.State.Version))
	}
	g.lastVersion = g.bot.discord.State.Version

	settings, planDiff, err := g.getSyncSettings(g.ctx, g.job)
	if err != nil {
		return fmt.Errorf("failed to get guild sync settings")
	}
	if planDiff == nil {
		planDiff = &users.DiscordSyncChanges{}
	}

	if err := g.bot.discord.RequestGuildMembers(g.guild.ID, "", 0, "", false); err != nil {
		g.logger.Error("failed to request guild members. %w", zap.Error(err))
	}

	errs := multierr.Combine()
	if settings.IsStatusLogEnabled() {
		if err := g.sendStartStatusLog(settings.StatusLogSettings.ChannelId); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	base := modules.NewBaseModule(g.logger.Named("module").With(zap.String("job", g.job), zap.String("discord_guild_id", g.guild.ID)),
		g.bot.db, g.bot.discord, g.guild, g.job, g.bot.cfg, g.bot.appCfg, g.bot.enricher, settings)

	state := &types.State{
		GuildID: g.guild.ID,
		Users:   types.Users{},
	}
	logs := []*discordgo.MessageEmbed{}
	for _, module := range g.modules {
		g.logger.Debug("running discord guild module", zap.String("dc_module", module))

		m, err := modules.GetModule(module, base)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("%s: %w", module, err))
			continue
		}

		s, mLogs, err := m.Plan(g.ctx)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("%s: %w", module, err))
			continue
		}

		state.Merge(s)
		logs = append(logs, mLogs...)
	}

	plan, ls, err := state.Calculate(g.ctx, g.bot.discord, settings.DryRun)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("error during plan calculation. %w", err))
		return errs
	}
	logs = append(logs, ls...)
	plan.DryRun = settings.DryRun

	// Encode plan as yaml for our "change list"
	b := bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(plan); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to encode plan to yaml for diff. %w", err))
	}
	planDiff.Add(&users.DiscordSyncChange{
		Time: timestamp.Now(),
		Plan: b.String(),
	})

	if !plan.DryRun {
		pLogs, err := plan.Apply(g.ctx, g.bot.discord)
		logs = append(logs, pLogs...)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("error during plan apply. %w", err))
		}
	}

	if settings.IsStatusLogEnabled() {
		if err := g.sendStatusLog(settings.StatusLogSettings.ChannelId, logs); err != nil {
			errs = multierr.Append(errs, err)
		}

		if err := g.sendEndStatusLog(settings.StatusLogSettings.ChannelId, time.Since(start), errs); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if err := g.setLastSyncInterval(g.ctx, g.job, planDiff); err != nil {
		g.logger.Error("error setting job props last sync time", zap.Error(err))
		errs = multierr.Append(errs, err)
	}

	g.logger.Info("completed sync run", zap.Duration("duration", time.Since(start)))

	return errs
}

func (g *Guild) Stop() {
	g.cancel()
}

func (g *Guild) sendStartStatusLog(channelId string) error {
	channel, err := g.bot.discord.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	if _, err := g.bot.discord.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  "Starting sync...",
		Author: embeds.EmbedAuthor,
		Color:  embeds.ColorInfo,
		Footer: embeds.EmbedFooterVersion,
	}); err != nil {
		return fmt.Errorf("failed to send status log start embed. %w", err)
	}

	return nil
}

func (g *Guild) sendStatusLog(channelId string, logs []*discordgo.MessageEmbed) error {
	if len(logs) == 0 {
		return nil
	}

	channel, err := g.bot.discord.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	// Split logs embeds into chunks
	for i := 0; i < len(logs); i += discordLogsEmbedChunkSize {
		end := i + discordLogsEmbedChunkSize

		if end > len(logs) {
			end = len(logs)
		}

		if _, err := g.bot.discord.ChannelMessageSendEmbeds(channel.ID, logs[i:end]); err != nil {
			return fmt.Errorf("failed to send status log embeds. %w", err)
		}
	}

	return nil
}

func (g *Guild) sendEndStatusLog(channelId string, duration time.Duration, errs error) error {
	channel, err := g.bot.discord.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	logs := []*discordgo.MessageEmbed{}
	if errs != nil {
		logs = append(logs, &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Errors during sync",
			Description: fmt.Sprintf("Following errors occured during sync:\n```\n%v\n```", errs),
			Author:      embeds.EmbedAuthor,
			Color:       embeds.ColorError,
		})
	}

	logs = append(logs, &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Sync completed!",
		Description: fmt.Sprintf("Completed in %s.", duration),
		Author:      embeds.EmbedAuthor,
		Color:       embeds.ColorSuccess,
		Footer:      embeds.EmbedFooterVersion,
	})

	if _, err := g.bot.discord.ChannelMessageSendEmbeds(channel.ID, logs); err != nil {
		return fmt.Errorf("failed to send status log completed embeds. %w", err)
	}

	return nil
}

func (g *Guild) getSyncSettings(ctx context.Context, job string) (*users.DiscordSyncSettings, *users.DiscordSyncChanges, error) {
	stmt := tJobProps.
		SELECT(
			tJobProps.DiscordSyncSettings,
			tJobProps.DiscordSyncChanges,
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		).
		LIMIT(1)

	var dest users.JobProps
	if err := stmt.QueryContext(ctx, g.bot.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	// Make sure the defaults are set
	dest.Default(job)

	return dest.DiscordSyncSettings, dest.DiscordSyncChanges, nil
}

func (g *Guild) setLastSyncInterval(ctx context.Context, job string, pDiff *users.DiscordSyncChanges) error {
	tJobProps := table.FivenetJobProps
	stmt := tJobProps.
		UPDATE(
			tJobProps.DiscordLastSync,
			tJobProps.DiscordSyncChanges,
		).
		SET(
			jet.CURRENT_TIMESTAMP(),
			pDiff,
		).
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		)

	if _, err := stmt.ExecContext(ctx, g.bot.db); err != nil {
		return fmt.Errorf("failed to update job last sync data. %w", err)
	}

	return nil
}
