package centrum

import (
	"testing"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	unit := createUnitForTest(t, srv, ctx, "Alpha-Deleted")
	require.NoError(t, trackerStub.SetUserMappingForUser(ctx, 1, &unit.Id))

	require.NoError(t, srv.units.Delete(ctx, unit.GetId()))

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
	insertAssignmentRowForTest(t, db, unit.GetId(), 2)
	require.NoError(t, srv.units.SyncUnitMembership(ctx, unit.GetId()))
	deleteAssignmentRowForTest(t, db, unit.GetId(), 2)

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

func TestSyncUnitMembershipDoesNotRestoreOffDutyUserMapping(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Off-Duty-Recovery")
	insertAssignmentRowForTest(t, db, unit.GetId(), 1)
	require.NoError(t, trackerStub.SetUserMappingForUser(ctx, 1, &unit.Id))
	delete(trackerStub.markers, 1)

	require.NoError(t, srv.units.SyncUnitMembership(ctx, unit.GetId()))

	_, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, 1, unitAssignmentCountForTest(t, db, unit.GetId(), 1))
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

func TestSyncUserUnitMappingRemovesAssignmentForUserOnDifferentJob(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Job-Change")
	seedAssignmentForTest(t, db, trackerStub, unit.GetId(), 1)
	trackerStub.markers[1].Job = "police"

	require.NoError(t, srv.units.SyncUserUnitMapping(ctx, 1))

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Nil(t, mapping.UnitId)
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, false)
}

func TestReconcileUserJobChangeRemovesOnlyCrossJobAssignment(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Job-Change-Event")
	seedAssignmentForTest(t, db, trackerStub, unit.GetId(), 1)

	removed, err := srv.units.ReconcileUserJobChange(ctx, 1, "ambulance")
	require.NoError(t, err)
	assert.False(t, removed)
	assert.Equal(t, 1, unitAssignmentCountForTest(t, db, unit.GetId(), 1))

	removed, err = srv.units.ReconcileUserJobChange(ctx, 1, "police")
	require.NoError(t, err)
	assert.True(t, removed)
	assert.Zero(t, unitAssignmentCountForTest(t, db, unit.GetId(), 1))
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, false)
}

func TestJobChangeRemovesUnitAndDispatcherStateTogether(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 1,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Combined-Job-Change")
	seedAssignmentForTest(t, db, trackerStub, unit.GetId(), 1)
	require.NoError(t, srv.units.SyncUnitMembership(ctx, unit.GetId()))
	require.NoError(t, srv.dispatchers.SetUserState(ctx, "ambulance", 1, true))

	trackerStub.markers[1].Job = "police"

	removed, err := srv.units.ReconcileUserJobChange(ctx, 1, "police")
	require.NoError(t, err)
	assert.True(t, removed)
	require.NoError(t, srv.dispatchers.SetUserState(ctx, "ambulance", 1, false))

	assert.Zero(t, unitAssignmentCountForTest(t, db, unit.GetId(), 1))
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, false)

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, int64(0), mapping.GetUnitId())

	dispatchers, err := srv.dispatchers.Get(ctx, "ambulance")
	require.NoError(t, err)
	assert.Empty(t, dispatchers.GetDispatchers())
}

func TestJobChangeReconciliationIsIdempotentAndPreservesCurrentJobState(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 1,
	})

	unit := createUnitForTest(t, srv, ctx, "Alpha-Idempotent-Job-Change")
	seedAssignmentForTest(t, db, trackerStub, unit.GetId(), 1)
	require.NoError(t, srv.units.SyncUnitMembership(ctx, unit.GetId()))

	removed, err := srv.units.ReconcileUserJobChange(ctx, 1, "ambulance")
	require.NoError(t, err)
	assert.False(t, removed)
	assert.Equal(t, 1, unitAssignmentCountForTest(t, db, unit.GetId(), 1))

	require.NoError(t, srv.dispatchers.SetUserState(ctx, "ambulance", 1, true))
	require.NoError(t, srv.dispatchers.SetUserState(ctx, "ambulance", 1, true))
	dispatchers, err := srv.dispatchers.Get(ctx, "ambulance")
	require.NoError(t, err)
	assert.Len(t, dispatchers.GetDispatchers(), 1)
}

