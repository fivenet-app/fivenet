package perms

import (
	"os"
	"testing"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/internal/modules"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/testdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestMain(m *testing.M) {
	// Enable ESX compatibility for database tables
	tables.EnableESXCompat()

	code := m.Run()
	os.Exit(code)
}

func TestBasicPerms(t *testing.T) {
	dbServer := servers.NewDBServer(t, true)
	natsServer := servers.NewNATSServer(t, true)

	ctx := t.Context()

	var ps perms.Permissions
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			fx.Invoke(func(p perms.Permissions) {
				ps = p
			}),
		)...,
	)
	require.NotNil(t, app)

	app.RequireStart()
	defer app.RequireStop()
	assert.NotNil(t, ps)

	userInfo := &userinfo.UserInfo{
		UserId:    testdata.Users[0].GetUserId(),
		Job:       testdata.Users[0].GetJob(),
		JobGrade:  testdata.Users[0].GetJobGrade(),
		Superuser: true,
	}
	// 1. Superuser can do everything
	can := ps.Can(userInfo, "citizens.CitizensService", "ListCitizens")
	assert.True(t, can, "Superuser has all permissions (citizens.CitizensService/ListCitizens)")
	can = ps.Can(userInfo, "documents.DocumentsService", "DeleteComment")
	assert.True(t, can, "Superuser has all permissions (documents.DocumentsService/DeleteComment)")

	// 2. Non-superuser (ambulance, 17) - Some basic tests that role permissions and attributes are working
	userInfo.Superuser = false
	can = ps.Can(userInfo, "citizens.CitizensService", "ListCitizens")
	assert.True(t, can, "User should have permission `citizens.CitizensService/ListCitizens`")
	can = ps.Can(userInfo, "jobs.TimeclockService", "ListTimeclock")
	assert.True(t, can, "User should have permission `jobs.TimeclockService/ListTimeclock`")

	can = ps.Can(userInfo, "settings.LawsService", "CreateOrUpdateLawBook")
	assert.False(
		t,
		can,
		"User should not have permission `settings.LawsService/CreateOrUpdateLawBook`",
	)
	can = ps.Can(userInfo, "settings.LawsService", "DeleteLawBook")
	assert.False(t, can, "User should not have permission `settings.SettingsService/DeleteLawBook`")

	attributes, err := ps.FlattenRoleAttributes(userInfo.GetJob(), userInfo.GetJobGrade())
	require.NoError(t, err, "FlattenRoleAttributes should not return an error")
	assert.NotEmpty(t, attributes, "FlattenRoleAttributes should return non-empty attributes")
	// Check if the expected flattened role attributes are returned
	assert.Len(t, attributes, 29, "FlattenRoleAttributes should return 29 attributes")
	for _, attr := range []string{
		"citizens-citizensservice-getuser-jobs-ambulance", "livemap-livemapservice-stream-players-ambulance", "mailer-mailerservice-createorupdateemail-fields-job",
		"jobs-conductservice-listconductentries-access-own", "jobs-conductservice-listconductentries-access-all", "wiki-wikiservice-createpage-fields-public",
		"qualifications-qualificationsservice-deletequalification-access-own", "qualifications-qualificationsservice-deletequalification-access-lower_rank",
		"qualifications-qualificationsservice-deletequalification-access-same_rank", "qualifications-qualificationsservice-deletequalification-access-any",
		"livemap-livemapservice-stream-markers-ambulance", "citizens-citizensservice-listcitizens-fields-userprops-bloodtype", "citizens-citizensservice-listcitizens-fields-phonenumber",
		"qualifications-qualificationsservice-updatequalification-access-own", "qualifications-qualificationsservice-updatequalification-access-lower_rank",
		"qualifications-qualificationsservice-updatequalification-access-same_rank", "qualifications-qualificationsservice-updatequalification-access-any",
		"livemap-livemapservice-deletemarker-access-own", "livemap-livemapservice-deletemarker-access-lower_rank", "livemap-livemapservice-deletemarker-access-same_rank",
		"livemap-livemapservice-deletemarker-access-any", "calendar-calendarservice-createcalendar-fields-job", "calendar-calendarservice-createcalendar-fields-public",
		"livemap-livemapservice-createorupdatemarker-access-own", "livemap-livemapservice-createorupdatemarker-access-lower_rank",
		"livemap-livemapservice-createorupdatemarker-access-same_rank", "livemap-livemapservice-createorupdatemarker-access-any",
		"qualifications-qualificationsservice-createqualification-fields-public", "completor-completorservice-completedocumentcategories-jobs-ambulance",
	} {
		assert.Contains(
			t,
			attributes,
			attr,
			"Make sure ambulance job has the expected attributes set for ambulance rank 17",
		)
	}
	for _, attr := range []string{"livemap-livemapservice-stream-players-doj", "livemap-livemapservice-stream-players-police"} {
		assert.NotContains(
			t,
			attributes,
			attr,
			"livemap-livemapservice-stream-players for doj and police should not be in the list of attributes for ambulance rank 17",
		)
	}

	role, err := ps.GetRoleByJobAndGrade(ctx, userInfo.GetJob(), userInfo.GetJobGrade())
	require.NoError(t, err, "GetRoleByJobAndGrade should not return an error")
	require.NotNil(t, role, "GetRoleByJobAndGrade should return non-nil role")

	rolePerms, err := ps.GetEffectiveRolePermissions(ctx, role.GetId())
	require.NoError(t, err, "GetEffectiveRolePermissions should not return an error")
	assert.Len(t, rolePerms, 43, "GetEffectiveRolePermissions should return 43 perms")

	// 3. Non-superuser (ambulance, 20) - should have more attributes than ambulance, 17 (further player locations access)
	userInfo = &userinfo.UserInfo{
		Job:      testdata.Users[1].GetJob(),
		JobGrade: testdata.Users[1].GetJobGrade(),
	}
	attributes, err = ps.FlattenRoleAttributes(userInfo.GetJob(), userInfo.GetJobGrade())
	require.NoError(t, err, "FlattenRoleAttributes should not return an error")
	assert.Len(t, attributes, 31, "FlattenRoleAttributes should now return 31 attributes")
	for _, attr := range []string{"livemap-livemapservice-stream-players-ambulance", "livemap-livemapservice-stream-players-doj", "livemap-livemapservice-stream-players-police"} {
		assert.Contains(
			t,
			attributes,
			attr,
			"livemap-livemapservice-stream-players for ambulance, doj and police should be in the list of attributes for ambulance rank 20",
		)
	}

	// 4. unemployed user - should not have any attributes but default perms
	userInfo = &userinfo.UserInfo{
		Job:      "unemployed",
		JobGrade: 0,
	}
	// Retrieve the default role id (technically this is a bit of a hack, but to make `GetEffectiveRolePermissions` work)
	role, err = ps.GetRoleByJobAndGrade(ctx, perms.DefaultRoleJob, userInfo.GetJobGrade())
	require.NoError(t, err, "GetRoleByJobAndGrade should not return an error")
	require.NotNil(t, role, "GetRoleByJobAndGrade should return non-nil role")
	rolePerms, err = ps.GetEffectiveRolePermissions(ctx, role.GetId())
	require.NoError(t, err, "GetEffectiveRolePermissions should not return an error")
	assert.Len(t, rolePerms, 4, "GetEffectiveRolePermissions should return 4 perms")
	attributes, err = ps.FlattenRoleAttributes(userInfo.GetJob(), userInfo.GetJobGrade())
	require.NoError(t, err, "FlattenRoleAttributes should not return an error")
	assert.Empty(t, attributes, "FlattenRoleAttributes should now return 0 attributes")
}
