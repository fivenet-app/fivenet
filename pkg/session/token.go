package session

import (
	"errors"
	"strings"

	"github.com/galexrt/arpanet/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type CitizenInfoClaims struct {
	AccountID            uint   `json:"accid"`
	Username             string `json:"usrname"`
	ActiveCharID         uint64 `json:"actcharid"`
	ActiveCharIdentifier string `json:"actcharident"`

	jwt.RegisteredClaims
}

var Tokens *TokenManager

type TokenManager struct {
	jwtSigningKey []byte
}

func NewTokenManager() *TokenManager {
	return &TokenManager{
		jwtSigningKey: []byte(strings.TrimSpace(config.C.JWT.Secret)),
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
