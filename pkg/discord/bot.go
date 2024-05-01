package discord

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/discord/commands"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/server/admin"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var (
	tJobProps = table.FivenetJobProps.AS("jobprops")
)

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
	}, []string{"job", "status"})

	metricGuildsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "discord_bot",
		Name:      "guilds_total_count",
		Help:      "Total count of Discord guilds being ready.",
	})

	metricGuildsReady = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "discord_bot",
		Name:      "guilds_ready_count",
		Help:      "Count of Discord guilds being ready.",
	})
)

type BotParams struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	Enricher  *mstlystcdata.Enricher
	Config    *config.Config
	AppConfig appconfig.IConfig
}

type Bot struct {
	ctx      context.Context
	logger   *zap.Logger
	tracer   trace.Tracer
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	cfg      *config.Discord
	appCfg   appconfig.IConfig

	cmds *commands.Cmds

	wg sync.WaitGroup

	id           string
	discord      *discordgo.Session
	activeGuilds *xsync.MapOf[string, *Guild]
}

func NewBot(p BotParams) (*Bot, error) {
	if !p.Config.Discord.Enabled {
		return nil, nil
	}

	// Create a new Discord session using the provided login information.
	discord, err := discordgo.New("Bot " + p.Config.Discord.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating discord session. %w", err)
	}
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences

	cmds, err := commands.New(p.Logger, discord, p.Config)
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

		cmds: cmds,

		wg: sync.WaitGroup{},

		discord:      discord,
		activeGuilds: xsync.NewMapOf[string, *Guild](),
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := b.start(ctx); err != nil {
			return err
		}

		go func() {
			b.Sync()
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		// Stop all guilds and discord session
		if err := b.stop(); err != nil {
			return err
		}

		cancel()

		return nil
	}))

	return b, nil
}

func (b *Bot) start(ctx context.Context) error {
	var ready atomic.Bool
	b.discord.AddHandler(func(discord *discordgo.Session, r *discordgo.Ready) {
		b.logger.Info(fmt.Sprintf("Ready with %d guilds", len(r.Guilds)))
		ready.Store(true)
	})

	b.discord.AddHandler(func(s *discordgo.Session, i *discordgo.GuildCreate) {
		b.logger.Info("discord server joined", zap.String("discord_guild_id", i.ID))
	})

	if err := b.discord.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	if b.cfg.Commands.Enabled {
		if err := b.cmds.RegisterGlobalCommands(); err != nil {
			return fmt.Errorf("failed to register global commands. %w", err)
		}
	}

	for {
		if b.discord.State.Ready.Version > 0 && ready.Load() {
			return b.setBotPresence()
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("discord client failed to get ready, version %d", b.discord.State.Version)

		case <-time.After(750 * time.Millisecond):
		}
	}
}

func (b *Bot) setBotPresence() error {
	if b.cfg.Presence.GameStatus != nil {
		if err := b.discord.UpdateGameStatus(0, *b.cfg.Presence.GameStatus); err != nil {
			return err
		}
	} else if b.cfg.Presence.ListeningStatus != nil {
		if err := b.discord.UpdateListeningStatus(*b.cfg.Presence.ListeningStatus); err != nil {
			return err
		}
	} else if b.cfg.Presence.StreamingStatus != nil {
		url := ""
		if b.cfg.Presence.StreamingStatusUrl != nil {
			url = *b.cfg.Presence.StreamingStatusUrl
		}
		if err := b.discord.UpdateStreamingStatus(0, *b.cfg.Presence.StreamingStatus, url); err != nil {
			return err
		}
	} else if b.cfg.Presence.WatchStatus != nil {
		if err := b.discord.UpdateWatchStatus(0, *b.cfg.Presence.WatchStatus); err != nil {
			return err
		}
	}
	b.logger.Info("bot presence has been set")

	return nil
}

func (b *Bot) refreshBotGuilds() error {
	usr, err := b.discord.User("@me")
	if err != nil {
		return fmt.Errorf("error obtaining account details: %w", err)
	}
	b.id = usr.ID

	return nil
}

