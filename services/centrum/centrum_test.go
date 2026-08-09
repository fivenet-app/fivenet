package centrum

import (
	"context"
	"database/sql"
	"maps"
	"os"
	"testing"

	centrumres "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/access"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbtracker "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/tracker"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/internal/modules"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
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

func (t *centrumJoinUnitTestTracker) GetUserMapping(userId int32) (*pbtracker.UserMapping, bool, error) {
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

	dbServer := servers.NewDBServer(t, true)
	natsServer := servers.NewNATSServer(t, true)
	trackerStub := newCentrumJoinUnitTestTracker()

	var srv *Server
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			userinfo.RetrieverModule,
			fx.Provide(notifi.New),
			fx.Provide(grpcSrvModule),
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

func insertUnitRowForTest(t *testing.T, db *sql.DB, job, name string) int64 {
	t.Helper()

	res, err := db.Exec(
		`INSERT INTO fivenet_centrum_units (job, name, initials, color) VALUES (?, ?, ?, ?)`,
		job,
		name,
		name[:1],
		"#445566",
	)
	require.NoError(t, err)

	id, err := res.LastInsertId()
	require.NoError(t, err)
	return id
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

func seedDispatchAccessForTest(
	t *testing.T,
	srv *Server,
	ctx context.Context,
	userJob string,
	targetJob string,
) {
	t.Helper()

	_, err := srv.settings.Update(ctx, targetJob, &centrumsettings.Settings{
		Job: targetJob,
		OfferedAccess: &centrumaccess.CentrumAccess{
			Jobs: []*centrumaccess.CentrumJobAccess{
				{
					SourceJob:  userJob,
					Job:        targetJob,
					Access:     centrumaccess.CentrumAccessLevel_CENTRUM_ACCESS_LEVEL_DISPATCH,
					AcceptedAt: timestamp.Now(),
				},
			},
		},
	})
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

	targetUnitID := insertUnitRowForTest(t, db, "police", "Bravo-Target")

	resp, err := srv.JoinUnit(ctx, &pbcentrum.JoinUnitRequest{
		UnitId: &targetUnitID,
	})
	require.ErrorIs(t, err, errorscentrum.ErrUnitPermDenied)
	assert.Nil(t, resp)

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Equal(t, currentUnit.GetId(), mapping.GetUnitId())
	assert.Equal(t, 1, unitAssignmentCountForTest(t, db, currentUnit.GetId(), 1))
	assert.Equal(t, 0, unitAssignmentCountForTest(t, db, targetUnitID, 1))
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
		UserId:   21,
		Job:      "ambulance",
		JobGrade: 20,
	})
	seedDispatchAccessForTest(t, srv, creatorCtx, "ambulance", "police")

	autoDispatch := createDispatchForTest(t, srv, creatorCtx)
	assert.Equal(t, []string{"ambulance"}, autoDispatch.GetJobs().GetJobStrings())

	dispatch := createDispatchForTest(t, srv, creatorCtx, "ambulance", "police")
	assert.ElementsMatch(t, []string{"ambulance", "police"}, dispatch.GetJobs().GetJobStrings())

	denyCtx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   22,
		Job:      "fire",
		JobGrade: 10,
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

	updateCtx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   21,
		Job:      "ambulance",
		JobGrade: 20,
	})
	updateResp, err := srv.UpdateDispatch(updateCtx, &pbcentrum.UpdateDispatchRequest{
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
