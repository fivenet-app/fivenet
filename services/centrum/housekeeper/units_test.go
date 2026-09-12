package housekeeper

import (
	"context"
	"errors"
	"maps"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adhocore/gronx"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/tracker"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	centrummetrics "github.com/fivenet-app/fivenet/v2026/services/centrum/metrics"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

type unitAssignmentWatchSourceStub struct {
	store *store.Store[centrumunits.Unit, *centrumunits.Unit]
	ready chan struct{}
}

func (s *unitAssignmentWatchSourceStub) WatchAll(
	ctx context.Context,
) (store.IKVWatcher[centrumunits.Unit, *centrumunits.Unit], error) {
	close(s.ready)
	return s.store.WatchAll(ctx)
}

func TestDispatchHousekeeperSchedules(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name     string
		schedule string
		want     time.Time
	}{
		{
			name:     "targeted cancellation recovery",
			schedule: cancelOldDispatchesSchedule,
			want:     time.Date(2026, time.January, 1, 0, 5, 0, 0, time.UTC),
		},
		{
			name:     "kv recovery audit",
			schedule: deleteOldDispatchesKVSchedule,
			want:     time.Date(2026, time.January, 1, 2, 15, 0, 0, time.UTC),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			next, err := gronx.NextTickAfter(
				test.schedule,
				time.Date(2026, time.January, 1, 0, 0, 1, 0, time.UTC),
				false,
			)
			require.NoError(t, err)
			assert.Equal(t, test.want, next)
		})
	}
}

func TestDeleteOldDispatchesSelectsOldestFirst(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectQuery("(?s)ORDER BY.*created_at.*ASC.*id.*ASC.*LIMIT \\\\?").
		WithArgs(75).
		WillReturnRows(sqlmock.NewRows([]string{"dispatch_id"}))

	h := &Housekeeper{db: db}
	deleted, err := h.deleteOldDispatches(t.Context())
	require.NoError(t, err)
	assert.Zero(t, deleted)
	require.NoError(t, mock.ExpectationsWereMet())
}

type unitAssignmentsStub struct {
	jobByUser map[int32]bool
	removed   []int32
}

func (s *unitAssignmentsStub) UserInJob(
	_ context.Context,
	_ qrm.DB,
	_ string,
	userID int32,
) (bool, error) {
	return s.jobByUser[userID], nil
}

func (s *unitAssignmentsStub) UpdateUnitAssignments(
	_ context.Context,
	_ string,
	_ *int32,
	_ int64,
	_ []int32,
	toRemove []int32,
) error {
	s.removed = append(s.removed, toRemove...)
	return nil
}

type housekeeperTrackerStub struct {
	mappings        map[int32]*tracker.UserMapping
	onDuty          map[int32]bool
	markers         map[int32]*livemapmarkers.UserMarker
	userMarkerStore *store.Store[livemapmarkers.UserMarker, *livemapmarkers.UserMarker]
	watchReady      chan struct{}
}

func (t *housekeeperTrackerStub) ListTrackedJobs() []string { return nil }

func (t *housekeeperTrackerStub) GetUserMarkerById(id int32) (*livemapmarkers.UserMarker, bool) {
	if t.markers != nil {
		marker, ok := t.markers[id]
		return marker, ok && marker != nil
	}
	return &livemapmarkers.UserMarker{UserId: id, Job: "ambulance"}, true
}

func (t *housekeeperTrackerStub) IsUserOnDuty(userId int32) bool {
	return t.onDuty[userId]
}

func (t *housekeeperTrackerStub) Subscribe(
	ctx context.Context,
) (store.IKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker], error) {
	if t.userMarkerStore != nil {
		watcher, err := t.userMarkerStore.WatchAll(ctx)
		if t.watchReady != nil {
			close(t.watchReady)
		}
		return watcher, err
	}

	return nil, nil
}

func (t *housekeeperTrackerStub) GetFilteredUserMarkers(
	_ *permissionsattributes.JobGradeList,
	_ *pbuserinfo.UserInfo,
) []*livemapmarkers.UserMarker {
	return nil
}

