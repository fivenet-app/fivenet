package housekeeper

import (
	"context"
	"maps"
	"testing"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/tracker"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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
	mappings map[int32]*tracker.UserMapping
	onDuty   map[int32]bool
}

func (t *housekeeperTrackerStub) ListTrackedJobs() []string { return nil }

func (t *housekeeperTrackerStub) GetUserMarkerById(id int32) (*livemapmarkers.UserMarker, bool) {
	return &livemapmarkers.UserMarker{UserId: id}, true
}

func (t *housekeeperTrackerStub) IsUserOnDuty(userId int32) bool {
	return t.onDuty[userId]
}

func (t *housekeeperTrackerStub) Subscribe(
	_ context.Context,
) (store.IKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker], error) {
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
