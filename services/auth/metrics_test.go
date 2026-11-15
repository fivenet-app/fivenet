package auth

import (
	"context"
	"testing"

	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestLoginMissingPasswordRecordsFailureMetric(t *testing.T) {
	t.Parallel()

	metrics := newAuthMetrics(nil)

	srv := newAuthErrorTestServer(t, &refreshAccountSessionStore{
		getLoginAccountByUsernameFn: func(_ context.Context, _ string) (*model.FivenetAccounts, error) {
			return &model.FivenetAccounts{ID: 1}, nil
		},
	})
	srv.metrics = metrics

	_, err := srv.Login(t.Context(), &pbauth.LoginRequest{})
	require.Error(t, err)
	require.InDelta(
		t,
		float64(1),
		testutil.ToFloat64(
			metrics.loginAttempts.WithLabelValues("failure", "account_missing_password"),
		),
		0,
	)
}
