package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const jwtTokenTestSecret = "secret-jwt-token-secret-for-testing"

var basicCitizenInfoClaim = &CitizenInfoClaims{
	AccID:    123456,
	Username: "example-username",
	CharID:   987654,
	RegisteredClaims: jwt.RegisteredClaims{
		Subject: "example-subject",
	},
}

// Even though is kinda a duplicate of go JWT lib, I want to make sure we don't have
// issues parsing our custom claims structure.
func TestToken(t *testing.T) {
	tm := NewTokenMgr(jwtTokenTestSecret)
	assert.NotNil(t, tm)
	claims := basicCitizenInfoClaim
	token, err := tm.NewWithClaims(claims)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Make sure we can parse the created token
	parsedClaims, err := tm.ParseWithClaims(token)
	require.NoError(t, err)
	// Custom Claim struct
	assert.Equal(t, claims.AccID, parsedClaims.AccID)
	assert.Equal(t, claims.Username, parsedClaims.Username)
	assert.Equal(t, claims.CharID, parsedClaims.CharID)
	// RegisteredClaims
	assert.Equal(t, claims.Subject, parsedClaims.Subject)

	// Make sure we can't parse the generated token if it is only part of it
	parsedClaims, err = tm.ParseWithClaims(token[5:50])
	require.Error(t, err)
	assert.Nil(t, parsedClaims)

	// Parse a random JWT token from https://jwt.io/ site
	randomToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	parsedClaims, err = tm.ParseWithClaims(randomToken)
	require.Error(t, err)
	assert.Nil(t, parsedClaims)
}
