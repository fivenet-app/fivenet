package auth

import (
	"context"
	"os"
	"testing"

	"github.com/galexrt/fivenet/internal/tests/proto"
	"github.com/galexrt/fivenet/internal/tests/servers"
	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	servers.TestDBServer.Setup()
	servers.TestNATSServer.Setup()

	code := m.Run()

	servers.TestDBServer.Stop()
	servers.TestNATSServer.Stop()

	os.Exit(code)
}

func TestFullAuthFlow(t *testing.T) {
	defer servers.TestDBServer.Reset()

	db := servers.TestDBServer.DB()

	ctx := context.Background()
	logger := zap.NewNop()
	tp := tracesdk.NewTracerProvider()
	ui := userinfo.NewMockUserInfoRetriever(map[int32]*userinfo.UserInfo{})
	tm := auth.NewTokenMgr("")

	eventus, err := events.NewEventus(logger, servers.TestNATSServer.GetURL())
	assert.NoError(t, err)

	p, err := perms.New(ctx, logger, db, tp, eventus)
	assert.NoError(t, err)
	defer p.Stop()

	aud := &audit.Noop{}
	c, err := mstlystcdata.NewCache(ctx, logger, tp, db)
	assert.NoError(t, err)
	enricher := mstlystcdata.NewEnricher(c)
	srv := NewServer(db, auth.NewGRPCAuth(ui, tm), tm, p, enricher, aud, ui)

	client, _, cancel := NewTestAuthServiceClient(srv)
	defer cancel()

	// First login without credentials
	loginReq := &LoginRequest{}
	loginReq.Username = ""
	loginReq.Password = ""
	res, err := client.Login(ctx, loginReq)
	assert.Error(t, err)
	assert.Nil(t, res)
	proto.CompareGRPCError(t, InvalidLoginErr, err)

	// Login with invalid credentials
	loginReq.Username = "non-existant-username"
	loginReq.Password = "non-existant-password"
	res, err = client.Login(ctx, loginReq)
	assert.Error(t, err)
	assert.Nil(t, res)
	proto.CompareGRPCError(t, InvalidLoginErr, err)

	// user-3: Login with valid account that has one char
	loginReq.Username = "user-3"
	loginReq.Password = "password"
	res, err = client.Login(ctx, loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	if res == nil {
		assert.FailNow(t, "user-3: Login with valid account failed, response is nil")
	}
	assert.NotEmpty(t, res.GetToken())

	// user-3: Create authenticated metadate and get characters (only has one char)
	md := metadata.New(map[string]string{"Authorization": "Bearer " + res.GetToken()})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq := &GetCharactersRequest{}
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
	res, err = client.Login(ctx, loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	if res == nil {
		assert.FailNow(t, "user-1: Login with valid account failed, response is nil")
	}
	assert.NotEmpty(t, res.GetToken())

	// user-1: Create authenticated metadate and get characters
	md = metadata.New(map[string]string{"Authorization": "Bearer " + res.GetToken()})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq = &GetCharactersRequest{}
	getCharsRes, err = client.GetCharacters(ctx, getCharsReq)
	assert.NoError(t, err)
	assert.NotNil(t, getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(t, "user-1: Empty char list returned for valid account that should have 2 chars")
	}
	assert.Len(t, getCharsRes.GetChars(), 2)

	// user-1: Choose an invalid character
	chooseCharReq := &ChooseCharacterRequest{}
	chooseCharReq.CharId = 2 // Char id 2 is `user-2`'s char
	chooseCharRes, err := client.ChooseCharacter(ctx, chooseCharReq)
	assert.Error(t, err)
	assert.Nil(t, chooseCharRes)
	proto.CompareGRPCError(t, UnableToChooseCharErr, err)

	// user-1: Choose valid character but we don't have permissions on the role
	err = p.UpdateRolePermissions(ctx, 1, perms.AddPerm{
		Id:  1,
		Val: false,
	})
	assert.NoError(t, err)
	// Remove perm from ambulance rank 1 role as `user-1` is a medic
	err = p.UpdateRolePermissions(ctx, 2, perms.AddPerm{
		Id:  1,
		Val: false,
	})
	assert.NoError(t, err)
	chooseCharReq.CharId = 1
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	assert.Error(t, err)
	assert.Nil(t, chooseCharRes)
	proto.CompareGRPCError(t, UnableToChooseCharErr, err)

	// user-1: Choose valid character, now we add a permssion
	// Perm ID 1 is `AuthService.ChooseCharacter`
	err = p.UpdateRolePermissions(ctx, 2, perms.AddPerm{
		Id:  1,
		Val: true,
	})
	assert.NoError(t, err)
	chooseCharReq.CharId = 1
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	assert.NoError(t, err)
	assert.NotNil(t, chooseCharRes)
}
