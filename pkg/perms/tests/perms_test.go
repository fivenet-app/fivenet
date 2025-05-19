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
	_ = ctx

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

	superuserUI := &userinfo.UserInfo{
		UserId:    testdata.Users[0].UserId,
		Job:       testdata.Users[0].Job,
		JobGrade:  testdata.Users[0].JobGrade,
		SuperUser: true,
	}
	// 1. Superuser can do everything
	can := ps.Can(superuserUI, "CitizenStoreService", "ListCitizens")
	assert.True(t, can, "Superuser has all permissions (CitizenStoreService.ListCitizens)")
	can = ps.Can(superuserUI, "DocStoreService", "DeleteComment")
	assert.True(t, can, "Superuser has all permissions (DocStoreService.DeleteComment)")

	// 2. Non-superuser can do only some things
	ui := &userinfo.UserInfo{
		UserId:   testdata.Users[0].UserId,
		Job:      testdata.Users[0].Job,
		JobGrade: testdata.Users[0].JobGrade,
	}
	can = ps.Can(ui, "CitizenStoreService", "ListCitizens")
	assert.True(t, can, "User should have permission to CitizenStoreService.ListCitizens")
	can = ps.Can(ui, "CitizenStoreService", "SetUserProps")
	assert.True(t, can, "User should have permission to CitizenStoreService.SetUserProps")

	can = ps.Can(ui, "RectorService", "UpdateRolePerms")
	assert.False(t, can, "User should not have permission to RectorService.UpdateRolePerms")

	attributes, err := ps.FlattenRoleAttributes(ui.Job, ui.JobGrade)
	assert.NoError(t, err, "FlattenRoleAttributes should not return an error")
	assert.NotEmpty(t, attributes, "FlattenRoleAttributes should return non-empty attributes")
	// Check if the SetUserProps attributes contain specific permissions
	assert.Contains(t, attributes, "citizenstoreservice-setuserprops-fields-userprops-wanted")
	assert.Contains(t, attributes, "citizenstoreservice-setuserprops-fields-userprops-job")
	// Check if jobgradelist attributes are present
	assert.Contains(t, attributes, "livemapperservice-stream-players-ambulance")
	assert.Contains(t, attributes, "livemapperservice-stream-players-police")
	assert.NotContains(t, attributes, "livemapperservice-stream-players-doj")

	attr, err := ps.AttrJobGradeList(ui, "", "", "")
	assert.NoError(t, err, "AttrJobGradeList should not return an error")
	assert.NotEmpty(t, attr, "AttrJobGradeList should return non-empty attributes")
}
