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

	specStore, err := store.New[centrumunits.UnitSpec, *centrumunits.UnitSpec](
		ctx,
		zap.NewNop(),
		js,
		"test_units_orphaned",
		store.WithKVPrefix[centrumunits.UnitSpec, *centrumunits.UnitSpec]("spec"),
		store.WithLocks[centrumunits.UnitSpec, *centrumunits.UnitSpec](nil),
	)
	require.NoError(t, err)
	statusStore, err := store.New[centrumunits.UnitStatus, *centrumunits.UnitStatus](
		ctx,
		zap.NewNop(),
		js,
		"test_units_orphaned",
		store.WithKVPrefix[centrumunits.UnitStatus, *centrumunits.UnitStatus]("status"),
		store.WithLocks[centrumunits.UnitStatus, *centrumunits.UnitStatus](nil),
	)
	require.NoError(t, err)
	membersStore, err := store.New[centrumunits.UnitAssignments, *centrumunits.UnitAssignments](
		ctx,
		zap.NewNop(),
		js,
		"test_units_orphaned",
		store.WithKVPrefix[centrumunits.UnitAssignments, *centrumunits.UnitAssignments]("members"),
		store.WithLocks[centrumunits.UnitAssignments, *centrumunits.UnitAssignments](nil),
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
	require.NoError(t, specStore.Start(storeCtx, true))
	require.NoError(t, statusStore.Start(storeCtx, true))
	require.NoError(t, membersStore.Start(storeCtx, true))
	require.NoError(t, jobStore.Start(storeCtx, true))

	orphaned := &centrumunits.Unit{Id: 7, Job: "police"}
	valid := &centrumunits.Unit{Id: 8, Job: "police"}
	require.NoError(t, specStore.Put(ctx, "7", unitSpecFromUnit(orphaned)))
	require.NoError(t, specStore.Put(ctx, "8", unitSpecFromUnit(valid)))
	require.NoError(
		t,
		membersStore.Put(ctx, "7", &centrumunits.UnitAssignments{UnitId: 7, Job: "police"}),
	)
	require.NoError(
		t,
		membersStore.Put(ctx, "8", &centrumunits.UnitAssignments{UnitId: 8, Job: "police"}),
	)
	require.NoError(t, jobStore.Put(ctx, "police.7", &common.IDMapping{Id: 7}))
	require.NoError(t, jobStore.Put(ctx, "police.8", &common.IDMapping{Id: 8}))

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	units := &UnitDB{
		db:           db,
		specStore:    specStore,
		statusStore:  statusStore,
		membersStore: membersStore,
		jobMapping:   jobStore,
	}
	require.NoError(t, units.removeOrphanedProjections(ctx, []*centrumunits.Unit{valid}))

	_, err = specStore.Get("7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	_, err = membersStore.Get("7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	_, err = jobStore.Get("police.7")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
	got, err := specStore.Get("8")
	require.NoError(t, err)
	assert.Equal(t, int64(8), got.GetId())
	_, err = jobStore.Get("police.8")
	require.NoError(t, err)
}

func TestUnitComponentsComposeIndependently(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()
	newStore := func() (*UnitDB, context.CancelFunc) {
		storeCtx, cancel := context.WithCancel(ctx)
		spec, err := store.New[centrumunits.UnitSpec](
			storeCtx,
			zap.NewNop(),
			js,
			"test_unit_components",
			store.WithKVPrefix[centrumunits.UnitSpec, *centrumunits.UnitSpec]("spec"),
			store.WithLocks[centrumunits.UnitSpec](nil),
		)
		require.NoError(t, err)
		status, err := store.New[centrumunits.UnitStatus](
			storeCtx,
			zap.NewNop(),
			js,
			"test_unit_components",
			store.WithKVPrefix[centrumunits.UnitStatus, *centrumunits.UnitStatus]("status"),
			store.WithLocks[centrumunits.UnitStatus](nil),
		)
		require.NoError(t, err)
		members, err := store.New[centrumunits.UnitAssignments](
			storeCtx,
			zap.NewNop(),
			js,
			"test_unit_components",
			store.WithKVPrefix[centrumunits.UnitAssignments, *centrumunits.UnitAssignments](
				"members",
			),
			store.WithLocks[centrumunits.UnitAssignments](nil),
		)
		require.NoError(t, err)
		for _, component := range []interface {
			Start(context.Context, bool) error
		}{spec, status, members} {
			require.NoError(t, component.Start(storeCtx, true))
		}
		return &UnitDB{specStore: spec, statusStore: status, membersStore: members}, cancel
	}

	units, cancel := newStore()
	t.Cleanup(cancel)
	key := "42"
	require.NoError(
		t,
		units.specStore.Put(
			ctx,
			key,
			&centrumunits.UnitSpec{Id: 42, Job: "police", Name: "Adam", SortOrder: 3},
		),
	)
	require.NoError(
		t,
		units.statusStore.Put(
			ctx,
			key,
			&centrumunits.UnitStatus{
				Id:     1,
				UnitId: 42,
				Status: centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE,
			},
		),
	)
	require.NoError(
		t,
		units.membersStore.Put(
			ctx,
			key,
			&centrumunits.UnitAssignments{
				UnitId: 42,
				Job:    "police",
				Users:  []*centrumunits.UnitAssignment{{UnitId: 42, UserId: 7}},
			},
		),
	)

	// A status write must not alter the independently stored spec or members.
	require.NoError(
		t,
		units.statusStore.Put(
			ctx,
			key,
			&centrumunits.UnitStatus{
				Id:     2,
				UnitId: 42,
				Status: centrumunits.StatusUnit_STATUS_UNIT_BUSY,
			},
		),
	)
	unit, err := units.Get(ctx, 42)
	require.NoError(t, err)
	assert.Equal(t, "Adam", unit.GetName())
	assert.Equal(t, centrumunits.StatusUnit_STATUS_UNIT_BUSY, unit.GetStatus().GetStatus())
	require.Len(t, unit.GetUsers(), 1)
	assert.Equal(t, int32(7), unit.GetUsers()[0].GetUserId())

	// Replacing membership must preserve the independently stored status/spec.
	require.NoError(
		t,
		units.membersStore.Put(
			ctx,
			key,
			&centrumunits.UnitAssignments{
				UnitId: 42,
				Job:    "police",
				Users:  []*centrumunits.UnitAssignment{{UnitId: 42, UserId: 9}},
			},
		),
	)
	unit, err = units.Get(ctx, 42)
	require.NoError(t, err)
	assert.Equal(t, "Adam", unit.GetName())
	assert.Equal(t, centrumunits.StatusUnit_STATUS_UNIT_BUSY, unit.GetStatus().GetStatus())
	require.Len(t, unit.GetUsers(), 1)
	assert.Equal(t, int32(9), unit.GetUsers()[0].GetUserId())
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
