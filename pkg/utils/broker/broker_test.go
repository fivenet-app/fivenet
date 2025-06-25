package broker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrokerStart(t *testing.T) {
	type testMessage struct {
		ID   int
		Data string
	}

	broker := New[testMessage]()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// Start the broker in a separate goroutine
	go broker.Start(ctx)

	// Test subscription
	sub1 := broker.Subscribe()
	sub2 := broker.Subscribe()

	assert.Equal(t, int64(2), broker.SubCount())

	// Test publishing
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		select {
		case msg := <-sub1:
			require.Equal(t, 1, msg.ID, "unexpected ID received on sub1")
			require.Equal(t, "test", msg.Data, "unexpected Data received on sub1")

		case <-time.After(1 * time.Second):
			t.Error("timeout waiting for message on sub1")
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case msg := <-sub2:
			require.Equal(t, 1, msg.ID, "unexpected ID received on sub2")
			require.Equal(t, "test", msg.Data, "unexpected Data received on sub2")

		case <-time.After(1 * time.Second):
			t.Error("timeout waiting for message on sub2")
		}
	}()

	broker.Publish(testMessage{ID: 1, Data: "test"})
	wg.Wait()

	// Test unsubscribe and ensure subscriber count is correct
	broker.Unsubscribe(sub1)
	// Wait for the broker to process the unsubscribe
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int64(1), broker.SubCount(), "expected 1 subscriber after unsubscribing")

	// Ensure unsubscribed channel is closed
	select {
	case _, ok := <-sub1:
		if ok {
			t.Error("expected sub1 to be closed")
		}

	default:
		t.Error("expected sub1 to be closed, but it is still open")
	}

	// Test context cancellation and wait a moment for the broker to shut down
	cancel()
	time.Sleep(100 * time.Millisecond)

	select {
	case _, ok := <-sub2:
		if ok {
			t.Error("expected sub2 to be closed after context cancellation")
		}
	default:
		// sub2 should be closed
	}
}
