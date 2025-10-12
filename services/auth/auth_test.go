package auth

import (
	"net/http"
	"os"
	"strings"
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

	ctx := t.Context()

	assert := assert.New(t)
	require := require.New(t)

	clientConn, grpcSrvModule, err := modules.TestGRPCServer(ctx)
	require.NoError(err)

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
	assert.NotNil(app)

	app.RequireStart()
	defer app.RequireStop()
	assert.NotNil(srv)

	client := pbauth.NewAuthServiceClient(clientConn)

	// First login without credentials
	loginReq := &pbauth.LoginRequest{}
	loginReq.Username = ""
	loginReq.Password = ""
	res, err := client.Login(ctx, loginReq)
	require.Error(err)
	assert.Nil(res)
	proto.CompareGRPCStatusCode(t, codes.InvalidArgument, err)

	// Login with invalid credentials
	loginReq.Username = "non-existent-username"
	loginReq.Password = "non-existent-password"
	res, err = client.Login(ctx, loginReq)
	require.Error(err)
	assert.Nil(res)
	proto.CompareGRPCError(t, errorsauth.ErrInvalidLogin, err)

	// user-3: Login with valid account that has one char
	loginReq.Username = "user-3"
	loginReq.Password = "password"
	mdUser1 := metadata.New(map[string]string{})
	res, err = client.Login(ctx, loginReq, grpc.Header(&mdUser1))
	require.NoError(err)
	assert.NotNil(res)
	if res == nil {
		assert.FailNow("user-3: Login with valid account failed, response is nil")
	}
	cookies := mdUser1.Get("set-cookie")
	// Make sure we have both cookies
	require.Len(cookies, 2, "Expected 2 cookies to be set")
	foundAuthed := -1
	foundToken := -1
	for i, c := range cookies {
		if strings.HasPrefix(c, "fivenet_authed=") {
			foundAuthed = i
		}
		if strings.HasPrefix(c, "fivenet_token=") {
			foundToken = i
		}
	}
	require.True(foundAuthed != -1, "Expected a cookie starting with 'fivenet_authed='")
	require.True(foundToken != -1, "Expected a cookie starting with 'fivenet_token='")

	cookie, err := http.ParseSetCookie(cookies[foundToken])
	require.NoError(err)
	userToken := cookie.Value
	assert.NotEmpty(userToken)

	// user-3: Create authenticated metadate and get characters (only has one char)
	md := metadata.New(map[string]string{"Authorization": "Bearer " + userToken})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq := &pbauth.GetCharactersRequest{}
	getCharsRes, err := client.GetCharacters(ctx, getCharsReq)
	require.NoError(err)
	assert.NotNil(getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(
			"user-3: Empty char list returned for valid account that should have 2 chars",
		)
	}
	assert.Len(getCharsRes.GetChars(), 1)

	// user-1: Login with valid account (2 chars)
	loginReq.Username = "user-1"
	loginReq.Password = "password"
	mdUser2 := metadata.New(map[string]string{})
	res, err = client.Login(ctx, loginReq, grpc.Header(&mdUser2))
	require.NoError(err)
	assert.NotNil(res)
	if res == nil {
		assert.FailNow("user-1: Login with valid account failed, response is nil")
	}
	cookies = mdUser2.Get("set-cookie")
	// Make sure we have both cookies
	require.Len(cookies, 2, "Expected 2 cookies to be set")
	foundAuthed = -1
	foundToken = -1
	for i, c := range cookies {
		if strings.HasPrefix(c, "fivenet_authed=") {
			foundAuthed = i
		}
		if strings.HasPrefix(c, "fivenet_token=") {
			foundToken = i
		}
	}
	require.True(foundAuthed != -1, "Expected a cookie starting with 'fivenet_authed='")
	require.True(foundToken != -1, "Expected a cookie starting with 'fivenet_token='")

	cookie, err = http.ParseSetCookie(cookies[foundToken])
	require.NoError(err)
	userToken = cookie.Value
	assert.NotEmpty(userToken)

	// user-1: Create authenticated metadate and get characters
	md = metadata.New(map[string]string{"Authorization": "Bearer " + userToken})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq = &pbauth.GetCharactersRequest{}
	getCharsRes, err = client.GetCharacters(ctx, getCharsReq)
	require.NoError(err)
	assert.NotNil(getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(
			"user-1: Empty char list returned for valid account that should have 2 chars",
		)
	}
	assert.Len(getCharsRes.GetChars(), 2)

	// user-1: Choose an invalid character
	chooseCharReq := &pbauth.ChooseCharacterRequest{}
	chooseCharReq.CharId = 2 // Char id 2 is `user-2`'s char
	chooseCharRes, err := client.ChooseCharacter(ctx, chooseCharReq)
	require.Error(err)
	assert.Nil(chooseCharRes)
	proto.CompareGRPCError(t, errorsauth.ErrUnableToChooseChar, err)

	role, err := srv.ps.GetRoleByJobAndGrade(ctx, "ambulance", 1)
	require.NoError(err)
	require.NotNil(role)

	perm, err := srv.ps.GetPermission(
		ctx,
		permsauth.AuthServicePerm,
		permsauth.AuthServiceChooseCharacterPerm,
	)
	require.NoError(err)
	assert.NotNil(perm)

	// user-1: Choose valid character, the job role doesn't have permissions but the **default permissions** should still allow us to login
	err = srv.ps.RemovePermissionsFromRole(ctx, role.GetId(), perm.GetId())
	require.NoError(err)
	// Disable choose char perm but the **default permissions** will still allow us to login
	err = srv.ps.UpdateRolePermissions(ctx, role.GetId(), perms.AddPerm{
		Id:  perm.GetId(),
		Val: false,
	})
	require.NoError(err)
	chooseCharReq.CharId = 1
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	require.NoError(err)
	assert.NotNil(chooseCharRes)

	// user-1: Choose valid character, now we allow "choose char" perm for the job role
	err = srv.ps.UpdateRolePermissions(ctx, role.GetId(), perms.AddPerm{
		Id:  perm.GetId(),
		Val: true,
	})
	require.NoError(err)
	chooseCharReq.CharId = 1
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	require.NoError(err)
	assert.NotNil(chooseCharRes)
	if chooseCharRes != nil {
		assert.NotNil(chooseCharRes.GetChar())
	}
}
