// Package retry provides small helpers for bounded retry loops.
package retry

import (
	"context"
	"math/rand/v2"
	"time"
)

// BackoffConfig controls exponential retry delays.
type BackoffConfig struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       float64
}

// Backoff produces jittered exponential delays.
type Backoff struct {
	config BackoffConfig
	delay  time.Duration
	jitter func(time.Duration) time.Duration
}

// NewBackoff creates a backoff using config. A zero multiplier defaults to 2.
func NewBackoff(config BackoffConfig) *Backoff {
	if config.Multiplier == 0 {
		config.Multiplier = 2
	}

	b := &Backoff{config: config}
	b.jitter = func(delay time.Duration) time.Duration {
		return jitterDelay(delay, config.Jitter)
	}
	return b
}

// Next returns the next delay in the exponential backoff sequence.
func (b *Backoff) Next() time.Duration {
	if b.delay == 0 {
		b.delay = b.config.InitialDelay
		if b.config.MaxDelay > 0 && b.delay > b.config.MaxDelay {
			b.delay = b.config.MaxDelay
		}
	} else if b.delay < b.config.MaxDelay {
		b.delay = time.Duration(float64(b.delay) * b.config.Multiplier)
		if b.delay > b.config.MaxDelay {
			b.delay = b.config.MaxDelay
		}
	}

	delay := b.jitter(b.delay)
	if b.config.MaxDelay > 0 && delay > b.config.MaxDelay {
		return b.config.MaxDelay
	}

	return delay
}

// Reset starts the delay sequence from its initial delay.
func (b *Backoff) Reset() {
	b.delay = 0
}

// Wait blocks for the next backoff delay or returns when ctx is canceled.
func (b *Backoff) Wait(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	timer := time.NewTimer(b.Next())
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func jitterDelay(delay time.Duration, jitter float64) time.Duration {
	if delay <= 0 || jitter <= 0 {
		return delay
	}

	//nolint:gosec // A non-cryptographic jitter source is sufficient here.
	factor := 1 + (rand.Float64()*2-1)*jitter
	return time.Duration(float64(delay) * factor)
}

// Do invokes fn up to attempts times. It retries only errors accepted by
// shouldRetry and returns the number of retries performed alongside the final
// error. attempts less than one still invokes fn once.
func Do(attempts int, fn func() error, shouldRetry func(error) bool) (int, error) {
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for attempt := range attempts {
		err := fn()
		if err == nil || !shouldRetry(err) {
			return attempt, err
		}
		lastErr = err
	}

	return attempts - 1, lastErr
}

// DoWithBackoff is Do with a context-aware wait before each retry.
func DoWithBackoff(
	ctx context.Context,
	attempts int,
	backoff *Backoff,
	fn func() error,
	shouldRetry func(error) bool,
) (int, error) {
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for attempt := range attempts {
		err := fn()
		if err == nil || !shouldRetry(err) {
			return attempt, err
		}
		lastErr = err

		if attempt == attempts-1 {
			break
		}
		if err := backoff.Wait(ctx); err != nil {
			return attempt, err
		}
	}

	return attempts - 1, lastErr
}
