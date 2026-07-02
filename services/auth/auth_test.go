package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"

	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	permsauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth/perms"
	"github.com/fivenet-app/fivenet/v2026/internal/modules"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/proto"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/servers"
	grpcserver "github.com/fivenet-app/fivenet/v2026/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	errorsauth "github.com/fivenet-app/fivenet/v2026/services/auth/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestFullAuthFlow(t *testing.T) {
	t.Parallel()
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
	mdUser3 := metadata.New(map[string]string{})
	res, err = client.Login(ctx, loginReq, grpc.Header(&mdUser3))
	require.NoError(err)
	assert.NotNil(res)
	if res == nil {
		assert.FailNow("user-3: Login with valid account failed, response is nil")
	}
	cookies := mdUser3.Get("set-cookie")
	// Make sure we have both cookies
	require.Len(cookies, 2, "Expected 2 cookies to be set")
	foundAuthed := -1
	foundToken := -1
	for i, c := range cookies {
		if strings.HasPrefix(c, auth.AuthedCookieName+"=") {
			foundAuthed = i
		}
		if strings.HasPrefix(c, auth.AccCookieName+"=") {
			foundToken = i
		}
	}
	require.NotEqual(
		-1,
		foundAuthed,
		"Expected a cookie starting with '"+auth.AuthedCookieName+"='",
	)
	require.NotEqual(-1, foundToken, "Expected a cookie starting with '"+auth.AccCookieName+"='")

	cookie, err := http.ParseSetCookie(cookies[foundToken])
	require.NoError(err)
	accountToken := cookie.Value
	assert.NotEmpty(accountToken)

	// user-3: Create authenticated metadate and get characters (only has one char)
	mdUser3Outgoing := metadata.New(map[string]string{"Authorization": "Bearer " + accountToken})
	ctx = metadata.NewOutgoingContext(ctx, mdUser3Outgoing)
	getCharsReq := &pbauth.GetCharactersRequest{}
	getCharsRes, err := client.GetCharacters(ctx, getCharsReq)
	require.NoError(err)
	assert.NotNil(getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(
			"user-3: Empty getchar response returned for valid account that should have 2 chars in total",
		)
	}
	assert.Len(getCharsRes.GetChars(), 1)

	// user-1: Login with valid account (2 chars)
	loginReq.Username = "user-1"
	loginReq.Password = "password"
	mdUser1 := metadata.New(map[string]string{})
	res, err = client.Login(ctx, loginReq, grpc.Header(&mdUser1))
	require.NoError(err)
	assert.NotNil(res)
	if res == nil {
		assert.FailNow("user-1: Login with valid account failed, response is nil")
	}
	cookies = mdUser1.Get("set-cookie")
	// Make sure we have both cookies
	require.Len(cookies, 2, "Expected 2 cookies to be set")
	foundAuthed = -1
	foundToken = -1
	for i, c := range cookies {
		if strings.HasPrefix(c, auth.AuthedCookieName+"=") {
			foundAuthed = i
		}
		if strings.HasPrefix(c, auth.AccCookieName+"=") {
			foundToken = i
		}
	}
	require.NotEqual(
		-1,
		foundAuthed,
		"Expected a cookie starting with '"+auth.AuthedCookieName+"='",
	)
	require.NotEqual(-1, foundToken, "Expected a cookie starting with '"+auth.AccCookieName+"='")

	cookie, err = http.ParseSetCookie(cookies[foundToken])
	require.NoError(err)
	accountToken = cookie.Value
	assert.NotEmpty(accountToken)

	// user-1: Create authenticated metadate and get characters
	mdUser1Outgoing := metadata.New(map[string]string{"Cookie": cookie.String()})
	ctx = metadata.NewOutgoingContext(ctx, mdUser1Outgoing)
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
	chooseCharReq := &pbauth.ChooseCharacterRequest{
		CharId: 2, // Char id 2 is `user-2`'s char
	}
	chooseCharRes, err := client.ChooseCharacter(ctx, chooseCharReq)
	require.Error(err)
	assert.Nil(chooseCharRes)
	proto.CompareGRPCError(t, errorsauth.ErrUnableToChooseChar, err)

	role, err := srv.ps.GetRoleByJobAndGrade(ctx, "ambulance", 1)
	require.NoError(err)
	require.NotNil(role)

	perm, err := srv.ps.GetPermission(
		ctx,
		permsauth.Namespace,
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
	if chooseCharRes != nil {
		assert.NotEmpty(chooseCharRes.GetToken())
		assert.Equal("user-1", chooseCharRes.GetUsername())
	}

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
		assert.NotEmpty(chooseCharRes.GetToken())
		assert.Equal("user-1", chooseCharRes.GetUsername())
	}

	accInfoResp, err := client.GetAccountInfo(ctx, &pbauth.GetAccountInfoRequest{})
	require.NoError(err)
	assert.NotNil(accInfoResp)
	if accInfoResp != nil {
		assert.NotNil(accInfoResp)
		if accInfoResp.GetAccount() != nil {
			acc := accInfoResp.GetAccount()
			assert.Equal("user-1", acc.GetUsername())
		}
	}
}

