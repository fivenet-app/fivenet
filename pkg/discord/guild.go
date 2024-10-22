package discord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/discord/modules"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
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
	id  discord.GuildID

	logger *zap.Logger
	bot    *Bot
	guild  discord.Guild

	base    *modules.BaseModule
	modules []modules.Module

	settings *atomic.Pointer[users.DiscordSyncSettings]
	events   *broker.Broker[interface{}]
}

func NewGuild(c context.Context, b *Bot, guild discord.Guild, job string) (*Guild, error) {
	ctx, cancel := context.WithCancel(c)

	events := broker.New[interface{}]()
	go events.Start(ctx)

	g := &Guild{
		mutex:  sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,

		job: job,
		id:  guild.ID,

		logger:  b.logger.Named("guild").With(zap.String("job", job), zap.Uint64("discord_guild_id", uint64(guild.ID))),
		bot:     b,
		guild:   guild,
		modules: []modules.Module{},

		settings: &atomic.Pointer[users.DiscordSyncSettings]{},
		events:   events,
	}

	settings, _, err := g.getSyncSettings(c)
	if err != nil {
		return nil, err
	}
	g.settings.Store(settings)

	g.base = modules.NewBaseModule(ctx, g.logger.Named("module"),
		g.bot.db, g.bot.discord, g.guild, g.job, g.bot.cfg, g.bot.appCfg, g.bot.enricher,
		settings,
	)

	ms := []string{"qualifications"}
	if b.cfg.GroupSync.Enabled {
		ms = append(ms, "groupsync")
	}
	if b.cfg.UserInfoSync.Enabled {
		ms = append(ms, "userinfo")
	}

	errs := multierr.Combine()
	for _, module := range ms {
		g.logger.Debug("getting discord guild module", zap.String("dc_module", module))

		m, err := modules.GetModule(module, g.base, g.events)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("%s: %w", module, err))
			continue
		}

		g.modules = append(g.modules, m)
	}

	return g, errs
}

func (g *Guild) Run() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	start := time.Now()

	// Get sync settings on run start
	settings, planDiff, err := g.getSyncSettings(g.ctx)
	if err != nil {
		return fmt.Errorf("failed to get guild sync settings")
	}
	g.base.SetSettings(settings)
	if planDiff == nil {
		planDiff = &users.DiscordSyncChanges{}
	}

	if _, err := g.bot.discord.Members(g.guild.ID); err != nil {
		g.logger.Error("failed to request guild members. %w", zap.Error(err))
	}

	errs := multierr.Combine()
	channelId := discord.NullChannelID
	if settings.IsStatusLogEnabled() {
		chId, err := strconv.ParseUint(settings.StatusLogSettings.ChannelId, 10, 64)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to parse status log channel id. %w", err))
		}

		channelId = discord.ChannelID(chId)
		if channelId != discord.NullChannelID {
			if err := g.sendStartStatusLog(discord.ChannelID(channelId)); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
	}

	// Run modules
	state := &types.State{
		GuildID: g.guild.ID,
		Users:   types.Users{},
	}
	logs := []discord.Embed{}
	for _, module := range g.modules {
		g.logger.Debug("running discord guild module", zap.String("dc_module", module.GetName()))

		s, mLogs, err := module.Plan(g.ctx)
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

	if channelId != discord.NullChannelID {
		if err := g.sendStatusLog(channelId, logs); err != nil {
			errs = multierr.Append(errs, err)
		}

		if err := g.sendEndStatusLog(channelId, time.Since(start), errs); err != nil {
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

func (g *Guild) sendStartStatusLog(channelId discord.ChannelID) error {
	channel, err := g.bot.discord.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	if _, err := g.bot.discord.SendEmbeds(channel.ID, discord.Embed{
		Type:   discord.NormalEmbed,
		Title:  "Starting sync...",
		Author: embeds.EmbedAuthor,
		Color:  embeds.ColorInfo,
		Footer: embeds.EmbedFooterVersion,
	}); err != nil {
		return fmt.Errorf("failed to send status log start embed. %w", err)
	}

	return nil
}

func (g *Guild) sendStatusLog(channelId discord.ChannelID, logs []discord.Embed) error {
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

		if _, err := g.bot.discord.SendEmbeds(channel.ID, logs[i:end]...); err != nil {
			return fmt.Errorf("failed to send status log embeds. %w", err)
		}
	}

	return nil
}

func (g *Guild) sendEndStatusLog(channelId discord.ChannelID, duration time.Duration, errs error) error {
	channel, err := g.bot.discord.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	logs := []discord.Embed{}
	if errs != nil {
		logs = append(logs, discord.Embed{
			Title:       "Errors during sync",
			Description: fmt.Sprintf("Following errors occured during sync:\n```\n%v\n```", errs),
			Author:      embeds.EmbedAuthor,
			Color:       embeds.ColorError,
		})
	}

	logs = append(logs, discord.Embed{
		Title:       "Sync completed!",
		Description: fmt.Sprintf("Completed in %s.", duration),
		Author:      embeds.EmbedAuthor,
		Color:       embeds.ColorSuccess,
		Footer:      embeds.EmbedFooterVersion,
	})

	if _, err := g.bot.discord.SendEmbeds(channel.ID, logs...); err != nil {
		return fmt.Errorf("failed to send status log completed embeds. %w", err)
	}

	return nil
}

func (g *Guild) getSyncSettings(ctx context.Context) (*users.DiscordSyncSettings, *users.DiscordSyncChanges, error) {
	stmt := tJobProps.
		SELECT(
			tJobProps.DiscordSyncSettings,
			tJobProps.DiscordSyncChanges,
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(jet.String(g.job)),
		).
		LIMIT(1)

	var dest users.JobProps
	if err := stmt.QueryContext(ctx, g.bot.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	// Make sure the defaults are set
	dest.Default(g.job)

	g.settings.Store(dest.DiscordSyncSettings)

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