func (b *Bot) Sync() error {
	b.logger.Info("sleeping 5 seconds before running first discord sync")
	time.Sleep(5 * time.Second)

	func() {
		ctx, span := b.tracer.Start(b.ctx, "discord_bot")
		defer span.End()

		if err := b.runSync(ctx); err != nil {
			b.logger.Error("failed to sync roles", zap.Error(err))
		}
	}()

	for {
		select {
		case <-b.ctx.Done():
			return nil

		case <-time.After(b.appCfg.Get().Discord.SyncInterval.AsDuration()):
			b.logger.Info("running discord sync")
			func() {
				ctx, span := b.tracer.Start(b.ctx, "discord_bot")
				defer span.End()

				if err := b.runSync(ctx); err != nil {
					b.logger.Error("failed to sync discord", zap.Error(err))
				}
			}()
		}
	}
}

// getGuilds Each guild is effectively associated with a Job via the JobProps
func (b *Bot) getGuilds(ctx context.Context) error {
	if err := b.refreshBotGuilds(); err != nil {
		return err
	}

	guildsDB, err := b.getJobGuildsFromDB(ctx)
	if err != nil {
		return err
	}

	if len(guildsDB) == 0 {
		b.logger.Debug("no job discord guild connections found")
		return nil
	}

	for job, guildID := range guildsDB {
		var found *discordgo.Guild
		if !slices.ContainsFunc(b.discord.State.Ready.Guilds, func(in *discordgo.Guild) bool {
			if in.ID == guildID {
				found = in
				return true
			}
			return false
		}) {
			// Make sure to stop any active stuff with the previously active guild
			g, ok := b.activeGuilds.Load(guildID)
			if ok {
				err := g.Stop()
				b.activeGuilds.Delete(guildID)
				if err != nil {
					return err
				}
			}

			continue
		}

		if found == nil {
			return fmt.Errorf("didn't find bot being in guild %s (job: %s)", guildID, job)
		}

		if _, ok := b.activeGuilds.Load(guildID); ok {
			continue
		}

		g, err := NewGuild(b, found, job)
		if err != nil {
			return err
		}
		b.activeGuilds.Store(g.ID, g)
	}

	return nil
}

func (b *Bot) getJobGuildsFromDB(ctx context.Context) (map[string]string, error) {
	stmt := tJobProps.
		SELECT(
			tJobProps.Job.AS("guild.job"),
			tJobProps.DiscordGuildID.AS("guild.id"),
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.DiscordGuildID.IS_NOT_NULL(),
		)

	var dest []*Guild
	if err := stmt.QueryContext(ctx, b.db, &dest); err != nil {
		return nil, err
	}

	guilds := map[string]string{}
	for _, g := range dest {
		guilds[g.Job] = g.ID
	}

	return guilds, nil
}

func (b *Bot) runSync(ctx context.Context) error {
	if err := b.getGuilds(ctx); err != nil {
		return err
	}

	totalCount := float64(0)
	readyCount := float64(0)
	b.activeGuilds.Range(func(key string, value *Guild) bool {
		totalCount++
		if !value.guild.Unavailable {
			readyCount++
		}

		return true
	})

	metricGuildsTotal.Add(totalCount)
	metricGuildsReady.Add(readyCount)

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
				logger := b.logger.With(zap.String("job", g.Job), zap.String("discord_guild_id", g.ID))

				if err := g.Run(); err != nil {
					logger.Error("error during sync", zap.Error(err))
					errs = multierr.Append(errs, err)

					metricLastSync.WithLabelValues(g.Job, "failed").SetToCurrentTime()
				} else {
					metricLastSync.WithLabelValues(g.Job, "success").SetToCurrentTime()
				}

				if err := b.setLastSyncInterval(ctx, g.Job); err != nil {
					logger.Error("error setting job props last sync time", zap.Error(err))
					errs = multierr.Append(errs, err)
				}
			}(guild)
		}
	}()

	b.activeGuilds.Range(func(_ string, guild *Guild) bool {
		workChannel <- guild
		return true
	})

	close(workChannel)

	b.wg.Wait()

	return errs
}

func (b *Bot) stop() error {
	var e error
	b.activeGuilds.Range(func(key string, guild *Guild) bool {
		if err := guild.Stop(); err != nil {
			e = err
			return false
		}

		return true
	})

	b.activeGuilds.Clear()

	if e != nil {
		return e
	}

	return b.discord.Close()
}

func (b *Bot) setLastSyncInterval(ctx context.Context, job string) error {
	tJobProps := table.FivenetJobProps
	stmt := tJobProps.
		UPDATE(
			tJobProps.DiscordLastSync,
		).
		SET(
			tJobProps.DiscordLastSync.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		)

	if _, err := stmt.ExecContext(ctx, b.db); err != nil {
		return err
	}

	return nil
}
