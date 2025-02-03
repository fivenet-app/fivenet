package discord

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/discord/commands"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/admin"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var tJobProps = table.FivenetJobProps.AS("jobprops")

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("discord_bot")
}

var BotModule = fx.Module("discord_bot",
	fx.Provide(
		NewBot,
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
)

type BotParams struct {
	fx.In

	LC         fx.Lifecycle
	Shutdowner fx.Shutdowner

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	JS        *events.JSWrapper
	DB        *sql.DB
	Enricher  *mstlystcdata.Enricher
	Config    *config.Config
	AppConfig appconfig.IConfig
	I18n      *lang.I18n
	Perms     perms.Permissions
}

type Bot struct {
	types.BotState

	ctx      context.Context
	logger   *zap.Logger
	tracer   trace.Tracer
	js       *events.JSWrapper
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	cfg      *config.Discord
	appCfg   appconfig.IConfig
	i18n     *lang.I18n
	perms    perms.Permissions

	cmds *commands.Cmds

	wg sync.WaitGroup

	syncTimer *time.Timer

	dc           *state.State
	activeGuilds *xsync.MapOf[discord.GuildID, *Guild]
}

func NewBot(p BotParams) (*Bot, error) {
	if !p.Config.Discord.Enabled {
		return nil, nil
	}

	// Create a new Discord session using the provided login information.
	state := state.New("Bot " + p.Config.Discord.Token)
	state.AddIntents(gateway.IntentGuildMembers)
	state.AddIntents(gateway.IntentGuildPresences)
	state.AddIntents(gateway.IntentGuildIntegrations)

	cmds, err := commands.New(commands.Params{
		Logger: p.Logger,
		S:      state,
		Cfg:    p.Config,
		I18n:   p.I18n,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating commands for discord bot. %w", err)
	}

	cancelCtx, cancel := context.WithCancel(context.Background())

	b := &Bot{
		ctx:      cancelCtx,
		logger:   p.Logger,
		tracer:   p.TP.Tracer("discord_bot"),
		js:       p.JS,
		db:       p.DB,
		enricher: p.Enricher,
		cfg:      &p.Config.Discord,
		appCfg:   p.AppConfig,
		i18n:     p.I18n,
		perms:    p.Perms,

		cmds: cmds,

		wg: sync.WaitGroup{},

		dc:           state,
		activeGuilds: xsync.NewMapOf[discord.GuildID, *Guild](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerStreams(ctxStartup, b.js); err != nil {
			return err
		}

		go func() {
			if err := state.Connect(cancelCtx); err != nil {
				p.Logger.Error("failed to connect to discord gateway", zap.Error(err))
				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					b.logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
				}
			}
		}()

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
					b.handleAppConfigUpdate(cfg)
				}
			}
		}()

		go b.syncLoop()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		// Stop all guilds and discord session
		if err := b.stop(); err != nil {
			return err
		}

		state.Close()

		cancel()

		return nil
	}))

	return b, nil
}

func (b *Bot) handleAppConfigUpdate(cfg *appconfig.Cfg) {
	b.setBotPresence(cfg.Discord.BotPresence)

	if b.syncTimer != nil {
		b.syncTimer.Reset(cfg.Discord.SyncInterval.AsDuration())
	}
}

func (b *Bot) start(ctx context.Context) error {
	var ready atomic.Bool

	b.dc.AddHandler(func(ev *gateway.ReadyEvent) {
		b.logger.Info(fmt.Sprintf("connected to gateway, ready with %d guilds", len(ev.Guilds)), zap.String("me", ev.User.Tag()))
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
				return fmt.Errorf("failed to obtain bot account details: %w", err)
			}

			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("discord client failed to get ready, version %d", b.dc.Ready().Version)

		case <-time.After(750 * time.Millisecond):
		}
	}

	if b.cfg.Commands.Enabled {
		if err := b.cmds.RegisterCommands(commands.CommandParams{
			DB:       b.db,
			L:        b.i18n,
			BotState: b,
			Perms:    b.perms,
		}); err != nil {
			return fmt.Errorf("failed to register commands. %w", err)
		}
	}

	b.handleAppConfigUpdate(b.appCfg.Get())

	return nil
}

