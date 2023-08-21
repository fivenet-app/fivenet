package discord

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
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
}

func NewBot(p BotParams) *Bot {
	ctx, cancel := context.WithCancel(context.Background())

	b := &Bot{
		ctx:      ctx,
		logger:   p.Logger,
		db:       p.DB,
		enricher: p.Enricher,
		token:    p.Config.Discord.Bot.Token,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := b.start(); err != nil {
			return err
		}

		go b.SyncRoles()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return b.discord.Close()
	}))

	return b
}

func (b *Bot) refreshGuilds() error {
	usr, err := b.discord.User("@me")
	if err != nil {
		return fmt.Errorf("error obtaining account details: %w", err)
	}
	b.id = usr.ID

	return nil
}

func (b *Bot) start() error {
	// Create a new Discord session using the provided login information.
	discord, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	b.discord = discord

	if err := b.refreshGuilds(); err != nil {
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

func (b *Bot) SyncRoles() error {
	for {
		select {
		case <-b.ctx.Done():
			return nil

		case <-time.After(2 * time.Second):
			if err := b.syncRoles(); err != nil {
				b.logger.Error("failed to sync roles", zap.Error(err))
			}
		}
	}
}

type Guild struct {
	Job string `alias:"job"`
	ID  uint64 `alias:"id"`
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
	guilds, err := b.getGuildsFromDB()
	if err != nil {
		return err
	}

	if err := b.createJobRoles(guilds); err != nil {
		return err
	}

	for _, guild := range guilds {
		stmt := tOauth2Accs.
			SELECT(
				tOauth2Accs.ExternalID,
			).
			FROM(
				tOauth2Accs.
					INNER_JOIN(tAccs,
						tAccs.ID.EQ(tOauth2Accs.AccountID),
					).
					INNER_JOIN(tUsers,
						tUsers.Identifier.LIKE(jet.CONCAT(jet.String("char%:"), tAccs.License)),
					),
			).
			WHERE(jet.AND(
				tOauth2Accs.Provider.EQ(jet.String("discord")),
				tUsers.Job.EQ(jet.String(guild.Job)),
			))

		var dest []struct {
			ExternalID uint64
			Job        string
			JobGrade   int32
		}
		if err := stmt.QueryContext(b.ctx, b.db, &dest); err != nil {
			return err
		}

		for _, user := range dest {
			if err := b.setUserJobRole(user.ExternalID, user.Job, user.JobGrade); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bot) createJobRoles(guilds []*Guild) error {
	for _, g := range guilds {
		guild, err := b.discord.Guild(strconv.Itoa(int(g.ID)))
		if err != nil {
			return err
		}

		job := b.enricher.GetJobByName(g.Job)
		if job == nil {
			b.logger.Error("unknown job for discord guild")
			continue
		}

		for i := len(job.Grades) - 1; i >= 0; i-- {
			grade := job.Grades[i]
			name := fmt.Sprintf(RoleFormat, grade.Grade, grade.Label)

			if utils.InSliceFunc(guild.Roles, func(in *discordgo.Role) bool {
				return in.Name == name
			}) {
				continue
			}

			if _, err := b.discord.GuildRoleCreate(guild.ID, &discordgo.RoleParams{
				Name: name,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bot) setUserJobRole(ext uint64, job string, grade int32) error {
	// TODO set user's role in discord
	//b.discord.GuildMember()

	return nil
}

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("discord_bot")
}
