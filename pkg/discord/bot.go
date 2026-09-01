package discord

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	botWorkerCount = 3

	metricsSubsystem = "discord_bot"
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("discord.bot")
}

var BotModule = fx.Module("discord.bot",
	fx.Provide(
		New,
	),
	fx.Decorate(wrapLogger),
)

type botMetrics struct {
	lastSync     *prometheus.GaugeVec
	guildsTotal  prometheus.Gauge
	syncDuration *prometheus.GaugeVec
}

var (
	botMetricsOnce sync.Once
	botMetricsInst *botMetrics
)

func getBotMetrics() *botMetrics {
	botMetricsOnce.Do(func() {
		botMetricsInst = &botMetrics{
			lastSync: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "last_sync",
				Help:      "Last time sync has completed.",
			}, []string{admin.MetricsJobNameLabel, "status"}),
			guildsTotal: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "guilds_total_count",
				Help:      "Total count of Discord guilds being ready.",
			}),
			syncDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "sync_duration_seconds",
				Help:      "Duration of the last sync operation in seconds.",
			}, []string{admin.MetricsJobNameLabel}),
		}

		prometheus.MustRegister(
			botMetricsInst.lastSync,
			botMetricsInst.guildsTotal,
			botMetricsInst.syncDuration,
		)
	})

	return botMetricsInst
}

type BotParams struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	JS        *events.JSWrapper
	DB        *sql.DB
	Enricher  mstlystcdata.IEnricher
	Config    *config.Config
	AppConfig appconfig.IConfig
	Perms     perms.Permissions
	I18n      i18n.Ii18n

	Discord *state.State
}

