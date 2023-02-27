package session

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserInfoClaims struct {
	License string `json:"ident"`

	jwt.RegisteredClaims
}
