package auth

import (
	"context"
	"os"
	"testing"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms/mock"
	"github.com/galexrt/arpanet/tests/dbmanager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestMain(m *testing.M) {
	dbmanager.TestDBManager.Setup()

	code := m.Run()

	dbmanager.TestDBManager.Stop()

	os.Exit(code)
}

func TestFullAuthFlow(t *testing.T) {
	defer dbmanager.TestDBManager.Reset()

	tm := auth.NewTokenManager("")
	p := mock.NewMock()
	srv := NewServer(dbmanager.TestDBManager.DB(), auth.NewGRPCAuth(tm), tm, p)

	client, _, cancel := NewTestAuthServiceClient(srv)
	defer cancel()

	ctx := context.Background()

	// First login without credentials
	loginReq := &LoginRequest{}
	loginReq.Username = ""
	loginReq.Password = ""
	res, err := client.Login(ctx, loginReq)
	assert.Error(t, err)
	assert.Nil(t, res)
	s := status.FromContextError(err)
	assert.Equal(t, codes.Unauthenticated, s.Code())
	assert.Equal(t, InvalidLoginErr, err)

	// Login with invalid credentials
	loginReq.Username = "non-existant-username"
	loginReq.Password = "non-existant-password"
	res, err = client.Login(ctx, loginReq)
	assert.Error(t, err)
	assert.Nil(t, res)
	s = status.FromContextError(err)
	assert.Equal(t, codes.Unauthenticated, s.Code())
	assert.Equal(t, InvalidLoginErr, err)

	// user-1: Login with valid account that has one char
	loginReq.Username = "user-1"
	loginReq.Password = "test-password"
	res, err = client.Login(ctx, loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	if res == nil {
		assert.FailNow(t, "user-1: Login with valid account failed, response is nil")
	}
	assert.NotEmpty(t, res.GetToken())

	// user-1: Create authenticated metadate and get characters (only has one char)
	md := metadata.New(map[string]string{"Authorization": "Bearer " + res.GetToken()})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq := &GetCharactersRequest{}
	getCharsRes, err := client.GetCharacters(ctx, getCharsReq)
	assert.NoError(t, err)
	assert.NotNil(t, getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(t, "user-1: Empty char list returned for valid account that should have 2 chars")
	}
	assert.Len(t, getCharsRes.GetChars(), 1)

	// user-2: Login with valid account (2 chars)
	loginReq.Username = "user-2"
	loginReq.Password = "test-password"
	res, err = client.Login(ctx, loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	if res == nil {
		assert.FailNow(t, "user-2: Login with valid account failed, response is nil")
	}
	assert.NotEmpty(t, res.GetToken())

	// user-2: Create authenticated metadate and get characters
	md = metadata.New(map[string]string{"Authorization": "Bearer " + res.GetToken()})
	ctx = metadata.NewOutgoingContext(ctx, md)
	getCharsReq = &GetCharactersRequest{}
	getCharsRes, err = client.GetCharacters(ctx, getCharsReq)
	assert.NoError(t, err)
	assert.NotNil(t, getCharsRes)
	if getCharsRes == nil {
		assert.FailNow(t, "user-2: Empty char list returned for valid account that should have 2 chars")
	}
	assert.Len(t, getCharsRes.GetChars(), 2)

	// user-2: Choose an invalid character
	chooseCharReq := &ChooseCharacterRequest{}
	chooseCharReq.CharId = 1 // Not user-2's char
	chooseCharRes, err := client.ChooseCharacter(ctx, chooseCharReq)
	assert.NoError(t, err)
	assert.Nil(t, chooseCharRes)
	s = status.FromContextError(err)
	assert.Equal(t, codes.OutOfRange, s.Code())
	assert.Equal(t, "user-2: That's not your character!", s.Message())

	// user-2: Choose valid character but we don't have permissions
	chooseCharReq.CharId = 2
	chooseCharRes, err = client.ChooseCharacter(ctx, chooseCharReq)
	assert.Error(t, err)
	assert.Nil(t, chooseCharRes)
	s = status.FromContextError(err)
	assert.Equal(t, codes.OutOfRange, s.Code())
	assert.Equal(t, UnableToChooseCharErr, err)
}
