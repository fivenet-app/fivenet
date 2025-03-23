package locks

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/internal/tests/servers"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	if err := servers.TestNATSServer.Setup(); err != nil {
		fmt.Println("failed to setup nats test server: %w", err)
		return
	}
	defer servers.TestNATSServer.Stop()

	code := m.Run()

	os.Exit(code)
}

func getNatsClient(t *testing.T, bucket string) *Locks {
	js, err := servers.TestNATSServer.GetJS()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:   bucket,
		History:  3,
		Storage:  jetstream.MemoryStorage,
		TTL:      4 * LockTimeout,
		Replicas: 1,
	})
	require.NoError(t, err)

	n, err := New(zap.NewNop(), kv, bucket, 6*time.Second)
	require.NoError(t, err)
	return n
}

func TestNats_LockUnlock(t *testing.T) {
	n := getNatsClient(t, "basic")

	lockKey := path.Join("acme", "example.com", "sites", "example.com")

	err := n.Lock(context.Background(), lockKey)
	if err != nil {
		t.Errorf("Unlock() error = %v", err)
	}

	err = n.Unlock(context.Background(), lockKey)
	if err != nil {
		t.Errorf("Unlock() error = %v", err)
	}
}

func TestNats_MultipleLocks(t *testing.T) {
	lockKey := path.Join("acme", "example.com", "sites", "example.com")

	n1 := getNatsClient(t, "basic")
	n2 := getNatsClient(t, "basic")
	n3 := getNatsClient(t, "basic")

	err := n1.Lock(context.Background(), lockKey)
	if err != nil {
		t.Errorf("Lock() error = %v", err)
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		n1.Unlock(context.Background(), lockKey)
	}()

	err = n2.Lock(context.Background(), lockKey)
	if err != nil {
		t.Errorf("Lock() error = %v", err)
	}

	n2.Unlock(context.Background(), lockKey)

	time.Sleep(100 * time.Millisecond)
	err = n3.Lock(context.Background(), lockKey)
	if err != nil {
		t.Errorf("Lock() error = %v", err)
	}

	n3.Unlock(context.Background(), lockKey)

	tracker := int32(0)
	wg := sync.WaitGroup{}
	for i := range 500 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-time.After(time.Duration(200+rand.Float64()*(2000-200+1)) * time.Millisecond)
			n := getNatsClient(t, "basic")
			connName := fmt.Sprintf("nats-%d", i)

			err := n.Lock(context.Background(), lockKey)
			if err != nil {
				t.Errorf("Lock() %s error = %v: %d", connName, err, n.getRev("LOCK."+lockKey))
			}

			v := atomic.AddInt32(&tracker, 1)
			if v != 1 {
				panic("Had a concurrent lock")
			}

			t.Logf("worker %d has the lock (%v)", i, v)

			atomic.AddInt32(&tracker, -1)

			err = n.Unlock(context.Background(), lockKey)
			if err != nil {
				t.Errorf("Unlock() %s error = %v: %d", connName, err, n.getRev("LOCK."+lockKey))
			}
		}(i)
	}

	wg.Wait()
}