type Bot struct {
	discordtypes.BotState

	logger   *zap.Logger
	tracer   trace.Tracer
	js       *events.JSWrapper
	db       *sql.DB
	enricher mstlystcdata.IEnricher
	dcCfg    *config.Discord
	authCfg  *config.Auth
	appCfg   appconfig.IConfig
	perms    perms.Permissions
	i18n     i18n.Ii18n

	publicURL          string
	oauth2ProviderName string
	metrics            *botMetrics

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
	ctxCancel, cancel := context.WithCancel(context.Background())

	oauth2ProviderName := "discord"
	if provider := p.Config.OAuth2.GetProviderByType(
		config.OAuth2ProviderDiscord,
	); provider != nil {
		oauth2ProviderName = provider.Name
	}

	b := &Bot{
		logger:   p.Logger,
		tracer:   p.TP.Tracer("discord.bot"),
		js:       p.JS,
		db:       p.DB,
		enricher: p.Enricher,
		dcCfg:    &p.Config.Discord,
		authCfg:  &p.Config.Auth,
		appCfg:   p.AppConfig,
		perms:    p.Perms,
		i18n:     p.I18n,

		publicURL:          p.Config.HTTP.PublicURL,
		oauth2ProviderName: oauth2ProviderName,
		metrics:            getBotMetrics(),

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
					case <-ctxCancel.Done():
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

								b.metrics.lastSync.WithLabelValues(guild.job, "failed").
									SetToCurrentTime()
							} else {
								b.metrics.lastSync.WithLabelValues(guild.job, "success").
									SetToCurrentTime()
							}
						}()

						b.metrics.syncDuration.WithLabelValues(guild.job).
							Set(elapsed.Seconds())
					}
				}
			})
		}

		if err := registerStreams(ctxStartup, b.js); err != nil {
			return err
		}

		if err := b.start(ctxCancel); err != nil {
			return err
		}

		// Handle app config updates
		go func() {
			configUpdateCh := b.appCfg.Subscribe()
			for {
				select {
				case <-ctxCancel.Done():
					b.appCfg.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					if cfg == nil {
						continue
					}
					b.handleAppConfigUpdate(ctxCancel, cfg)
				}
			}
		}()

		go b.syncLoop(ctxCancel)

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

	b.dc.AddHandler(b.handlePrivateMessage)

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
		b.logger.Info("running discord sync", zap.Bool("dry_run", b.dcCfg.DryRun))
		func() {
			ctx, span := b.tracer.Start(ctx, "discord.bot")
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

	b.metrics.guildsTotal.Set(totalCount)

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

func (b *Bot) GetStatus() discordtypes.BotStatus {
	status := discordtypes.BotStatus{}
	if interval := b.syncTime.Load(); interval != nil {
		status.SyncInterval = *interval
	}

	for _, guild := range b.activeGuilds.All() {
		status.Guilds = append(status.Guilds, guild.status())
	}
	sort.Slice(status.Guilds, func(i, j int) bool {
		return status.Guilds[i].GuildID < status.Guilds[j].GuildID
	})

	return status
}

func (b *Bot) GetStatusForGuild(guildID discord.GuildID) (discordtypes.BotStatus, bool) {
	guild, ok := b.activeGuilds.Load(guildID)
	if !ok || guild == nil {
		return discordtypes.BotStatus{}, false
	}

	status := discordtypes.BotStatus{}
	if interval := b.syncTime.Load(); interval != nil {
		status.SyncInterval = *interval
	}
	status.Guilds = []discordtypes.GuildStatus{guild.status()}
	return status, true
}

func (b *Bot) IsUserConfigAdmin(ctx context.Context, userID discord.UserID) (bool, error) {
	account := &struct {
		License string                  `alias:"license"`
		Groups  *accounts.AccountGroups `alias:"groups"`
	}{}
	stmt := table.FivenetAccountsOauth2.
		SELECT(
			table.FivenetAccounts.License.AS("license"),
			table.FivenetAccounts.Groups.AS("groups"),
		).
		FROM(
			table.FivenetAccountsOauth2.
				INNER_JOIN(
					table.FivenetAccounts,
					table.FivenetAccounts.ID.EQ(table.FivenetAccountsOauth2.AccountID),
				),
		).
		WHERE(mysql.AND(
			table.FivenetAccountsOauth2.Provider.EQ(mysql.String("discord")),
			table.FivenetAccountsOauth2.ExternalID.EQ(
				mysql.String(strconv.FormatUint(uint64(userID), 10)),
			),
			table.FivenetAccounts.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, b.db, account); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	_, _, configGroups, configUsers := userinfo.EffectiveAdminLists(
		b.authCfg.JobAdminGroups, b.authCfg.JobAdminUsers,
		b.authCfg.ConfigAdminGroups, b.authCfg.ConfigAdminUsers,
		b.appCfg,
	)
	return account.Groups.ContainsAnyGroup(configGroups) ||
		slices.Contains(configUsers, account.License), nil
}

func (b *Bot) DebugUser(
	ctx context.Context,
	guildID discord.GuildID,
	userID discord.UserID,
) (*discordtypes.UserDebug, error) {
	result := &discordtypes.UserDebug{DiscordUserID: userID}
	guild, ok := b.activeGuilds.Load(guildID)
	if !ok || guild == nil {
		return result, nil
	}
	result.MemberFound = func() bool { _, err := b.dc.Member(guildID, userID); return err == nil }()
	result.SyncRunning = guild.IsRunning()
	result.LastSync = guild.status().LastSync
	result.GroupSyncEnabled = b.dcCfg.GroupSync.Enabled
	result.UserInfoEnabled = b.dcCfg.UserInfoSync.Enabled

	var row struct {
		Groups   *accounts.AccountGroups `alias:"groups"`
		Enabled  bool                    `alias:"enabled"`
		UserID   *int32                  `alias:"user_id"`
		Job      *string                 `alias:"job"`
		JobGrade *int32                  `alias:"job_grade"`
	}
	tAccounts := table.FivenetAccounts
	tUsers := table.FivenetUser
	stmt := table.FivenetAccountsOauth2.
		SELECT(
			tAccounts.Groups.AS("groups"),
			tAccounts.Enabled.AS("enabled"),
			tUsers.ID.AS("user_id"),
			tUsers.Job.AS("job"),
			tUsers.JobGrade.AS("job_grade"),
		).
		FROM(table.FivenetAccountsOauth2.
			INNER_JOIN(tAccounts, tAccounts.ID.EQ(table.FivenetAccountsOauth2.AccountID)).
			LEFT_JOIN(table.FivenetUserAccounts, table.FivenetUserAccounts.AccountID.EQ(tAccounts.ID)).
			LEFT_JOIN(tUsers, tUsers.ID.EQ(table.FivenetUserAccounts.UserID)),
		).
		WHERE(mysql.AND(
			table.FivenetAccountsOauth2.Provider.EQ(mysql.String("discord")),
			table.FivenetAccountsOauth2.ExternalID.EQ(
				mysql.String(strconv.FormatUint(uint64(userID), 10)),
			),
		)).
		LIMIT(1)
	if err := stmt.QueryContext(ctx, b.db, &row); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, nil
		}
		return nil, err
	}
	result.DiscordLinked, result.AccountFound, result.AccountActive = true, true, row.Enabled
	if row.UserID != nil {
		result.UserID = *row.UserID
	}
	if row.Job != nil {
		result.Job = *row.Job
	}
	if row.JobGrade != nil {
		result.JobGrade = *row.JobGrade
	}
	result.PartOfJob = result.Job == guild.job
	result.Groups = append(result.Groups, row.Groups.GetGroups()...)

	if b.dcCfg.GroupSync.Enabled {
		for group, mapping := range b.dcCfg.GroupSync.Mapping {
			if slices.ContainsFunc(
				result.Groups,
				func(v string) bool { return strings.EqualFold(strings.TrimSpace(v), group) },
			) {
				result.MappedGroups = append(result.MappedGroups, group)
				if !mapping.NotSameJob || !result.PartOfJob {
					result.ExpectedRoles = append(result.ExpectedRoles, mapping.RoleName)
				}
			}
		}
	}
	member, err := b.dc.Member(guildID, userID)
	if err == nil {
		roles, roleErr := b.dc.Roles(guildID)
		if roleErr == nil {
			for _, role := range roles {
				if slices.Contains(member.RoleIDs, role.ID) {
					result.ActualRoles = append(result.ActualRoles, role.Name)
				}
			}
		}
	}
	for _, expected := range result.ExpectedRoles {
		if !slices.Contains(result.ActualRoles, expected) {
			result.MissingRoles = append(result.MissingRoles, expected)
		}
	}
	return result, nil
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

func (b *Bot) handlePrivateMessage(ev *gateway.MessageCreateEvent) {
	// Ignore messages from bots and webhooks
	if ev.Author.Bot || ev.WebhookID.IsValid() {
		return
	}
	// Only react to private messages, ignore guild messages
	if ev.GuildID.IsValid() {
		return
	}

	locale := b.i18n.GetFallbackLanguage()
	if ev.Author.Locale != "" {
		locale = ev.Author.Locale
	}
	t := b.i18n.Translator(locale)

	if _, err := b.dc.SendMessageComplex(ev.ChannelID, api.SendMessageData{
		Embeds: []discord.Embed{
			{
				Type:        discord.NormalEmbed,
				Title:       t("discord.messages.private_message.title", nil),
				Description: t("discord.messages.private_message.desc", nil),
				Author:      embeds.EmbedAuthor,
				Color:       embeds.ColorInfo,
				Footer:      embeds.EmbedFooterFiveNet,
			},
		},
		Components: discord.Components(
			&discord.ActionRowComponent{
				&discord.ButtonComponent{
					Label:    t("discord.messages.private_message.help_button", nil),
					Style:    discord.SecondaryButtonStyle(),
					CustomID: "help",
				},
				&discord.ButtonComponent{
					Label: t("discord.commands.fivenet.open_link", nil),
					Style: discord.LinkButtonStyle(b.publicURL),
				},
			},
		),
	}); err != nil {
		b.logger.Error(
			"failed to send message via direct message",
			zap.Uint64("discord_user_id", uint64(ev.Author.ID)),
			zap.Error(err),
		)
	}
}
