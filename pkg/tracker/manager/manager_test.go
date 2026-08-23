package manager

import (
	"context"
	"database/sql"
	"os"
	"testing"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/internal/modules"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/units"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/protobuf/proto"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func newTrackerManagerForTest(t *testing.T) (*Manager, *sql.DB, *tracker.TestTracker, func()) {
	t.Helper()

	ctx := t.Context()
	dbServer := servers.NewDBServer(ctx, t, true)
	natsServer := servers.NewNATSServer(t, true)

	var manager *Manager
	var trackerStub *tracker.TestTracker
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			fx.Provide(tracker.NewForTests),
			fx.Provide(dispatchers.New),
			fx.Provide(helpers.New),
			fx.Provide(settings.New),
			fx.Provide(units.New),
			fx.Provide(New),
			fx.Provide(access.NewCentrumUnitsSubjectObjectAccess),

			fx.Invoke(func(t tracker.ITracker) {
				trackerStub = t.(*tracker.TestTracker)
			}),
			fx.Invoke(func(m *Manager) {
				manager = m
			}),
		)...,
	)
	require.NotNil(t, app)

	app.RequireStart()

	db, err := dbServer.DB()
	require.NoError(t, err)
	require.NotNil(t, manager)
	require.NotNil(t, trackerStub)

	return manager, db, trackerStub, func() {
		app.RequireStop()
	}
}

