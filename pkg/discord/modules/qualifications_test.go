package modules

import (
	"sync/atomic"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jobssettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/settings"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestQualificationsPlanRoles(t *testing.T) {
	t.Parallel()
	var settings atomic.Pointer[jobssettings.DiscordSyncSettings]
	settings.Store(&jobssettings.DiscordSyncSettings{QualificationsRoleFormat: "%abbr% - %name%"})

	g := &QualificationsSync{BaseModule: &BaseModule{
		logger:   zaptest.NewLogger(t),
		job:      "police",
		settings: &settings,
	}}

	customFormat := "Qualification %abbr% | %name%"
	entries := []*qualificationsEntry{
		{
			ID:                 1,
			Abbreviation:       "EMT",
			QualificationTitle: "Empty Role",
			DiscordSettings:    &qualifications.QualificationDiscordSettings{},
		},
		{
			ID:                 2,
			Abbreviation:       "EMT",
			QualificationTitle: "Advanced",
			DiscordSettings: &qualifications.QualificationDiscordSettings{
				RoleName: new("Advanced"),
			},
		},
		{
			ID:                 3,
			Abbreviation:       "K9",
			QualificationTitle: "Handler",
			DiscordSettings: &qualifications.QualificationDiscordSettings{
				RoleName:   new("Handler"),
				RoleFormat: &customFormat,
			},
		},
	}

	roles, logs, err := g.planRoles(entries)
	require.NoError(t, err)

	require.Len(t, logs, 1)
	assert.Contains(t, logs[0].Title, "Empty role name")
	assert.Contains(t, logs[0].Title, "ID: 1")

	require.Len(t, roles, 2)
	assert.Equal(t, "EMT - Advanced", roles[2].Name)
	assert.Equal(t, "Qualifications-2", roles[2].Module)
	assert.Equal(t, "police", roles[2].Job)

	assert.Equal(t, "Qualification K9 | Handler", roles[3].Name)
	assert.Equal(t, "Qualifications-3", roles[3].Module)
}

func TestQualificationsPlanUsersMergesRolesForSameUser(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	for range 2 {
		mock.ExpectQuery("SELECT .*fivenet_qualifications_results.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"qualification_user_mapping.external_id",
				"qualification_user_mapping.user_id",
				"jobs.job",
				"jobs.grade",
			}).
				AddRow("12345", int32(7), "police", int32(1)))
	}

	g := &QualificationsSync{BaseModule: &BaseModule{db: db, job: "police"}}

	roleA := &discordtypes.Role{Name: "Role A", Module: "Qualifications-1", Job: "police"}
	roleB := &discordtypes.Role{Name: "Role B", Module: "Qualifications-2", Job: "police"}
	users, logs, err := g.planUsers(t.Context(), map[int64]*discordtypes.Role{
		1: roleA,
		2: roleB,
	})

	require.NoError(t, err)
	assert.Empty(t, logs)
	require.Len(t, users, 1)

	user := users[12345]
	require.NotNil(t, user)
	require.NotNil(t, user.Roles)
	assert.ElementsMatch(t, []*discordtypes.Role{roleA, roleB}, user.Roles.Sum)
}

func TestQualificationsQueryAndPlanUsersForQualificationSkipsInvalidExternalID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	}()

	mock.ExpectQuery("SELECT .*fivenet_qualifications_results.*").
		WillReturnRows(sqlmock.NewRows([]string{
			"qualification_user_mapping.external_id",
			"qualification_user_mapping.user_id",
			"jobs.job",
			"jobs.grade",
		}).
			AddRow("12345", int32(1), "police", int32(1)).
			AddRow("invalid-id", int32(2), "police", int32(2)))

	g := &QualificationsSync{BaseModule: &BaseModule{db: db, job: "police"}}
	users := discordtypes.Users{}
	role := &discordtypes.Role{Name: "Role", Module: "Qualifications-1", Job: "police"}

	err = g.queryAndPlanUsersForQualification(t.Context(), 1, role, &users)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse user oauth2 external id")

	require.Len(t, users, 1)
	user := users[12345]
	require.NotNil(t, user)
	require.NotNil(t, user.Roles)
	assert.ElementsMatch(t, []*discordtypes.Role{role}, user.Roles.Sum)
}
