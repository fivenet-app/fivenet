package dmv

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	pbdmv "github.com/fivenet-app/fivenet/gen/go/proto/services/dmv"
	"github.com/fivenet-app/fivenet/internal/modules"
	"github.com/fivenet-app/fivenet/internal/tests/servers"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	grpcserver "github.com/fivenet-app/fivenet/pkg/grpc"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	if err := servers.TestDBServer.Setup(); err != nil {
		fmt.Printf("failed to setup mysql test server. %v\n", err)
		return
	}
	defer servers.TestDBServer.Stop()

	if err := servers.TestNATSServer.Setup(); err != nil {
		fmt.Println("failed to setup nats test server: %w", err)
		return
	}
	defer servers.TestNATSServer.Stop()

	// All tests assume esx compat mode
	tables.EnableESXCompat()

	code := m.Run()

	os.Exit(code)
}

func TestListVehicles(t *testing.T) {
	defer servers.TestDBServer.Reset()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	clientConn, grpcSrvModule, err := modules.TestGRPCServer(ctx)
	require.NoError(t, err)

	ui := userinfo.NewMockUserInfoRetriever(map[int32]*userinfo.UserInfo{
		1: {
			AccountId: 1,
			Enabled:   true,
			UserId:    1,
			License:   "3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4",
			Group:     "user",
			Job:       "ambulance",
			JobGrade:  17,
		},
	})

	var srv *Server
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			fx.Provide(func() userinfo.UserInfoRetriever {
				return ui
			}),
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

	client := pbdmv.NewDMVServiceClient(clientConn)

	tm := auth.NewTokenMgr("")
	token, err := tm.NewWithClaims(&auth.CitizenInfoClaims{
		AccID:    1,
		Username: "user-1",
		CharID:   1,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "fivenet",
			Subject:  "3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4",
			ID:       "1",
			Audience: []string{"fivenet"},
		},
	})
	require.NoError(t, err)

	md := metadata.New(map[string]string{"Authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// List Vehicles
	pageSize := int64(5)
	listReq := &pbdmv.ListVehiclesRequest{
		Pagination: &database.PaginationRequest{
			Offset:   0,
			PageSize: &pageSize,
		},
	}
	res, err := client.ListVehicles(ctx, listReq)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Pagination)
	assert.Len(t, res.Vehicles, 5)

	app.RequireStop()
}
