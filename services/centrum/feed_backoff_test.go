package centrum

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/retry"
	"github.com/stretchr/testify/require"
)

func TestFeedRestartBackoffGrowsAndCaps(t *testing.T) {
	t.Parallel()

	config := feedRestartBackoff
	config.Jitter = 0
	backoff := retry.NewBackoff(config)

	require.Equal(t, time.Second, backoff.Next())
	require.Equal(t, 2*time.Second, backoff.Next())
	require.Equal(t, 4*time.Second, backoff.Next())
	require.Equal(t, 8*time.Second, backoff.Next())
	require.Equal(t, 10*time.Second, backoff.Next())
	require.Equal(t, 10*time.Second, backoff.Next())
}

func TestFeedRestartBackoffReset(t *testing.T) {
	t.Parallel()

	config := feedRestartBackoff
	config.Jitter = 0
	backoff := retry.NewBackoff(config)

	backoff.Next()
	backoff.Next()
	backoff.Reset()

	require.Equal(t, time.Second, backoff.Next())
}

func TestJitterFeedRestartDelayStaysWithinBounds(t *testing.T) {
	t.Parallel()

	const delay = 10 * time.Second
	minJitter := time.Duration(float64(delay) * (1 - feedRestartJitter))
	maxJitter := time.Duration(float64(delay) * (1 + feedRestartJitter))

	for range 100 {
		backoff := retry.NewBackoff(feedRestartBackoff)
		var jittered time.Duration
		for range 5 {
			jittered = backoff.Next()
		}
		require.GreaterOrEqual(t, jittered, minJitter)
		require.LessOrEqual(t, jittered, maxJitter)
	}
}
