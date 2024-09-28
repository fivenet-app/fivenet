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

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/discord/commands"
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

const (
	disconnectRestartCountThreshold = 5
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

	metricGuildsReady = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "discord_bot",
		Name:      "guilds_ready_count",
		Help:      "Count of Discord guilds being ready.",
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

	id              string
	disconnectCount atomic.Uint64
	discord         *discordgo.Session
	activeGuilds    *xsync.MapOf[string, *Guild]
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

		disconnectCount: atomic.Uint64{},
		discord:         discord,
		activeGuilds:    xsync.NewMapOf[string, *Guild](),
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := b.start(ctx); err != nil {
			return err
		}

		// Cause discord client disconnects to cause bot restart after so many tries
		b.discord.AddHandler(func(discord *discordgo.Session, r *discordgo.Disconnect) {
			b.logger.Warn("discord client disconnected")

			b.disconnectCount.Add(1)

			if b.disconnectCount.Load() > disconnectRestartCountThreshold {
				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					b.logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
				}
			}
		})

		go func() {
			b.logger.Info("sleeping 5 seconds before running first discord sync")
			time.Sleep(5 * time.Second)

			if err := b.sync(); err != nil {
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

	b.discord.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		b.logger.Info("discord server joined", zap.String("discord_guild_id", g.ID))
	})

	if err := b.discord.Open(); err != nil {
		return fmt.Errorf("error opening discord connection: %w", err)
	}

	for {
		if b.discord.State.Ready.Version > 0 && ready.Load() {
			if err := b.refreshBotUserGuilds(); err != nil {
				return fmt.Errorf("failed to refresh bot user guilds. %w", err)
			}

			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("discord client failed to get ready, version %d", b.discord.State.Version)

		case <-time.After(750 * time.Millisecond):
		}
	}

	if err := b.setBotPresence(); err != nil {
		return fmt.Errorf("failed to set bot presence. %w", err)
	}

	if b.cfg.Commands.Enabled {
		if err := b.cmds.RegisterGlobalCommands(); err != nil {
			return fmt.Errorf("failed to register global commands. %w", err)
		}
	}

	return nil
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

func (b *Bot) refreshBotUserGuilds() error {
	usr, err := b.discord.User("@me")
	if err != nil {
		return fmt.Errorf("error obtaining account details: %w", err)
	}
	b.id = usr.ID

	return nil
}

func (b *Bot) sync() error {
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

	for job, guildID := range jobGuilds {
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
				g.Stop()

				b.activeGuilds.Delete(guildID)
			}

			continue
		}

		if found == nil {
			b.logger.Warn("didn't find bot being in guild", zap.String("discord_guild_id", guildID), zap.String("job", job))
			continue
		}

		if _, ok := b.activeGuilds.Load(guildID); ok {
			continue
		}

		g, err := NewGuild(b.ctx, b, found, job)
		if err != nil {
			return err
		}
		b.activeGuilds.Store(g.id, g)
	}

	return nil
}

func (b *Bot) getJobGuildsFromDB(ctx context.Context) (map[string]string, error) {
	stmt := tJobProps.
		SELECT(
			tJobProps.Job.AS("job"),
			tJobProps.DiscordGuildID.AS("id"),
		).
		FROM(tJobProps).
		WHERE(jet.AND(
			tJobProps.DiscordGuildID.IS_NOT_NULL(),
			tJobProps.Job.EQ(jet.String("mechanic")),
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

	guilds := map[string]string{}
	for _, g := range dest {
		guilds[g.Job] = g.ID
	}

	return guilds, nil
}

func (b *Bot) runSync(ctx context.Context) error {
	if err := b.getGuilds(ctx); err != nil {
		return fmt.Errorf("failed to get guilds. %w", err)
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

	metricGuildsTotal.Set(totalCount)
	metricGuildsReady.Set(readyCount)

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
				logger := b.logger.With(zap.String("job", g.job), zap.String("discord_guild_id", g.id))

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

	b.activeGuilds.Range(func(_ string, guild *Guild) bool {
		workChannel <- guild
		return true
	})

	close(workChannel)

	b.wg.Wait()

	return errs
}

func (b *Bot) stop() error {
	errs := multierr.Combine()
	b.activeGuilds.Range(func(key string, guild *Guild) bool {
		guild.Stop()

		return true
	})

	b.activeGuilds.Clear()

	if errs != nil {
		return errs
	}

	return b.discord.Close()
}
