package centrum

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
)

func newDispatchAssignmentFixture(t *testing.T) (*Server, *sql.DB, int64, int64) {
	t.Helper()

	srv, db, _ := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	unit := createUnitForTest(t, srv, ctx, "Assignment-Test")
	require.NoError(
		t,
		srv.units.UpdateUnitAssignments(
			ctx,
			"ambulance",
			new(int32(1)),
			unit.GetId(),
			[]int32{1},
			nil,
		),
	)
	dispatch := createDispatchForTest(t, srv, ctx, "ambulance")
	return srv, db, unit.GetId(), dispatch.GetId()
}

func assignmentKeyForTest(dispatchID, unitID int64) string {
	return fmt.Sprintf("assignment.%d.%d", dispatchID, unitID)
}

func assignmentExpiryForTest(t *testing.T, db *sql.DB, dispatchID, unitID int64) sql.NullTime {
	t.Helper()

	var expiresAt sql.NullTime
	err := db.QueryRow(
		"SELECT expires_at FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?",
		dispatchID,
		unitID,
	).Scan(&expiresAt)
	require.NoError(t, err)
	return expiresAt
}

func takeDispatchForTest(
	t *testing.T,
	srv *Server,
	dispatchID int64,
	resp centrumdispatches.TakeDispatchResp,
) error {
	t.Helper()
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	_, err := srv.TakeDispatch(ctx, &pbcentrum.TakeDispatchRequest{
		DispatchIds: []int64{dispatchID},
		Resp:        resp,
	})
	return err
}

func TestAssignDispatchCreatesThirtySecondPendingAssignment(t *testing.T) {
	t.Parallel()

	srv, _, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)

	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)

	dsp, err := srv.dispatches.Get(ctx, dispatchID)
	require.NoError(t, err)
	require.Len(t, dsp.GetUnits(), 1)
	require.NotNil(t, dsp.GetUnits()[0].GetExpiresAt())
	require.InDelta(t, 30, time.Until(dsp.GetUnits()[0].GetExpiresAt().AsTime()).Seconds(), 2)

	_, err = srv.dispatches.IdleStore().Get(ctx, assignmentKeyForTest(dispatchID, unitID))
	require.NoError(t, err)
}

func TestTakeDispatchAcceptsUnassignedDispatch(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)

	require.NoError(
		t,
		takeDispatchForTest(
			t,
			srv,
			dispatchID,
			centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED,
		),
	)

	require.False(t, assignmentExpiryForTest(t, db, dispatchID, unitID).Valid)
	dsp, err := srv.dispatches.Get(t.Context(), dispatchID)
	require.NoError(t, err)
	require.Len(t, dsp.GetUnits(), 1)
	require.Equal(t, unitID, dsp.GetUnits()[0].GetUnitId())
}

func TestTakeDispatchRejectsMissingAssignmentOnAssignedDispatch(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)

	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)
	_, err = db.Exec(
		"DELETE FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?",
		dispatchID,
		unitID,
	)
	require.NoError(t, err)

	err = takeDispatchForTest(
		t,
		srv,
		dispatchID,
		centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED,
	)
	require.ErrorIs(t, err, errorscentrum.ErrNotPartOfDispatch)
}

func TestTakeDispatchConcurrentSelfTakeCreatesOneAssignment(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	start := make(chan struct{})
	errs := make(chan error, 2)

	for range 2 {
		go func() {
			<-start
			errs <- takeDispatchForTest(
				t,
				srv,
				dispatchID,
				centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED,
			)
		}()
	}
	close(start)

	for range 2 {
		require.NoError(t, <-errs)
	}

	var count int
	require.NoError(
		t,
		db.QueryRow(
			"SELECT COUNT(*) FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?",
			dispatchID,
			unitID,
		).Scan(&count),
	)
	require.Equal(t, 1, count)
}

func TestForcedAssignDispatchCreatesAcceptedAssignmentWithoutExpiry(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)

	_, err := srv.AssignDispatch(ctx, &pbcentrum.AssignDispatchRequest{
		DispatchId: dispatchID,
		ToAdd:      []int64{unitID},
		Forced:     new(true),
	})
	require.NoError(t, err)

	dsp, err := srv.dispatches.Get(ctx, dispatchID)
	require.NoError(t, err)
	require.Len(t, dsp.GetUnits(), 1)
	require.Nil(t, dsp.GetUnits()[0].GetExpiresAt())
	require.False(t, assignmentExpiryForTest(t, db, dispatchID, unitID).Valid)
	_, err = srv.dispatches.IdleStore().Get(ctx, assignmentKeyForTest(dispatchID, unitID))
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}

func TestTakeDispatchAcceptanceCancelsTimerAndClearsExpiry(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)

	require.NoError(
		t,
		takeDispatchForTest(
			t,
			srv,
			dispatchID,
			centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED,
		),
	)
	require.False(t, assignmentExpiryForTest(t, db, dispatchID, unitID).Valid)
	dsp, err := srv.dispatches.Get(ctx, dispatchID)
	require.NoError(t, err)
	require.Nil(t, dsp.GetUnits()[0].GetExpiresAt())
	_, err = srv.dispatches.IdleStore().Get(ctx, assignmentKeyForTest(dispatchID, unitID))
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}

func TestTakeDispatchDeclineRemovesAssignmentAndTimer(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)

	require.NoError(
		t,
		takeDispatchForTest(
			t,
			srv,
			dispatchID,
			centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_DECLINED,
		),
	)
	var count int
	require.NoError(
		t,
		db.QueryRow("SELECT COUNT(*) FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?", dispatchID, unitID).
			Scan(&count),
	)
	require.Zero(t, count)
	_, err = srv.dispatches.IdleStore().Get(ctx, assignmentKeyForTest(dispatchID, unitID))
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}

