package centrum

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrummanager"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrumstate"
	"github.com/fivenet-app/fivenet/internal/modules"
	"github.com/fivenet-app/fivenet/internal/tests/servers"
	grpcserver "github.com/fivenet-app/fivenet/pkg/grpc"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	if err := servers.TestDBServer.Setup(); err != nil {
		fmt.Println("failed to setup mysql test server: %w", err)
		return
	}
	defer servers.TestDBServer.Stop()

	if err := servers.TestNATSServer.Setup(); err != nil {
		fmt.Println("failed to setup nats test server: %w", err)
		return
	}
	defer servers.TestNATSServer.Stop()

	code := m.Run()

	os.Exit(code)
}

func TestBasicCentrumFlow(t *testing.T) {
	defer servers.TestDBServer.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientConn, grpcSrvModule, err := modules.TestGRPCServer(ctx)
	require.NoError(t, err)

	var srv *Server
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			fx.Provide(tracker.NewForTests),
			centrumstate.StateModule,
			centrummanager.Module,
			fx.Provide(grpcSrvModule),
			fx.Provide(grpcserver.AsService(func(p Params) (*Server, error) {
				srv, err = NewServer(p)
				return srv, err
			})),

			fx.Invoke(func(*grpc.Server) {}),
		)...,
	)
	assert.NotNil(t, app)

	app.RequireStart()
	assert.NotNil(t, srv)

	client := NewCentrumServiceClient(clientConn)
	defer cancel()
	_ = ctx
	_ = client

	// TODO add basic test

	app.RequireStop()
}
