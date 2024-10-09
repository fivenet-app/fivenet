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
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
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
	DB        *sql.DB
	Enricher  *mstlystcdata.Enricher
	Config    *config.Config
	AppConfig appconfig.IConfig
	I18n      *lang.I18n
}

type Bot struct {
	ctx      context.Context
	logger   *zap.Logger
	tracer   trace.Tracer
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	cfg      *config.Discord
	appCfg   appconfig.IConfig
	i18n     *lang.I18n

	cmds *commands.Cmds

	wg sync.WaitGroup

	discord      *state.State
	activeGuilds *xsync.MapOf[discord.GuildID, *Guild]
}

func NewBot(p BotParams) (*Bot, error) {
	if !p.Config.Discord.Enabled {
		return nil, nil
	}

	// Create a new Discord session using the provided login information.
	state := state.New("Bot " + p.Config.Discord.Token)
	state.AddIntents(gateway.IntentGuilds)
	state.AddIntents(gateway.IntentGuildMembers)
	state.AddIntents(gateway.IntentGuildPresences)
	state.AddIntents(gateway.IntentGuildIntegrations)

	state.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := state.Me()
		p.Logger.Info("connected to gateway", zap.String("me", me.Tag()))
	})

	cmds, err := commands.New(p.Logger, state, p.Config, p.I18n)
	if err != nil {
		return nil, fmt.Errorf("error creating commands for discord bot. %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	b := &Bot{
		ctx:      ctx,
		logger:   p.Logger,
		tracer:   p.TP.Tracer("discord_bot"),
		db:       p.DB,
		enricher: p.Enricher,
		cfg:      &p.Config.Discord,
		appCfg:   p.AppConfig,
		i18n:     p.I18n,

		cmds: cmds,

		wg: sync.WaitGroup{},

		discord:      state,
		activeGuilds: xsync.NewMapOf[discord.GuildID, *Guild](),
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go func() {
			if err := state.Connect(ctx); err != nil {
				p.Logger.Error("failed to connect to discord gateway", zap.Error(err))
				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					b.logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
				}
			}
		}()

		if err := b.start(ctx); err != nil {
			return err
		}

		// Handle app config updates
		go func() {
			configUpdateCh := b.appCfg.Subscribe()
			for {
				select {
				case <-ctx.Done():
					b.appCfg.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					b.handleAppConfigUpdate(cfg)
				}
			}
		}()

		go func() {
			b.logger.Info("sleeping 5 seconds before running first discord sync")
			time.Sleep(5 * time.Second)

			if err := b.syncLoop(); err != nil {
				b.logger.Error("error from discord bot sync loop", zap.Error(err))
				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					b.logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
				}
			}
		}()

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
}

func (b *Bot) start(ctx context.Context) error {
	var ready atomic.Bool

	b.discord.AddHandler(func(r *gateway.ReadyEvent) {
		b.logger.Info(fmt.Sprintf("ready with %d guilds", len(r.Guilds)))
		ready.Store(true)
	})

	b.discord.AddHandler(func(g *gateway.GuildCreateEvent) {
		b.logger.Info("discord server joined", zap.Uint64("discord_guild_id", uint64(g.ID)))
	})

	go func() {
		if err := b.discord.Open(ctx); err != nil {
			b.logger.Error("error opening discord connection", zap.Error(err))
		}
	}()

	for {
		if b.discord.Ready().Version > 0 && ready.Load() {
			if _, err := b.discord.Me(); err != nil {
				return fmt.Errorf("failed to obtain bot account details: %w", err)
			}

			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("discord client failed to get ready, version %d", b.discord.Ready().Version)

		case <-time.After(750 * time.Millisecond):
		}
	}

	if b.cfg.Commands.Enabled {
		if err := b.cmds.RegisterCommands(); err != nil {
			return fmt.Errorf("failed to register commands. %w", err)
		}
	}

	b.handleAppConfigUpdate(b.appCfg.Get())

	return nil
}

func (b *Bot) syncLoop() error {
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
		select {
		case <-b.ctx.Done():
			return nil

		case <-time.After(syncInterval):
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

	if _, err := b.discord.Guilds(); err != nil {
		return fmt.Errorf("failed to get guild info from discord. %w", err)
	}

	guilds, err := b.discord.GuildStore.Guilds()
	if err != nil {
		return fmt.Errorf("failed to get guilds from store. %w", err)
	}

	for job, guildID := range jobGuilds {
		idx := slices.IndexFunc(guilds, func(in discord.Guild) bool {
			return in.ID == guildID
		})
		if idx == -1 {
			// Make sure to stop any active stuff with the previously active guild
			g, ok := b.activeGuilds.Load(guildID)
			if ok {
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

	return b.discord.Close()
}
