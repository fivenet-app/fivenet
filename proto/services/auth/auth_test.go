package auth

import (
	"context"
	"testing"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms/mock"
	"github.com/galexrt/arpanet/tests/dbmanager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	dbmanager.TestDBManager.Setup()

	m.Run()

	dbmanager.TestDBManager.Stop()
}

func TestLogin(t *testing.T) {
	defer dbmanager.TestDBManager.Reset()

	tm := auth.NewTokenManager("")
	p := mock.NewMock()
	srv := NewServer(dbmanager.TestDBManager.DB(), auth.NewGRPCAuth(tm), tm, p)

	client, _, cancel := NewTestAuthServiceClient(srv)
	defer cancel()

	ctx := context.Background()
	req := &LoginRequest{}

	// First login without credentials
	req.Username = ""
	req.Password = ""
	res, err := client.Login(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, res)
}
