package perms

import (
	"context"
	"os"
	"testing"

	"github.com/fivenet-app/fivenet/v2025/internal/modules"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/testdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
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

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

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
		UserId:    testdata.Users[0].UserId,
		Job:       testdata.Users[0].Job,
		JobGrade:  testdata.Users[0].JobGrade,
		Superuser: true,
	}
	// 1. Superuser can do everything
	can := ps.Can(userInfo, "CitizensService", "ListCitizens")
	assert.True(t, can, "Superuser has all permissions (CitizensService.ListCitizens)")
	can = ps.Can(userInfo, "DocumentsService", "DeleteComment")
	assert.True(t, can, "Superuser has all permissions (DocumentsService.DeleteComment)")

	// 2. Non-superuser (ambulance, 17) - Some basic tests that role permissions and attributes are working
	userInfo.Superuser = false
	can = ps.Can(userInfo, "CitizensService", "ListCitizens")
	assert.True(t, can, "User should have permission `CitizensService.ListCitizens`")
	can = ps.Can(userInfo, "TimeclockService", "ListTimeclock")
	assert.True(t, can, "User should have permission `TimeclockService.ListTimeclock`")

	can = ps.Can(userInfo, "LawsService", "CreateOrUpdateLawBook")
	assert.False(t, can, "User should not have permission `LawsService.CreateOrUpdateLawBook`")
	can = ps.Can(userInfo, "LawsService", "DeleteLawBook")
	assert.False(t, can, "User should not have permission `SettingsService.DeleteLawBook`")

	attributes, err := ps.FlattenRoleAttributes(userInfo.Job, userInfo.JobGrade)
	assert.NoError(t, err, "FlattenRoleAttributes should not return an error")
	assert.NotEmpty(t, attributes, "FlattenRoleAttributes should return non-empty attributes")
	// Check if the expected flattened role attributes are returned
	assert.Len(t, attributes, 29, "FlattenRoleAttributes should return 29 attributes")
	for _, attr := range []string{"citizensservice-getuser-jobs-ambulance", "livemapservice-stream-players-ambulance", "mailerservice-createorupdateemail-fields-job", "jobsconductservice-listconductentries-access-own", "jobsconductservice-listconductentries-access-all", "wikiservice-createpage-fields-public", "qualificationsservice-deletequalification-access-own", "qualificationsservice-deletequalification-access-lower_rank", "qualificationsservice-deletequalification-access-same_rank", "qualificationsservice-deletequalification-access-any", "livemapservice-stream-markers-ambulance", "citizensservice-listcitizens-fields-user_props-bloodtype", "citizensservice-listcitizens-fields-phonenumber", "qualificationsservice-updatequalification-access-own", "qualificationsservice-updatequalification-access-lower_rank", "qualificationsservice-updatequalification-access-same_rank", "qualificationsservice-updatequalification-access-any", "livemapservice-deletemarker-access-own", "livemapservice-deletemarker-access-lower_rank", "livemapservice-deletemarker-access-same_rank", "livemapservice-deletemarker-access-any", "calendarservice-createcalendar-fields-job", "calendarservice-createcalendar-fields-public", "livemapservice-createorupdatemarker-access-own", "livemapservice-createorupdatemarker-access-lower_rank", "livemapservice-createorupdatemarker-access-same_rank", "livemapservice-createorupdatemarker-access-any", "qualificationsservice-createqualification-fields-public", "completorservice-completedocumentcategories-jobs-ambulance"} {
		assert.Contains(t, attributes, attr, "Make sure ambulance job has the expected attributes set for ambulance rank 17")
	}
	for _, attr := range []string{"livemapservice-stream-players-doj", "livemapservice-stream-players-police"} {
		assert.NotContains(t, attributes, attr, "livemapservice-stream-players for doj and police should not be in the list of attributes for ambulance rank 17")
	}

	role, err := ps.GetRoleByJobAndGrade(ctx, userInfo.Job, userInfo.JobGrade)
	assert.NoError(t, err, "GetRoleByJobAndGrade should not return an error")
	require.NotNil(t, role, "GetRoleByJobAndGrade should return non-nil role")

	rolePerms, err := ps.GetEffectiveRolePermissions(ctx, role.ID)
	assert.NoError(t, err, "GetEffectiveRolePermissions should not return an error")
	assert.Len(t, rolePerms, 43, "GetEffectiveRolePermissions should return 43 perms")

	// 3. Non-superuser (ambulance, 20) - should have more attributes than ambulance, 17 (further player locations access)
	userInfo = &userinfo.UserInfo{
		Job:      testdata.Users[1].Job,
		JobGrade: testdata.Users[1].JobGrade,
	}
	attributes, err = ps.FlattenRoleAttributes(userInfo.Job, userInfo.JobGrade)
	assert.NoError(t, err, "FlattenRoleAttributes should not return an error")
	assert.Len(t, attributes, 31, "FlattenRoleAttributes should now return 31 attributes")
	for _, attr := range []string{"livemapservice-stream-players-ambulance", "livemapservice-stream-players-doj", "livemapservice-stream-players-police"} {
		assert.Contains(t, attributes, attr, "livemapservice-stream-players for ambulance, doj and police should be in the list of attributes for ambulance rank 20")
	}

	// 4. unemployed user - should not have any attributes but default perms
	userInfo = &userinfo.UserInfo{
		Job:      "unemployed",
		JobGrade: 0,
	}
	// Retrieve the default role id (technically this is a bit of a hack, but to make `GetEffectiveRolePermissions` work)
	role, err = ps.GetRoleByJobAndGrade(ctx, perms.DefaultRoleJob, userInfo.JobGrade)
	assert.NoError(t, err, "GetRoleByJobAndGrade should not return an error")
	require.NotNil(t, role, "GetRoleByJobAndGrade should return non-nil role")
	rolePerms, err = ps.GetEffectiveRolePermissions(ctx, role.ID)
	assert.NoError(t, err, "GetEffectiveRolePermissions should not return an error")
	assert.Len(t, rolePerms, 4, "GetEffectiveRolePermissions should return 4 perms")
	attributes, err = ps.FlattenRoleAttributes(userInfo.Job, userInfo.JobGrade)
	assert.NoError(t, err, "FlattenRoleAttributes should not return an error")
	assert.Len(t, attributes, 0, "FlattenRoleAttributes should now return 0 attributes")
}
