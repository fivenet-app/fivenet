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
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRemoveOrphanedProjections(t *testing.T) {
	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	ctx := t.Context()

	dispatchStore, err := store.New[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](
		ctx,
		zap.NewNop(),
		js,
		"test_dispatches_orphaned",
		store.WithKVPrefix[centrumdispatches.Dispatch, *centrumdispatches.Dispatch]("id"),
		store.WithLocks[centrumdispatches.Dispatch, *centrumdispatches.Dispatch](nil),
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

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, dispatchStore.Start(storeCtx, true))
	require.NoError(t, jobStore.Start(storeCtx, true))

	dispatch := &centrumdispatches.Dispatch{
		Id:   7,
		Jobs: &centrum.JobList{Jobs: []*centrum.JobListEntry{{Name: "police"}}},
	}
	require.NoError(t, dispatchStore.Put(ctx, "7", dispatch))
	require.NoError(t, jobStore.Put(ctx, "police.7", &common.IDMapping{Id: 7}))

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectQuery("SELECT.*fivenet_centrum_dispatches.*WHERE.*id").
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	dispatches := &DispatchDB{
		db:         db,
		store:      dispatchStore,
		jobMapping: jobStore,
	}
	require.NoError(t, dispatches.removeOrphanedProjections(ctx))
	require.NoError(t, mock.ExpectationsWereMet())

	_, err = dispatchStore.Get("7")
	assert.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	_, err = jobStore.Get("police.7")
	assert.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}
