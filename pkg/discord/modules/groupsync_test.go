package modules

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGroupSyncPlanUser(t *testing.T) {
	type fields struct {
		job     string
		cfg     *config.Discord
		dbSetup func(t *testing.T) (*sql.DB, func())
	}

	type args struct {
		user  *groupSyncUser
		roles map[string]*discordtypes.Role
	}

	tests := []struct {
		name string

		fields fields
		args   args

		wantUserNil    bool
		wantRoles      []*discordtypes.Role
		wantErr        bool
		wantLogTitle   string
		wantSameJobSet bool
	}{
		{
			name: "returns nil when user has no mapped groups",
			fields: fields{
				cfg: &config.Discord{
					GroupSync: config.DiscordGroupSync{
						Mapping: map[string]config.DiscordGroupRole{
							"supporter": {RoleName: "Supporter"},
						},
					},
				},
			},
			args: args{
				user: &groupSyncUser{
					ExternalID: "12345",
					Groups:     &accounts.AccountGroups{Groups: []string{"unknown"}},
				},
				roles: map[string]*discordtypes.Role{"supporter": {ID: 1, Name: "Supporter"}},
			},
			wantUserNil: true,
		},
		{
			name: "returns embed when mapped role cannot be found",
			fields: fields{
				cfg: &config.Discord{
					GroupSync: config.DiscordGroupSync{
						Mapping: map[string]config.DiscordGroupRole{
							"supporter": {RoleName: "Supporter"},
						},
					},
				},
			},
			args: args{
				user: &groupSyncUser{
					ExternalID: "12345",
					Groups:     &accounts.AccountGroups{Groups: []string{"supporter"}},
				},
				roles: map[string]*discordtypes.Role{"otherrole": {ID: 2, Name: "OtherRole"}},
			},
			wantUserNil:  true,
			wantLogTitle: "Group Sync: Failed to find dc role for group Supporter",
		},
		{
			name: "returns error when external id is not numeric",
			fields: fields{
				cfg: &config.Discord{
					GroupSync: config.DiscordGroupSync{
						Mapping: map[string]config.DiscordGroupRole{
							"supporter": {RoleName: "Supporter"},
						},
					},
				},
			},
			args: args{
				user: &groupSyncUser{
					ExternalID: "not-a-number",
					Groups:     &accounts.AccountGroups{Groups: []string{"supporter"}},
				},
				roles: map[string]*discordtypes.Role{"supporter": {ID: 1, Name: "Supporter"}},
			},
			wantUserNil: true,
			wantErr:     true,
		},
		{
			name: "returns user with all mapped roles",
			fields: fields{
				cfg: &config.Discord{
					GroupSync: config.DiscordGroupSync{
						Mapping: map[string]config.DiscordGroupRole{
							"supporter": {RoleName: "Supporter"},
							"donator":   {RoleName: "Donator"},
						},
					},
				},
			},
			args: args{
				user: &groupSyncUser{
					ExternalID: "12345",
					Groups:     &accounts.AccountGroups{Groups: []string{"supporter", "donator"}},
				},
				roles: map[string]*discordtypes.Role{
					"supporter": {ID: 1, Name: "Supporter"},
					"donator":   {ID: 2, Name: "Donator"},
				},
			},
			wantUserNil: false,
			wantRoles: []*discordtypes.Role{
				{ID: 1, Name: "Supporter"},
				{ID: 2, Name: "Donator"},
			},
		},
		{
			name: "returns user with all mapped roles (single group mapping)",
			fields: fields{
				cfg: &config.Discord{
					GroupSync: config.DiscordGroupSync{
						Mapping: map[string]config.DiscordGroupRole{
							"supporter": {RoleName: "Supporter"},
						},
					},
				},
			},
			args: args{
				user: &groupSyncUser{
					ExternalID: "12345",
					Groups:     &accounts.AccountGroups{Groups: []string{"supporter"}},
				},
				roles: map[string]*discordtypes.Role{
					"supporter": {ID: 1, Name: "Supporter"},
				},
			},
			wantUserNil: false,
			wantRoles: []*discordtypes.Role{
				{ID: 1, Name: "Supporter"},
			},
		},
		{
			name: "not same job enabled and user without matching job still gets role",
			fields: fields{
				job: "police",
				cfg: &config.Discord{
					GroupSync: config.DiscordGroupSync{
						Mapping: map[string]config.DiscordGroupRole{
							"supporter": {RoleName: "Supporter", NotSameJob: true},
						},
					},
				},
				dbSetup: func(t *testing.T) (*sql.DB, func()) {
					db, mock, err := sqlmock.New()
					require.NoError(t, err)
					mock.ExpectQuery("SELECT .*fivenet_user_jobs.*").
						// Third argument is the query limit.
						WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 10).
						WillReturnRows(sqlmock.NewRows([]string{"job", "grade"}))

					cleanup := func() {
						assert.NoError(t, mock.ExpectationsWereMet())
						_ = db.Close()
					}
					return db, cleanup
				},
			},
			args: args{
				user: &groupSyncUser{
					ExternalID: "12345",
					UserID:     42,
					Groups:     &accounts.AccountGroups{Groups: []string{"supporter"}},
				},
				roles: map[string]*discordtypes.Role{"supporter": {ID: 1, Name: "Supporter"}},
			},
			wantUserNil: false,
			wantRoles: []*discordtypes.Role{
				{ID: 1, Name: "Supporter"},
			},
		},
		{
			name: "not same job check errors and user is skipped",
			fields: fields{
				job: "police",
				cfg: &config.Discord{
					GroupSync: config.DiscordGroupSync{
						Mapping: map[string]config.DiscordGroupRole{
							"supporter": {RoleName: "Supporter", NotSameJob: true},
						},
					},
				},
				dbSetup: func(t *testing.T) (*sql.DB, func()) {
					db, mock, err := sqlmock.New()
					require.NoError(t, err)
					mock.ExpectQuery("SELECT .*fivenet_user_jobs.*").
						// Third argument is the query limit.
						WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 10).
						WillReturnError(errors.New("db error"))

					cleanup := func() {
						assert.NoError(t, mock.ExpectationsWereMet())
						_ = db.Close()
					}
					return db, cleanup
				},
			},
			args: args{
				user: &groupSyncUser{
					ExternalID: "12345",
					UserID:     42,
					Groups:     &accounts.AccountGroups{Groups: []string{"supporter"}},
				},
				roles: map[string]*discordtypes.Role{"supporter": {ID: 1, Name: "Supporter"}},
			},
			wantUserNil: true,
		},
	}

	logger := zaptest.NewLogger(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := &BaseModule{
				logger: logger,
				cfg:    tt.fields.cfg,
				job:    tt.fields.job,
			}

			var cleanup func()
			if tt.fields.dbSetup != nil {
				db, c := tt.fields.dbSetup(t)
				base.db = db
				cleanup = c
			}
			if cleanup != nil {
				defer cleanup()
			}

			g := &GroupSync{BaseModule: base}

			user, logs, err := g.planUser(context.Background(), tt.args.user, tt.args.roles)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantUserNil {
				assert.Nil(t, user)
			} else {
				require.NotNil(t, user)
				assert.Equal(t, discord.UserID(12345), user.ID)
				require.NotNil(t, user.Roles)
				assert.ElementsMatch(t, tt.wantRoles, user.Roles.Sum)
			}

			if tt.wantLogTitle != "" {
				require.NotEmpty(t, logs)
				assert.Equal(t, tt.wantLogTitle, logs[0].Title)
			} else {
				assert.Empty(t, logs)
			}

			assert.Equal(t, tt.wantSameJobSet, tt.args.user.SameJob)
		})
	}
}

