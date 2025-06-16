package centrum

import (
	"context"
	"os"
	"testing"

	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/internal/modules"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrummanager"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumstate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	// All tests assume esx compat mode for now
	tables.EnableESXCompat()

	code := m.Run()
	os.Exit(code)
}

func TestBasicCentrumFlow(t *testing.T) {
	dbServer := servers.NewDBServer(t, true)
	natsServer := servers.NewNATSServer(t, true)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	clientConn, grpcSrvModule, err := modules.TestGRPCServer(ctx)
	require.NoError(t, err)

	var srv *Server
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			fx.Provide(modules.TestUserInfoRetriever),
			fx.Provide(tracker.NewForTests),
			centrumstate.StateModule,
			centrummanager.Module,
			fx.Provide(grpcSrvModule),
			fx.Provide(func(p Params) (Result, error) {
				r, err := NewServer(p)
				srv = r.Server
				return r, err
			}),

			fx.Invoke(func(*grpc.Server) {}),
		)...,
	)
	assert.NotNil(t, app)

	app.RequireStart()
	defer app.RequireStop()
	assert.NotNil(t, srv)

	client := pbcentrum.NewCentrumServiceClient(clientConn)
	defer cancel()
	_ = ctx
	_ = client
}
