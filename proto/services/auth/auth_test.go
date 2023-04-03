package auth

import (
	"context"
	"os"
	"testing"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms/mock"
	"github.com/galexrt/fivenet/tests/dbmanager"
	"github.com/galexrt/fivenet/tests/proto"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	dbmanager.TestDBManager.Setup()

	code := m.Run()

	dbmanager.TestDBManager.Stop()

	os.Exit(code)
}

func TestFullAuthFlow(t *testing.T) {
	defer dbmanager.TestDBManager.Reset()

	db := dbmanager.TestDBManager.DB()

	ctx := context.Background()
	tm := auth.NewTokenManager("")
	p := mock.NewMock()
	c, err := mstlystcdata.NewCache(ctx, zap.NewNop(), db)
	assert.NoError(t, err)
	enricher := mstlystcdata.NewEnricher(c)
	srv := NewServer(db, auth.NewGRPCAuth(tm), tm, p, enricher)

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
	chooseCharReq.CharId = 1 // Char id 1 is not `user-1`'s char
	chooseCharRes, err := client.ChooseCharacter(ctx, chooseCharReq)
	assert.Error(t, err)
	assert.Nil(t, chooseCharRes)
	proto.CompareGRPCError(t, UnableToChooseCharErr, err)

	// user-1: Choose valid character but we don't have permissions
	chooseCharReq.CharId = 2
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	assert.Error(t, err)
	assert.Nil(t, chooseCharRes)
	proto.CompareGRPCError(t, UnableToChooseCharErr, err)

	// user-1: Choose valid character, now we add a permssion
	p.AddUserPerm(1, "test123-perm")
	chooseCharReq.CharId = 2
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	assert.Error(t, err)
	assert.Nil(t, chooseCharRes)
	proto.CompareGRPCError(t, UnableToChooseCharErr, err)
}