func TestGroupSyncPlanRoles(t *testing.T) {
	t.Parallel()

	g := &GroupSync{
		BaseModule: &BaseModule{
			logger: zaptest.NewLogger(t),
			cfg: &config.Discord{
				GroupSync: config.DiscordGroupSync{
					Mapping: map[string]config.DiscordGroupRole{
						"Group.One": {RoleName: "Shared Role", Color: "#112233"},
						"group.two": {RoleName: "Shared Role", Color: "#445566"},
						"  group.three  ": {RoleName: "Unique Role"},
					},
				},
			},
		},
	}

	roles := g.planRoles()

	require.Len(t, roles, 3)
	require.Contains(t, roles, "group.one")
	require.Contains(t, roles, "group.two")
	require.Contains(t, roles, "group.three")

	// Group keys with the same configured role name should point to the same planned role.
	assert.Same(t, roles["group.one"], roles["group.two"])
	assert.Equal(t, "Shared Role", roles["group.one"].Name)

	assert.NotSame(t, roles["group.one"], roles["group.three"])
	assert.Equal(t, "Unique Role", roles["group.three"].Name)
}

func TestGroupSyncPlanRolesSharedRolePermissions(t *testing.T) {
	t.Parallel()

	permissions := int64(discord.PermissionKickMembers)

	g := &GroupSync{
		BaseModule: &BaseModule{
			logger: zaptest.NewLogger(t),
			cfg: &config.Discord{
				GroupSync: config.DiscordGroupSync{
					Mapping: map[string]config.DiscordGroupRole{
						"group.one": {RoleName: "Shared Role"},
						"group.two": {RoleName: "Shared Role", Permissions: &permissions},
					},
				},
			},
		},
	}

	roles := g.planRoles()

	require.Len(t, roles, 2)
	require.Contains(t, roles, "group.one")
	require.Contains(t, roles, "group.two")
	assert.Same(t, roles["group.one"], roles["group.two"])

	require.NotNil(t, roles["group.one"].Permissions)
	assert.Equal(t, discord.Permissions(permissions), *roles["group.one"].Permissions)
}
