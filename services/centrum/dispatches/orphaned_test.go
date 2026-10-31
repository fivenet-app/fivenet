package dispatches

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	centrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRemoveOrphanedProjections(t *testing.T) {
	t.Parallel()

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	ctx := t.Context()

	dispatchStore, err := store.New[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec](
		ctx,
		zap.NewNop(),
		js,
		"test_dispatches_orphaned",
		store.WithKVPrefix[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec]("spec"),
		store.WithLocks[centrumdispatches.DispatchSpec, *centrumdispatches.DispatchSpec](nil),
	)
	require.NoError(t, err)
	jobStore, err := store.New[common.IDMapping, *common.IDMapping](
		ctx,
		zap.NewNop(),
		js,
		"test_dispatches_orphaned",
		store.WithKVPrefix[common.IDMapping, *common.IDMapping]("job"),
		store.WithLocks[common.IDMapping, *common.IDMapping](nil),
	)
	require.NoError(t, err)
	assignmentsStore, err := store.New[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments](
		ctx,
		zap.NewNop(),
		js,
		"test_dispatches_orphaned",
		store.WithKVPrefix[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments](
			"assignments",
		),
		store.WithLocks[centrumdispatches.DispatchAssignments, *centrumdispatches.DispatchAssignments](
			nil,
		),
	)
	require.NoError(t, err)
	statusStore, err := store.New[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus](
		ctx,
		zap.NewNop(),
		js,
		"test_dispatches_orphaned",
		store.WithKVPrefix[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus](
			"status",
		),
		store.WithLocks[centrumdispatches.DispatchStatus, *centrumdispatches.DispatchStatus](nil),
	)
	require.NoError(t, err)

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, dispatchStore.Start(storeCtx, true))
	require.NoError(t, jobStore.Start(storeCtx, true))
	require.NoError(t, assignmentsStore.Start(storeCtx, true))
	require.NoError(t, statusStore.Start(storeCtx, true))

	dispatch := &centrumdispatches.DispatchSpec{
		Id:   7,
		Jobs: &centrum.JobList{Jobs: []*centrum.JobListEntry{{Name: "police"}}},
	}
	require.NoError(t, dispatchStore.Put(ctx, "7", dispatch))
	require.NoError(t, jobStore.Put(ctx, "police.7", &common.IDMapping{Id: 7}))
	// Simulate an interrupted job-list update: cleanup must not leave an
	// obsolete feed index behind just because it is absent from the spec.
	require.NoError(t, jobStore.Put(ctx, "fire.7", &common.IDMapping{Id: 7}))

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectQuery("SELECT.*fivenet_centrum_dispatches.*WHERE.*id").
		WithArgs(int64(7)).
		WillReturnError(qrm.ErrNoRows)

	dispatches := &DispatchDB{
		db:               db,
		specStore:        dispatchStore,
		assignmentsStore: assignmentsStore,
		statusStore:      statusStore,
		jobMapping:       jobStore,
	}
	require.NoError(t, dispatches.removeOrphanedProjections(ctx))
	require.NoError(t, mock.ExpectationsWereMet())

	_, err = dispatchStore.Get("7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	_, err = jobStore.Get("police.7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	_, err = jobStore.Get("fire.7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}
