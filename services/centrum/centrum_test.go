package centrum

import (
	"context"
	"database/sql"
	"maps"
	"os"
	"slices"
	"testing"

	centrumres "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/access"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbtracker "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/tracker"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/internal/modules"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	trackerpkg "github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/helpers"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/settings"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/units"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

type centrumJoinUnitTestTracker struct {
	markers  map[int32]*livemapmarkers.UserMarker
	mappings map[int32]*pbtracker.UserMapping
}

func newCentrumJoinUnitTestTracker() *centrumJoinUnitTestTracker {
	return &centrumJoinUnitTestTracker{
		markers: map[int32]*livemapmarkers.UserMarker{
			1: {
				UserId: 1,
				Hidden: false,
			},
		},
		mappings: map[int32]*pbtracker.UserMapping{},
	}
}

func (t *centrumJoinUnitTestTracker) ListTrackedJobs() []string {
	return nil
}

func (t *centrumJoinUnitTestTracker) GetUserMarkerById(
	id int32,
) (*livemapmarkers.UserMarker, bool) {
	marker, ok := t.markers[id]
	return marker, ok
}

func (t *centrumJoinUnitTestTracker) IsUserOnDuty(userId int32) bool {
	_, ok := t.markers[userId]
	return ok
}

func (t *centrumJoinUnitTestTracker) Subscribe(
	_ context.Context,
) (store.IKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker], error) {
	return nil, nil
}

func (t *centrumJoinUnitTestTracker) GetFilteredUserMarkers(
	_ *permissionsattributes.JobGradeList,
	_ *pbuserinfo.UserInfo,
) []*livemapmarkers.UserMarker {
	return nil
}

func (t *centrumJoinUnitTestTracker) GetUserMapping(
	userId int32,
) (*pbtracker.UserMapping, bool, error) {
	mapping, ok := t.mappings[userId]
	if !ok {
		return nil, false, nil
	}

	return mapping, true, nil
}

func (t *centrumJoinUnitTestTracker) SetUserMapping(
	ctx context.Context,
	mapping *pbtracker.UserMapping,
) error {
	_ = ctx
	if mapping == nil {
		return nil
	}

	t.mappings[mapping.GetUserId()] = mapping
	return nil
}

func (t *centrumJoinUnitTestTracker) SetUserMappingForUser(
	ctx context.Context,
	userId int32,
	unitId *int64,
) error {
	return t.SetUserMapping(ctx, &pbtracker.UserMapping{
		UserId: userId,
		UnitId: unitId,
	})
}

func (t *centrumJoinUnitTestTracker) UnsetUnitIDForUser(ctx context.Context, userId int32) error {
	return t.SetUserMappingForUser(ctx, userId, nil)
}

func (t *centrumJoinUnitTestTracker) DeleteUserMapping(_ context.Context, userId int32) error {
	delete(t.mappings, userId)
	return nil
}

func (t *centrumJoinUnitTestTracker) ListUserMappings(
	_ context.Context,
) (map[int32]*pbtracker.UserMapping, error) {
	out := make(map[int32]*pbtracker.UserMapping, len(t.mappings))
	maps.Copy(out, t.mappings)

	return out, nil
}

func newCentrumJoinUnitTestServer(
	t *testing.T,
) (*Server, *sql.DB, *centrumJoinUnitTestTracker) {
	t.Helper()

	ctx := t.Context()

	_, grpcSrvModule, err := modules.TestGRPCServer(ctx)
	require.NoError(t, err)

	dbServer := servers.NewDBServer(ctx, t, true)
	natsServer := servers.NewNATSServer(t, true)
	trackerStub := newCentrumJoinUnitTestTracker()
	jobsCatalog := mstlystcdata.NewTestJobs(map[string]*jobs.Job{
		"ambulance": {
			Name:  "ambulance",
			Label: "LSMD",
			Grades: []*jobs.JobGrade{
				{
					JobName: func() *string {
						jobName := "ambulance"
						return &jobName
					}(),
					Grade: 1,
					Label: "Rank 1",
				},
			},
		},
		"police": {
			Name:  "police",
			Label: "LSPD",
			Grades: []*jobs.JobGrade{
				{
					JobName: func() *string {
						jobName := "police"
						return &jobName
					}(),
					Grade: 1,
					Label: "Rank 1",
				},
			},
		},
		"doj": {
			Name:  "doj",
			Label: "DOJ",
			Grades: []*jobs.JobGrade{
				{
					JobName: func() *string {
						jobName := "doj"
						return &jobName
					}(),
					Grade: 1,
					Label: "Rank 1",
				},
			},
		},
		"unemployed": {
			Name:  "unemployed",
			Label: "Unemployed",
			Grades: []*jobs.JobGrade{
				{
					JobName: func() *string {
						jobName := "unemployed"
						return &jobName
					}(),
					Grade: 1,
					Label: "Rank 1",
				},
			},
		},
	})

	var srv *Server
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			userinfo.RetrieverModule,
			fx.Provide(notifi.New),
			fx.Provide(grpcSrvModule),
			fx.Decorate(func(_ mstlystcdata.IJobs) mstlystcdata.IJobs { return jobsCatalog }),
			fx.Provide(func() trackerpkg.ITracker {
				return trackerStub
			}),
			fx.Provide(helpers.New),
			fx.Provide(settings.New),
			fx.Provide(dispatchers.New),
			fx.Provide(units.New),
			fx.Provide(dispatches.New),
			fx.Provide(func(p Params) Result {
				r := NewServer(p)
				srv = r.Server
				return r
			}),
			fx.Invoke(func(*grpc.Server) {}),
		)...,
	)
	app.RequireStart()

	t.Cleanup(func() {
		app.RequireStop()
		dbServer.Stop()
	})

	db, err := dbServer.DB()
	require.NoError(t, err)
	require.NotNil(t, srv)

	return srv, db, trackerStub
}

