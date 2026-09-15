package centrum

import (
	"fmt"
	"slices"
	"testing"
	"time"

	centrumres "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	require.True(t, resp.GetUpdated())
	require.NotNil(t, resp.GetStatus())
	assert.Positive(t, resp.GetStatus().GetId())
}

func TestDispatchAssignmentRemovalUpdatesDurableStateAndProjection(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})
	trackerStub.markers[2] = &livemapmarkers.UserMarker{UserId: 2, Job: "ambulance"}
	upsertUserJobForTest(t, db, 1, "ambulance", 17)
	upsertUserJobForTest(t, db, 2, "ambulance", 17)

	firstUnit := createUnitForTest(t, srv, ctx, "Alpha-Dispatch-Expiry")
	secondUnit := createUnitForTest(t, srv, ctx, "Bravo-Dispatch-Expiry")
	require.NoError(
		t,
		srv.units.UpdateUnitAssignments(ctx, "ambulance", nil, firstUnit.GetId(), []int32{1}, nil),
	)
	require.NoError(
		t,
		srv.units.UpdateUnitAssignments(ctx, "ambulance", nil, secondUnit.GetId(), []int32{2}, nil),
	)

	dispatch := createDispatchForTest(t, srv, ctx)
	job := "ambulance"
	expiresAt := time.Now().Add(time.Minute)
	require.NoError(t, srv.dispatches.UpdateAssignments(
		ctx,
		&job,
		nil,
		dispatch.GetId(),
		[]int64{firstUnit.GetId(), secondUnit.GetId()},
		nil,
		expiresAt,
	))

	firstTimerKey := fmt.Sprintf("assignment.%d.%d", dispatch.GetId(), firstUnit.GetId())
	secondTimerKey := fmt.Sprintf("assignment.%d.%d", dispatch.GetId(), secondUnit.GetId())
	_, err := srv.dispatches.IdleStore().Get(ctx, firstTimerKey)
	require.NoError(t, err)
	_, err = srv.dispatches.IdleStore().Get(ctx, secondTimerKey)
	require.NoError(t, err)

	require.NoError(t, srv.dispatches.UpdateAssignments(
		ctx,
		&job,
		nil,
		dispatch.GetId(),
		nil,
		[]int64{firstUnit.GetId()},
		time.Time{},
	))

	var assignmentCount int
	require.NoError(t, db.QueryRow(
		`SELECT COUNT(*) FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?`,
		dispatch.GetId(),
		firstUnit.GetId(),
	).Scan(&assignmentCount))
	assert.Zero(t, assignmentCount)

	projected, err := srv.dispatches.Get(ctx, dispatch.GetId())
	require.NoError(t, err)
	require.Len(t, projected.GetUnits(), 1)
	assert.Equal(t, secondUnit.GetId(), projected.GetUnits()[0].GetUnitId())

	activity, err := srv.ListDispatchActivity(
		ctx,
		&pbcentrum.ListDispatchActivityRequest{Id: dispatch.GetId()},
	)
	require.NoError(t, err)
	assert.True(
		t,
		slices.ContainsFunc(
			activity.GetActivity(),
			func(status *centrumdispatches.DispatchStatus) bool {
				return status.GetUnitId() == firstUnit.GetId() &&
					status.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED
			},
		),
	)

	_, err = srv.dispatches.IdleStore().Get(ctx, firstTimerKey)
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	_, err = srv.dispatches.IdleStore().Get(ctx, secondTimerKey)
	require.NoError(t, err)
}
