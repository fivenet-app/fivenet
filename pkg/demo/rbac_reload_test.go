package demo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/internal/modules"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestDemoSeedRBACReloadsPermsCache(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	dbServer := servers.NewDBServer(ctx, t, true)
	natsServer := servers.NewNATSServer(t, true)

	db, err := dbServer.DB()
	require.NoError(t, err)

	targetJob := PoliceJob
	require.NoError(t, cleanupDemoRBACForJob(ctx, db, targetJob))

	var loadedPerms perms.Permissions
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			fx.Invoke(func(p perms.Permissions) {
				loadedPerms = p
			}),
		)...,
	)
	require.NotNil(t, app)

	app.RequireStart()
	t.Cleanup(app.RequireStop)

	require.NotNil(t, loadedPerms)

	demoCfg := &config.Config{
		Demo: config.Demo{
			Enabled:   true,
			TargetJob: targetJob,
		},
	}
	d := &Demo{
		logger: zap.NewNop(),
		db:     db,
		cfg:    demoCfg,
		perms:  loadedPerms,
	}
	d.initRandomizers()

	highestGrade, ok, err := d.lookupHighestJobGrade(ctx, targetJob)
	require.NoError(t, err)
	require.True(t, ok, "expected demo target job to have at least one grade")

	require.NoError(t, d.seedDemoCatalog(ctx))

	role, err := loadedPerms.GetRoleByJobAndGrade(ctx, targetJob, highestGrade)
	require.NoError(t, err)
	require.NotNil(t, role)

	rolePerms, err := loadedPerms.GetEffectiveRolePermissions(ctx, role.GetId())
	require.NoError(t, err)
	require.NotEmpty(t, rolePerms)

	userPerms, err := loadedPerms.GetPermissionsOfUser(&userinfo.UserInfo{
		Job:      targetJob,
		JobGrade: highestGrade,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userPerms)

	rolePermGuardSet := map[string]struct{}{}
	for _, p := range rolePerms {
		rolePermGuardSet[p.GetGuardName()] = struct{}{}
	}

	userPermGuardSet := map[string]struct{}{}
	for _, p := range userPerms {
		userPermGuardSet[p.GetGuardName()] = struct{}{}
	}

	for guard := range rolePermGuardSet {
		assert.Contains(t, userPermGuardSet, guard)
	}
}

func cleanupDemoRBACForJob(ctx context.Context, db *sql.DB, job string) error {
	stmts := []string{
		`DELETE rp
		 FROM fivenet_rbac_roles_permissions rp
		 INNER JOIN fivenet_rbac_roles r ON r.id = rp.role_id
		 WHERE r.job = ?`,
		`DELETE ra
		 FROM fivenet_rbac_roles_attrs ra
		 INNER JOIN fivenet_rbac_roles r ON r.id = ra.role_id
		 WHERE r.job = ?`,
		`DELETE FROM fivenet_rbac_job_permissions WHERE job = ?`,
		`DELETE FROM fivenet_rbac_job_attrs WHERE job = ?`,
		`DELETE FROM fivenet_rbac_roles WHERE job = ?`,
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		if _, err := tx.ExecContext(ctx, stmt, job); err != nil {
			return err
		}
	}

	return tx.Commit()
}