func createUnitForTest(
	t *testing.T,
	srv *Server,
	ctx context.Context,
	name string,
) *centrumunits.Unit {
	t.Helper()

	resp, err := srv.CreateOrUpdateUnit(ctx, &pbcentrum.CreateOrUpdateUnitRequest{
		Unit: &centrumunits.Unit{
			Name:     name,
			Initials: name[:1],
			Color:    "#112233",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetUnit())

	return resp.GetUnit()
}

func seedAssignmentForTest(
	t *testing.T,
	db *sql.DB,
	tracker *centrumJoinUnitTestTracker,
	unitID int64,
	userID int32,
) {
	t.Helper()

	_, err := db.Exec(
		`INSERT INTO fivenet_centrum_units_users (unit_id, user_id) VALUES (?, ?)`,
		unitID,
		userID,
	)
	require.NoError(t, err)

	require.NoError(t, tracker.SetUserMappingForUser(t.Context(), userID, &unitID))
}

func unitAssignmentCountForTest(t *testing.T, db *sql.DB, unitID int64, userID int32) int {
	t.Helper()

	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM fivenet_centrum_units_users WHERE unit_id = ? AND user_id = ?`,
		unitID,
		userID,
	).Scan(&count)
	require.NoError(t, err)

	return count
}

func insertAssignmentRowForTest(t *testing.T, db *sql.DB, unitID int64, userID int32) {
	t.Helper()

	_, err := db.Exec(
		`INSERT INTO fivenet_centrum_units_users (unit_id, user_id) VALUES (?, ?)`,
		unitID,
		userID,
	)
	require.NoError(t, err)
}

func moveAssignmentRowForTest(t *testing.T, db *sql.DB, unitID int64, userID int32) {
	t.Helper()

	_, err := db.Exec(
		`UPDATE fivenet_centrum_units_users SET unit_id = ? WHERE user_id = ?`,
		unitID,
		userID,
	)
	require.NoError(t, err)
}

func deleteAssignmentRowForTest(t *testing.T, db *sql.DB, unitID int64, userID int32) {
	t.Helper()

	_, err := db.Exec(
		`DELETE FROM fivenet_centrum_units_users WHERE unit_id = ? AND user_id = ?`,
		unitID,
		userID,
	)
	require.NoError(t, err)
}

func assertUnitCacheHasUser(
	t *testing.T,
	srv *Server,
	ctx context.Context,
	unitID int64,
	userID int32,
	want bool,
) {
	t.Helper()

	unit, err := srv.units.Get(ctx, unitID)
	require.NoError(t, err)
	require.NotNil(t, unit)

	assert.Equal(
		t,
		want,
		slices.ContainsFunc(unit.GetUsers(), func(in *centrumunits.UnitAssignment) bool {
			return in.GetUserId() == userID
		}),
	)
}

func seedDispatchAccessForTest(
	t *testing.T,
	srv *Server,
	ctx context.Context,
	userJob string,
	targetJob string,
) {
	t.Helper()

	acceptedAt := timestamp.Now()

	// Update settings to allow userJob <-> targetJob DISPATCH access to each other
	_, err := srv.settings.Update(ctx, userJob, &centrumsettings.Settings{
		Job:     userJob,
		Enabled: true,
		Access: &centrumaccess.CentrumAccess{
			Jobs: []*centrumaccess.CentrumJobAccess{
				{
					SourceJob: userJob,
					Job:       targetJob,
					Access:    centrumaccess.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_DISPATCH,
				},
			},
		},
	})
	require.NoError(t, err)
	_, err = srv.settings.Update(ctx, targetJob, &centrumsettings.Settings{
		Job:     targetJob,
		Enabled: true,
		Access: &centrumaccess.CentrumAccess{
			Jobs: []*centrumaccess.CentrumJobAccess{
				{
					SourceJob: targetJob,
					Job:       userJob,
					Access:    centrumaccess.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_DISPATCH,
				},
			},
		},
	})
	require.NoError(t, err)

	// Accept the userJob/targetJob access now
	targetJobSettings, err := srv.settings.Get(ctx, targetJob)
	require.NoError(t, err)
	require.NotNil(t, targetJobSettings)
	targetJobSettings.Access.Jobs[0].AcceptedAt = acceptedAt
	_, err = srv.settings.Update(ctx, targetJob, targetJobSettings)
	require.NoError(t, err)

	userJobSettings, err := srv.settings.Get(ctx, userJob)
	require.NoError(t, err)
	require.NotNil(t, userJobSettings)
	userJobSettings.Access.Jobs[0].AcceptedAt = acceptedAt
	_, err = srv.settings.Update(ctx, userJob, userJobSettings)
	require.NoError(t, err)
}

func createDispatchForTest(
	t *testing.T,
	srv *Server,
	ctx context.Context,
	jobs ...string,
) *centrumdispatches.Dispatch {
	t.Helper()

	dispatchJobs := &centrumres.JobList{}
	if len(jobs) > 0 {
		dispatchJobs.Jobs = make([]*centrumres.JobListEntry, 0, len(jobs))
		for _, job := range jobs {
			dispatchJobs.Jobs = append(dispatchJobs.Jobs, &centrumres.JobListEntry{Name: job})
		}
	}

	resp, err := srv.CreateDispatch(ctx, &pbcentrum.CreateDispatchRequest{
		Dispatch: &centrumdispatches.Dispatch{
			Message: "initial dispatch",
			Jobs:    dispatchJobs,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetDispatch())

	return resp.GetDispatch()
}

func TestSyncUserUnitMappingRepairsMissingTrackerMapping(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Sync")
	insertAssignmentRowForTest(t, db, unit.GetId(), 1)

	_, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.False(t, ok)

	require.NoError(t, srv.units.SyncUserUnitMapping(ctx, 1))

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Equal(t, unit.GetId(), mapping.GetUnitId())
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, true)
}

func TestSyncUserUnitMappingClearsRemovedUser(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Remove")
	seedAssignmentForTest(t, db, trackerStub, unit.GetId(), 1)
	require.NoError(t, srv.units.SyncUnitMembership(ctx, unit.GetId()))
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, true)

	deleteAssignmentRowForTest(t, db, unit.GetId(), 1)

	require.NoError(t, srv.units.SyncUserUnitMapping(ctx, 1))

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Nil(t, mapping.UnitId)
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, false)
}

func TestSyncUserUnitMappingRefreshesOldAndNewUnitOnMove(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	oldUnit := createUnitForTest(t, srv, ctx, "Alpha-Old")
	newUnit := createUnitForTest(t, srv, ctx, "Bravo-New")
	seedAssignmentForTest(t, db, trackerStub, oldUnit.GetId(), 1)
	staleOldUnitID := oldUnit.GetId()
	require.NoError(t, trackerStub.SetUserMappingForUser(ctx, 2, &staleOldUnitID))
	require.NoError(t, srv.units.SyncUnitMembership(ctx, oldUnit.GetId()))
	assertUnitCacheHasUser(t, srv, ctx, oldUnit.GetId(), 1, true)
	assertUnitCacheHasUser(t, srv, ctx, newUnit.GetId(), 1, false)

	moveAssignmentRowForTest(t, db, newUnit.GetId(), 1)

	require.NoError(t, srv.units.SyncUserUnitMapping(ctx, 1))

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Equal(t, newUnit.GetId(), mapping.GetUnitId())
	assertUnitCacheHasUser(t, srv, ctx, oldUnit.GetId(), 1, false)
	assertUnitCacheHasUser(t, srv, ctx, newUnit.GetId(), 1, true)

	_, ok, err = trackerStub.GetUserMapping(2)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestSyncUnitMembershipClearsStaleMappingForMissingUnit(t *testing.T) {
	t.Parallel()

	srv, _, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	missingUnitID := int64(999_999)
	require.NoError(t, trackerStub.SetUserMappingForUser(ctx, 1, &missingUnitID))

	require.NoError(t, srv.units.SyncUnitMembership(ctx, missingUnitID))

	_, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestSyncUnitMembershipPreservesValidMappingAndClearsStaleOnes(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Membership")
	insertAssignmentRowForTest(t, db, unit.GetId(), 1)

	staleUnitID := unit.GetId()
	require.NoError(t, trackerStub.SetUserMappingForUser(ctx, 2, &staleUnitID))

	require.NoError(t, srv.units.SyncUnitMembership(ctx, unit.GetId()))

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Equal(t, unit.GetId(), mapping.GetUnitId())
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, true)

	_, ok, err = trackerStub.GetUserMapping(2)
	require.NoError(t, err)
	assert.False(t, ok)
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 2, false)
}

func TestSyncUserUnitMappingDeletesMappingForOffDutyUserWithoutAssignment(t *testing.T) {
	t.Parallel()

	srv, _, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	delete(trackerStub.markers, 1)
	unitID := int64(123)
	require.NoError(t, trackerStub.SetUserMappingForUser(ctx, 1, &unitID))

	require.NoError(t, srv.units.SyncUserUnitMapping(ctx, 1))

	_, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestRemoveUnitAssignmentsDoesNotWriteTrackerMapping(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Mapping-Delete")
	seedAssignmentForTest(t, db, trackerStub, unit.GetId(), 1)
	require.NoError(t, srv.units.SyncUnitMembership(ctx, unit.GetId()))

	delete(trackerStub.mappings, 1)

	creatorID := int32(1)
	require.NoError(t, srv.units.RemoveUnitAssignments(
		ctx,
		"",
		&creatorID,
		unit.GetId(),
		[]int32{1},
	))

	_, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, 0, unitAssignmentCountForTest(t, db, unit.GetId(), 1))
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, false)
}

func TestJoinUnitKeepsCurrentUnitWhenTargetValidationFails(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	currentUnit := createUnitForTest(t, srv, ctx, "Alpha-Current")
	seedAssignmentForTest(t, db, trackerStub, currentUnit.GetId(), 1)

	policeCtx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   2,
		Job:      "police",
		JobGrade: 17,
	})
	targetUnit := createUnitForTest(t, srv, policeCtx, "Bravo-Target")

	resp, err := srv.JoinUnit(ctx, &pbcentrum.JoinUnitRequest{
		UnitId: &targetUnit.Id,
	})
	require.ErrorIs(t, err, errorscentrum.ErrUnitPermDenied)
	assert.Nil(t, resp)

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Equal(t, currentUnit.GetId(), mapping.GetUnitId())
	assert.Equal(t, 1, unitAssignmentCountForTest(t, db, currentUnit.GetId(), 1))
	assert.Equal(t, 0, unitAssignmentCountForTest(t, db, targetUnit.GetId(), 1))
}

func TestJoinUnitMovesUserAfterValidationSucceeds(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	currentUnit := createUnitForTest(t, srv, ctx, "Alpha-Current")
	seedAssignmentForTest(t, db, trackerStub, currentUnit.GetId(), 1)

	targetUnit := createUnitForTest(t, srv, ctx, "Bravo-Target")

	resp, err := srv.JoinUnit(ctx, &pbcentrum.JoinUnitRequest{
		UnitId: &targetUnit.Id,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetUnit())
	assert.Equal(t, targetUnit.GetId(), resp.GetUnit().GetId())

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Equal(t, targetUnit.GetId(), mapping.GetUnitId())
	assert.Equal(t, 0, unitAssignmentCountForTest(t, db, currentUnit.GetId(), 1))
	assert.Equal(t, 1, unitAssignmentCountForTest(t, db, targetUnit.GetId(), 1))

	sameResp, err := srv.JoinUnit(ctx, &pbcentrum.JoinUnitRequest{
		UnitId: &targetUnit.Id,
	})
	require.NoError(t, err)
	require.NotNil(t, sameResp)
	require.NotNil(t, sameResp.GetUnit())
	assert.Equal(t, targetUnit.GetId(), sameResp.GetUnit().GetId())
}

func TestJoinUnitLeavePathRemovesCurrentUnit(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	currentUnit := createUnitForTest(t, srv, ctx, "Alpha-Current")
	seedAssignmentForTest(t, db, trackerStub, currentUnit.GetId(), 1)

	resp, err := srv.JoinUnit(ctx, &pbcentrum.JoinUnitRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Nil(t, resp.GetUnit())

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Nil(t, mapping.UnitId)
	assert.Equal(t, 0, unitAssignmentCountForTest(t, db, currentUnit.GetId(), 1))
}

func TestJoinUnitOffDutyDeletesStaleMapping(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Off-Duty")
	seedAssignmentForTest(t, db, trackerStub, unit.GetId(), 1)
	deleteAssignmentRowForTest(t, db, unit.GetId(), 1)
	delete(trackerStub.markers, 1)

	resp, err := srv.JoinUnit(ctx, &pbcentrum.JoinUnitRequest{})
	require.ErrorIs(t, err, errorscentrum.ErrNotOnDuty)
	assert.Nil(t, resp)

	_, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestCreateDispatchRejectsUnauthorizedJobs(t *testing.T) {
	t.Parallel()

	srv, _, _ := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   11,
		Job:      "fire",
		JobGrade: 10,
	})

	_, err := srv.CreateDispatch(ctx, &pbcentrum.CreateDispatchRequest{
		Dispatch: &centrumdispatches.Dispatch{
			Message: "blocked",
			Jobs: &centrumres.JobList{
				Jobs: []*centrumres.JobListEntry{
					{Name: "fire"},
					{Name: "police"},
				},
			},
		},
	})
	require.ErrorIs(t, err, errorscentrum.ErrDispatchJobPermDenied)
}

func TestCreateAndUpdateDispatchAuthorization(t *testing.T) {
	t.Parallel()

	srv, _, _ := newCentrumJoinUnitTestServer(t)

	creatorCtx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 20,
	})
	seedDispatchAccessForTest(t, srv, creatorCtx, "ambulance", "police")

	autoDispatch := createDispatchForTest(t, srv, creatorCtx)
	assert.Equal(t, []string{"ambulance"}, autoDispatch.GetJobs().GetJobStrings())

	dispatch := createDispatchForTest(t, srv, creatorCtx, "ambulance", "police")
	assert.ElementsMatch(t, []string{"ambulance", "police"}, dispatch.GetJobs().GetJobStrings())

	denyCtx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   3,
		Job:      "doj",
		JobGrade: 16,
	})
	_, err := srv.UpdateDispatch(denyCtx, &pbcentrum.UpdateDispatchRequest{
		Dispatch: &centrumdispatches.Dispatch{
			Id:      dispatch.GetId(),
			Message: "unauthorized change",
			Jobs: &centrumres.JobList{
				Jobs: []*centrumres.JobListEntry{
					{Name: "fire"},
				},
			},
		},
	})
	require.ErrorIs(t, err, errorscentrum.ErrNotPartOfDispatch)
}

func TestUpdateDispatchAllowsDispatcher(t *testing.T) {
	t.Parallel()

	srv, _, _ := newCentrumJoinUnitTestServer(t)

	creatorCtx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 20,
	})
	seedDispatchAccessForTest(t, srv, creatorCtx, "ambulance", "police")

	dispatch := createDispatchForTest(t, srv, creatorCtx, "ambulance", "police")

	_, err := srv.TakeControl(creatorCtx, &pbcentrum.TakeControlRequest{
		Signon: true,
	})
	require.NoError(t, err)

	updateResp, err := srv.UpdateDispatch(creatorCtx, &pbcentrum.UpdateDispatchRequest{
		Dispatch: &centrumdispatches.Dispatch{
			Id:      dispatch.GetId(),
			Message: "authorized change",
			Jobs:    dispatch.GetJobs(),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.NotNil(t, updateResp.GetDispatch())
	assert.Equal(t, "authorized change", updateResp.GetDispatch().GetMessage())
}

func TestUpdateDispatchStatusAllowsMissingTrackerMapping(t *testing.T) {
	t.Parallel()

	srv, _, trackerStub := newCentrumJoinUnitTestServer(t)

	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 20,
	})
	dispatch := createDispatchForTest(t, srv, ctx)
	jobGrade := int32(20)
	trackerStub.markers[1] = &livemapmarkers.UserMarker{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: &jobGrade,
		Hidden:   false,
	}
	require.NoError(t, srv.dispatchers.SetUserState(ctx, "ambulance", 1, true))
	delete(trackerStub.mappings, 1)

	resp, err := srv.UpdateDispatchStatus(ctx, &pbcentrum.UpdateDispatchStatusRequest{
		DispatchId: dispatch.GetId(),
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_EN_ROUTE,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
}
