package discord

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
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
