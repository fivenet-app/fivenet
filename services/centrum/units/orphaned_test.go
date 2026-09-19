package units

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRemoveOrphanedProjections(t *testing.T) {
	t.Parallel()

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	ctx := t.Context()

	unitStore, err := store.New[centrumunits.Unit, *centrumunits.Unit](
		ctx,
		zap.NewNop(),
		js,
		"test_units_orphaned",
		store.WithKVPrefix[centrumunits.Unit, *centrumunits.Unit]("id"),
		store.WithLocks[centrumunits.Unit, *centrumunits.Unit](nil),
	)
	require.NoError(t, err)
	jobStore, err := store.New[common.IDMapping, *common.IDMapping](
		ctx,
		zap.NewNop(),
		js,
		"test_units_orphaned",
		store.WithKVPrefix[common.IDMapping, *common.IDMapping]("job"),
		store.WithLocks[common.IDMapping, *common.IDMapping](nil),
	)
	require.NoError(t, err)

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, unitStore.Start(storeCtx, true))
	require.NoError(t, jobStore.Start(storeCtx, true))

	orphaned := &centrumunits.Unit{Id: 7, Job: "police"}
	valid := &centrumunits.Unit{Id: 8, Job: "police"}
	require.NoError(t, unitStore.Put(ctx, "7", orphaned))
	require.NoError(t, unitStore.Put(ctx, "8", valid))
	require.NoError(t, jobStore.Put(ctx, "police.7", &common.IDMapping{Id: 7}))
	require.NoError(t, jobStore.Put(ctx, "police.8", &common.IDMapping{Id: 8}))

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	units := &UnitDB{
		db:         db,
		store:      unitStore,
		jobMapping: jobStore,
	}
	require.NoError(t, units.removeOrphanedProjections(ctx, []*centrumunits.Unit{valid}))

	_, err = unitStore.Get("7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	_, err = jobStore.Get("police.7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	got, err := unitStore.Get("8")
	require.NoError(t, err)
	assert.Equal(t, int64(8), got.GetId())
	_, err = jobStore.Get("police.8")
	assert.NoError(t, err)
}

func TestUpdateReturnsNoRowsForMissingUnit(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT.*fivenet_centrum_units.*").
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectRollback()

	units := &UnitDB{db: db}
	_, err = units.Update(t.Context(), 1, &centrumunits.Unit{Id: 42, Job: "police"})
	require.ErrorIs(t, err, sql.ErrNoRows, "expected sql.ErrNoRows, got %v", err)
	require.NoError(t, mock.ExpectationsWereMet())
}
