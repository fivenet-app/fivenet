package discord

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGuildRunRejectsSyncDuringCooldown(t *testing.T) {
	t.Parallel()
	g := &Guild{
		initiated: func() (v atomic.Bool) { v.Store(true); return }(),
		lastSync:  time.Now(),
	}

	require.ErrorIs(t, g.Run(false), ErrSyncCooldownTime)
	require.False(t, g.IsRunning())
}

func TestGuildSetLastSyncIntervalUpdatesTimestamp(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = db.Close()
	})
	mock.ExpectExec("UPDATE .*fivenet_job_props.*").WillReturnResult(sqlmock.NewResult(0, 1))

	g := &Guild{
		bot:    &Bot{db: db},
		job:    "police",
		logger: zaptest.NewLogger(t),
	}
	before := time.Now()
	require.NoError(t, g.setLastSyncInterval(context.Background(), "police", nil))
	require.True(t, g.lastSync.After(before) || g.lastSync.Equal(before))
}
