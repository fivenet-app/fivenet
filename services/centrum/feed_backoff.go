package centrum

import (
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/retry"
)

const (
	feedRestartInitialDelay = time.Second
	feedRestartMaxDelay     = 10 * time.Second
	feedRestartJitter       = 0.2
)

var feedRestartBackoff = retry.BackoffConfig{
	InitialDelay: feedRestartInitialDelay,
	MaxDelay:     feedRestartMaxDelay,
	Multiplier:   2,
	Jitter:       feedRestartJitter,
}
