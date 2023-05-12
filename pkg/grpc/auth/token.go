package auth

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpireTime  = 25 * time.Hour
	TokenRenewalTime = 24 * time.Hour
	TokenMaxRenews   = 5
)

type CitizenInfoClaims struct {
	AccountID    uint64 `json:"accid"`
	Username     string `json:"usrnm"`
	CharID       int32  `json:"chrid"`
	CharJob      string `json:"chrjb"`
	CharJobGrade int32  `json:"chrjbg"`
	RenewedCount int32  `json:"renwc"`

	jwt.RegisteredClaims
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
		claims.CharID = activeChar.UserId
		claims.CharJob = activeChar.Job
		claims.CharJobGrade = activeChar.JobGrade
	} else {
		claims.CharID = 0
		claims.CharJob = ""
		claims.CharJobGrade = 0
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
