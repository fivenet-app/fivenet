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
	t.Parallel()
	assert, require := assert.New(t), require.New(t)

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

	assert.Equal(int64(2), broker.SubCount())

	// Test publishing
	var wg sync.WaitGroup

	wg.Go(func() {
		select {
		case msg := <-sub1:
			require.Equal(1, msg.ID, "unexpected ID received on sub1")
			require.Equal("test", msg.Data, "unexpected Data received on sub1")

		case <-time.After(1 * time.Second):
			t.Error("timeout waiting for message on sub1")
		}
	})

	wg.Go(func() {
		select {
		case msg := <-sub2:
			require.NotNil(msg)
			require.Equal(1, msg.ID, "unexpected ID received on sub2")
			require.Equal("test", msg.Data, "unexpected Data received on sub2")

		case <-time.After(1 * time.Second):
			t.Error("timeout waiting for message on sub2")
		}
	})

	broker.Publish(testMessage{ID: 1, Data: "test"})
	wg.Wait()

	// Test unsubscribe and ensure subscriber count is correct
	broker.Unsubscribe(sub1)
	// Wait for the broker to process the unsubscribe
	time.Sleep(100 * time.Millisecond)
	assert.Equal(int64(1), broker.SubCount(), "expected 1 subscriber after unsubscribing")

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

func TestBrokerClosesSlowResyncSubscriber(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	broker := NewWithResyncOnSlowSubscriber[int](1)
	go broker.Start(ctx)

	sub := broker.Subscribe()
	broker.Publish(1)
	broker.Publish(2)
	broker.Publish(3)

	assert.Equal(t, 1, <-sub)
	_, ok := <-sub
	assert.False(t, ok)
	assert.Eventually(
		t,
		func() bool { return broker.SubCount() == 0 },
		time.Second,
		10*time.Millisecond,
	)
}

func TestBrokerOperationsAfterShutdownReturn(t *testing.T) {
	t.Parallel()

	b := New[int]()
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		b.Start(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("broker did not stop")
	}

	finished := make(chan struct{})
	go func() {
		ch := b.Subscribe()
		b.Publish(1)
		b.Unsubscribe(ch)
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("broker operation blocked after shutdown")
	}
}

func TestBrokerActiveSubscribersCloseOnShutdown(t *testing.T) {
	t.Parallel()

	b := New[int]()
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		b.Start(ctx)
		close(done)
	}()

	sub := b.Subscribe()
	assert.Equal(t, int64(1), b.SubCount())
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("broker did not stop")
	}
	select {
	case _, ok := <-sub:
		assert.False(t, ok)
	case <-time.After(time.Second):
		t.Fatal("subscriber was not closed")
	}
	assert.Equal(t, int64(0), b.SubCount())
}
