package livemapper

import (
	"context"
	"os"
	"testing"

	pblivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper"
	"github.com/fivenet-app/fivenet/internal/modules"
	"github.com/fivenet-app/fivenet/internal/tests/servers"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	grpcserver "github.com/fivenet-app/fivenet/pkg/grpc"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/services/centrum/centrumbrokers"
	"github.com/fivenet-app/fivenet/services/centrum/centrummanager"
	"github.com/fivenet-app/fivenet/services/centrum/centrumstate"
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
			fx.Provide(centrumbrokers.New),
			fx.Provide(tracker.NewForTests),
			centrumstate.StateModule,
			centrummanager.Module,
			fx.Provide(grpcSrvModule),
			fx.Provide(grpcserver.AsService(func(p Params) *Server {
				srv = NewServer(p)
				return srv
			})),

			fx.Invoke(func(*grpc.Server) {}),
		)...,
	)
	assert.NotNil(t, app)

	app.RequireStart()
	defer app.RequireStop()
	assert.NotNil(t, srv)

	client := pblivemapper.NewLivemapperServiceClient(clientConn)
	defer cancel()
	_ = ctx
	_ = client
}
