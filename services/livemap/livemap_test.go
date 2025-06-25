package livemap

import (
	"os"
	"testing"

	pblivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap"
	"github.com/fivenet-app/fivenet/v2025/internal/modules"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	grpcserver "github.com/fivenet-app/fivenet/v2025/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
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

	ctx := t.Context()

	clientConn, grpcSrvModule, err := modules.TestGRPCServer(ctx)
	require.NoError(t, err)

	var srv *Server
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			fx.Provide(modules.TestUserInfoRetriever),
			fx.Provide(tracker.NewForTests),
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

	client := pblivemap.NewLivemapServiceClient(clientConn)
	_ = ctx
	_ = client
}