func TestTakeDispatchTwoSecondGraceBoundary(t *testing.T) {
	t.Parallel()

	t.Run("within grace", func(t *testing.T) {
		t.Parallel()

		srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
		ctx := auth.ContextWithUserInfo(
			t.Context(),
			&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
		)
		_, err := srv.AssignDispatch(
			ctx,
			&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
		)
		require.NoError(t, err)
		_, err = db.Exec(
			"UPDATE fivenet_centrum_dispatches_asgmts SET expires_at = DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 SECOND) WHERE dispatch_id = ? AND unit_id = ?",
			dispatchID,
			unitID,
		)
		require.NoError(t, err)
		require.NoError(
			t,
			takeDispatchForTest(
				t,
				srv,
				dispatchID,
				centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED,
			),
		)
		require.False(t, assignmentExpiryForTest(t, db, dispatchID, unitID).Valid)
	})

	t.Run("after grace", func(t *testing.T) {
		t.Parallel()

		srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
		ctx := auth.ContextWithUserInfo(
			t.Context(),
			&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
		)
		_, err := srv.AssignDispatch(
			ctx,
			&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
		)
		require.NoError(t, err)
		_, err = db.Exec(
			"UPDATE fivenet_centrum_dispatches_asgmts SET expires_at = DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 3 SECOND) WHERE dispatch_id = ? AND unit_id = ?",
			dispatchID,
			unitID,
		)
		require.NoError(t, err)
		err = takeDispatchForTest(
			t,
			srv,
			dispatchID,
			centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED,
		)
		require.ErrorIs(t, err, errorscentrum.ErrNotPartOfDispatch)
	})
}

func TestUpdateExpiredAssignmentsRemovesOnlyExpiredRows(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)

	_, err = db.Exec(
		"UPDATE fivenet_centrum_dispatches_asgmts SET expires_at = DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 3 SECOND) WHERE dispatch_id = ? AND unit_id = ?",
		dispatchID,
		unitID,
	)
	require.NoError(t, err)
	var expired bool
	require.NoError(t, db.QueryRow(
		"SELECT expires_at <= DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 2 SECOND) FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?",
		dispatchID,
		unitID,
	).Scan(&expired))
	require.True(t, expired)
	deleted, err := srv.dispatches.UpdateExpiredAssignments(
		ctx,
		new("ambulance"),
		dispatchID,
		[]int64{unitID},
	)
	require.NoError(t, err)
	require.Equal(t, 1, deleted)

	var count int
	require.NoError(t, db.QueryRow(
		"SELECT COUNT(*) FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?",
		dispatchID,
		unitID,
	).Scan(&count))
	require.Zero(t, count)
}

func TestUpdateExpiredAssignmentsSkipsAssignmentsWithinGracePeriod(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)

	_, err = db.Exec(
		"UPDATE fivenet_centrum_dispatches_asgmts SET expires_at = DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 SECOND) WHERE dispatch_id = ? AND unit_id = ?",
		dispatchID,
		unitID,
	)
	require.NoError(t, err)
	_, err = srv.dispatches.UpdateExpiredAssignments(ctx, new("ambulance"), dispatchID, []int64{unitID})
	require.NoError(t, err)
	require.True(t, assignmentExpiryForTest(t, db, dispatchID, unitID).Valid)
}

func TestDeleteExpiredAssignmentsRetainsValidRows(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)

	deleted, err := srv.dispatches.DeleteExpiredAssignments(ctx, dispatchID, []int64{unitID})
	require.NoError(t, err)
	require.Zero(t, deleted)
	require.True(t, assignmentExpiryForTest(t, db, dispatchID, unitID).Valid)

	err = takeDispatchForTest(t, srv, dispatchID, centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED)
	require.NoError(t, err)
	deleted, err = srv.dispatches.DeleteExpiredAssignments(ctx, dispatchID, []int64{unitID})
	require.NoError(t, err)
	require.Zero(t, deleted)
	require.False(t, assignmentExpiryForTest(t, db, dispatchID, unitID).Valid)
}

func TestEmptyUnitCleanupRemovesAssignmentAndTimer(t *testing.T) {
	t.Parallel()

	srv, db, unitID, dispatchID := newDispatchAssignmentFixture(t)
	ctx := auth.ContextWithUserInfo(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Job: "ambulance", JobGrade: 1},
	)
	_, err := srv.AssignDispatch(
		ctx,
		&pbcentrum.AssignDispatchRequest{DispatchId: dispatchID, ToAdd: []int64{unitID}},
	)
	require.NoError(t, err)

	require.NoError(
		t,
		srv.units.RemoveUnitAssignments(ctx, "ambulance", new(int32(1)), unitID, []int32{1}),
	)
	require.NoError(
		t,
		srv.dispatches.UpdateAssignments(
			ctx,
			new("ambulance"),
			nil,
			dispatchID,
			nil,
			[]int64{unitID},
			time.Time{},
		),
	)

	var count int
	require.NoError(
		t,
		db.QueryRow("SELECT COUNT(*) FROM fivenet_centrum_dispatches_asgmts WHERE dispatch_id = ? AND unit_id = ?", dispatchID, unitID).
			Scan(&count),
	)
	require.Zero(t, count)
	_, err = srv.dispatches.IdleStore().Get(ctx, assignmentKeyForTest(dispatchID, unitID))
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}
