package discordtypes

import (
	"context"
	"errors"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/stretchr/testify/require"
)

func TestUsersAddMergesUsersAndDeduplicatesRoles(t *testing.T) {
	role := &Role{ID: 10, Name: "Police"}
	users := Users{}
	users.Add(&User{ID: 42, Nickname: strptr("old"), Roles: &UserRoles{Sum: Roles{role}}})
	users.Add(&User{ID: 42, Nickname: strptr("new"), Roles: &UserRoles{Sum: Roles{role, {ID: 11, Name: "Admin"}}}})

	got := users[42]
	require.Equal(t, "new", *got.Nickname)
	require.Len(t, got.Roles.Sum, 2)
}

func TestUserJobHelpers(t *testing.T) {
	ambulance := &users.UserJob{}
	ambulance.SetJob("ambulance")
	ambulance.SetGrade(1)
	police := &users.UserJob{}
	police.SetJob("police")
	police.SetGrade(3)
	user := &User{Jobs: []*users.UserJob{ambulance, police}}

	require.True(t, user.HasJob("police"))
	require.False(t, user.HasJob("fire"))
	require.Equal(t, int32(3), user.GetJobInfo("police").GetGrade())
	require.Nil(t, user.GetJobInfo("fire"))
}

func TestUserAddRoleDeduplicatesByIDOrName(t *testing.T) {
	user := &User{ID: 1}
	user.AddRole(&Role{ID: 10, Name: "Police"})
	user.AddRole(&Role{ID: 10, Name: "Different name"})
	user.AddRole(&Role{ID: 11, Name: "Police"})

	require.Len(t, user.Roles.Sum, 1)
}

func TestStateMerge(t *testing.T) {
	called := false
	processor := func(context.Context, discord.GuildID, discord.Member, *User) ([]discord.Embed, error) {
		called = true
		return nil, nil
	}
	state := &State{GuildID: 123, Users: Users{42: {ID: 42, Roles: &UserRoles{}}}}
	state.Merge(&State{
		Roles:          Roles{{ID: 1, Name: "Police"}},
		Users:          Users{42: {ID: 42, Nickname: strptr("Officer")}},
		UserProcessors: []UserProcessorHandler{processor},
	})
	state.Merge(nil)

	require.Len(t, state.Roles, 1)
	require.Equal(t, "Officer", *state.Users[42].Nickname)
	require.Len(t, state.UserProcessors, 1)
	_, err := state.UserProcessors[0](context.Background(), 123, discord.Member{}, state.Users[42])
	require.NoError(t, err)
	require.True(t, called)
}

func TestStateCalculateFiltersUnchangedUsersAndIgnoresBots(t *testing.T) {
	role := &Role{ID: 9, Name: "Keep"}
	state := &State{
		GuildID: 123,
		Roles:   Roles{role},
		Users: Users{
			1: {ID: 1, Roles: &UserRoles{Sum: Roles{role}}},
		},
	}
	dc := newTestDiscordState(t, 123, []discord.Member{
		{User: discord.User{ID: 1}, RoleIDs: []discord.RoleID{9}},
		{User: discord.User{ID: 2, Bot: true}},
	})

	plan, _, err := state.Calculate(context.Background(), dc, false)
	require.NoError(t, err)
	require.Empty(t, plan.Users)
}

func TestStateCalculateAddsAndRemovesManagedRoles(t *testing.T) {
	keep := &Role{ID: 9, Name: "Keep"}
	remove := &Role{ID: 10, Name: "Remove"}
	state := &State{
		GuildID: 123,
		Roles:   Roles{keep, remove},
		Users: Users{
			1: {ID: 1, Roles: &UserRoles{Sum: Roles{keep}}},
		},
	}
	dc := newTestDiscordState(t, 123, []discord.Member{{User: discord.User{ID: 1}, RoleIDs: []discord.RoleID{10}}})

	plan, _, err := state.Calculate(context.Background(), dc, false)
	require.NoError(t, err)
	require.Len(t, plan.Users, 1)
	require.Equal(t, discord.RoleID(9), plan.Users[0].Roles.ToAdd[0].ID)
	require.Equal(t, discord.RoleID(10), plan.Users[0].Roles.ToRemove[0].ID)
}

func TestStateCalculatePropagatesProcessorErrors(t *testing.T) {
	wantErr := errors.New("processor failed")
	state := &State{
		GuildID: 123,
		UserProcessors: []UserProcessorHandler{func(context.Context, discord.GuildID, discord.Member, *User) ([]discord.Embed, error) {
			return nil, wantErr
		}},
	}
	dc := newTestDiscordState(t, 123, []discord.Member{{User: discord.User{ID: 1}}})

	_, _, err := state.Calculate(context.Background(), dc, false)
	require.ErrorIs(t, err, wantErr)
}

func strptr(value string) *string { return &value }

func newTestDiscordState(t *testing.T, guildID discord.GuildID, members []discord.Member) *state.State {
	t.Helper()
	dc := state.NewWithIntents("", gateway.IntentGuildMembers|gateway.IntentGuilds)
	for i := range members {
		require.NoError(t, dc.Cabinet.MemberSet(guildID, &members[i], false))
	}
	for _, role := range []discord.Role{
		{ID: 9, Name: "Keep"},
		{ID: 10, Name: "Remove"},
	} {
		require.NoError(t, dc.Cabinet.RoleSet(guildID, &role, false))
	}
	return dc
}
