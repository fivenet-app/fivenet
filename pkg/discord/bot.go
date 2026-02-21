package discord

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const botWorkerCount = 3

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("discord.bot")
}

var BotModule = fx.Module("discord.bot",
	fx.Provide(
		New,
	),
	fx.Decorate(wrapLogger),
)

var (
	metricLastSync = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "discord_bot",
		Name:      "last_sync",
		Help:      "Last time sync has completed.",
	}, []string{"job_name", "status"})

	metricGuildsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "discord_bot",
		Name:      "guilds_total_count",
		Help:      "Total count of Discord guilds being ready.",
	})

	metricSyncDuration = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "discord_bot",
		Name:      "sync_duration_seconds",
		Help:      "Duration of the last sync operation in seconds.",
	}, []string{"job_name"})
)

type BotParams struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	JS        *events.JSWrapper
	DB        *sql.DB
	Enricher  *mstlystcdata.Enricher
	Config    *config.Config
	AppConfig appconfig.IConfig
	Perms     perms.Permissions

	Discord *state.State
}

type Bot struct {
	discordtypes.BotState

	logger   *zap.Logger
	tracer   trace.Tracer
	js       *events.JSWrapper
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	cfg      *config.Discord
	appCfg   appconfig.IConfig
	perms    perms.Permissions

	oauth2ProviderName string

	wg     sync.WaitGroup
	workCh chan *Guild

	syncTimer *time.Timer
	syncTime  atomic.Pointer[time.Duration]

	enabled      bool
	dc           *state.State
	activeGuilds *xsync.Map[discord.GuildID, *Guild]
}

type Result struct {
	fx.Out

	Bot      *Bot
	BotState discordtypes.BotState
}

func New(p BotParams) Result {
	cancelCtx, cancel := context.WithCancel(context.Background())

	oauth2ProviderName := "discord"
	if provider := p.Config.OAuth2.GetProviderByType(config.OAuth2ProviderDiscord); provider != nil {
		oauth2ProviderName = provider.Name
	}

	b := &Bot{
		logger:   p.Logger,
		tracer:   p.TP.Tracer("discord.bot"),
		js:       p.JS,
		db:       p.DB,
		enricher: p.Enricher,
		cfg:      &p.Config.Discord,
		appCfg:   p.AppConfig,
		perms:    p.Perms,

		oauth2ProviderName: oauth2ProviderName,

		wg:     sync.WaitGroup{},
		workCh: make(chan *Guild, 3),

		enabled:      p.Config.Discord.Enabled && p.Config.Discord.Sync,
		dc:           p.Discord,
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		// Discord bot or sync not enabled
		if !b.enabled {
			return nil
		}

		// Setup sync timer
		syncInterval := b.appCfg.Get().Discord.GetSyncInterval().AsDuration()
		b.syncTime.Store(&syncInterval)
		b.syncTimer = time.NewTimer(syncInterval)

		// Start bot workers
		for range botWorkerCount {
			b.wg.Go(func() {
				for {
					select {
					case <-cancelCtx.Done():
						return

					case guild := <-b.workCh:
						var elapsed time.Duration

						func() {
							logger := b.logger.With(
								zap.String("job", guild.job),
								zap.Uint64("discord_guild_id", uint64(guild.gid)),
							)

							start := time.Now()
							defer func() {
								elapsed = time.Since(start)
								// Recover from a panic and set err accordingly
								if e := recover(); e != nil {
									var err error
									if er, ok := e.(error); ok {
										err = fmt.Errorf("recovered from panic. %w", er)
									} else {
										//nolint:errorlint // `er` is not guaranteed to be an error type, so we want it to be treated as a "string" here.
										err = fmt.Errorf("recovered from panic. %v", er)
									}

									logger.Error(
										"discord guild sync panic",
										zap.Error(err),
										zap.StackSkip("stacktrace", 2),
									)
								}
							}()

							// Ignore the cooldown for the periodic sync
							if err := guild.Run(true); err != nil {
								logger.Error("error during discord sync", zap.Error(err))

								metricLastSync.WithLabelValues(guild.job, "failed").
									SetToCurrentTime()
							} else {
								metricLastSync.WithLabelValues(guild.job, "success").SetToCurrentTime()
							}
						}()

						metricSyncDuration.With(prometheus.Labels{"job_name": guild.job}).
							Set(elapsed.Seconds())
					}
				}
			})
		}

		if err := registerStreams(ctxStartup, b.js); err != nil {
			return err
		}

		if err := b.start(cancelCtx); err != nil {
			return err
		}

		// Handle app config updates
		go func() {
			configUpdateCh := b.appCfg.Subscribe()
			for {
				select {
				case <-cancelCtx.Done():
					b.appCfg.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					if cfg == nil {
						continue
					}
					b.handleAppConfigUpdate(cancelCtx, cfg)
				}
			}
		}()

		go b.syncLoop(cancelCtx)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		// Stop all guilds and discord session
		if err := b.stop(); err != nil {
			return err
		}

		cancel()

		b.wg.Wait()

		return nil
	}))

	return Result{
		Bot:      b,
		BotState: b,
	}
}

