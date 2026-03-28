package modules

import (
	"errors"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	jobssettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/settings"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestConstructUserNickname(t *testing.T) {
	tests := []struct {
		name      string
		firstname string
		lastname  string
		prefix    string
		suffix    string
		expected  string
	}{
		{
			name:      "Basic name without prefix or suffix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "",
			suffix:    "",
			expected:  "John Doe",
		},
		{
			name:      "Name with prefix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "Dr.",
			suffix:    "",
			expected:  "Dr. John Doe",
		},
		{
			name:      "Name with suffix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "",
			suffix:    "PhD",
			expected:  "John Doe PhD",
		},
		{
			name:      "Name with prefix and suffix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "Dr.",
			suffix:    "PhD",
			expected:  "Dr. John Doe PhD",
		},
		{
			name:      "Name exceeding max length",
			firstname: "Jonathan",
			lastname:  "Doe-Smith-Jackson",
			prefix:    "Dr.",
			suffix:    "PhD",
			expected:  "Dr. J. Doe-Smith-Jackson PhD",
		},
		{
			name:      "Name with only prefix exceeding max length",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "A very longp",
			suffix:    "",
			expected:  "A very longp John Doe",
		},
		{
			name:      "Name with only suffix exceeding max length",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "",
			suffix:    "A very longs",
			expected:  "John Doe A very longs",
		},
		{
			name:      "Empty lastname",
			firstname: "Max",
			lastname:  "",
			prefix:    "Dr.",
			suffix:    "PhD",
			expected:  "Dr. Max PhD",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := &UserInfo{}
			result := g.constructUserNickname(
				test.firstname,
				test.lastname,
				test.prefix,
				test.suffix,
			)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

var jobGrades = map[string]map[int32]*discordtypes.Role{
	"police": {
		// Employee role for police job
		0: {ID: 110, Name: "Police"},
		// Job grade roles for police job
		1: {ID: 1, Name: "Grade1"},
		2: {ID: 2, Name: "Grade2"},
		3: {ID: 3, Name: "Grade3"},
	},
}

func TestPlanUser(t *testing.T) {
	type fields struct {
		job            string
		employeeRole   *discordtypes.Role
		unemployedRole *discordtypes.Role
		absenceRole    *discordtypes.Role
		jobGradeRoles  map[int32]*discordtypes.Role
		groupRoles     map[string]*discordtypes.Role
	}
	type args struct {
		u        *userRoleMapping
		settings *jobssettings.DiscordSyncSettings
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mockDiscord    func(*testing.T) IState
		wantUserNil    bool
		wantErr        bool
		wantLogTitle   string
		wantRolesCount int

		expectedUser *discordtypes.User
	}{
		{
			name: "discord.Member returns not found, returns log",
			fields: fields{
				job: "police",
			},
			args: args{
				u: &userRoleMapping{
					ExternalID: "12345",
					Firstname:  "John",
					Lastname:   "Doe",
					Jobs: []*users.UserJob{
						{Job: "police", Grade: 1},
					},
				},
				settings: &jobssettings.DiscordSyncSettings{},
			},
			mockDiscord: func(t *testing.T) IState {
				return &mockDiscord{
					memberFunc: func(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error) {
						return nil, &httputil.HTTPError{Status: http.StatusNotFound}
					},
				}
			},
			wantUserNil:  true,
			wantErr:      false,
			wantLogTitle: "UserInfo: Employee not found on Discord: John Doe",
		},
		{
			name: "discord.Member returns error, returns error",
			fields: fields{
				job: "police",
			},
			args: args{
				u: &userRoleMapping{
					ExternalID: "12345",
					Firstname:  "John",
					Lastname:   "Doe",
					Jobs: []*users.UserJob{
						{Job: "police", Grade: 1},
					},
				},
				settings: &jobssettings.DiscordSyncSettings{},
			},
			mockDiscord: func(t *testing.T) IState {
				return &mockDiscord{
					memberFunc: func(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error) {
						return nil, errors.New("api error")
					},
				}
			},
			wantUserNil: true,
			wantErr:     true,
		},
		{
			name: "nickname sync enabled, nickname is set",
			fields: fields{
				job: "police",
			},
			args: args{
				u: &userRoleMapping{
					ExternalID: "12345",
					Firstname:  "Jane",
					Lastname:   "Smith",
					Jobs: []*users.UserJob{
						{Job: "police", Grade: 1},
					},
				},
				settings: func() *jobssettings.DiscordSyncSettings {
					s := &jobssettings.DiscordSyncSettings{}
					s.UserInfoSyncSettings = &jobssettings.UserInfoSyncSettings{
						SyncNicknames: true,
					}
					return s
				}(),
			},
			mockDiscord: func(t *testing.T) IState {
				return &mockDiscord{
					memberFunc: func(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error) {
						return &discord.Member{
							User: discord.User{ID: 12345},
							Nick: "",
						}, nil
					},
				}
			},
			wantUserNil: false,
			wantErr:     false,
		},
		{
			name: "user already has the correct nickname and roles, no changes needed",
			fields: fields{
				job:           "police",
				employeeRole:  jobGrades["police"][0],
				jobGradeRoles: jobGrades["police"],
			},
			args: args{
				u: &userRoleMapping{
					ExternalID: "12345",
					Firstname:  "Jane",
					Lastname:   "Smith",
					Jobs: []*users.UserJob{
						{Job: "police", Grade: 1},
					},
				},
				settings: func() *jobssettings.DiscordSyncSettings {
					s := &jobssettings.DiscordSyncSettings{}
					s.UserInfoSyncSettings = &jobssettings.UserInfoSyncSettings{
						SyncNicknames:       true,
						EmployeeRoleEnabled: true,
					}
					return s
				}(),
			},
			mockDiscord: func(t *testing.T) IState {
				return &mockDiscord{
					memberFunc: func(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error) {
						return &discord.Member{
							User: discord.User{ID: 12345},
							Nick: "Jane Smith",
							RoleIDs: []discord.RoleID{
								110, // Employee role
								1,   // Job grade role for grade 1
							},
						}, nil
					},
				}
			},
			wantUserNil: false,
			wantErr:     false,
			expectedUser: &discordtypes.User{
				ID: 12345,
				// Nickname is empty because no update is needed in this test
				Roles: &discordtypes.UserRoles{
					Sum: discordtypes.Roles{
						jobGrades["police"][1], // Grade 1 role
						jobGrades["police"][0], // Employee role
					},
				},
				Jobs: []*users.UserJob{
					{
						Job:   "police",
						Grade: 1,
					},
				},
			},
		},
		{
			name: "user has incorrect roles, roles to add and remove are calculated",
			fields: fields{
				job:           "police",
				employeeRole:  jobGrades["police"][0],
				jobGradeRoles: jobGrades["police"],
			},
			args: args{
				u: &userRoleMapping{
					ExternalID: "12345",
					Firstname:  "Jane",
					Lastname:   "Smith",
					Jobs: []*users.UserJob{
						{Job: "police", Grade: 2},
					},
				},
				settings: func() *jobssettings.DiscordSyncSettings {
					s := &jobssettings.DiscordSyncSettings{}
					s.UserInfoSyncSettings = &jobssettings.UserInfoSyncSettings{
						SyncNicknames:       true,
						EmployeeRoleEnabled: true,
					}
					return s
				}(),
			},
			mockDiscord: func(t *testing.T) IState {
				return &mockDiscord{
					memberFunc: func(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error) {
						return &discord.Member{
							User: discord.User{ID: 12345},
							Nick: "Jane Smith",
							RoleIDs: []discord.RoleID{
								110, // Employee role
								1,   // Incorrect job grade role (grade 1 instead of grade 2)
							},
						}, nil
					},
				}
			},
			wantUserNil: false,
			wantErr:     false,
			expectedUser: &discordtypes.User{
				ID: 12345,
				// Nickname is empty because no update is needed in this test
				Roles: &discordtypes.UserRoles{
					Sum: discordtypes.Roles{
						jobGrades["police"][0], // Employee role
						jobGrades["police"][2], // Correct job grade role for grade 2
					},
				},
				Jobs: []*users.UserJob{
					{
						Job:   "police",
						Grade: 2,
					},
				},
			},
		},
	}

	logger := zaptest.NewLogger(t)
	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDiscord := tt.mockDiscord(t)
			g := &UserInfo{
				BaseModule: &BaseModule{
					logger:  logger,
					discord: mockDiscord,
					guild:   discord.Guild{ID: 1, OwnerID: 999},
					job:     tt.fields.job,
				},
				employeeRole:   tt.fields.employeeRole,
				absenceRole:    tt.fields.absenceRole,
				groupRoles:     tt.fields.groupRoles,
				jobGradeRoles:  jobGrades[tt.fields.job],
				unemployedRole: tt.fields.unemployedRole,
			}

			user, logs, err := g.planUser(tt.args.u, tt.args.settings)
			if tt.wantErr {
				require.Error(err, tt.name+": expected an error but got none")
			} else {
				require.NoError(err, tt.name+": unexpected error: %v", err)
			}

			if tt.wantUserNil {
				assert.Nil(user, tt.name+": expected user to be nil")
			} else {
				assert.NotNil(user, tt.name+": expected user to not be nil")

				if tt.expectedUser != nil {
					assert.Equal(tt.expectedUser.ID, user.ID, tt.name+": expected user id to match")
					assert.Equal(
						tt.expectedUser.Nickname,
						user.Nickname,
						tt.name+": expected nickname to match",
					)
					if tt.expectedUser.Roles != nil {
						assert.NotNil(user.Roles, tt.name+": expected user roles to not be nil")
						assert.ElementsMatch(
							tt.expectedUser.Roles.Sum,
							user.Roles.Sum,
							tt.name+": expected roles sum to match",
						)
						assert.ElementsMatch(
							tt.expectedUser.Roles.ToAdd,
							user.Roles.ToAdd,
							tt.name+": expected roles to add to match",
						)
						assert.ElementsMatch(
							tt.expectedUser.Roles.ToRemove,
							user.Roles.ToRemove,
							tt.name+": expected roles to remove to match",
						)
					} else {
						assert.Nil(user.Roles, tt.name+": expected user roles to be nil")
					}
					assert.Equal(
						tt.expectedUser.Jobs,
						user.Jobs,
						tt.name+": expected jobs to match",
					)
				}
			}

			if tt.wantLogTitle != "" {
				// logs may be empty if user is nil and not found, so allow empty logs if wantLogTitle is set to empty string
				assert.NotEmpty(logs, tt.name+": expected logs to not be empty")
				if len(logs) > 0 {
					assert.Equal(tt.wantLogTitle, logs[0].Title, tt.name+": log title mismatch")
				}
			}
		})
	}
}

// mockDiscord is a mock implementation of the Discord state for testing purposes.
type mockDiscord struct {
	*state.State

	memberFunc func(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error)
}

func (m *mockDiscord) Member(
	guildID discord.GuildID,
	userID discord.UserID,
) (*discord.Member, error) {
	if m.memberFunc != nil {
		return m.memberFunc(guildID, userID)
	}
	return &discord.Member{User: discord.User{ID: userID}}, nil
}
