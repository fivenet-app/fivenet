package auth

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpireTime  = 96 * time.Hour
	TokenRenewalTime = 48 * time.Hour
)

type CitizenInfoClaims struct {
	jwt.RegisteredClaims

	AccID    uint64 `json:"accid"`
	Username string `json:"usr"`
	CharID   int32  `json:"chrid"`
}

type TokenMgr struct {
	jwtSigningKey []byte
}

func NewTokenMgr(jwtSecret string) *TokenMgr {
	return &TokenMgr{
		jwtSigningKey: []byte(strings.TrimSpace(jwtSecret)),
	}
}

func (t *TokenMgr) NewWithClaims(claims *CitizenInfoClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSigningKey)
}

func (t *TokenMgr) ParseWithClaims(tokenString string) (*CitizenInfoClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CitizenInfoClaims{},
		func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return "", errors.New("failed to verify jwt token method")
			}
			return t.jwtSigningKey, nil
		},
	)
	if err != nil {
		return nil, errors.New("failed to parse jwt token")
	}

	claims, ok := token.Claims.(*CitizenInfoClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("failed to parse token claims")
}

func BuildTokenClaimsFromAccount(
	account *model.FivenetAccounts,
	activeChar *users.User,
) *CitizenInfoClaims {
	claims := &CitizenInfoClaims{
		AccID:    account.ID,
		Username: *account.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "fivenet",
			Subject:  account.License,
			ID:       strconv.FormatUint(account.ID, 10),
			Audience: []string{"fivenet"},
		},
	}
	SetTokenClaimsTimes(claims)

	if activeChar != nil {
		claims.CharID = activeChar.GetUserId()
	} else {
		claims.CharID = 0
	}

	return claims
}

func SetTokenClaimsTimes(claims *CitizenInfoClaims) {
	now := time.Now()
	// A usual scenario is to set the expiration time relative to the current time
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(TokenExpireTime))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
}
