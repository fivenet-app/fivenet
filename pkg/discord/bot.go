package discord

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tJobProps   = table.FivenetJobProps
	tOauth2Accs = table.FivenetOauth2Accounts
	tAccs       = table.FivenetAccounts
	tUsers      = table.Users.AS("users")
)

const RoleFormat = "[%02d] %s"

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("discord_bot")
}

var BotModule = fx.Module("discord_bot",
	fx.Provide(
		NewBot,
	),
	fx.Decorate(wrapLogger),
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
	token    string

	id      string
	discord *discordgo.Session

	guildsMutex  sync.Mutex
	joinedGuilds []*discordgo.Guild

	syncInterval  time.Duration
	roleFormat    string
	nicknameRegex *regexp.Regexp
}

func NewBot(p BotParams) (*Bot, error) {
	if !p.Config.Discord.Bot.Enabled {
		return nil, nil
	}

	nicknameRegex, err := regexp.Compile(p.Config.Discord.Bot.NicknameRegex)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	b := &Bot{
		ctx:      ctx,
		logger:   p.Logger,
		db:       p.DB,
		enricher: p.Enricher,
		token:    p.Config.Discord.Bot.Token,

		syncInterval:  p.Config.Discord.Bot.SyncInterval,
		roleFormat:    p.Config.Discord.Bot.RoleFormat,
		nicknameRegex: nicknameRegex,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := b.start(p.Config.Discord.Bot.Token); err != nil {
			return err
		}

		go b.SyncRoles()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return b.discord.Close()
	}))

	return b, nil
}

func (b *Bot) start(token string) error {
	// Create a new Discord session using the provided login information.
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	b.discord = discord

	if err := b.refreshBotGuilds(); err != nil {
		return err
	}

	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		b.guildsMutex.Lock()
		defer b.guildsMutex.Unlock()

		b.joinedGuilds = discord.State.Guilds
		b.logger.Info(fmt.Sprintf("Ready with %d guilds", len(b.joinedGuilds)))
	})

	if err := discord.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

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

func (b *Bot) SyncRoles() error {
	for {
		select {
		case <-b.ctx.Done():
			return nil

		case <-time.After(b.syncInterval):
			if err := b.syncRoles(); err != nil {
				b.logger.Error("failed to sync roles", zap.Error(err))
			}
		}
	}
}

func (b *Bot) getGuilds() ([]*Guild, error) {
	if err := b.refreshBotGuilds(); err != nil {
		return nil, err
	}

	guilds, err := b.getGuildsFromDB()
	if err != nil {
		return nil, err
	}

	for i := len(guilds) - 1; i >= 0; i-- {
		var found *discordgo.Guild
		if !utils.InSliceFunc(b.joinedGuilds, func(in *discordgo.Guild) bool {
			if in.ID == guilds[i].ID {
				found = in
				return true
			}
			return false
		}) {
			guilds = utils.RemoveFromSlice(guilds, i)
		} else {
			if found == nil {
				return nil, fmt.Errorf("didn't find bot being in guild %s", guilds[i].ID)
			}
			guilds[i].init(b, found)
		}
	}

	return guilds, nil
}

func (b *Bot) getGuildsFromDB() ([]*Guild, error) {
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

	return dest, nil
}

func (b *Bot) syncRoles() error {
	guilds, err := b.getGuilds()
	if err != nil {
		return err
	}

	// Each guild is effectively associated with a Job via the JobProps
	for _, g := range guilds {
		if err := g.createJobRoles(); err != nil {
			return err
		}

		if err := g.syncUserInfo(); err != nil {
			return err
		}
	}

	return nil
}