func (t *housekeeperTrackerStub) GetUserMapping(
	userId int32,
) (*tracker.UserMapping, bool, error) {
	mapping, ok := t.mappings[userId]
	if !ok {
		return nil, false, nil
	}

	return mapping, true, nil
}

func (t *housekeeperTrackerStub) SetUserMapping(
	_ context.Context,
	mapping *tracker.UserMapping,
) error {
	if mapping != nil {
		t.mappings[mapping.GetUserId()] = mapping
	}
	return nil
}

func (t *housekeeperTrackerStub) SetUserMappingForUser(
	ctx context.Context,
	userId int32,
	unitId *int64,
) error {
	return t.SetUserMapping(ctx, &tracker.UserMapping{UserId: userId, UnitId: unitId})
}

func (t *housekeeperTrackerStub) UnsetUnitIDForUser(_ context.Context, userId int32) error {
	delete(t.mappings, userId)
	return nil
}

func (t *housekeeperTrackerStub) DeleteUserMapping(_ context.Context, userId int32) error {
	delete(t.mappings, userId)
	return nil
}

func (t *housekeeperTrackerStub) ListUserMappings(
	_ context.Context,
) (map[int32]*tracker.UserMapping, error) {
	out := make(map[int32]*tracker.UserMapping, len(t.mappings))
	maps.Copy(out, t.mappings)
	return out, nil
}

func TestCheckAndUpdateUnitUsersRemovesCrossJobUsers(t *testing.T) {
	t.Parallel()

	unitID := int64(42)
	stub := &unitAssignmentsStub{
		jobByUser: map[int32]bool{
			1: true,
			2: false,
		},
	}
	trackerStub := &housekeeperTrackerStub{
		mappings: map[int32]*tracker.UserMapping{
			1: {UserId: 1, UnitId: &unitID},
			2: {UserId: 2, UnitId: &unitID},
		},
		onDuty: map[int32]bool{
			1: true,
			2: true,
		},
	}

	h := &Housekeeper{
		logger:          zap.NewNop(),
		db:              nil,
		tracker:         trackerStub,
		unitAssignments: stub,
	}
	unit := &centrumunits.Unit{
		Id:    unitID,
		Job:   "ambulance",
		Users: []*centrumunits.UnitAssignment{{UserId: 1}, {UserId: 2}},
	}

	found, removed, err := h.checkAndUpdateUnitUsers(t.Context(), unit)
	require.NoError(t, err)
	assert.Equal(t, []int32{1}, found)
	assert.Equal(t, 1, removed)
	assert.Equal(t, []int32{2}, stub.removed)
}

func TestWatchUnitAssignmentsRecoversAfterCleanupFailure(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	unitStore, err := store.New[centrumunits.Unit](
		ctx,
		zap.NewNop(),
		js,
		"test_unit_assignment_watch",
		store.WithLocks[centrumunits.Unit](nil),
	)
	require.NoError(t, err)

	source := &unitAssignmentWatchSourceStub{
		store: unitStore,
		ready: make(chan struct{}),
	}
	attempts := make(chan int, 2)
	cleanupCalls := 0
	h := &Housekeeper{
		logger:                    zap.NewNop(),
		metrics:                   centrummetrics.Get(),
		unitAssignmentWatchSource: source,
		removeEmptyUnit: func(_ context.Context, _ *centrumunits.Unit) (int, error) {
			cleanupCalls++
			attempts <- cleanupCalls
			if cleanupCalls == 1 {
				return 0, errors.New("temporary cleanup failure")
			}

			return 1, nil
		},
	}

	done := make(chan error, 1)
	go func() {
		done <- h.watchUnitAssignments(ctx)
	}()
	<-source.ready

	emptyUnit := &centrumunits.Unit{Id: 42}
	require.NoError(t, unitStore.Put(ctx, "unit.42", emptyUnit))
	require.Equal(t, 1, <-attempts)

	// Deletion events must not invoke assignment cleanup.
	require.NoError(t, unitStore.Delete(ctx, "unit.42"))
	select {
	case call := <-attempts:
		t.Fatalf("cleanup invoked for delete event, call %d", call)
	case <-time.After(100 * time.Millisecond):
	}

	// A later PUT must still be handled after the prior cleanup failure.
	require.NoError(t, unitStore.Put(ctx, "unit.42", emptyUnit))
	require.Equal(t, 2, <-attempts)

	cancel()
	require.NoError(t, <-done)
}