func TestRefreshUserLocations(t *testing.T) {
	t.Parallel()

	manager, db, _, stop := newTrackerManagerForTest(t)
	defer stop()

	ctx := t.Context()
	// Run the refreshUserLocations method to make sure the database state has been loaded
	err := manager.refreshUserLocations(ctx, true)
	require.NoError(t, err)

	list := manager.userLocStore.List()
	assert.Empty(t, list)

	// Insert user locations
	require.NoError(
		t,
		insertCitizenLocations(
			ctx,
			db,
			1,
			"ambulance",
			3,
			1.0,
			1.0,
			false,
		),
	)
	require.NoError(
		t,
		insertCitizenLocations(
			ctx,
			db,
			2,
			"ambulance",
			3,
			1.0,
			1.0,
			true,
		),
	)

	require.NoError(t, manager.refreshUserLocations(ctx, false))

	// The cache refresh is synchronous here, so the store should be populated immediately.
	list = manager.userLocStore.List()
	assert.Len(t, list, 2)

	user1, err := manager.userLocStore.Get(userMarkerKey(int32(1), "ambulance", 3))
	require.NoError(t, err)
	assert.NotNil(t, user1)
	assert.InEpsilon(t, 1.0, user1.GetX(), 0.0001)
	assert.InEpsilon(t, 1.0, user1.GetY(), 0.0001)

	user2, err := manager.userLocStore.Get(userMarkerKey(int32(2), "ambulance", 3))
	require.NoError(t, err)
	assert.NotNil(t, user2)
	assert.InEpsilon(t, 1.0, user2.GetX(), 0.0001)
	assert.InEpsilon(t, 1.0, user2.GetY(), 0.0001)

	staleUser1 := proto.Clone(user1).(*livemapmarkers.UserMarker)
	staleUser1.SetJob("police")
	staleUser1.SetJobGrade(4)
	require.NoError(
		t,
		manager.userLocStore.Put(ctx, userMarkerKey(int32(1), "police", 4), staleUser1),
	)

	removed, err := manager.cleanupUserIDs(ctx, map[int32]any{
		1: nil,
		2: nil,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, removed)

	user1, err = manager.userLocStore.Get(userMarkerKey(int32(1), "ambulance", 3))
	require.NoError(t, err)
	assert.NotNil(t, user1)

	staleUser1, err = manager.userLocStore.Get(userMarkerKey(int32(1), "police", 4))
	require.Error(t, err)
	assert.Nil(t, staleUser1)

	// Update user location (no event is sent for updates)
	require.NoError(
		t,
		insertCitizenLocations(
			ctx,
			db,
			2,
			"ambulance",
			3,
			5.0,
			5.0,
			true,
		),
	)

	require.NoError(t, manager.refreshUserLocations(ctx, false))

	user1, err = manager.userLocStore.Get(userMarkerKey(int32(1), "ambulance", 3))
	require.NoError(t, err)
	assert.NotNil(t, user1)
	assert.InEpsilon(t, 1.0, user1.GetX(), 0.0001)
	assert.InEpsilon(t, 1.0, user1.GetY(), 0.0001)

	user2, err = manager.userLocStore.Get(userMarkerKey(int32(2), "ambulance", 3))
	require.NoError(t, err)
	assert.NotNil(t, user2)
	assert.InEpsilon(t, 5.0, user2.GetX(), 0.0001)
	assert.InEpsilon(t, 5.0, user2.GetY(), 0.0001)

	require.NoError(t, removeUserLocations(ctx, db))

	require.NoError(t, manager.refreshUserLocations(ctx, false))

	list = manager.userLocStore.List()
	assert.Empty(t, list)

	user1, err = manager.userLocStore.Get(userMarkerKey(int32(1), "ambulance", 3))
	require.Error(t, err)
	assert.Nil(t, user1)

	user2, err = manager.userLocStore.Get(userMarkerKey(int32(2), "ambulance", 3))
	require.Error(t, err)
	assert.Nil(t, user2)

	// Check that a snapshot entry exists in the KeyValue store
	kv, err := manager.js.KeyValue(ctx, tracker.BucketUserLoc)
	require.NoError(t, err)
	assert.NotNil(t, kv)
}

func TestRefreshUserLocationsRemovesOldJobDispatcherOnJobChange(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	manager, db, trackerStub, stop := newTrackerManagerForTest(t)
	defer stop()

	require.NoError(t, insertCitizenLocations(ctx, db, 1, "ambulance", 3, 1.0, 1.0, false))
	require.NoError(t, manager.refreshUserLocations(ctx, true))
	trackerStub.SeedUserMarker(&livemapmarkers.UserMarker{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: proto.Int32(3),
		Hidden:   false,
	})

	marker, err := manager.userByIDStore.Get(tracker.UserIdKey(1))
	require.NoError(t, err)
	require.NotNil(t, marker)
	require.Equal(t, "ambulance", marker.GetJob())

	require.NoError(t, manager.dispatchers.SetUserState(ctx, "ambulance", 1, true))

	dispatchersForAmbulance, err := manager.dispatchers.Get(ctx, "ambulance")
	require.NoError(t, err)
	require.Len(t, dispatchersForAmbulance.GetDispatchers(), 1)

	require.NoError(t, dbUpdateUserJob(ctx, db, 1, "army"))
	require.NoError(t, insertCitizenLocations(ctx, db, 1, "army", 3, 1.0, 1.0, false))

	require.NoError(t, manager.refreshUserLocations(ctx, false))

	marker, err = manager.userByIDStore.Get(tracker.UserIdKey(1))
	require.NoError(t, err)
	require.NotNil(t, marker)
	require.Equal(t, "army", marker.GetJob())

	dispatchersForAmbulance, err = manager.dispatchers.Get(ctx, "ambulance")
	require.NoError(t, err)
	assert.Empty(t, dispatchersForAmbulance.GetDispatchers())

	oldKey := userMarkerKey(1, "ambulance", 3)
	_, err = manager.userLocStore.Get(oldKey)
	require.Error(t, err)

	newKey := userMarkerKey(1, "army", 3)
	refreshedMarker, err := manager.userLocStore.Get(newKey)
	require.NoError(t, err)
	require.NotNil(t, refreshedMarker)
	require.Equal(t, "army", refreshedMarker.GetJob())
}

func TestCleanupUserIDsDeletesStaleLocationKeyWhenMarkerMissing(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	manager, db, _, stop := newTrackerManagerForTest(t)
	defer stop()

	staleUser := &livemapmarkers.UserMarker{
		UserId:   42,
		Job:      "police",
		JobGrade: new(int32(4)),
	}
	staleKey := userMarkerKey(staleUser.GetUserId(), staleUser.GetJob(), staleUser.GetJobGrade())
	require.NoError(t, manager.userLocStore.Put(ctx, staleKey, staleUser))

	require.NoError(t, dbInsertTestUser(ctx, db, 42, "police"))
	require.NoError(t, dbInsertDispatcher(ctx, db, "police", 42))
	require.NoError(t, manager.dispatchers.LoadFromDB(ctx, "police"))

	removed, err := manager.cleanupUserIDs(ctx, map[int32]any{
		staleUser.GetUserId(): nil,
	})
	require.NoError(t, err)
	require.Equal(t, 1, removed)

	got, err := manager.userLocStore.Get(staleKey)
	require.Error(t, err)
	require.Nil(t, got)

	dispatchersForPolice, err := manager.dispatchers.Get(ctx, "police")
	require.NoError(t, err)
	assert.Empty(t, dispatchersForPolice.GetDispatchers())
}

func insertCitizenLocations(
	ctx context.Context,
	db *sql.DB,
	userId int32,
	job string,
	grade int32,
	x float64,
	y float64,
	hidden bool,
) error {
	stmt := tLocs.
		INSERT(
			tLocs.UserID,
			tLocs.Job,
			tLocs.JobGrade,
			tLocs.X,
			tLocs.Y,
			tLocs.Hidden,
		).
		VALUES(
			userId,
			job,
			grade,
			x,
			y,
			hidden,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tLocs.Job.SET(mysql.StringExp(mysql.Raw("VALUES(`job`)"))),
			tLocs.JobGrade.SET(mysql.IntExp(mysql.Raw("VALUES(`job_grade`)"))),
			tLocs.X.SET(mysql.FloatExp(mysql.Raw("VALUES(`x`)"))),
			tLocs.Y.SET(mysql.FloatExp(mysql.Raw("VALUES(`y`)"))),
			tLocs.Hidden.SET(mysql.BoolExp(mysql.Raw("VALUES(`hidden`)"))),
		)

	_, err := stmt.ExecContext(ctx, db)

	return err
}

func dbUpdateUserJob(ctx context.Context, db *sql.DB, userID int32, job string) error {
	_, err := db.ExecContext(
		ctx,
		"UPDATE fivenet_user SET job = ? WHERE id = ?",
		job,
		userID,
	)

	return err
}

func dbInsertTestUser(ctx context.Context, db *sql.DB, userID int32, job string) error {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO fivenet_user (id, license, identifier, job, job_grade, firstname, lastname, dateofbirth, sex, height, phone_number, disabled, visum, playtime, created_at, updated_at) VALUES (?, '', CONCAT('test:', ?), 'user', ?, 1, 'Test', 'User', '01.01.2000', 'm', 180, '0000000', 0, 0, 0, CURRENT_TIMESTAMP(3), CURRENT_TIMESTAMP(3))",
		userID,
		userID,
		job,
	)

	return err
}

func dbInsertDispatcher(ctx context.Context, db *sql.DB, job string, userID int32) error {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO fivenet_centrum_dispatchers (job, user_id) VALUES (?, ?)",
		job,
		userID,
	)

	return err
}

func removeUserLocations(ctx context.Context, db *sql.DB) error {
	stmt := tLocs.
		DELETE().
		WHERE(tLocs.UserID.IS_NOT_NULL().OR(tLocs.UserID.IS_NULL())).
		LIMIT(10000)

	_, err := stmt.ExecContext(ctx, db)

	return err
}
