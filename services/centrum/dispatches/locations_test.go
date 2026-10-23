package dispatches

import (
	"sync"
	"testing"

	centrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords"
	"github.com/stretchr/testify/require"
)

func TestUpdateDispatchLocationReconcilesJobsAndTerminalStatus(t *testing.T) {
	t.Parallel()

	db := &DispatchDB{
		dispatchLocationsMutex:  &sync.Mutex{},
		dispatchLocations:       make(map[string]*coords.Coords[*centrumdispatches.Dispatch]),
		dispatchLocationEntries: make(map[int64]map[string]*centrumdispatches.Dispatch),
	}

	dispatch := &centrumdispatches.Dispatch{
		Id:   7,
		Jobs: &centrum.JobList{Jobs: []*centrum.JobListEntry{{Name: "police"}}},
		Status: &centrumdispatches.DispatchStatus{
			Status: centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW,
		},
	}
	require.NoError(t, db.UpdateDispatchLocation(dispatch))
	require.True(
		t,
		db.GetLocations("police").Has(dispatch, centrumdispatches.DispatchPointMatchFn(7)),
	)

	updated := &centrumdispatches.Dispatch{
		Id:   7,
		Jobs: &centrum.JobList{Jobs: []*centrum.JobListEntry{{Name: "fire"}}},
		Status: &centrumdispatches.DispatchStatus{
			Status: centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW,
		},
	}
	require.NoError(t, db.UpdateDispatchLocation(updated))
	require.False(
		t,
		db.GetLocations("police").Has(dispatch, centrumdispatches.DispatchPointMatchFn(7)),
	)
	require.True(t, db.GetLocations("fire").Has(updated, centrumdispatches.DispatchPointMatchFn(7)))

	updated.Status.Status = centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED
	require.NoError(t, db.UpdateDispatchLocation(updated))
	require.False(
		t,
		db.GetLocations("fire").Has(updated, centrumdispatches.DispatchPointMatchFn(7)),
	)
}
