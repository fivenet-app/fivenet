package centrum

import (
	"context"
	"database/sql"
	"errors"
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
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	trackerpkg "github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
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
	markers    map[int32]*livemapmarkers.UserMarker
	mappings   map[int32]*pbtracker.UserMapping
	mappingErr error
}

func newCentrumJoinUnitTestTracker() *centrumJoinUnitTestTracker {
	return &centrumJoinUnitTestTracker{
		markers: map[int32]*livemapmarkers.UserMarker{
			1: {
				UserId: 1,
				Job:    "ambulance",
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
	marker, ok := t.markers[userId]
	return ok && marker != nil && !marker.GetHidden()
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
	if t.mappingErr != nil {
		return nil, false, t.mappingErr
	}

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
					JobName: new("ambulance"),
					Grade:   1,
					Label:   "Rank 1",
				},
			},
		},
		"police": {
			Name:  "police",
			Label: "LSPD",
			Grades: []*jobs.JobGrade{
				{
					JobName: new("police"),
					Grade:   1,
					Label:   "Rank 1",
				},
			},
		},
		"doj": {
			Name:  "doj",
			Label: "DOJ",
			Grades: []*jobs.JobGrade{
				{
					JobName: new("doj"),
					Grade:   1,
					Label:   "Rank 1",
				},
			},
		},
		"unemployed": {
			Name:  "unemployed",
			Label: "Unemployed",
			Grades: []*jobs.JobGrade{
				{
					JobName: new("unemployed"),
					Grade:   1,
					Label:   "Rank 1",
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
			fx.Provide(access.NewCentrumUnitsSubjectObjectAccess),
			fx.Provide(access.NewJobGroupsSubjectObjectAccess),
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

func TestSendLatestStateFailsWhenOwnUnitMappingCannotBeRead(t *testing.T) {
	t.Parallel()

	srv, _, trackerStub := newCentrumJoinUnitTestServer(t)
	trackerStub.mappingErr = errors.New("tracker unavailable")

	err := srv.sendLatestState(
		t.Context(),
		&testCentrumStreamServer{ctx: t.Context()},
		&pbuserinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
		&centrumsettings.EffectiveAccess{},
		nil,
	)
	require.ErrorContains(t, err, "failed to get own unit mapping")
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

func upsertUserJobForTest(
	t *testing.T,
	db *sql.DB,
	userID int32,
	job string,
	grade int32,
) {
	t.Helper()

	_, err := db.Exec(
		`INSERT INTO fivenet_user_jobs (user_id, job, grade, is_primary) VALUES (?, ?, ?, TRUE)
		 ON DUPLICATE KEY UPDATE job = VALUES(job), grade = VALUES(grade), is_primary = VALUES(is_primary)`,
		userID,
		job,
		grade,
	)
	require.NoError(t, err)
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