func TestDispatcherJobsToRemove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		previous userDutyContext
		seen     bool
		current  userDutyContext
		want     []string
	}{
		{
			name:    "initial visible marker",
			current: userDutyContext{job: "ambulance"},
		},
		{
			name:     "job change",
			previous: userDutyContext{job: "ambulance"},
			seen:     true,
			current:  userDutyContext{job: "police"},
			want:     []string{"ambulance"},
		},
		{
			name:     "becomes hidden",
			previous: userDutyContext{job: "ambulance"},
			seen:     true,
			current:  userDutyContext{job: "ambulance", hidden: true},
			want:     []string{"ambulance"},
		},
		{
			name:    "initial hidden marker",
			current: userDutyContext{job: "ambulance", hidden: true},
			want:    []string{"ambulance"},
		},
		{
			name:     "marker deleted after job change",
			previous: userDutyContext{job: "police"},
			seen:     true,
			current:  userDutyContext{job: "ambulance", hidden: true},
			want:     []string{"police", "ambulance"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				test.want,
				dispatcherJobsToRemove(test.previous, test.seen, test.current),
			)
		})
	}
}

func TestUserMarkerKeyContext(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name    string
		key     string
		wantID  int32
		wantJob string
		wantErr bool
	}{
		{
			name:    "standard key",
			key:     "ambulance.3.42",
			wantID:  42,
			wantJob: "ambulance",
		},
		{
			name:    "job containing dots",
			key:     "state.police.3.42",
			wantID:  42,
			wantJob: "state.police",
		},
		{
			name:    "missing user id",
			key:     "ambulance.3",
			wantErr: true,
		},
		{
			name:    "invalid user id",
			key:     "ambulance.3.user",
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userID, context, err := userMarkerKeyContext(test.key)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.wantID, userID)
			assert.Equal(t, test.wantJob, context.job)
			assert.True(t, context.hidden)
		})
	}
}

func TestUserMarkerDeleteReconcilesParsedUserID(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	userMarkerStore, err := store.New[livemapmarkers.UserMarker](
		ctx,
		zap.NewNop(),
		js,
		"test_user_marker_delete_watch",
		store.WithLocks[livemapmarkers.UserMarker](nil),
	)
	require.NoError(t, err)

	trackerStub := &housekeeperTrackerStub{
		mappings:        map[int32]*tracker.UserMapping{},
		onDuty:          map[int32]bool{},
		markers:         map[int32]*livemapmarkers.UserMarker{},
		userMarkerStore: userMarkerStore,
		watchReady:      make(chan struct{}),
	}
	syncedUsers := make(chan int32, 2)
	dispatcherRemovals := make(chan struct {
		job    string
		userID int32
		active bool
	}, 1)
	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
		tracer:  noop.NewTracerProvider().Tracer("test"),
		tracker: trackerStub,
		syncUserUnitMapping: func(_ context.Context, userID int32) error {
			syncedUsers <- userID
			return nil
		},
		setDispatcherState: func(_ context.Context, job string, userID int32, active bool) error {
			dispatcherRemovals <- struct {
				job    string
				userID int32
				active bool
			}{job: job, userID: userID, active: active}
			return nil
		},
	}

	done := make(chan error, 1)
	go func() {
		done <- h.watchUserChanges(ctx, map[int32]userDutyContext{})
	}()
	<-trackerStub.watchReady

	const userID int32 = 42
	require.NoError(t, userMarkerStore.Put(ctx, "ambulance.3.42", &livemapmarkers.UserMarker{
		UserId: userID,
		Job:    "ambulance",
	}))
	require.Equal(t, userID, <-syncedUsers)

	require.NoError(t, userMarkerStore.Delete(ctx, "ambulance.3.42"))
	require.Equal(t, userID, <-syncedUsers)
	removal := <-dispatcherRemovals
	assert.Equal(t, "ambulance", removal.job)
	assert.Equal(t, userID, removal.userID)
	assert.False(t, removal.active)

	cancel()
	require.NoError(t, <-done)
}

