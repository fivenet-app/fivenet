package store

import (
	"context"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/tests"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestBasicStoreCreateAndUse(t *testing.T) {
	t.Parallel()
	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	defer func() {
		if err := shutdown(); err != nil {
			t.Error(err)
		}
	}()

	logger := zaptest.NewLogger(t)
	ctx := t.Context()

	bucket := "test1"
	store, err := New[tests.SimpleObject](ctx, logger, js, bucket)
	require.NoError(t, err)
	assert.NotNil(t, store)

	storeCtx, storeCancel := context.WithCancel(ctx)
	t.Cleanup(storeCancel)
	err = store.Start(storeCtx, false)
	require.NoError(t, err)

	// Check if the Jetstream KV was auto-created
	kv, err := js.KeyValue(ctx, bucket)
	require.NoError(t, err)
	assert.NotNil(t, kv)

	// Retrieve a non-existent key
	val, err := store.Get("non-existent-key")
	require.Error(t, err)
	assert.Nil(t, val)

	// Create and ensure two values are stored
	first := &tests.SimpleObject{
		Field1: "First",
		Field2: false,
	}
	err = store.Put(ctx, "first", first)
	require.NoError(t, err)

	second := &tests.SimpleObject{
		Field1: "Second",
		Field2: true,
	}
	err = store.Put(ctx, "second", second)
	require.NoError(t, err)

	changed, err := store.PutIfChanged(ctx, "first", first)
	require.NoError(t, err)
	assert.False(t, changed)

	changed, err = store.PutIfChanged(ctx, "first", &tests.SimpleObject{Field1: "Updated"})
	require.NoError(t, err)
	assert.True(t, changed)

	keys := store.Keys("")
	assert.Len(t, keys, 2)

	list := store.List()
	assert.Len(t, list, 2)

	// Retrieved values are **always clones** so compare exported values
	firstRetrieved, err := store.Get("first")
	require.NoError(t, err)
	assert.NotNil(t, firstRetrieved)

	assert.EqualExportedValues(t, firstRetrieved, first)

	// Check if Get returns the correct result for a "locally cached" value
	secondRetrieved, err := store.Get("second")
	require.NoError(t, err)
	assert.NotNil(t, secondRetrieved)

	assert.EqualExportedValues(t, secondRetrieved, second)

	// Make sure that ComputeUpdate works as expected
	newField1Val := "Hello World!"
	err = store.ComputeUpdate(
		ctx,
		"second",
		func(key string, existing *tests.SimpleObject) (*tests.SimpleObject, bool, error) {
			existing.Field1 = newField1Val
			existing.Field2 = false
			return existing, true, nil
		},
	)
	require.NoError(t, err)

	secondRetrieved, err = store.Get("second")
	require.NoError(t, err)
	assert.NotNil(t, secondRetrieved)

	if secondRetrieved != nil {
		assert.Equal(t, newField1Val, secondRetrieved.GetField1())
		assert.False(t, secondRetrieved.GetField2())
	}

	// Check that range runs the callback 2 times
	runCount := 0
	store.Range(func(key string, value *tests.SimpleObject) bool {
		runCount++
		return true
	})
	assert.Equal(t, 2, runCount)

	// Check that deleting and finally clearing, removes values as expected
	err = store.Delete(ctx, "first")
	require.NoError(t, err)

	list = store.List()
	assert.Len(t, list, 1)

	err = store.Clear(ctx)
	require.NoError(t, err)

	list = store.List()
	assert.Empty(t, list)
}

func TestRangeWithKVPrefixUsesCallerKeys(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx := t.Context()
	store, err := New[tests.SimpleObject](
		ctx,
		zaptest.NewLogger(t),
		js,
		"range_prefix",
		WithKVPrefix[tests.SimpleObject]("id"),
	)
	require.NoError(t, err)

	storeCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	require.NoError(t, store.Start(storeCtx, false))
	require.NoError(t, store.Put(ctx, "42", &tests.SimpleObject{Field1: "prefixed"}))

	values := map[string]*tests.SimpleObject{}
	store.Range(func(key string, value *tests.SimpleObject) bool {
		values[key] = value
		return true
	})

	require.Len(t, values, 1)
	require.Equal(t, "prefixed", values["42"].GetField1())
}

