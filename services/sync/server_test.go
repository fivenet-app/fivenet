package sync

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/stretchr/testify/require"
)

func TestAuthFuncOverrideMarksAPITokenAuth(t *testing.T) {
	t.Parallel()

	server := &Server{
		tokens: []string{"sync-token"},
	}

	ctx := auth.SetTokenInGRPCContext(t.Context(), "sync-token")
	out, err := server.AuthFuncOverride(ctx, "/services.sync.SyncService/SendData")
	require.NoError(t, err)

	kind, ok := auth.GetAuthKindFromContext(out)
	require.True(t, ok)
	require.Equal(t, auth.AuthKindAPIToken, kind)
}