func TestUserMarkerDeleteRetainsCanonicalDispatcherState(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	userMarkerStore, err := store.New[livemapmarkers.UserMarker](
		ctx,
		zap.NewNop(),
		js,
		"test_user_marker_stale_job_delete_watch",
		store.WithLocks[livemapmarkers.UserMarker](nil),
	)
	require.NoError(t, err)

	const userID int32 = 42
	trackerStub := &housekeeperTrackerStub{
		mappings: map[int32]*tracker.UserMapping{},
		onDuty:   map[int32]bool{userID: true},
		markers: map[int32]*livemapmarkers.UserMarker{
			userID: {UserId: userID, Job: "ambulance"},
		},
		userMarkerStore: userMarkerStore,
		watchReady:      make(chan struct{}),
	}
	dispatcherRemovals := make(chan string, 2)
	syncedUsers := make(chan int32, 2)
	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
		tracer:  noop.NewTracerProvider().Tracer("test"),
		tracker: trackerStub,
		syncUserUnitMapping: func(_ context.Context, id int32) error {
			syncedUsers <- id
			return nil
		},
		setDispatcherState: func(_ context.Context, job string, _ int32, active bool) error {
			assert.False(t, active)
			dispatcherRemovals <- job
			return nil
		},
	}
	// Seed the stale location key before subscribing; it represents state left
	// behind by an earlier tracker refresh.
	require.NoError(t, userMarkerStore.Put(ctx, "police.4.42", &livemapmarkers.UserMarker{
		UserId: userID,
		Job:    "police",
	}))

	done := make(chan error, 1)
	go func() { done <- h.watchUserChanges(ctx, map[int32]userDutyContext{}) }()
	<-trackerStub.watchReady
	// This post-subscription update is a barrier that also establishes the
	// live canonical context in the watcher.
	require.NoError(t, userMarkerStore.Put(ctx, "ambulance.3.42", &livemapmarkers.UserMarker{
		UserId: userID,
		Job:    "ambulance",
	}))
	require.Equal(t, userID, <-syncedUsers)

	// Cleanup removes a leftover police key after ambulance became the
	// canonical marker.
	require.NoError(t, userMarkerStore.Delete(ctx, "police.4.42"))

	assert.Equal(t, "police", <-dispatcherRemovals)
	select {
	case job := <-dispatcherRemovals:
		t.Fatalf("unexpected dispatcher removal for current marker job %q", job)
	case <-time.After(100 * time.Millisecond):
	}

	cancel()
	require.NoError(t, <-done)
}