func (b *Bot) handleAppConfigUpdate(ctx context.Context, cfg *appconfig.Cfg) {
	b.setBotPresence(ctx, cfg.Discord.GetBotPresence())

	// Only reset sync timer when interval has changed
	currentSyncTime := b.syncTime.Load()
	if currentSyncTime == nil || *currentSyncTime != cfg.Discord.GetSyncInterval().AsDuration() {
		newSyncTime := cfg.Discord.GetSyncInterval().AsDuration()
		b.syncTime.Store(&newSyncTime)
		b.syncTimer.Reset(newSyncTime)
	}
}

func (b *Bot) start(ctx context.Context) error {
	var ready atomic.Bool

	b.dc.AddHandler(func(ev *gateway.ReadyEvent) {
		b.logger.Info(
			fmt.Sprintf("connected to gateway, ready with %d guilds", len(ev.Guilds)),
			zap.String("me", ev.User.Tag()),
		)
		ready.Store(true)
	})

	b.dc.AddHandler(func(ev *gateway.GuildCreateEvent) {
		b.logger.Info("discord server joined", zap.Uint64("discord_guild_id", uint64(ev.ID)))
	})

	b.dc.AddHandler(func(ev *gateway.GuildMemberAddEvent) {
		g, ok := b.activeGuilds.Load(ev.GuildID)
		if !ok {
			return
		}

		g.events.Publish(ev)
	})

	for {
		if b.dc.Ready().Version > 0 && ready.Load() {
			if _, err := b.dc.Me(); err != nil {
				return fmt.Errorf("failed to obtain bot account details. %w", err)
			}

			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"discord client failed to get ready in time, version %d",
				b.dc.Ready().Version,
			)

		case <-time.After(750 * time.Millisecond):
		}
	}

	b.handleAppConfigUpdate(ctx, b.appCfg.Get())

	return nil
}

func (b *Bot) syncLoop(ctx context.Context) {
	for {
		b.logger.Info("running discord sync", zap.Bool("dry_run", b.cfg.DryRun))
		func() {
			ctx, span := b.tracer.Start(ctx, "discord_bot")
			defer span.End()

			if err := b.runSync(ctx); err != nil {
				b.logger.Error("failed to sync discord", zap.Error(err))
			}
		}()

		if syncTime := b.syncTime.Load(); syncTime != nil {
			b.syncTimer.Reset(*syncTime)
		} else {
			// Fallback to sane value
			b.syncTimer.Reset(10 * time.Minute)
		}

		select {
		case <-ctx.Done():
			return

		case <-b.syncTimer.C:
		}
	}
}

