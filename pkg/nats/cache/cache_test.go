package cache

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCachePutGetAndDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	kv := newFakeKV("test")
	c, err := New[wrapperspb.StringValue, *wrapperspb.StringValue](ctx, zap.NewNop(), nil, kv.bucket,
		withKeyValue[wrapperspb.StringValue, *wrapperspb.StringValue](kv),
	)
	require.NoError(t, err)

	require.NoError(t, c.Start(ctx, true))

	original := &wrapperspb.StringValue{Value: "hello"}
	require.NoError(t, c.Put(ctx, "foo/bar", original))

	require.Eventually(t, func() bool { return c.Has("foo/bar") }, time.Second, 10*time.Millisecond)

	got, err := c.Get("foo/bar")
	require.NoError(t, err)
	require.Equal(t, original.Value, got.Value)
	require.NotSame(t, original, got, "cached value should be cloned on read")

	got.Value = "changed"
	latest, err := c.Get("foo/bar")
	require.NoError(t, err)
	require.Equal(t, "hello", latest.Value, "mutating returned value must not mutate cache")

	require.NoError(t, c.Delete(ctx, "foo/bar"))
	require.Eventually(t, func() bool { return !c.Has("foo/bar") }, time.Second, 10*time.Millisecond)
}

func TestCacheTTLExpiry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	kv := newFakeKV("ttl")
	ttl := 30 * time.Millisecond
	c, err := New[wrapperspb.StringValue, *wrapperspb.StringValue](ctx, zap.NewNop(), nil, kv.bucket,
		withKeyValue[wrapperspb.StringValue, *wrapperspb.StringValue](kv),
		withTTL[wrapperspb.StringValue, *wrapperspb.StringValue](ttl),
	)
	require.NoError(t, err)

	require.NoError(t, c.Start(ctx, true))
	require.NoError(t, c.Put(ctx, "short", &wrapperspb.StringValue{Value: "live"}))

	require.Eventually(t, func() bool { return c.Has("short") }, time.Second, 10*time.Millisecond)

	time.Sleep(ttl + 20*time.Millisecond)
	require.False(t, c.Has("short"))
	_, err = c.Get("short")
	require.ErrorIs(t, err, jetstream.ErrKeyNotFound)
}

func TestCacheKeysListAndRange(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	kv := newFakeKV("keys")
	c, err := New[wrapperspb.StringValue, *wrapperspb.StringValue](ctx, zap.NewNop(), nil, kv.bucket,
		withKeyValue[wrapperspb.StringValue, *wrapperspb.StringValue](kv),
		withPrefix[wrapperspb.StringValue, *wrapperspb.StringValue]("pref"),
		withIgnoredKeys[wrapperspb.StringValue, *wrapperspb.StringValue]("ignored"),
	)
	require.NoError(t, err)

	require.NoError(t, c.Start(ctx, true))

	require.NoError(t, c.Put(ctx, "first", &wrapperspb.StringValue{Value: "one"}))
	require.NoError(t, c.Put(ctx, "ignored", &wrapperspb.StringValue{Value: "skip"}))
	require.NoError(t, c.Put(ctx, "second", &wrapperspb.StringValue{Value: "two"}))

	require.Eventually(t, func() bool { return c.Has("first") && c.Has("second") }, time.Second, 10*time.Millisecond)

	keys := c.Keys("")
	require.ElementsMatch(t, []string{"first", "second"}, keys)

	list := c.List()
	require.Len(t, list, 2)

	count := 0
	c.Range(func(key string, value *wrapperspb.StringValue) bool {
		count++
		return false
	})
	require.Equal(t, 1, count, "range should stop when callback returns false")

	require.NoError(t, c.Clear(ctx))
	require.Eventually(t, func() bool { return !c.Has("first") && !c.Has("second") }, time.Second, 10*time.Millisecond)
}

