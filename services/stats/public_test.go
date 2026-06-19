package stats

import (
	"testing"

	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPublicStatsEmptyCache(t *testing.T) {
	t.Parallel()

	srv := &Server{worker: &worker{}}

	resp, err := srv.GetPublicStats(t.Context(), &pbstats.GetPublicStatsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Nil(t, resp.GetStats())
}
