package auth

import (
	"testing"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const jwtTokenTestSecret = "secret-jwt-token-secret-for-testing"

var testUserCombinedClaim = &authclaims.CombinedClaims{
	AccID:    123456,
	Username: "example-username",
	UserID:   987654,
	RegisteredClaims: jwt.RegisteredClaims{
		Subject: "example-subject",
	},
}

// Even though is kinda a duplicate of go JWT lib, I want to make sure we don't have
// issues parsing our custom claims structure.
func TestToken(t *testing.T) {
	t.Parallel()
	tm := NewTokenMgr(jwtTokenTestSecret)
	assert.NotNil(t, tm)
	claims := testUserCombinedClaim
	token, err := tm.FromCombinedClaims(claims)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Make sure we can parse the created token
	parsedAccClaims, err := tm.ParseAccToken(token)
	require.NoError(t, err)
	// Custom Claim struct
	assert.Equal(t, claims.AccID, parsedAccClaims.AccID)
	assert.Equal(t, claims.Username, parsedAccClaims.Username)
	assert.Equal(t, claims.Subject, parsedAccClaims.Subject)
	// RegisteredClaims
	assert.Equal(t, claims.Subject, parsedAccClaims.Subject)

	// Make sure we can't parse the generated token if it is only part of it
	parsedAccClaims, err = tm.ParseAccToken(token[5:50])
	require.Error(t, err)
	assert.Nil(t, parsedAccClaims)

	// Parse a random JWT token from https://jwt.io/ site to provoke an error
	randomToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	parsedAccClaims, err = tm.ParseAccToken(randomToken)
	require.Error(t, err)
	assert.Nil(t, parsedAccClaims)
}

func TestMapAccountToClaims(t *testing.T) {
	t.Parallel()

	account := &accounts.Account{
		Id:       123,
		Username: "example-user",
		License:  "license-123",
		Groups: &accounts.AccountGroups{
			Groups: []string{"bootstrap-admin"},
		},
	}

	claims := MapAccountToClaims(account, true)
	require.NotNil(t, claims)
	assert.True(t, claims.CanBeSuperuser)
	assert.Equal(t, account.GetGroups().GetGroups(), claims.Groups)
	assert.Equal(t, account.GetUsername(), claims.Username)
	assert.Equal(t, account.GetLicense(), claims.Subject)
}

func TestParseUserTokenReturnsWrappedExpiryError(t *testing.T) {
	t.Parallel()

	tm := NewTokenMgr(jwtTokenTestSecret)
	expiredClaims := &authclaims.UserInfoClaims{
		AccID:  1,
		UserID: 2,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Subject:   "license",
			Audience:  []string{"fivenet"},
			Issuer:    "fivenet",
		},
	}

	token, err := tm.FromUserClaims(expiredClaims)
	require.NoError(t, err)

	parsed, err := tm.ParseUserToken(token)
	require.Error(t, err)
	assert.Nil(t, parsed)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}