// getGuilds Each guild is effectively associated with a Job via the JobProps.
func (b *Bot) getGuilds(ctx context.Context) error {
	jobGuilds, err := b.getJobGuildsFromDB(ctx)
	if err != nil {
		return err
	}

	if len(jobGuilds) == 0 {
		b.logger.Debug("no job discord guild connections found")
		return nil
	}

	guilds, err := b.dc.Guilds()
	if err != nil {
		return fmt.Errorf("failed to get guilds from dc state. %w", err)
	}

	for _, guildInfo := range jobGuilds {
		idx := slices.IndexFunc(guilds, func(in discord.Guild) bool {
			return in.ID == guildInfo.GuildID
		})
		if idx == -1 {
			// Make sure to stop any active stuff with the previously active guild
			if g, ok := b.activeGuilds.Load(guildInfo.GuildID); ok {
				g.Stop()

				b.activeGuilds.Delete(guildInfo.GuildID)
			}

			b.logger.Warn(
				"didn't find bot in guild (anymore?)",
				zap.Uint64("discord_guild_id", uint64(guildInfo.GuildID)),
				zap.String("job", guildInfo.Job),
			)
			continue
		}

		// Check if the guild is already existing and therefore active
		if _, ok := b.activeGuilds.Load(guildInfo.GuildID); ok {
			continue
		}

		g, err := NewGuild(
			ctx,
			b,
			guilds[idx],
			guildInfo.Job,
			guildInfo.LastSync.AsTime(),
			b.oauth2ProviderName,
		)
		if err != nil {
			return err
		}
		b.activeGuilds.Store(guildInfo.GuildID, g)
	}

	return nil
}

type jobGuild struct {
	Job      string               `alias:"job"       sql:"primary_key"`
	GuildID  discord.GuildID      `alias:"guild_id"`
	LastSync *timestamp.Timestamp `alias:"last_sync"`
}

func (b *Bot) getJobGuildsFromDB(ctx context.Context) ([]*jobGuild, error) {
	tJobProps := table.FivenetJobProps.AS("job_props")

	stmt := tJobProps.
		SELECT(
			tJobProps.Job.AS("jobguild.job"),
			tJobProps.DiscordGuildID.AS("jobguild.guild_id"),
			tJobProps.DiscordLastSync.AS("jobguild.last_sync"),
		).
		FROM(tJobProps).
		WHERE(mysql.AND(
			tJobProps.DiscordGuildID.IS_NOT_NULL(),
		))

	var guilds []*jobGuild
	if err := stmt.QueryContext(ctx, b.db, &guilds); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return guilds, nil
}

func (b *Bot) runSync(ctx context.Context) error {
	// Discord bot or sync not enabled
	if !b.enabled {
		b.logger.Warn("skipping discord sync since bot or sync is not enabled")
		return nil
	}

	if err := b.getGuilds(ctx); err != nil {
		return fmt.Errorf("failed to get guilds. %w", err)
	}

	totalCount := float64(0)

	metricGuildsTotal.Set(totalCount)

	errs := multierr.Combine()

	// Submit guilds to sync via work channel
	for _, guild := range b.activeGuilds.All() {
		b.workCh <- guild
	}

	return errs
}

func (b *Bot) stop() error {
	errs := multierr.Combine()
	for _, guild := range b.activeGuilds.All() {
		guild.Stop()
	}

	b.activeGuilds.Clear()

	if errs != nil {
		return errs
	}

	return b.dc.Close()
}

// State helpers for commands and modules

func (b *Bot) GetJobFromGuildID(guildId discord.GuildID) (string, bool) {
	guild, ok := b.activeGuilds.Load(guildId)
	if !ok || guild == nil {
		return "", false
	}

	return guild.job, true
}

func (b *Bot) RunSync(guildID discord.GuildID) (bool, error) {
	// Submit guild to sync queue via work channel
	guild, ok := b.activeGuilds.Load(guildID)
	if !ok {
		return false, errors.New("no active guild found for guild ID")
	}

	b.workCh <- guild

	return false, nil
}

func (b *Bot) IsUserGuildAdmin(
	ctx context.Context,
	channelId discord.ChannelID,
	userId discord.UserID,
) (bool, error) {
	perms, err := b.dc.Permissions(channelId, userId)
	if err != nil {
		return false, err
	}

	return perms.Has(discord.PermissionAdministrator), nil
}
