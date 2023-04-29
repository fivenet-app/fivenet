package auth

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/fivenet/proto/resources/users"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpireTime  = 16 * time.Hour
	TokenRenewalTime = 4 * time.Hour
	TokenMaxRenews   = 5
)

type CitizenInfoClaims struct {
	AccountID          uint64 `json:"accid"`
	Username           string `json:"usrnm"`
	ActiveCharID       int32  `json:"chrid"`
	ActiveCharJob      string `json:"chrjb"`
	ActiveCharJobGrade int32  `json:"chrjbg"`
	RenewedCount       int32  `json:"renwc"`

	jwt.RegisteredClaims
}

type TokenManager struct {
	jwtSigningKey []byte
}

func NewTokenManager(jwtSecret string) *TokenManager {
	return &TokenManager{
		jwtSigningKey: []byte(strings.TrimSpace(jwtSecret)),
	}
}

func (t *TokenManager) NewWithClaims(claims *CitizenInfoClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSigningKey)
}

func (t *TokenManager) ParseWithClaims(tokenString string) (*CitizenInfoClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CitizenInfoClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", errors.New("failed to verify jwt token method")
		}
		return t.jwtSigningKey, nil
	})
	if err != nil {
		return nil, errors.New("failed to parse jwt token")
	}

	claims, ok := token.Claims.(*CitizenInfoClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("failed to parse token claims")
}

func BuildTokenClaimsFromAccount(account *model.FivenetAccounts, activeChar *users.User) *CitizenInfoClaims {
	claims := &CitizenInfoClaims{
		AccountID:    account.ID,
		Username:     *account.Username,
		RenewedCount: 0,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "fivenet",
			Subject:  account.License,
			ID:       strconv.FormatUint(uint64(account.ID), 10),
			Audience: []string{"fivenet"},
		},
	}
	SetTokenClaimsTimes(claims)

	if activeChar != nil {
		claims.ActiveCharID = activeChar.UserId
		claims.ActiveCharJob = activeChar.Job
		claims.ActiveCharJobGrade = activeChar.JobGrade
	} else {
		claims.ActiveCharID = 0
		claims.ActiveCharJob = ""
		claims.ActiveCharJobGrade = 0
	}

	return claims
}

func SetTokenClaimsTimes(claims *CitizenInfoClaims) {
	now := time.Now()
	// A usual scenario is to set the expiration time relative to the current time
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(TokenExpireTime))
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
	claims.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
}