//nolint:paralleltest,tparallel // Subtests must run before the parent deletes the shared test data.
func TestPrefixedStoresAreIsolatedInSharedBucket(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx := t.Context()
	newStore := func(prefix string) *Store[tests.SimpleObject, *tests.SimpleObject] {
		t.Helper()
		store, err := New[tests.SimpleObject](
			ctx,
			zaptest.NewLogger(t),
			js,
			"prefix_isolation",
			WithKVPrefix[tests.SimpleObject](prefix),
		)
		require.NoError(t, err)
		require.NoError(t, store.Start(ctx, false))
		return store
	}

	ids := newStore("id")
	jobs := newStore("job")
	require.NoError(t, ids.Put(ctx, "42", &tests.SimpleObject{Field1: "unit"}))
	require.NoError(t, jobs.Put(ctx, "42", &tests.SimpleObject{Field1: "job"}))

	for _, test := range []struct {
		name  string
		store *Store[tests.SimpleObject, *tests.SimpleObject]
		value string
	}{
		{name: "id", store: ids, value: "unit"},
		{name: "job", store: jobs, value: "job"},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, []string{"42"}, test.store.Keys(""))
			list := test.store.List()
			require.Len(t, list, 1)
			require.Equal(t, test.value, list[0].GetField1())
			filtered := test.store.ListFiltered("", nil)
			require.Len(t, filtered, 1)
			require.Equal(t, test.value, filtered[0].GetField1())
		})
	}

	require.NoError(t, ids.Delete(ctx, "42"))
	require.False(t, ids.Has("42"))
	remaining, err := jobs.Get("42")
	require.NoError(t, err)
	require.Equal(t, "job", remaining.GetField1())
}

func TestRangeWithKVPrefixRespectsIgnoredKeysAndClonesValues(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx := t.Context()
	store, err := New[tests.SimpleObject](
		ctx,
		zaptest.NewLogger(t),
		js,
		"range_prefix_ignored",
		WithKVPrefix[tests.SimpleObject]("id"),
		WithIgnoredKeys[tests.SimpleObject]("ignored"),
	)
	require.NoError(t, err)
	require.NoError(t, store.Start(ctx, false))
	require.NoError(t, store.Put(ctx, "visible", &tests.SimpleObject{Field1: "original"}))
	require.NoError(t, store.Put(ctx, "ignored", &tests.SimpleObject{Field1: "hidden"}))

	calls := 0
	store.Range(func(key string, value *tests.SimpleObject) bool {
		calls++
		require.Equal(t, "visible", key)
		value.Field1 = "mutated"
		return false
	})
	require.Equal(t, 1, calls, "Range must stop when the callback returns false")

	stored, err := store.Get("visible")
	require.NoError(t, err)
	require.Equal(t, "original", stored.GetField1(), "Range must provide a clone")
}

func TestWatchAllClosesUpdatesWhenCanceled(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	defer func() { require.NoError(t, shutdown()) }()

	ctx, cancel := context.WithCancel(t.Context())
	store, err := New[tests.SimpleObject](ctx, zaptest.NewLogger(t), js, "watch_close")
	require.NoError(t, err)

	watcher, err := store.WatchAll(ctx)
	require.NoError(t, err)
	cancel()

	select {
	case _, ok := <-watcher.Updates():
		require.False(t, ok)
	case <-time.After(time.Second):
		t.Fatal("watcher updates channel did not close after cancellation")
	}
}

func TestReconcileWatcherSnapshotRemovesOnlyMissingKeys(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	store, err := New[tests.SimpleObject](
		t.Context(),
		zaptest.NewLogger(t),
		js,
		"watch_reconcile",
		WithIgnoredKeys[tests.SimpleObject]("ignored"),
	)
	require.NoError(t, err)

	store.data.Store("retained", &tests.SimpleObject{Field1: "retained"})
	store.data.Store("stale", &tests.SimpleObject{Field1: "stale"})
	store.data.Store("ignored", &tests.SimpleObject{Field1: "ignored"})

	store.reconcileWatcherSnapshot(t.Context(), map[string]struct{}{"retained": {}})

	assert.True(t, store.Has("retained"))
	assert.False(t, store.Has("stale"))
	assert.True(t, store.Has("ignored"))
}