func TestCacheIgnoredKeysAreSkipped(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	kv := newFakeKV("ignored")
	c, err := New[wrapperspb.StringValue, *wrapperspb.StringValue](ctx, zap.NewNop(), nil, kv.bucket,
		withKeyValue[wrapperspb.StringValue, *wrapperspb.StringValue](kv),
		withIgnoredKeys[wrapperspb.StringValue, *wrapperspb.StringValue]("skip"),
	)
	require.NoError(t, err)

	require.NoError(t, c.Start(ctx, true))

	// Simulate external KV write that should be ignored by cache watcher.
	_, err = kv.Put(ctx, "skip", []byte("payload"))
	require.NoError(t, err)

	// Cache should not track ignored key even after watcher update.
	time.Sleep(20 * time.Millisecond)
	require.False(t, c.Has("skip"))
	require.Empty(t, c.Keys(""))
	require.Empty(t, c.List())
}

func TestCachePropagatesKVErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	putErr := errors.New("kv put failed")
	delErr := errors.New("kv delete failed")
	kv := newFailingKV("failures", putErr, delErr)

	c, err := New[wrapperspb.StringValue, *wrapperspb.StringValue](ctx, zap.NewNop(), nil, kv.bucket,
		withKeyValue[wrapperspb.StringValue, *wrapperspb.StringValue](kv),
	)
	require.NoError(t, err)

	require.NoError(t, c.Start(ctx, true))

	err = c.Put(ctx, "bad", &wrapperspb.StringValue{Value: "x"})
	require.ErrorIs(t, err, putErr)
	require.False(t, c.Has("bad"))

	err = c.Delete(ctx, "bad")
	require.ErrorIs(t, err, delErr)
}

// --- Test helpers ---------------------------------------------------------------------------

func withKeyValue[T any, U protoutils.ProtoMessage[T]](kv jetstream.KeyValue) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.kv = kv
	}
}

func withPrefix[T any, U protoutils.ProtoMessage[T]](prefix string) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.prefix = prefix
	}
}

func withTTL[T any, U protoutils.ProtoMessage[T]](ttl time.Duration) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.ttl = &ttl
	}
}

func withIgnoredKeys[T any, U protoutils.ProtoMessage[T]](keys ...string) Option[T, U] {
	return func(c *Cache[T, U]) {
		c.ignoredKeys = keys
	}
}

type fakeKV struct {
	bucket   string
	mu       sync.Mutex
	watchers []*fakeWatcher
	data     map[string][]byte
	revision uint64

	putErr    error
	deleteErr error
}

func newFakeKV(bucket string) *fakeKV {
	return &fakeKV{bucket: bucket, data: make(map[string][]byte)}
}

func newFailingKV(bucket string, putErr, deleteErr error) *fakeKV {
	return &fakeKV{
		bucket:    bucket,
		data:      make(map[string][]byte),
		putErr:    putErr,
		deleteErr: deleteErr,
	}
}

func (kv *fakeKV) Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	val, ok := kv.data[key]
	if !ok {
		return nil, jetstream.ErrKeyNotFound
	}

	return &fakeEntry{
		bucket:   kv.bucket,
		key:      key,
		value:    val,
		revision: kv.revision,
		created:  time.Now(),
		op:       jetstream.KeyValuePut,
	}, nil
}

func (kv *fakeKV) GetRevision(ctx context.Context, key string, revision uint64) (jetstream.KeyValueEntry, error) {
	return nil, errors.New("not implemented")
}

func (kv *fakeKV) Put(ctx context.Context, key string, value []byte) (uint64, error) {
	if kv.putErr != nil {
		return 0, kv.putErr
	}

	kv.mu.Lock()
	kv.revision++
	rev := kv.revision
	kv.data[key] = value
	watchers := append([]*fakeWatcher(nil), kv.watchers...)
	kv.mu.Unlock()

	entry := &fakeEntry{
		bucket:   kv.bucket,
		key:      key,
		value:    value,
		revision: rev,
		created:  time.Now(),
		op:       jetstream.KeyValuePut,
	}

	for _, w := range watchers {
		w.send(entry)
	}

	return rev, nil
}

func (kv *fakeKV) PutString(ctx context.Context, key string, value string) (uint64, error) {
	return kv.Put(ctx, key, []byte(value))
}

