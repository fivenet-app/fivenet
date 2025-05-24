package auth

import (
	"context"
	"net/http"
	"os"
	"testing"

	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	permsauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth/perms"
	"github.com/fivenet-app/fivenet/v2025/internal/modules"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/proto"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/servers"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	grpcserver "github.com/fivenet-app/fivenet/v2025/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	errorsauth "github.com/fivenet-app/fivenet/v2025/services/auth/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	// All tests assume esx compat mode
	tables.EnableESXCompat()

	code := m.Run()
	os.Exit(code)
}

func TestFullAuthFlow(t *testing.T) {
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

	client := pbauth.NewAuthServiceClient(clientConn)

	// First login without credentials
	loginReq := &pbauth.LoginRequest{}
	loginReq.Username = ""
	loginReq.Password = ""
	res, err := client.Login(ctx, loginReq)
	assert.Error(t, err)
	assert.Nil(t, res)
	proto.CompareGRPCStatusCode(t, codes.InvalidArgument, err)

	// Login with invalid credentials
	loginReq.Username = "non-existant-username"
	loginReq.Password = "non-existant-password"
	res, err = client.Login(ctx, loginReq)
	assert.Error(t, err)
	assert.Nil(t, res)
	proto.CompareGRPCError(t, errorsauth.ErrInvalidLogin, err)

	// user-3: Login with valid account that has one char
	loginReq.Username = "user-3"
	loginReq.Password = "password"
	mdUser1 := metadata.New(map[string]string{})
	res, err = client.Login(ctx, loginReq, grpc.Header(&mdUser1))
	assert.NoError(t, err)
	assert.NotNil(t, res)
	if res == nil {
		assert.FailNow(t, "user-3: Login with valid account failed, response is nil")
	}
	cookies := mdUser1.Get("set-cookie")
	cookie, err := http.ParseSetCookie(cookies[0])
	require.NoError(t, err)
	userToken := cookie.Value
	assert.NotEmpty(t, userToken)

	// user-3: Create authenticated metadate and get characters (only has one char)
	md := metadata.New(map[string]string{"Authorization": "Bearer " + userToken})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq := &pbauth.GetCharactersRequest{}
	getCharsRes, err := client.GetCharacters(ctx, getCharsReq)
	assert.NoError(t, err)
	assert.NotNil(t, getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(t, "user-3: Empty char list returned for valid account that should have 2 chars")
	}
	assert.Len(t, getCharsRes.GetChars(), 1)

	// user-1: Login with valid account (2 chars)
	loginReq.Username = "user-1"
	loginReq.Password = "password"
	mdUser2 := metadata.New(map[string]string{})
	res, err = client.Login(ctx, loginReq, grpc.Header(&mdUser2))
	assert.NoError(t, err)
	assert.NotNil(t, res)
	if res == nil {
		assert.FailNow(t, "user-1: Login with valid account failed, response is nil")
	}
	cookies = mdUser2.Get("set-cookie")
	cookie, err = http.ParseSetCookie(cookies[0])
	require.NoError(t, err)
	userToken = cookie.Value
	assert.NotEmpty(t, userToken)

	// user-1: Create authenticated metadate and get characters
	md = metadata.New(map[string]string{"Authorization": "Bearer " + userToken})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq = &pbauth.GetCharactersRequest{}
	getCharsRes, err = client.GetCharacters(ctx, getCharsReq)
	assert.NoError(t, err)
	assert.NotNil(t, getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(t, "user-1: Empty char list returned for valid account that should have 2 chars")
	}
	assert.Len(t, getCharsRes.GetChars(), 2)

	// user-1: Choose an invalid character
	chooseCharReq := &pbauth.ChooseCharacterRequest{}
	chooseCharReq.CharId = 2 // Char id 2 is `user-2`'s char
	chooseCharRes, err := client.ChooseCharacter(ctx, chooseCharReq)
	assert.Error(t, err)
	assert.Nil(t, chooseCharRes)
	proto.CompareGRPCError(t, errorsauth.ErrUnableToChooseChar, err)

	role, err := srv.ps.GetRoleByJobAndGrade(ctx, "ambulance", 1)
	assert.NoError(t, err)
	require.NotNil(t, role)

	perm, err := srv.ps.GetPermission(ctx, permsauth.AuthServicePerm, permsauth.AuthServiceChooseCharacterPerm)
	assert.NoError(t, err)
	assert.NotNil(t, perm)

	// user-1: Choose valid character, the job role doesn't have permissions but the **default permissions** should still allow us to login
	err = srv.ps.RemovePermissionsFromRole(ctx, role.ID, perm.Id)
	assert.NoError(t, err)
	// Disable choose char perm but the **default permissions** will still allow us to login
	err = srv.ps.UpdateRolePermissions(ctx, role.ID, perms.AddPerm{
		Id:  perm.Id,
		Val: false,
	})
	assert.NoError(t, err)
	chooseCharReq.CharId = 1
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	assert.NoError(t, err)
	assert.NotNil(t, chooseCharRes)

	// user-1: Choose valid character, now we allow "choose char" perm for the job role
	err = srv.ps.UpdateRolePermissions(ctx, role.ID, perms.AddPerm{
		Id:  perm.Id,
		Val: true,
	})
	assert.NoError(t, err)
	chooseCharReq.CharId = 1
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	assert.NoError(t, err)
	assert.NotNil(t, chooseCharRes)
	if chooseCharRes != nil {
		assert.NotNil(t, chooseCharRes.Char)
	}
}