func (b *Bot) syncLoop() {
	for {
		b.logger.Info("running discord sync")
		func() {
			ctx, span := b.tracer.Start(b.ctx, "discord_bot")
			defer span.End()

			if err := b.runSync(ctx); err != nil {
				b.logger.Error("failed to sync discord", zap.Error(err))
			}
		}()

		syncInterval := b.appCfg.Get().Discord.SyncInterval.AsDuration()
		if b.syncTimer == nil {
			b.syncTimer = time.NewTimer(syncInterval)
		} else {
			b.syncTimer.Reset(syncInterval)
		}

		select {
		case <-b.ctx.Done():
			return

		case <-b.syncTimer.C:
		}
	}
}

// getGuilds Each guild is effectively associated with a Job via the JobProps
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

	for job, guildID := range jobGuilds {
		idx := slices.IndexFunc(guilds, func(in discord.Guild) bool {
			return in.ID == guildID
		})
		if idx == -1 {
			// Make sure to stop any active stuff with the previously active guild
			if g, ok := b.activeGuilds.Load(guildID); ok {
				g.Stop()

				b.activeGuilds.Delete(guildID)
			}

			b.logger.Warn("didn't find bot in guild (anymore?)", zap.Uint64("discord_guild_id", uint64(guildID)), zap.String("job", job))
			continue
		}

		if _, ok := b.activeGuilds.Load(guildID); ok {
			continue
		}

		g, err := NewGuild(b.ctx, b, guilds[idx], job)
		if err != nil {
			return err
		}
		b.activeGuilds.Store(guildID, g)
	}

	return nil
}

func (b *Bot) getJobGuildsFromDB(ctx context.Context) (map[string]discord.GuildID, error) {
	stmt := tJobProps.
		SELECT(
			tJobProps.Job.AS("job"),
			tJobProps.DiscordGuildID.AS("id"),
		).
		FROM(tJobProps).
		WHERE(jet.AND(
			tJobProps.DiscordGuildID.IS_NOT_NULL(),
		))

	var dest []*struct {
		Job string `alias:"job"`
		ID  string `alias:"id"`
	}
	if err := stmt.QueryContext(ctx, b.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	guilds := map[string]discord.GuildID{}
	for _, g := range dest {
		id, err := strconv.ParseUint(g.ID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse guild id %s as uint64. %w", g.ID, err)
		}

		guilds[g.Job] = discord.GuildID(id)
	}

	return guilds, nil
}

func (b *Bot) runSync(ctx context.Context) error {
	if err := b.getGuilds(ctx); err != nil {
		return fmt.Errorf("failed to get guilds. %w", err)
	}

	totalCount := float64(0)

	metricGuildsTotal.Set(totalCount)

	errs := multierr.Combine()

	// Run at max 3 syncs at once
	workChannel := make(chan *Guild, 3)

	// Retrieve guilds via channel
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		for guild := range workChannel {
			b.wg.Add(1)
			go func(g *Guild) {
				defer b.wg.Done()
				logger := b.logger.With(zap.String("job", g.job), zap.Uint64("discord_guild_id", uint64(g.id)))

				if err := g.Run(); err != nil {
					logger.Error("error during sync", zap.Error(err))
					errs = multierr.Append(errs, err)

					metricLastSync.WithLabelValues(g.job, "failed").SetToCurrentTime()
				} else {
					metricLastSync.WithLabelValues(g.job, "success").SetToCurrentTime()
				}
			}(guild)
		}
	}()

	b.activeGuilds.Range(func(_ discord.GuildID, guild *Guild) bool {
		workChannel <- guild
		return true
	})

	close(workChannel)

	b.wg.Wait()

	return errs
}

func (b *Bot) stop() error {
	errs := multierr.Combine()
	b.activeGuilds.Range(func(key discord.GuildID, guild *Guild) bool {
		guild.Stop()

		return true
	})

	b.activeGuilds.Clear()

	if errs != nil {
		return errs
	}

	return b.dc.Close()
}

func (b *Bot) GetJobFromGuildID(guildId discord.GuildID) (string, bool) {
	guild, ok := b.activeGuilds.Load(guildId)
	if !ok || guild == nil {
		return "", false
	}

	return guild.job, true
}
