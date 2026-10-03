package centrum

import (
	"math/rand"
	"time"
)

const (
	feedRestartInitialDelay = time.Second
	feedRestartMaxDelay     = 10 * time.Second
	feedRestartJitter       = 0.2
)

type feedRestartBackoff struct {
	delay  time.Duration
	jitter func(time.Duration) time.Duration
}

func newFeedRestartBackoff() feedRestartBackoff {
	return feedRestartBackoff{jitter: jitterFeedRestartDelay}
}

func (b *feedRestartBackoff) next() time.Duration {
	if b.jitter == nil {
		b.jitter = jitterFeedRestartDelay
	}

	if b.delay == 0 {
		b.delay = feedRestartInitialDelay
	} else if b.delay < feedRestartMaxDelay {
		b.delay *= 2
		if b.delay > feedRestartMaxDelay {
			b.delay = feedRestartMaxDelay
		}
	}

	return b.jitter(b.delay)
}

func (b *feedRestartBackoff) reset() {
	b.delay = 0
}

func jitterFeedRestartDelay(delay time.Duration) time.Duration {
	//nolint:gosec // A non-cryptographic jitter source is sufficient here.
	jitter := (rand.Float64()*2 - 1) * feedRestartJitter
	return time.Duration(float64(delay) * (1 + jitter))
}