func TestUpdateUnitStatusPersistsChangedReasonForSameStatus(t *testing.T) {
	t.Parallel()

	srv, _, _ := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})
	unit := createUnitForTest(t, srv, ctx, "Alpha-Status")
	firstReason := "available for patrol"

	first, updated, err := srv.units.UpdateStatus(ctx, unit.GetId(), &centrumunits.UnitStatus{
		UnitId: unit.GetId(),
		Status: centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE,
		Reason: &firstReason,
	})
	require.NoError(t, err)
	require.True(t, updated)
	require.NotNil(t, first)

	secondReason := "available for transport"
	second, updated, err := srv.units.UpdateStatus(ctx, unit.GetId(), &centrumunits.UnitStatus{
		UnitId: unit.GetId(),
		Status: centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE,
		Reason: &secondReason,
	})
	require.NoError(t, err)
	require.True(t, updated)
	require.NotNil(t, second)
	assert.Greater(t, second.GetId(), first.GetId())
	assert.Equal(t, "available for transport", second.GetReason())
}

func TestUpdateUnitStatusResponseContainsPersistedStatus(t *testing.T) {
	t.Parallel()

	srv, _, _ := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})
	unit := createUnitForTest(t, srv, ctx, "Alpha-Status-Response")
	userID := int32(1)
	require.NoError(
		t,
		srv.units.UpdateUnitAssignments(
			ctx,
			"ambulance",
			&userID,
			unit.GetId(),
			[]int32{userID},
			nil,
		),
	)

	first, err := srv.UpdateUnitStatus(ctx, &pbcentrum.UpdateUnitStatusRequest{
		UnitId: unit.GetId(),
		Status: centrumunits.StatusUnit_STATUS_UNIT_BUSY,
	})
	require.NoError(t, err)
	require.True(t, first.GetUpdated())
	require.NotNil(t, first.GetStatus())
	require.Equal(t, centrumunits.StatusUnit_STATUS_UNIT_BUSY, first.GetStatus().GetStatus())

	second, err := srv.UpdateUnitStatus(ctx, &pbcentrum.UpdateUnitStatusRequest{
		UnitId: unit.GetId(),
		Status: centrumunits.StatusUnit_STATUS_UNIT_BUSY,
	})
	require.NoError(t, err)
	assert.False(t, second.GetUpdated())
	require.NotNil(t, second.GetStatus())
	assert.Equal(t, first.GetStatus().GetId(), second.GetStatus().GetId())
}

func TestUpdateUnitAssignmentsDropsUsersOutsideUnitJob(t *testing.T) {
	t.Parallel()

	srv, db, trackerStub := newCentrumJoinUnitTestServer(t)
	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId:   1,
		Job:      "ambulance",
		JobGrade: 17,
	})

	trackerStub.markers[4] = &livemapmarkers.UserMarker{
		UserId: 4,
		Hidden: false,
	}
	upsertUserJobForTest(t, db, 1, "ambulance", 17)

	unit := createUnitForTest(t, srv, ctx, "Alpha-Assignments")

	require.NoError(t, srv.units.UpdateUnitAssignments(
		ctx,
		"ambulance",
		nil,
		unit.GetId(),
		[]int32{1, 4},
		nil,
	))

	assert.Equal(t, 1, unitAssignmentCountForTest(t, db, unit.GetId(), 1))
	assert.Equal(t, 0, unitAssignmentCountForTest(t, db, unit.GetId(), 4))
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 1, true)
	assertUnitCacheHasUser(t, srv, ctx, unit.GetId(), 4, false)

	mapping, ok, err := trackerStub.GetUserMapping(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, mapping)
	assert.Equal(t, unit.GetId(), mapping.GetUnitId())

	_, ok, err = trackerStub.GetUserMapping(4)
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
