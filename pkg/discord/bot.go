package discord

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/server/metrics"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tJobProps = table.FivenetJobProps
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
	lastSync = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metrics.Namespace,
		Subsystem: "discord_bot",
		Name:      "last_sync",
		Help:      "Last time sync has completed.",
	}, []string{"job", "status"})
)

type BotParams struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	DB       *sql.DB
	Enricher *mstlystcdata.Enricher
	Config   *config.Config
}

type Bot struct {
	ctx      context.Context
	logger   *zap.Logger
	db       *sql.DB
	enricher *mstlystcdata.Enricher
	cfg      *config.DiscordBot

	syncInterval time.Duration

	id           string
	discord      *discordgo.Session
	activeGuilds *xsync.MapOf[string, *Guild]
}

func NewBot(p BotParams) (*Bot, error) {
	if !p.Config.Discord.Bot.Enabled {
		return nil, nil
	}

	// Create a new Discord session using the provided login information.
	discord, err := discordgo.New("Bot " + p.Config.Discord.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences

	ctx, cancel := context.WithCancel(context.Background())
	b := &Bot{
		ctx:      ctx,
		logger:   p.Logger,
		db:       p.DB,
		enricher: p.Enricher,
		cfg:      &p.Config.Discord.Bot,

		discord:      discord,
		activeGuilds: xsync.NewMapOf[string, *Guild](),

		syncInterval: p.Config.Discord.Bot.SyncInterval,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := b.start(ctx); err != nil {
			return err
		}

		go func() {
			if err := b.setupSync(); err != nil {
				b.logger.Error("failed to set up sync for guilds")
				return
			}

			b.Sync()
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		// Stop all guilds and discord session
		return b.stop()
	}))

	return b, nil
}

func (b *Bot) start(ctx context.Context) error {
	b.discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		b.logger.Info(fmt.Sprintf("Ready with %d guilds", len(discord.State.Guilds)))
	})

	if err := b.discord.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	for {
		if b.discord.State.Ready.Version > 0 {
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
	b.logger.Info("Set bot presence")

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
	b.logger.Info("Sleeping 5 seconds before running first discord sync")
	time.Sleep(5 * time.Second)

	if err := b.runSync(); err != nil {
		b.logger.Error("failed to sync roles", zap.Error(err))
	}

	for {
		select {
		case <-b.ctx.Done():
			return nil

		case <-time.After(b.syncInterval):
			b.logger.Info("Running Discord Sync")
			if err := b.runSync(); err != nil {
				b.logger.Error("failed to sync roles", zap.Error(err))
			}
		}
	}
}

// getGuilds Each guild is effectively associated with a Job via the JobProps
func (b *Bot) getGuilds() error {
	if err := b.refreshBotGuilds(); err != nil {
		return err
	}

	guildsDB, err := b.getGuildsFromDB()
	if err != nil {
		return err
	}

	for job, guildID := range guildsDB {
		var found *discordgo.Guild
		if !utils.InSliceFunc(b.discord.State.Guilds, func(in *discordgo.Guild) bool {
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

func (b *Bot) getGuildsFromDB() (map[string]string, error) {
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
	if err := stmt.QueryContext(b.ctx, b.db, &dest); err != nil {
		return nil, err
	}

	guilds := map[string]string{}
	for _, g := range dest {
		guilds[g.Job] = g.ID
	}

	return guilds, nil
}

func (b *Bot) setupSync() error {
	if err := b.getGuilds(); err != nil {
		return err
	}

	var e error
	b.activeGuilds.Range(func(key string, guild *Guild) bool {
		if err := guild.Setup(); err != nil {
			e = err
			return false
		}

		return true
	})

	return e
}

func (b *Bot) runSync() error {
	if err := b.getGuilds(); err != nil {
		return err
	}

	b.activeGuilds.Range(func(key string, guild *Guild) bool {
		if err := guild.Run(); err != nil {
			b.logger.Error("error during sync", zap.String("job", guild.Job), zap.String("discord_guild_id", guild.ID),
				zap.Error(err))
			lastSync.WithLabelValues(guild.Job, "failed").SetToCurrentTime()
			return false
		}

		lastSync.WithLabelValues(guild.Job, "success").SetToCurrentTime()
		return true
	})

	return nil
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
