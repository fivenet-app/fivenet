package locks

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/galexrt/fivenet/internal/tests/servers"
	"github.com/nats-io/nats.go"
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

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: bucket,
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
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			//<-time.After(time.Duration(200+mrand.Float64()*(2000-200+1)) * time.Millisecond)
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

			// fmt.Printf("Worker %d has the lock (%v)\n", i, v)

			atomic.AddInt32(&tracker, -1)

			err = n.Unlock(context.Background(), lockKey)
			if err != nil {
				t.Errorf("Unlock() %s error = %v: %d", connName, err, n.getRev("LOCK."+lockKey))
			}
		}(i)
	}

	wg.Wait()
}