func TestUserMarkerDeleteRetainsSameJobDispatcherStateAcrossGradeChange(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	userMarkerStore, err := store.New[livemapmarkers.UserMarker](
		ctx,
		zap.NewNop(),
		js,
		"test_user_marker_stale_grade_delete_watch",
		store.WithLocks[livemapmarkers.UserMarker](nil),
	)
	require.NoError(t, err)

	const userID int32 = 42
	trackerStub := &housekeeperTrackerStub{
		mappings: map[int32]*tracker.UserMapping{},
		onDuty:   map[int32]bool{userID: true},
		markers: map[int32]*livemapmarkers.UserMarker{
			userID: {UserId: userID, Job: "police"},
		},
		userMarkerStore: userMarkerStore,
		watchReady:      make(chan struct{}),
	}
	dispatcherRemovals := make(chan string, 1)
	syncedUsers := make(chan int32, 2)
	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
		tracer:  noop.NewTracerProvider().Tracer("test"),
		tracker: trackerStub,
		syncUserUnitMapping: func(_ context.Context, id int32) error {
			syncedUsers <- id
			return nil
		},
		setDispatcherState: func(_ context.Context, job string, _ int32, _ bool) error {
			dispatcherRemovals <- job
			return nil
		},
	}
	// The old-grade key is stale location state, not a new duty transition.
	require.NoError(t, userMarkerStore.Put(ctx, "police.3.42", &livemapmarkers.UserMarker{
		UserId: userID,
		Job:    "police",
	}))

	done := make(chan error, 1)
	go func() { done <- h.watchUserChanges(ctx, map[int32]userDutyContext{}) }()
	<-trackerStub.watchReady
	require.NoError(t, userMarkerStore.Put(ctx, "police.4.42", &livemapmarkers.UserMarker{
		UserId: userID,
		Job:    "police",
	}))
	require.Equal(t, userID, <-syncedUsers)

	require.NoError(t, userMarkerStore.Delete(ctx, "police.3.42"))

	select {
	case job := <-dispatcherRemovals:
		t.Fatalf("unexpected dispatcher removal for active job %q", job)
	case <-time.After(100 * time.Millisecond):
	}

	cancel()
	require.NoError(t, <-done)
}

func TestCancelDispatchRetainsProjectionWhenStatusPersistenceFails(t *testing.T) {
	t.Parallel()

	deleted := false
	h := &Housekeeper{
		logger:  zap.NewNop(),
		metrics: centrummetrics.Get(),
		updateDispatchStatus: func(
			context.Context,
			int64,
			*centrumdispatches.DispatchStatus,
		) (*centrumdispatches.DispatchStatus, error) {
			return nil, errors.New("database unavailable")
		},
		deleteDispatch: func(context.Context, int64, bool) error {
			deleted = true
			return nil
		},
	}

	h.cancelDispatch(t.Context(), &centrumdispatches.Dispatch{Id: 42})
	assert.False(t, deleted)
}

func TestHandleProjectionCleanup(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name           string
		createdAt      time.Time
		wantReschedule bool
		wantDelete     bool
	}{
		{
			name:           "early tombstone is rescheduled",
			createdAt:      time.Now(),
			wantReschedule: true,
		},
		{
			name:       "expired projection is deleted",
			createdAt:  time.Now().Add(-dispatches.ProjectionTTL),
			wantDelete: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rescheduled := false
			deleted := false
			h := &Housekeeper{
				logger:  zap.NewNop(),
				metrics: centrummetrics.Get(),
				getDispatchProjection: func(context.Context, int64) (*centrumdispatches.Dispatch, error) {
					return &centrumdispatches.Dispatch{
						Id:        42,
						CreatedAt: timestamp.New(test.createdAt),
					}, nil
				},
				scheduleProjectionCleanup: func(context.Context, int64, *timestamp.Timestamp) error {
					rescheduled = true
					return nil
				},
				deleteDispatch: func(context.Context, int64, bool) error {
					deleted = true
					return nil
				},
			}

			h.handleProjectionCleanup(t.Context(), 42)
			assert.Equal(t, test.wantReschedule, rescheduled)
			assert.Equal(t, test.wantDelete, deleted)
		})
	}
}

func TestRecoveryBacklogReporting(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name        string
		items       []int
		limit       int
		wantItems   []int
		wantBacklog bool
	}{
		{
			name:      "within limit",
			items:     []int{1, 2},
			limit:     2,
			wantItems: []int{1, 2},
		},
		{
			name:        "one extra item signals backlog",
			items:       []int{1, 2, 3},
			limit:       2,
			wantItems:   []int{1, 2},
			wantBacklog: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			items, backlog := capWorkItems(test.items, test.limit)
			assert.Equal(t, test.wantItems, items)
			assert.Equal(t, test.wantBacklog, backlog)

			data := &cron.GenericCronData{Attributes: map[string]string{}}
			setBacklogAttribute(data, "backlog_remaining", backlog)
			assert.Equal(
				t,
				strconv.FormatBool(test.wantBacklog),
				data.GetAttribute("backlog_remaining"),
			)
		})
	}
}
