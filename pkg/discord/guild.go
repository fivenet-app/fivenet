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
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/modules"
	discordtypes "github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	minimumSyncWaitTime       = 30 * time.Second
	discordLogsEmbedChunkSize = 10
	maxErrs                   = 12
)

var ErrSyncCooldownTime = errors.New("guild still in sync cooldown time")

type Guild struct {
	initiated atomic.Bool
	running   atomic.Bool

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	job      string
	gid      discord.GuildID
	lastSync time.Time

	logger *zap.Logger
	bot    *Bot
	guild  discord.Guild

	base    *modules.BaseModule
	modules []modules.Module

	settings *atomic.Pointer[jobs.DiscordSyncSettings]
	events   *broker.Broker[any]
}

func NewGuild(
	ctx context.Context,
	b *Bot,
	guild discord.Guild,
	job string,
	lastSync time.Time,
) (*Guild, error) {
	ctx, cancel := context.WithCancel(ctx)

	events := broker.New[any]()
	go events.Start(ctx)

	g := &Guild{
		initiated: atomic.Bool{},
		running:   atomic.Bool{},

		mu:     sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,

		job: job,
		gid: guild.ID,

		logger: b.logger.Named("guild").
			With(zap.String("job", job), zap.Uint64("discord_guild_id", uint64(guild.ID))),
		bot:     b,
		guild:   guild,
		modules: []modules.Module{},

		settings: &atomic.Pointer[jobs.DiscordSyncSettings]{},
		events:   events,
	}

	settings, _, err := g.getSyncSettings(ctx)
	if err != nil {
		return nil, err
	}
	g.settings.Store(settings)

	g.base = modules.NewBaseModule(ctx, g.logger.Named("module"),
		g.bot.db, g.bot.dc, g.guild, g.job, g.bot.cfg, g.bot.appCfg, g.bot.enricher,
		settings,
	)

	ms := []string{}
	ms = append(ms, "qualifications")
	if b.cfg.GroupSync.Enabled {
		ms = append(ms, "groupsync")
	}
	if b.cfg.UserInfoSync.Enabled {
		ms = append(ms, "userinfo")
	}

	g.logger.Debug("getting discord guild modules", zap.Strings("dc_modules", ms))
	errs := multierr.Combine()
	for _, module := range ms {
		m, err := modules.GetModule(module, g.base, g.events)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("%s. %w", module, err))
			continue
		}

		g.modules = append(g.modules, m)
	}

	return g, errs
}

func (g *Guild) warmupStore() error {
	members, err := g.bot.dc.AllMembers(g.gid)
	if err != nil {
		return fmt.Errorf("failed to get guild members. %w", err)
	}
	for _, member := range members {
		if err := g.bot.dc.MemberSet(g.gid, &member, false); err != nil {
			return err
		}
	}

	return nil
}

func (g *Guild) IsRunning() bool {
	return g.running.Load()
}

func (g *Guild) Run(ignoreCooldown bool) error {
	if !g.initiated.Load() {
		if err := g.warmupStore(); err != nil {
			g.logger.Warn("error during store warmup", zap.Error(err))
		}

		g.initiated.Store(true)
	}

	// If the sync is already running, return
	if g.running.Load() {
		return nil
	}

	// Acquire lock and set running "flag"
	g.mu.Lock()
	defer g.mu.Unlock()
	g.running.Store(true)
	defer g.running.Store(false)

	// Make sure that the minimum wait time is over
	if !ignoreCooldown && time.Since(g.lastSync) <= minimumSyncWaitTime {
		return ErrSyncCooldownTime
	}

	start := time.Now()

	// Get sync settings on run start
	settings, planDiff, err := g.getSyncSettings(g.ctx)
	if err != nil {
		return errors.New("failed to get guild sync settings")
	}
	g.base.SetSettings(settings)
	if planDiff == nil {
		planDiff = &jobs.DiscordSyncChanges{}
	}

	if _, err := g.bot.dc.MembersAfter(g.guild.ID, 0, 0); err != nil {
		g.logger.Error("failed to request all guild members", zap.Error(err))
	}

	errs := multierr.Combine()
	channelId := discord.NullChannelID
	if settings.IsStatusLogEnabled() {
		chId, err := strconv.ParseUint(settings.GetStatusLogSettings().GetChannelId(), 10, 64)
		if err != nil {
			errs = multierr.Append(
				errs,
				fmt.Errorf("failed to parse status log channel id. %w", err),
			)
		}

		channelId = discord.ChannelID(chId)
		if channelId != discord.NullChannelID {
			if err := g.sendStartStatusLog(channelId); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
	}

	// Run modules
	state := &discordtypes.State{
		GuildID: g.guild.ID,
		Users:   discordtypes.Users{},
	}
	logs := []discord.Embed{}
	for _, module := range g.modules {
		g.logger.Debug("running discord guild module", zap.String("dc_module", module.GetName()))

		s, mLogs, err := module.Plan(g.ctx)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("%s. %w", module, err))
			continue
		}

		state.Merge(s)
		logs = append(logs, mLogs...)
	}

	// Allow config to force dry run mode
	plan, ls, err := state.Calculate(g.ctx, g.bot.dc, g.bot.cfg.DryRun || settings.GetDryRun())
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("error during plan calculation. %w", err))
		return errs
	}
	logs = append(logs, ls...)

	// Encode plan as yaml for our "change list"
	b := bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(plan); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to encode plan to yaml for diff. %w", err))
	}
	planDiff.Add(&jobs.DiscordSyncChange{
		Time: timestamp.Now(),
		Plan: b.String(),
	})

	if !plan.DryRun {
		pLogs, err := plan.Apply(g.ctx, g.bot.dc)
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
	channel, err := g.bot.dc.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	if _, err := g.bot.dc.SendEmbeds(channel.ID, discord.Embed{
		Type:        discord.NormalEmbed,
		Title:       "Starting sync...",
		Description: fmt.Sprintf("Dry run: %t", g.settings.Load().GetDryRun() || g.bot.cfg.DryRun),
		Author:      embeds.EmbedAuthor,
		Color:       embeds.ColorInfo,
		Footer:      embeds.EmbedFooterVersion,
	}); err != nil {
		return fmt.Errorf("failed to send status log start embed. %w", err)
	}

	return nil
}

