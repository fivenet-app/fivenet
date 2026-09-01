package discord

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestDebugUserUsesUserJobs(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})

	// The legacy fivenet_user job is deliberately different. DebugUser must
	// report the matching row from fivenet_user_jobs instead.
	mock.ExpectQuery("SELECT .*fivenet_user_jobs.*").
		WillReturnRows(sqlmock.NewRows([]string{"groups", "enabled", "user_id", "job", "job_grade"}).
			AddRow(nil, true, int32(42), "police", int32(3)))

	guildID := discord.GuildID(123)
	userID := discord.UserID(456)
	dc := state.NewWithIntents("", gateway.IntentGuildMembers|gateway.IntentGuilds)
	require.NoError(t, dc.Cabinet.MemberSet(guildID, &discord.Member{User: discord.User{ID: userID}}, false))

	jobGuild := &Guild{job: "police", gid: guildID}
	b := &Bot{
		db:           db,
		dc:           dc,
		dcCfg:        &config.Discord{},
		logger:       zaptest.NewLogger(t),
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	b.activeGuilds.Store(guildID, jobGuild)

	got, err := b.DebugUser(context.Background(), guildID, userID)
	require.NoError(t, err)
	require.True(t, got.DiscordLinked)
	require.True(t, got.AccountFound)
	require.Equal(t, int32(42), got.UserID)
	require.Equal(t, "police", got.Job)
	require.Equal(t, int32(3), got.JobGrade)
	require.True(t, got.PartOfJob)
}

func TestGetJobGuildsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_job_props.*").
		WillReturnRows(sqlmock.NewRows([]string{"jobguild.job", "jobguild.guild_id", "jobguild.last_sync"}).
			AddRow("police", int64(123), nil))

	b := &Bot{db: db}
	got, err := b.getJobGuildsFromDB(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "police", got[0].Job)
	require.Equal(t, discord.GuildID(123), got[0].GuildID)
}

func TestBotStatusAndGuildLookup(t *testing.T) {
	interval := 5 * time.Minute
	b := &Bot{
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	b.syncTime.Store(&interval)
	b.activeGuilds.Store(20, &Guild{gid: 20, job: "police"})
	b.activeGuilds.Store(10, &Guild{gid: 10, job: "ambulance"})

	job, ok := b.GetJobFromGuildID(10)
	require.True(t, ok)
	require.Equal(t, "ambulance", job)
	_, ok = b.GetJobFromGuildID(99)
	require.False(t, ok)

	status := b.GetStatus()
	require.Equal(t, interval, status.SyncInterval)
	require.Equal(t, discord.GuildID(10), status.Guilds[0].GuildID)
	require.Equal(t, discord.GuildID(20), status.Guilds[1].GuildID)

	status, ok = b.GetStatusForGuild(20)
	require.True(t, ok)
	require.Equal(t, "police", status.Guilds[0].Job)
	_, ok = b.GetStatusForGuild(99)
	require.False(t, ok)
}

func TestRunSyncQueuesExistingGuild(t *testing.T) {
	workCh := make(chan *Guild, 1)
	guild := &Guild{gid: 10, job: "police"}
	b := &Bot{activeGuilds: xsync.NewMap[discord.GuildID, *Guild](), workCh: workCh}
	b.activeGuilds.Store(10, guild)

	running, err := b.RunSync(10)
	require.NoError(t, err)
	require.False(t, running)
	require.Same(t, guild, <-workCh)

	_, err = b.RunSync(99)
	require.Error(t, err)
}

func TestGetJobGuildsFromDBPropagatesErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).WillReturnError(context.Canceled)

	b := &Bot{db: db}
	_, err = b.getJobGuildsFromDB(context.Background())
	require.ErrorIs(t, err, context.Canceled)
}

func TestGetGuildsWithNoConfiguredGuildsDoesNothing(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_job_props.*").WillReturnRows(
		sqlmock.NewRows([]string{"jobguild.job", "jobguild.guild_id", "jobguild.last_sync"}),
	)

	b := &Bot{
		db:           db,
		logger:       zaptest.NewLogger(t),
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	require.NoError(t, b.getGuilds(context.Background()))
	count := 0
	for range b.activeGuilds.All() {
		count++
	}
	require.Zero(t, count)
}

func TestIsUserConfigAdminMatchesConfiguredLicense(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_accounts.*").WillReturnRows(
		sqlmock.NewRows([]string{"license", "groups"}).AddRow("license-1", nil),
	)

	b := &Bot{
		db:      db,
		authCfg: &config.Auth{ConfigAdminUsers: []string{"license-1"}},
	}
	isAdmin, err := b.IsUserConfigAdmin(context.Background(), discord.UserID(456))
	require.NoError(t, err)
	require.True(t, isAdmin)
}

func TestIsUserConfigAdminReturnsFalseForMissingAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_accounts.*").WillReturnRows(
		sqlmock.NewRows([]string{"license", "groups"}),
	)

	b := &Bot{db: db, authCfg: &config.Auth{}}
	isAdmin, err := b.IsUserConfigAdmin(context.Background(), discord.UserID(456))
	require.NoError(t, err)
	require.False(t, isAdmin)
}

func TestDebugUserReportsInactiveAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_user_jobs.*").WillReturnRows(
		sqlmock.NewRows([]string{"groups", "enabled", "user_id", "job", "job_grade"}).
			AddRow(nil, false, int32(42), "police", int32(3)),
	)

	guildID := discord.GuildID(123)
	userID := discord.UserID(456)
	dc := state.NewWithIntents("", gateway.IntentGuildMembers)
	require.NoError(t, dc.Cabinet.MemberSet(guildID, &discord.Member{User: discord.User{ID: userID}}, false))
	b := &Bot{
		db:           db,
		dc:           dc,
		dcCfg:        &config.Discord{},
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	b.activeGuilds.Store(guildID, &Guild{job: "police", gid: guildID})

	got, err := b.DebugUser(context.Background(), guildID, userID)
	require.NoError(t, err)
	require.True(t, got.AccountFound)
	require.False(t, got.AccountActive)
}

func TestIsUserConfigAdminMatchesConfiguredGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_accounts.*").WillReturnRows(
		sqlmock.NewRows([]string{"license", "groups"}).AddRow("license-1", []byte(`["admins"]`)),
	)

	b := &Bot{db: db, authCfg: &config.Auth{ConfigAdminGroups: []string{"admins"}}}
	isAdmin, err := b.IsUserConfigAdmin(context.Background(), discord.UserID(456))
	require.NoError(t, err)
	require.True(t, isAdmin)
}

func TestIsUserConfigAdminPropagatesDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectQuery("SELECT .*fivenet_accounts.*").WillReturnError(context.Canceled)

	b := &Bot{db: db, authCfg: &config.Auth{}}
	_, err = b.IsUserConfigAdmin(context.Background(), discord.UserID(456))
	require.ErrorIs(t, err, context.Canceled)
}

func TestIsUserGuildAdminUsesCachedPermissions(t *testing.T) {
	guildID := discord.GuildID(123)
	channelID := discord.ChannelID(234)
	userID := discord.UserID(456)
	dc := state.NewWithIntents("", gateway.IntentGuildMembers|gateway.IntentGuilds)
	require.NoError(t, dc.Cabinet.GuildSet(&discord.Guild{ID: guildID}, false))
	require.NoError(t, dc.Cabinet.ChannelSet(&discord.Channel{ID: channelID, GuildID: guildID}, false))
	require.NoError(t, dc.Cabinet.MemberSet(guildID, &discord.Member{User: discord.User{ID: userID}, RoleIDs: []discord.RoleID{99}}, false))
	require.NoError(t, dc.Cabinet.RoleSet(guildID, &discord.Role{ID: 99, Permissions: discord.PermissionAdministrator}, false))

	b := &Bot{dc: dc}
	isAdmin, err := b.IsUserGuildAdmin(context.Background(), channelID, userID)
	require.NoError(t, err)
	require.True(t, isAdmin)

	require.NoError(t, dc.Cabinet.MemberSet(guildID, &discord.Member{User: discord.User{ID: userID}}, true))
	isAdmin, err = b.IsUserGuildAdmin(context.Background(), channelID, userID)
	require.NoError(t, err)
	require.False(t, isAdmin)
}

func TestIsUserGuildAdminReturnsStateError(t *testing.T) {
	b := &Bot{dc: state.NewWithIntents("", gateway.IntentGuildMembers|gateway.IntentGuilds)}
	_, err := b.IsUserGuildAdmin(context.Background(), discord.ChannelID(234), discord.UserID(456))
	require.Error(t, err)
}

func TestRunSyncSkipsDisabledBot(t *testing.T) {
	b := &Bot{enabled: false, logger: zaptest.NewLogger(t)}
	require.NoError(t, b.runSync(context.Background()))
}

func TestRunSyncQueuesActiveGuilds(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_job_props.*").WillReturnRows(
		sqlmock.NewRows([]string{"jobguild.job", "jobguild.guild_id", "jobguild.last_sync"}),
	)

	workCh := make(chan *Guild, 1)
	guild := &Guild{gid: 123, job: "police"}
	b := &Bot{
		enabled:      true,
		db:           db,
		logger:       zaptest.NewLogger(t),
		metrics:      getBotMetrics(),
		workCh:       workCh,
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	b.activeGuilds.Store(guild.gid, guild)

	require.NoError(t, b.runSync(context.Background()))
	require.Same(t, guild, <-workCh)
}