func TestChooseCharacterConfigAdminEligibility(t *testing.T) {
	t.Parallel()
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

	db, err := dbServer.DB()
	require.NoError(err)

	_, err = db.ExecContext(ctx, "UPDATE fivenet_accounts SET last_char = ? WHERE id = ?", 5, 1)
	require.NoError(err)

	client := pbauth.NewAuthServiceClient(clientConn)
	accountToken := loginAndGetAccountToken(t, client, ctx, "user-1")

	srv.jobAdminGroups = nil
	srv.jobAdminUsers = nil
	srv.configAdminGroups = nil
	srv.configAdminUsers = nil

	currentCfg := srv.appCfg.Get()
	require.NotNil(currentCfg)
	if currentCfg.GetAuth() == nil {
		currentCfg.Auth = &settings.Auth{}
	}
	currentCfg.GetAuth().SetLastCharLock(true)
	currentCfg.GetAuth().SetJobAdminGroups(nil)
	currentCfg.GetAuth().SetJobAdminUsers(nil)
	currentCfg.GetAuth().SetConfigAdminGroups(nil)
	currentCfg.GetAuth().SetConfigAdminUsers([]string{"3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4"})
	srv.appCfg.Set(currentCfg)

	chooseCtx := metadata.NewOutgoingContext(
		ctx,
		metadata.New(map[string]string{
			"Cookie": auth.AccCookieName + "=" + accountToken,
		}),
	)
	chooseHeaders := metadata.New(nil)
	chooseCharRes, err := client.ChooseCharacter(chooseCtx, &pbauth.ChooseCharacterRequest{
		CharId: 1,
	}, grpc.Header(&chooseHeaders))
	require.NoError(err)
	require.NotNil(chooseCharRes)
	require.NotNil(chooseCharRes.GetChar())
	assert.Equal(int32(1), chooseCharRes.GetChar().GetUserId())
	assert.Equal("user-1", chooseCharRes.GetUsername())
	assert.True(hasPermission(chooseCharRes.GetPermissions(), perms.PermCanBeSuperuser))

	updatedAccountToken := accountToken
	for _, cookieHeader := range chooseHeaders.Get("set-cookie") {
		cookie, err := http.ParseSetCookie(cookieHeader)
		require.NoError(err)
		if cookie.Name == auth.AccCookieName {
			updatedAccountToken = cookie.Value
			break
		}
	}

	parsedAccClaims, err := srv.tm.ParseAccToken(updatedAccountToken)
	require.NoError(err)
	assert.Equal(int64(1), parsedAccClaims.AccID)

	parsedUserClaims, err := srv.tm.ParseUserToken(chooseCharRes.GetToken())
	require.NoError(err)
	assert.Equal(int32(1), parsedUserClaims.UserID)

	impersonateCtx := metadata.NewOutgoingContext(
		ctx,
		metadata.New(map[string]string{
			"Cookie":        auth.AccCookieName + "=" + updatedAccountToken,
			"Authorization": "Bearer " + chooseCharRes.GetToken(),
		}),
	)
	impersonateRes, err := client.ImpersonateJob(impersonateCtx, &pbauth.ImpersonateJobRequest{
		JobGrade: 1,
	})
	require.NoError(err)
	require.NotNil(impersonateRes)
	assert.True(impersonateRes.GetState())
	assert.True(hasPermission(impersonateRes.GetPermissions(), perms.PermCanBeSuperuser))
}

func loginAndGetAccountToken(
	t *testing.T,
	client pbauth.AuthServiceClient,
	ctx context.Context,
	username string,
) string {
	t.Helper()

	loginReq := &pbauth.LoginRequest{
		Username: username,
		Password: "password",
	}
	md := metadata.New(map[string]string{})
	res, err := client.Login(ctx, loginReq, grpc.Header(&md))
	require.NoError(t, err)
	require.NotNil(t, res)

	cookies := md.Get("set-cookie")
	require.Len(t, cookies, 2, "Expected 2 cookies to be set")

	foundToken := -1
	for i, c := range cookies {
		if strings.HasPrefix(c, auth.AccCookieName+"=") {
			foundToken = i
			break
		}
	}
	require.NotEqual(t, -1, foundToken, "Expected a cookie starting with '"+auth.AccCookieName+"='")

	cookie, err := http.ParseSetCookie(cookies[foundToken])
	require.NoError(t, err)
	return cookie.Value
}

func hasPermission(
	perms []*permissionspermissions.Permission,
	want *permissionspermissions.Permission,
) bool {
	for _, perm := range perms {
		if perm == nil {
			continue
		}
		if perm.GetNamespace() == want.GetNamespace() &&
			perm.GetService() == want.GetService() &&
			perm.GetName() == want.GetName() {
			return true
		}
	}

	return false
}
