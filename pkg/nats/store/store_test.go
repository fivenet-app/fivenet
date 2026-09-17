package store

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/tests"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/proto"
)

type scriptedWatchKV struct {
	jetstream.KeyValue

	mu      sync.Mutex
	results []scriptedWatchResult
	calls   atomic.Int32
}

type capturedWatchKV struct {
	jetstream.KeyValue

	watchers chan jetstream.KeyWatcher
}

func (k *capturedWatchKV) Watch(
	ctx context.Context,
	keys string,
	opts ...jetstream.WatchOpt,
) (jetstream.KeyWatcher, error) {
	watcher, err := k.KeyValue.Watch(ctx, keys, opts...)
	if err == nil {
		k.watchers <- watcher
	}
	return watcher, err
}

type scriptedWatchResult struct {
	watcher jetstream.KeyWatcher
	err     error
}

func (s *scriptedWatchKV) Watch(
	context.Context,
	string,
	...jetstream.WatchOpt,
) (jetstream.KeyWatcher, error) {
	s.calls.Add(1)
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.results) == 0 {
		return nil, errors.New("unexpected watch")
	}

	result := s.results[0]
	s.results = s.results[1:]
	return result.watcher, result.err
}

type scriptedKeyWatcher struct {
	updates chan jetstream.KeyValueEntry
}

func newScriptedKeyWatcher(
	closeAfter bool,
	entries ...jetstream.KeyValueEntry,
) *scriptedKeyWatcher {
	w := &scriptedKeyWatcher{updates: make(chan jetstream.KeyValueEntry, len(entries))}
	for _, entry := range entries {
		w.updates <- entry
	}
	if closeAfter {
		close(w.updates)
	}
	return w
}

func (w *scriptedKeyWatcher) Updates() <-chan jetstream.KeyValueEntry { return w.updates }
func (*scriptedKeyWatcher) Stop() error                               { return nil }

type scriptedKeyValueEntry struct {
	key       string
	value     []byte
	revision  uint64
	operation jetstream.KeyValueOp
}

func newScriptedPutEntry(t *testing.T, key, value string, revision uint64) *scriptedKeyValueEntry {
	t.Helper()
	data, err := proto.Marshal(&tests.SimpleObject{Field1: value})
	require.NoError(t, err)
	return &scriptedKeyValueEntry{
		key:       key,
		value:     data,
		revision:  revision,
		operation: jetstream.KeyValuePut,
	}
}

func (*scriptedKeyValueEntry) Bucket() string                    { return "scripted" }
func (e *scriptedKeyValueEntry) Key() string                     { return e.key }
func (e *scriptedKeyValueEntry) Value() []byte                   { return e.value }
func (e *scriptedKeyValueEntry) Revision() uint64                { return e.revision }
func (*scriptedKeyValueEntry) Created() time.Time                { return time.Time{} }
func (*scriptedKeyValueEntry) Delta() uint64                     { return 0 }
func (e *scriptedKeyValueEntry) Operation() jetstream.KeyValueOp { return e.operation }

func newScriptedStore(
	t *testing.T,
	kv jetstream.KeyValue,
	opts ...Option[tests.SimpleObject, *tests.SimpleObject],
) *Store[tests.SimpleObject, *tests.SimpleObject] {
	t.Helper()
	opts = append(
		opts,
		WithJetstreamKV[tests.SimpleObject](kv),
		WithLocks[tests.SimpleObject](nil),
	)
	store, err := New[tests.SimpleObject](
		t.Context(),
		zaptest.NewLogger(t),
		nil,
		"scripted",
		opts...)
	require.NoError(t, err)
	return store
}