func (g *Guild) sendStatusLog(channelId discord.ChannelID, logs []discord.Embed) error {
	if len(logs) == 0 {
		return nil
	}

	channel, err := g.bot.dc.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	// Split log embeds into chunks
	for i := 0; i < len(logs); i += discordLogsEmbedChunkSize {
		end := min(i+discordLogsEmbedChunkSize, len(logs))

		if _, err := g.bot.dc.SendEmbeds(channel.ID, logs[i:end]...); err != nil {
			return fmt.Errorf("failed to send status log embeds. %w", err)
		}
	}

	return nil
}

func (g *Guild) sendEndStatusLog(
	channelId discord.ChannelID,
	duration time.Duration,
	errs error,
) error {
	channel, err := g.bot.dc.Channel(channelId)
	if err != nil {
		return fmt.Errorf("failed to get status log channel. %w", err)
	}

	logs := []discord.Embed{}
	if errs != nil {
		// Collect up to maxErrs error messages as fields in the embed
		var fields []discord.EmbedField
		errList := multierr.Errors(errs)
		totalErrs := len(errList)
		for i, err := range errList {
			if i >= maxErrs {
				break
			}

			fields = append(fields, discord.EmbedField{
				Name:   fmt.Sprintf("Error %d", i+1),
				Value:  err.Error(),
				Inline: false,
			})
		}

		logs = append(logs, discord.Embed{
			Title: "Errors during sync",
			Description: fmt.Sprintf(
				"Following errors occurred during sync (%d total):",
				totalErrs,
			),
			Author: embeds.EmbedAuthor,
			Color:  embeds.ColorError,
			Fields: fields,
			Footer: embeds.EmbedFooterVersion,
		})
	}

	logs = append(logs, discord.Embed{
		Title:       "Sync completed!",
		Description: fmt.Sprintf("Completed in %s.", duration),
		Author:      embeds.EmbedAuthor,
		Color:       embeds.ColorSuccess,
		Footer:      embeds.EmbedFooterVersion,
	})

	if _, err := g.bot.dc.SendEmbeds(channel.ID, logs...); err != nil {
		return fmt.Errorf("failed to send status log completed embeds. %w", err)
	}

	return nil
}

func (g *Guild) getSyncSettings(
	ctx context.Context,
) (*jobs.DiscordSyncSettings, *jobs.DiscordSyncChanges, error) {
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

	var dest jobs.JobProps
	if err := stmt.QueryContext(ctx, g.bot.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	// Make sure the defaults are set
	dest.Default(g.job)

	g.settings.Store(dest.GetDiscordSyncSettings())

	return dest.GetDiscordSyncSettings(), dest.GetDiscordSyncChanges(), nil
}

func (g *Guild) setLastSyncInterval(
	ctx context.Context,
	job string,
	pDiff *jobs.DiscordSyncChanges,
) error {
	t := time.Now()

	tJobProps := table.FivenetJobProps
	stmt := tJobProps.
		UPDATE(
			tJobProps.DiscordLastSync,
			tJobProps.DiscordSyncChanges,
		).
		SET(
			t,
			pDiff,
		).
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		)

	if _, err := stmt.ExecContext(ctx, g.bot.db); err != nil {
		return fmt.Errorf("failed to update job last sync data. %w", err)
	}

	g.lastSync = t

	return nil
}
