package perms

import (
	"os"
	"testing"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/internal/modules"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/testdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func guard(s string) string {
	return slug.Make(s)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestBasicPerms(t *testing.T) {
	t.Parallel()
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

	flattenRoleAttributes := func(roleAttrs []*permissionsattributes.RoleAttribute) []string {
		attributes := make([]string, 0, len(roleAttrs))
		for _, rAttr := range roleAttrs {
			attrKey := guard(
				rAttr.GetNamespace() + "." + rAttr.GetService() + "/" + rAttr.GetName() + "." + rAttr.GetKey(),
			)
			switch permissionsattributes.AttributeTypes(rAttr.GetType()) {
			case permissionsattributes.StringListAttributeType:
				for _, v := range rAttr.GetValue().GetStringList().GetStrings() {
					attributes = append(attributes, guard(attrKey+"."+v))
				}
			case permissionsattributes.JobListAttributeType:
				for _, v := range rAttr.GetValue().GetJobList().GetStrings() {
					attributes = append(attributes, guard(attrKey+"."+v))
				}
			case permissionsattributes.JobGradeListAttributeType:
				for v := range rAttr.GetValue().GetJobGradeList().GetJobs() {
					attributes = append(attributes, guard(attrKey+"."+v))
				}
			}
		}

		return attributes
	}

	userInfo := &userinfo.UserInfo{
		UserId:    testdata.Users[0].GetUserId(),
		Job:       testdata.Users[0].GetJob(),
		JobGrade:  testdata.Users[0].GetJobGrade(),
		Superuser: true,
	}
	// 1. Superuser can do everything
	can := ps.CanRaw(userInfo, "citizens", "CitizensService", "ListCitizens")
	assert.True(t, can, "Superuser has all permissions (citizens.CitizensService/ListCitizens)")
	can = ps.Can(userInfo, permsdocuments.CommentsService.DeleteComment.Perm)
	assert.True(t, can, "Superuser has all permissions (documents.CommentsService/DeleteComment)")
	// 1.1. Superuser doesn't have access to non-existing permissions
	can = ps.CanRaw(userInfo, "test123", "TestService", "DoSomething")
	assert.False(
		t,
		can,
		"Superuser has not access to non-existing permissions (test123.TestService/DoSomething)",
	)

	// 2. Non-superuser (ambulance, 17) - Some basic tests that role permissions and attributes are working
	userInfo.Superuser = false
	can = ps.CanRaw(userInfo, "citizens", "CitizensService", "ListCitizens")
	assert.True(t, can, "User should have permission `citizens.CitizensService/ListCitizens`")
	can = ps.CanRaw(userInfo, "jobs", "TimeclockService", "ListTimeclock")
	assert.True(t, can, "User should have permission `jobs.TimeclockService/ListTimeclock`")

	can = ps.CanRaw(userInfo, "settings", "LawsService", "CreateOrUpdateLawBook")
	assert.False(
		t,
		can,
		"User should not have permission `settings.LawsService/CreateOrUpdateLawBook`",
	)
	can = ps.CanRaw(userInfo, "settings", "LawsService", "DeleteLawBook")
	assert.False(t, can, "User should not have permission `settings.SettingsService/DeleteLawBook`")

	roleAttrs, err := ps.GetEffectiveRoleAttributes(ctx, userInfo.GetJob(), userInfo.GetJobGrade())
	require.NoError(t, err, "GetEffectiveRoleAttributes should not return an error")
	attributes := flattenRoleAttributes(roleAttrs)
	assert.NotEmpty(t, attributes, "GetEffectiveRoleAttributes should return non-empty attributes")
	// Check if the expected flattened role attributes are returned
	assert.Len(t, attributes, 29, "GetEffectiveRoleAttributes should return 29 attributes")
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
	assert.Len(t, rolePerms, 44, "GetEffectiveRolePermissions should return 44 perms")

	// 3. Non-superuser (ambulance, 20) - should have more attributes than ambulance, 17 (further player locations access)
	userInfo = &userinfo.UserInfo{
		Job:      testdata.Users[1].GetJob(),
		JobGrade: testdata.Users[1].GetJobGrade(),
	}
	roleAttrs, err = ps.GetEffectiveRoleAttributes(ctx, userInfo.GetJob(), userInfo.GetJobGrade())
	require.NoError(t, err, "GetEffectiveRoleAttributes should not return an error")
	attributes = flattenRoleAttributes(roleAttrs)
	assert.Len(t, attributes, 31, "GetEffectiveRoleAttributes should now return 31 attributes")
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
	assert.Len(t, rolePerms, 7, "GetEffectiveRolePermissions should return 7 perms")
	roleAttrs, err = ps.GetEffectiveRoleAttributes(ctx, userInfo.GetJob(), userInfo.GetJobGrade())
	require.NoError(t, err, "GetEffectiveRoleAttributes should not return an error")
	attributes = flattenRoleAttributes(roleAttrs)
	assert.Empty(t, attributes, "GetEffectiveRoleAttributes should now return 0 attributes")
}