func TestBasicStoreCreateAndUse(t *testing.T) {
	t.Parallel()
	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

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

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

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

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

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

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

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

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

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

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

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

func TestStoreRestartsWatcherBeforeInitialSnapshotCompletes(t *testing.T) {
	t.Parallel()

	kv := &scriptedWatchKV{results: []scriptedWatchResult{
		{
			// The first watcher applies a partial snapshot but closes before its
			// nil sentinel. The store must wait for the replacement snapshot.
			watcher: newScriptedKeyWatcher(true, newScriptedPutEntry(t, "partial", "old", 1)),
		},
		{
			watcher: newScriptedKeyWatcher(false,
				newScriptedPutEntry(t, "authoritative", "new", 2),
				nil,
			),
		},
	}}
	store := newScriptedStore(t, kv)

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	require.NoError(t, store.Start(ctx, true))

	assert.EqualValues(t, 2, kv.calls.Load())
	assert.False(t, store.Has("partial"))
	value, err := store.Get("authoritative")
	require.NoError(t, err)
	assert.Equal(t, "new", value.GetField1())
}

func TestStoreReconcilesStaleKeysAfterWatcherRestarts(t *testing.T) {
	t.Parallel()

	var remoteDeletes atomic.Int32
	kv := &scriptedWatchKV{results: []scriptedWatchResult{
		{
			// This is a completed initial snapshot, so Start may report ready.
			watcher: newScriptedKeyWatcher(true,
				newScriptedPutEntry(t, "retained", "current", 1),
				newScriptedPutEntry(t, "stale", "old", 2),
				nil,
			),
		},
		{
			// The replacement snapshot no longer contains stale.
			watcher: newScriptedKeyWatcher(false,
				newScriptedPutEntry(t, "retained", "current", 3),
				nil,
			),
		},
	}}
	store := newScriptedStore(t, kv,
		WithOnRemoteDeletedFn[tests.SimpleObject](
			func(context.Context, string, *tests.SimpleObject) error {
				remoteDeletes.Add(1)
				return nil
			},
		),
	)

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	require.NoError(t, store.Start(ctx, true))

	require.Eventually(t, func() bool {
		return kv.calls.Load() == 2 && !store.Has("stale")
	}, 2*time.Second, 10*time.Millisecond)
	assert.True(t, store.Has("retained"))
	assert.EqualValues(t, 1, remoteDeletes.Load())
}

func TestStoreRetriesFailedReplacementWatcher(t *testing.T) {
	t.Parallel()

	kv := &scriptedWatchKV{results: []scriptedWatchResult{
		{watcher: newScriptedKeyWatcher(true)},
		{err: errors.New("temporary watcher failure")},
		{watcher: newScriptedKeyWatcher(false, newScriptedPutEntry(t, "recovered", "yes", 1), nil)},
	}}
	store := newScriptedStore(t, kv)

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	require.NoError(t, store.Start(ctx, true))

	assert.EqualValues(t, 3, kv.calls.Load())
	value, err := store.Get("recovered")
	require.NoError(t, err)
	assert.Equal(t, "yes", value.GetField1())
}

func TestStoreStopsRetryingWatcherAfterCancellation(t *testing.T) {
	t.Parallel()

	kv := &scriptedWatchKV{results: []scriptedWatchResult{
		{watcher: newScriptedKeyWatcher(true)},
		{err: errors.New("watcher remains unavailable")},
		{err: errors.New("watcher remains unavailable")},
	}}
	store := newScriptedStore(t, kv)

	ctx, cancel := context.WithCancel(t.Context())
	startDone := make(chan error, 1)
	go func() { startDone <- store.Start(ctx, true) }()

	// The initial watcher closed before readiness, so Start remains blocked
	// while replacement watcher creation is retried.
	require.Eventually(
		t,
		func() bool { return kv.calls.Load() >= 2 },
		2*time.Second,
		10*time.Millisecond,
	)
	cancel()

	select {
	case err := <-startDone:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("store Start did not return after cancellation")
	}
}

func TestStoreRecoversAfterNATSWatcherStops(t *testing.T) {
	t.Parallel()

	js := nats.NewServer(t, nats.ServerOptions{InProcess: true}).GetJS()

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	baseKV, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  "watcher_recovery",
		Storage: jetstream.MemoryStorage,
	})
	require.NoError(t, err)
	initialData := map[string]string{
		"retained":  "old",
		"stale":     "old",
		"unrelated": "unchanged",
	}
	for key, value := range initialData {
		data, err := proto.Marshal(&tests.SimpleObject{Field1: value})
		require.NoError(t, err)
		_, err = baseKV.Put(ctx, key, data)
		require.NoError(t, err)
	}
	capturedKV := &capturedWatchKV{
		KeyValue: baseKV,
		watchers: make(chan jetstream.KeyWatcher, 2),
	}
	store, err := New[tests.SimpleObject](
		ctx,
		zaptest.NewLogger(t),
		js,
		"watcher_recovery",
		WithJetstreamKV[tests.SimpleObject](capturedKV),
		WithLocks[tests.SimpleObject](nil),
	)
	require.NoError(t, err)
	require.NoError(t, store.Start(ctx, true))
	for key, value := range initialData {
		actual, err := store.Get(key)
		require.NoError(t, err)
		assert.Equal(t, value, actual.GetField1())
	}

	firstWatcher := <-capturedKV.watchers
	require.NoError(t, firstWatcher.Stop())

	// These normal KV mutations occur while Store.Start waits to recreate the
	// real NATS watcher. Its replacement snapshot must update, delete, add,
	// and retain cached data correctly.
	data, err := proto.Marshal(&tests.SimpleObject{Field1: "new"})
	require.NoError(t, err)
	_, err = baseKV.Put(ctx, "retained", data)
	require.NoError(t, err)
	require.NoError(t, baseKV.Delete(ctx, "stale"))
	data, err = proto.Marshal(&tests.SimpleObject{Field1: "added during outage"})
	require.NoError(t, err)
	_, err = baseKV.Put(ctx, "added", data)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		retained, err := store.Get("retained")
		added, addedErr := store.Get("added")
		unrelated, unrelatedErr := store.Get("unrelated")
		return err == nil &&
			retained.GetField1() == "new" &&
			!store.Has("stale") &&
			addedErr == nil &&
			added.GetField1() == "added during outage" &&
			unrelatedErr == nil &&
			unrelated.GetField1() == "unchanged"
	}, 3*time.Second, 10*time.Millisecond)
	assert.Len(t, capturedKV.watchers, 1, "store must create a replacement watcher")
}