func TestDebugUserMapsGroupsAndActualRoles(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_user_jobs.*").WillReturnRows(
		sqlmock.NewRows([]string{"groups", "enabled", "user_id", "job", "job_grade"}).
			AddRow([]byte(`["admins"]`), true, int32(42), "police", int32(3)),
	)

	guildID := discord.GuildID(123)
	userID := discord.UserID(456)
	dc := state.NewWithIntents("", gateway.IntentGuildMembers|gateway.IntentGuilds)
	require.NoError(t, dc.Cabinet.MemberSet(guildID, &discord.Member{User: discord.User{ID: userID}, RoleIDs: []discord.RoleID{99}}, false))
	require.NoError(t, dc.Cabinet.RoleSet(guildID, &discord.Role{ID: 99, Name: "Admin"}, false))
	b := &Bot{
		db:           db,
		dc:           dc,
		dcCfg:        &config.Discord{GroupSync: config.DiscordGroupSync{Enabled: true, Mapping: map[string]config.DiscordGroupRole{"admins": {RoleName: "Admin"}}}},
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	b.activeGuilds.Store(guildID, &Guild{job: "police", gid: guildID})

	got, err := b.DebugUser(context.Background(), guildID, userID)
	require.NoError(t, err)
	require.Equal(t, []string{"admins"}, got.MappedGroups)
	require.Equal(t, []string{"Admin"}, got.ActualRoles)
	require.Empty(t, got.MissingRoles)
}

func TestHandlePrivateMessageIgnoresNonPrivateSources(t *testing.T) {
	t.Parallel()
	b := &Bot{}

	// Bot-authored and webhook messages must be ignored before any Discord
	// client or i18n dependency is accessed.
	b.handlePrivateMessage(&gateway.MessageCreateEvent{
		Author: discord.User{Bot: true},
	})
	b.handlePrivateMessage(&gateway.MessageCreateEvent{
		Author:    discord.User{ID: 1},
		WebhookID: discord.WebhookID(2),
	})
	b.handlePrivateMessage(&gateway.MessageCreateEvent{
		Author:  discord.User{ID: 1},
		GuildID: discord.GuildID(2),
	})
}

func TestGetGuildsCreatesConfiguredGuild(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	lastSync := timestamp.Now().AsTime()
	mock.ExpectQuery("SELECT .*fivenet_job_props.*").WillReturnRows(
		sqlmock.NewRows([]string{"jobguild.job", "jobguild.guild_id", "jobguild.last_sync"}).
			AddRow("police", int64(123), lastSync),
	)
	// NewGuild loads the job's sync settings. An absent row is valid because
	// getSyncSettings applies defaults.
	mock.ExpectQuery("SELECT .*fivenet_job_props.*").WillReturnRows(
		sqlmock.NewRows([]string{"discord_sync_settings", "discord_sync_changes"}),
	)

	guildID := discord.GuildID(123)
	dc := state.NewWithIntents("", gateway.IntentGuilds)
	require.NoError(t, dc.Cabinet.GuildSet(&discord.Guild{ID: guildID, Name: "Police"}, false))
	b := &Bot{
		db:           db,
		dc:           dc,
		dcCfg:        &config.Discord{},
		logger:       zaptest.NewLogger(t),
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}

	require.NoError(t, b.getGuilds(context.Background()))
	created, ok := b.activeGuilds.Load(guildID)
	require.True(t, ok)
	require.Equal(t, "police", created.job)
	require.Equal(t, guildID, created.gid)
	created.Stop()
}

func TestGetGuildsRemovesActiveGuildMissingFromDiscord(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_job_props.*").WillReturnRows(
		sqlmock.NewRows([]string{"jobguild.job", "jobguild.guild_id", "jobguild.last_sync"}).
			AddRow("police", int64(123), nil),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	active := &Guild{gid: 123, job: "police", cancel: cancel}
	dc := state.NewWithIntents("", gateway.IntentGuilds)
	// Having another cached guild makes State.Guilds return the cache without
	// falling back to a REST request, while guild 123 is absent from Discord.
	require.NoError(t, dc.Cabinet.GuildSet(&discord.Guild{ID: 999}, false))
	b := &Bot{
		db:           db,
		dc:           dc,
		logger:       zaptest.NewLogger(t),
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	b.activeGuilds.Store(123, active)

	require.NoError(t, b.getGuilds(ctx))
	_, ok := b.activeGuilds.Load(123)
	require.False(t, ok)
}

func TestGetGuildsRetainsActiveGuild(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_job_props.*").WillReturnRows(
		sqlmock.NewRows([]string{"jobguild.job", "jobguild.guild_id", "jobguild.last_sync"}).
			AddRow("police", int64(123), nil),
	)

	guildID := discord.GuildID(123)
	active := &Guild{gid: guildID, job: "police"}
	dc := state.NewWithIntents("", gateway.IntentGuilds)
	require.NoError(t, dc.Cabinet.GuildSet(&discord.Guild{ID: guildID}, false))
	b := &Bot{
		db:           db,
		dc:           dc,
		logger:       zaptest.NewLogger(t),
		activeGuilds: xsync.NewMap[discord.GuildID, *Guild](),
	}
	b.activeGuilds.Store(guildID, active)

	require.NoError(t, b.getGuilds(context.Background()))
	got, ok := b.activeGuilds.Load(guildID)
	require.True(t, ok)
	require.Same(t, active, got)
}
