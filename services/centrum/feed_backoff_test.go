package centrum

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFeedRestartBackoffGrowsAndCaps(t *testing.T) {
	t.Parallel()

	backoff := newFeedRestartBackoff()
	backoff.jitter = func(delay time.Duration) time.Duration { return delay }

	require.Equal(t, time.Second, backoff.next())
	require.Equal(t, 2*time.Second, backoff.next())
	require.Equal(t, 4*time.Second, backoff.next())
	require.Equal(t, 8*time.Second, backoff.next())
	require.Equal(t, 10*time.Second, backoff.next())
	require.Equal(t, 10*time.Second, backoff.next())
}

func TestFeedRestartBackoffReset(t *testing.T) {
	t.Parallel()

	backoff := newFeedRestartBackoff()
	backoff.jitter = func(delay time.Duration) time.Duration { return delay }

	backoff.next()
	backoff.next()
	backoff.reset()

	require.Equal(t, time.Second, backoff.next())
}

func TestFeedRestartBackoffZeroValueIsUsable(t *testing.T) {
	t.Parallel()

	var backoff feedRestartBackoff

	require.NotZero(t, backoff.next())
}

func TestJitterFeedRestartDelayStaysWithinBounds(t *testing.T) {
	t.Parallel()

	const delay = 10 * time.Second
	min := time.Duration(float64(delay) * (1 - feedRestartJitter))
	max := time.Duration(float64(delay) * (1 + feedRestartJitter))

	for range 100 {
		jittered := jitterFeedRestartDelay(delay)
		require.GreaterOrEqual(t, jittered, min)
		require.LessOrEqual(t, jittered, max)
	}
}
