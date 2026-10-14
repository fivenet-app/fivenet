package centrum

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWaitForReadyReturnsLoadError(t *testing.T) {
	t.Parallel()

	loadErr := errors.New("initial cache load failed")
	ready := make(chan struct{})
	close(ready)
	server := &Server{ready: ready, readyErr: loadErr}

	require.ErrorIs(t, server.waitForReady(t.Context()), loadErr)
}

func TestWaitForReadyHonorsContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	server := &Server{ready: make(chan struct{})}

	require.ErrorIs(t, server.waitForReady(ctx), context.Canceled)
}
