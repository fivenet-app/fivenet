package stats

import (
	"context"
	"testing"

	statsstore "github.com/fivenet-app/fivenet/v2026/stores/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockStatsStore struct {
	data Stats
	err  error
}

func (m *mockStatsStore) LoadPublicStats(ctx context.Context) (Stats, error) {
	return m.data, m.err
}

func TestWorkerLoadStatsRoundsValues(t *testing.T) {
	t.Parallel()

	w := &worker{
		logger: zap.NewNop(),
		store: &mockStatsStore{
			data: statsstore.Stats{
				"users_registered": {Value: new(int32(11))},
				"citizens_total":   {Value: new(int32(20))},
			},
		},
	}

	err := w.loadStats(t.Context())
	require.NoError(t, err)

	got := w.GetStats()
	require.NotNil(t, got)
	assert.Equal(t, int32(20), (*got)["users_registered"].GetValue())
	assert.Equal(t, int32(20), (*got)["citizens_total"].GetValue())
}

func TestWorkerLoadStatsStoresDataOnError(t *testing.T) {
	t.Parallel()

	w := &worker{
		logger: zap.NewNop(),
		store: &mockStatsStore{
			data: statsstore.Stats{
				"users_registered": {Value: new(int32(9))},
			},
			err: context.Canceled,
		},
	}

	err := w.loadStats(t.Context())
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)

	got := w.GetStats()
	require.NotNil(t, got)
	assert.Equal(t, int32(10), (*got)["users_registered"].GetValue())
}
