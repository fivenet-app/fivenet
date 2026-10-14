package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	t.Parallel()

	retryable := errors.New("retryable")
	calls := 0
	retries, err := Do(3, func() error {
		calls++
		if calls < 3 {
			return retryable
		}
		return nil
	}, func(err error) bool {
		return errors.Is(err, retryable)
	})

	require.NoError(t, err)
	assert.Equal(t, 3, calls)
	assert.Equal(t, 2, retries)
}

func TestDoStopsOnNonRetryableError(t *testing.T) {
	t.Parallel()

	nonRetryable := errors.New("non-retryable")
	calls := 0
	retries, err := Do(3, func() error {
		calls++
		return nonRetryable
	}, func(error) bool {
		return false
	})

	require.ErrorIs(t, err, nonRetryable)
	assert.Equal(t, 1, calls)
	assert.Zero(t, retries)
}

func TestDoExhaustsRetryableAttempts(t *testing.T) {
	t.Parallel()

	retryable := errors.New("retryable")
	calls := 0
	retries, err := Do(3, func() error {
		calls++
		return retryable
	}, func(err error) bool {
		return errors.Is(err, retryable)
	})

	require.ErrorIs(t, err, retryable)
	assert.Equal(t, 3, calls)
	assert.Equal(t, 2, retries)
}

func TestDoUsesOneAttemptForNonPositiveLimit(t *testing.T) {
	t.Parallel()

	calls := 0
	_, err := Do(0, func() error {
		calls++
		return nil
	}, func(error) bool {
		return true
	})

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestBackoffGrowsCapsAndResets(t *testing.T) {
	t.Parallel()

	backoff := NewBackoff(BackoffConfig{
		InitialDelay: 5 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2,
	})

	require.Equal(t, 5*time.Millisecond, backoff.Next())
	require.Equal(t, 10*time.Millisecond, backoff.Next())
	require.Equal(t, 20*time.Millisecond, backoff.Next())
	require.Equal(t, 40*time.Millisecond, backoff.Next())
	require.Equal(t, 80*time.Millisecond, backoff.Next())
	require.Equal(t, 100*time.Millisecond, backoff.Next())
	require.Equal(t, 100*time.Millisecond, backoff.Next())

	backoff.Reset()
	require.Equal(t, 5*time.Millisecond, backoff.Next())
}

func TestBackoffHardCapsJitteredDelay(t *testing.T) {
	t.Parallel()

	backoff := NewBackoff(BackoffConfig{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
	})
	backoff.jitter = func(delay time.Duration) time.Duration { return 2 * delay }

	require.Equal(t, 100*time.Millisecond, backoff.Next())
}

func TestBackoffWaitRespectsCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	backoff := NewBackoff(BackoffConfig{InitialDelay: time.Second})

	require.ErrorIs(t, backoff.Wait(ctx), context.Canceled)
}

func TestDoWithBackoffRetriesAndRespectsCancellation(t *testing.T) {
	t.Parallel()

	retryable := errors.New("retryable")
	calls := 0
	backoff := NewBackoff(BackoffConfig{InitialDelay: time.Millisecond})
	retries, err := DoWithBackoff(t.Context(), 3, backoff, func() error {
		calls++
		if calls < 3 {
			return retryable
		}
		return nil
	}, func(err error) bool {
		return errors.Is(err, retryable)
	})

	require.NoError(t, err)
	assert.Equal(t, 3, calls)
	assert.Equal(t, 2, retries)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	calls = 0
	_, err = DoWithBackoff(
		ctx,
		3,
		NewBackoff(BackoffConfig{InitialDelay: time.Second}),
		func() error {
			calls++
			return retryable
		},
		func(err error) bool {
			return errors.Is(err, retryable)
		},
	)

	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 1, calls)
}
