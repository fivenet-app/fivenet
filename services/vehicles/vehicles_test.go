package vehicles

import (
	"context"
	"os"
	"testing"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbvehicles "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2025/internal/modules"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	grpcserver "github.com/fivenet-app/fivenet/v2025/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/userinfo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	// All tests assume esx compat mode
	tables.EnableESXCompat()

	code := m.Run()
	os.Exit(code)
}

func TestListVehicles(t *testing.T) {
	dbServer := servers.NewDBServer(t, true)
	natsServer := servers.NewNATSServer(t, true)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	clientConn, grpcSrvModule, err := modules.TestGRPCServer(ctx)
	require.NoError(t, err)

	ui := userinfo.NewMockUserInfoRetriever(map[int32]*pbuserinfo.UserInfo{
		3: {
			AccountId: 3,
			Enabled:   true,
			UserId:    3,
			License:   "db7e039146d5bf1b6781e7bc1bef31f0bb1298ea",
			Group:     "user",
			Job:       "doj",
			JobGrade:  16,
		},
	})

	var srv *Server
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
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

	client := pbvehicles.NewVehiclesServiceClient(clientConn)

	tm := auth.NewTokenMgr("")
	token, err := tm.NewWithClaims(&auth.CitizenInfoClaims{
		AccID:    3,
		Username: "user-3",
		CharID:   3,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "fivenet",
			Subject:  "db7e039146d5bf1b6781e7bc1bef31f0bb1298ea",
			ID:       "3",
			Audience: []string{"fivenet"},
		},
	})
	require.NoError(t, err)

	md := metadata.New(map[string]string{"Authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// List Vehicles
	pageSize := int64(5)
	listReq := &pbvehicles.ListVehiclesRequest{
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
}
