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

	"github.com/fivenet-app/fivenet/v2025/internal/tests/servers"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func getNatsClient(
	ctx context.Context,
	js jetstream.JetStream,
	bucket string,
) (*Locks, error) {
	lBucket := fmt.Sprintf("%s_locks", bucket)
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:         lBucket,
		Description:    fmt.Sprintf("%s Locks", bucket),
		History:        1,
		MaxBytes:       -1,
		Storage:        jetstream.MemoryStorage,
		LimitMarkerTTL: 3 * time.Minute, // Set a limit marker TTL to avoid stale locks
	})
	if err != nil {
		return nil, err
	}

	n := NewWithKV(zap.NewNop(), kv, bucket, 6*time.Second)
	return n, nil
}

func TestNats_LockUnlock(t *testing.T) {
	natsServer := servers.NewNATSServer(t, true)
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

func TestNats_MultipleLocks(t *testing.T) {
	natsServer := servers.NewNATSServer(t, true)
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
	wg := sync.WaitGroup{}
	for i := range 500 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
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
		}(i)
	}

	wg.Wait()
}
