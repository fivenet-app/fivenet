package auth

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

const (
	// 4 days.
	TokenExpireTime = 96 * time.Hour
	// 2 days.
	TokenRenewalTime = 48 * time.Hour
)

var ErrFailedJWTVerify = errors.New("failed to verify jwt token method")

var TokenMgrModule = fx.Module("tokenMgr",
	fx.Provide(
		NewTokenMgrFromConfig,
	),
)

type TokenMgr struct {
	jwtSigningKey []byte
}

func NewTokenMgrFromConfig(cfg *config.Config) *TokenMgr {
	return NewTokenMgr(cfg.JWT.Secret)
}

func NewTokenMgr(jwtSecret string) *TokenMgr {
	return &TokenMgr{
		jwtSigningKey: []byte(strings.TrimSpace(jwtSecret)),
	}
}

func (t *TokenMgr) FromAccClaims(claims *authclaims.AccountInfoClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSigningKey)
}

func (t *TokenMgr) FromUserClaims(claims *authclaims.UserInfoClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSigningKey)
}

func (t *TokenMgr) FromCombinedClaims(claims *authclaims.CombinedClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSigningKey)
}

func (t *TokenMgr) ParseAccToken(tokenString string) (*authclaims.AccountInfoClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&authclaims.AccountInfoClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", ErrFailedJWTVerify
			}
			return t.jwtSigningKey, nil
		},
	)
	if err != nil {
		return nil, errors.New("failed to parse jwt acc token")
	}

	claims, ok := token.Claims.(*authclaims.AccountInfoClaims)
	if ok && token.Valid {
		// Ensure AccID is set
		if claims.AccID > 0 {
			return claims, nil
		}
	}

	return nil, errors.New("failed to parse acc token claims")
}

func (t *TokenMgr) ParseUserToken(tokenString string) (*authclaims.UserInfoClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&authclaims.UserInfoClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", ErrFailedJWTVerify
			}
			return t.jwtSigningKey, nil
		},
	)
	if err != nil {
		return nil, errors.New("failed to parse jwt user token")
	}

	claims, ok := token.Claims.(*authclaims.UserInfoClaims)
	if ok && token.Valid {
		// Ensure at least UserID is set "correctly"
		if claims.UserID > 0 {
			return claims, nil
		}
	}

	return nil, errors.New("failed to parse user token claims")
}

func (t *TokenMgr) ParseCombinedToken(tokenString string) (*authclaims.CombinedClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&authclaims.CombinedClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", errors.New("failed to verify jwt token method")
			}
			return t.jwtSigningKey, nil
		},
	)
	if err != nil {
		return nil, errors.New("failed to parse jwt user token")
	}

	claims, ok := token.Claims.(*authclaims.CombinedClaims)
	if ok && token.Valid {
		// Ensure at least AccID is set "correctly"
		if claims.AccID > 0 {
			return claims, nil
		}
	}

	return nil, errors.New("failed to parse user token claims")
}

func MapAccountToClaims(account *accounts.Account) *authclaims.AccountInfoClaims {
	accClaims := &authclaims.AccountInfoClaims{
		AccID:    account.Id,
		Username: account.Username,
		Groups:   account.GetGroups().GetGroups(),

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "fivenet",
			Subject:  account.License,
			ID:       strconv.FormatInt(account.Id, 10),
			Audience: []string{"fivenet"},
		},
	}
	setTokenClaimsTimes(&accClaims.RegisteredClaims)

	return accClaims
}

func MapUserToClaims(user *users.User) *authclaims.UserInfoClaims {
	userClaims := &authclaims.UserInfoClaims{
		UserID:   user.GetUserId(),
		Job:      &user.Job,
		JobGrade: &user.JobGrade,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "fivenet",
			ID:       strconv.FormatInt(int64(user.UserId), 10),
			Audience: []string{"fivenet"},
		},
	}
	setTokenClaimsTimes(&userClaims.RegisteredClaims)

	return userClaims
}

// Set the expiration time relative to the current (server) time for the given claims.
func setTokenClaimsTimes(claims *jwt.RegisteredClaims) {
	now := time.Now()
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(TokenExpireTime))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
}