func (kv *fakeKV) Create(ctx context.Context, key string, value []byte, opts ...jetstream.KVCreateOpt) (uint64, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if _, exists := kv.data[key]; exists {
		return 0, jetstream.ErrKeyExists
	}
	return kv.Put(ctx, key, value)
}

func (kv *fakeKV) Update(ctx context.Context, key string, value []byte, revision uint64) (uint64, error) {
	return kv.Put(ctx, key, value)
}

func (kv *fakeKV) Delete(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
	if kv.deleteErr != nil {
		return kv.deleteErr
	}

	kv.mu.Lock()
	kv.revision++
	rev := kv.revision
	delete(kv.data, key)
	watchers := append([]*fakeWatcher(nil), kv.watchers...)
	kv.mu.Unlock()

	entry := &fakeEntry{
		bucket:   kv.bucket,
		key:      key,
		revision: rev,
		created:  time.Now(),
		op:       jetstream.KeyValueDelete,
	}

	for _, w := range watchers {
		w.send(entry)
	}

	return nil
}

func (kv *fakeKV) Purge(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
	return kv.Delete(ctx, key, opts...)
}

func (kv *fakeKV) Watch(ctx context.Context, keys string, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	w := newFakeWatcher()

	kv.mu.Lock()
	kv.watchers = append(kv.watchers, w)
	kv.mu.Unlock()

	go func() {
		// Signal ready
		w.send(nil)
		<-ctx.Done()
		_ = w.Stop()
	}()

	return w, nil
}

func (kv *fakeKV) WatchAll(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	return kv.Watch(ctx, ">", opts...)
}

func (kv *fakeKV) WatchFiltered(ctx context.Context, keys []string, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	return kv.Watch(ctx, strings.Join(keys, ","), opts...)
}

func (kv *fakeKV) Keys(ctx context.Context, opts ...jetstream.WatchOpt) ([]string, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	keys := make([]string, 0, len(kv.data))
	for k := range kv.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func (kv *fakeKV) ListKeys(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyLister, error) {
	return nil, errors.New("not implemented")
}

func (kv *fakeKV) ListKeysFiltered(ctx context.Context, filters ...string) (jetstream.KeyLister, error) {
	return nil, errors.New("not implemented")
}

func (kv *fakeKV) History(ctx context.Context, key string, opts ...jetstream.WatchOpt) ([]jetstream.KeyValueEntry, error) {
	return nil, errors.New("not implemented")
}

func (kv *fakeKV) Bucket() string {
	return kv.bucket
}

func (kv *fakeKV) PurgeDeletes(ctx context.Context, opts ...jetstream.KVPurgeOpt) error {
	return errors.New("not implemented")
}

func (kv *fakeKV) Status(ctx context.Context) (jetstream.KeyValueStatus, error) {
	return nil, errors.New("not implemented")
}

type fakeWatcher struct {
	updates chan jetstream.KeyValueEntry
	closed  atomic.Bool
}

func newFakeWatcher() *fakeWatcher {
	return &fakeWatcher{updates: make(chan jetstream.KeyValueEntry, 32)}
}

func (w *fakeWatcher) Updates() <-chan jetstream.KeyValueEntry {
	return w.updates
}

func (w *fakeWatcher) Stop() error {
	if w.closed.Swap(true) {
		return nil
	}
	close(w.updates)
	return nil
}

func (w *fakeWatcher) send(entry jetstream.KeyValueEntry) {
	if w.closed.Load() {
		return
	}

	select {
	case w.updates <- entry:
	default:
	}
}

type fakeEntry struct {
	bucket   string
	key      string
	value    []byte
	revision uint64
	created  time.Time
	op       jetstream.KeyValueOp
}

func (e *fakeEntry) Bucket() string                  { return e.bucket }
func (e *fakeEntry) Key() string                     { return e.key }
func (e *fakeEntry) Value() []byte                   { return e.value }
func (e *fakeEntry) Revision() uint64                { return e.revision }
func (e *fakeEntry) Created() time.Time              { return e.created }
func (e *fakeEntry) Delta() uint64                   { return 0 }
func (e *fakeEntry) Operation() jetstream.KeyValueOp { return e.op }
