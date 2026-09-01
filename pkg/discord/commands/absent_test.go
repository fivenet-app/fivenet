package commands

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/stretchr/testify/require"
)

type absentBotState struct {
	discordtypes.BotState
	job string
	ok  bool
}

func (b absentBotState) GetJobFromGuildID(discord.GuildID) (string, bool) {
	return b.job, b.ok
}

func TestAbsentCommandRejectsInvalidDiscordContexts(t *testing.T) {
	t.Parallel()
	l, err := i18n.New()
	require.NoError(t, err)
	c := &AbsentCommand{l: l, b: absentBotState{}}

	resp := c.HandleCommand(context.Background(), cmdroute.CommandData{
		Event: &discord.InteractionEvent{GuildID: discord.NullGuildID},
	})
	require.Nil(t, resp)

	resp = c.HandleCommand(context.Background(), cmdroute.CommandData{
		Event: &discord.InteractionEvent{GuildID: 123},
	})
	require.NotNil(t, resp)
	require.NotEmpty(t, (*resp.Embeds)[0].Title)

	resp = c.HandleCommand(context.Background(), cmdroute.CommandData{
		Event: &discord.InteractionEvent{
			GuildID: 123,
			Member:  &discord.Member{},
		},
	})
	require.NotNil(t, resp)
	require.NotEmpty(t, (*resp.Embeds)[0].Title)
}

func TestAbsentCommandReportsUserWithoutMatchingJob(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectQuery("SELECT .*fivenet_user_jobs.*").WillReturnRows(
		sqlmock.NewRows([]string{"user_id", "job_grade"}),
	)

	l, err := i18n.New()
	require.NoError(t, err)
	c := &AbsentCommand{
		l:  l,
		db: db,
		b:  absentBotState{job: "police", ok: true},
	}
	resp := c.HandleCommand(context.Background(), cmdroute.CommandData{
		Event: &discord.InteractionEvent{
			GuildID: 123,
			Member:  &discord.Member{User: discord.User{ID: 456}},
		},
	})
	require.NotNil(t, resp)
	require.Contains(t, (*resp.Embeds)[0].Description, "unable to find")
}
