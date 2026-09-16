package locks

import (
	"context"
	"fmt"
	"math/rand/v2"
	"path"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	testnats "github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/servers"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func getNatsClient(
	ctx context.Context,
	js jetstream.JetStream,
	bucket string,
) (*Locks, error) {
	return getNatsClientWithTTL(ctx, js, bucket, 6*time.Second)
}

func getNatsClientWithTTL(
	ctx context.Context,
	js jetstream.JetStream,
	bucket string,
	maxLockAge time.Duration,
) (*Locks, error) {
	lBucket := fmt.Sprintf("%s_locks", bucket)
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:         lBucket,
		Description:    fmt.Sprintf("%s Locks", bucket),
		History:        1,
		MaxBytes:       -1,
		Storage:        jetstream.MemoryStorage,
		LimitMarkerTTL: 3 * maxLockAge, // Set a limit marker TTL to avoid stale locks
	})
	if err != nil {
		return nil, err
	}

	n := NewWithKV(zap.NewNop(), kv, bucket, maxLockAge)
	return n, nil
}

//nolint:paralleltest // This test is not safe to run in parallel due to shared state in the NATS server and lock keys.
func TestNats_LockUnlock(t *testing.T) {
	natsServer := servers.NewNATSServer(t)
	js := natsServer.GetJS()

	ctx := t.Context()
	n, err := getNatsClient(ctx, js, "basic")
	require.NoError(t, err)

	lockKey := path.Join("acme", "example.com", "sites", "example.com")

	err = n.Lock(ctx, lockKey)
	if err != nil {
		t.Errorf("Unlock() error = %v", err)
	}

	err = n.Unlock(ctx, lockKey)
	if err != nil {
		t.Errorf("Unlock() error = %v", err)
	}
}

//nolint:paralleltest // This test is not safe to run in parallel due to shared state in the NATS server and lock keys.
func TestNats_MultipleLocks(t *testing.T) {
	natsServer := servers.NewNATSServer(t)
	js := natsServer.GetJS()

	lockKey := path.Join("acme", "example.com", "sites", "example.com")

	ctx := t.Context()
	n1, err := getNatsClient(ctx, js, "basic")
	require.NoError(t, err)
	n2, err := getNatsClient(ctx, js, "basic")
	require.NoError(t, err)
	n3, err := getNatsClient(ctx, js, "basic")
	require.NoError(t, err)

	err = n1.Lock(ctx, lockKey)
	if err != nil {
		t.Errorf("Lock() error = %v", err)
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		n1.Unlock(ctx, lockKey)
	}()

	err = n2.Lock(ctx, lockKey)
	if err != nil {
		t.Errorf("Lock() error = %v", err)
	}

	n2.Unlock(ctx, lockKey)

	time.Sleep(100 * time.Millisecond)
	err = n3.Lock(ctx, lockKey)
	if err != nil {
		t.Errorf("Lock() error = %v", err)
	}

	n3.Unlock(ctx, lockKey)

	tracker := int32(0)
	var wg sync.WaitGroup
	for i := range 500 {
		wg.Go(func() {
			<-time.After(time.Duration(200+rand.Float64()*(2000-200+1)) * time.Millisecond)
			n, err := getNatsClient(ctx, js, "basic")
			require.NoError(t, err)
			connName := fmt.Sprintf("nats-%d", i)

			err = n.Lock(ctx, lockKey)
			if err != nil {
				t.Errorf("Lock() %s error = %v: %d", connName, err, n.getRev("LOCK."+lockKey))
			}

			v := atomic.AddInt32(&tracker, 1)
			if v != 1 {
				panic("Had a concurrent lock")
			}

			t.Logf("worker %d has the lock (%v)", i, v)

			atomic.AddInt32(&tracker, -1)

			err = n.Unlock(ctx, lockKey)
			if err != nil {
				t.Errorf("Unlock() %s error = %v: %d", connName, err, n.getRev("LOCK."+lockKey))
			}
		})
	}

	wg.Wait()
}

func TestNats_StaleUnlockCannotReleaseNewOwner(t *testing.T) {
	t.Parallel()

	natsServer := testnats.NewServer(t, testnats.ServerOptions{InProcess: true})
	js := natsServer.GetJetStream()
	ctx := t.Context()

	const maxLockAge = time.Second
	first, err := getNatsClientWithTTL(ctx, js, "stale_unlock", maxLockAge)
	require.NoError(t, err)
	second, err := getNatsClientWithTTL(ctx, js, "stale_unlock", maxLockAge)
	require.NoError(t, err)

	const key = "dispatch.42"
	locked, err := first.TryLock(ctx, key)
	require.NoError(t, err)
	require.True(t, locked)

	require.Eventually(t, func() bool {
		firstLocked, err := first.IsLocked(ctx, key)
		return err == nil && !firstLocked
	}, 3*time.Second, 20*time.Millisecond, "the first owner's per-key TTL should expire")

	locked, err = second.TryLock(ctx, key)
	require.NoError(t, err)
	require.True(t, locked)

	require.Error(
		t,
		first.Unlock(ctx, key),
		"a stale revision must not delete the new owner's lock",
	)
	locked, err = second.IsLocked(ctx, key)
	require.NoError(t, err)
	require.True(t, locked, "the second owner must retain its lock after stale unlock")
	require.NoError(t, second.Unlock(ctx, key))
}
