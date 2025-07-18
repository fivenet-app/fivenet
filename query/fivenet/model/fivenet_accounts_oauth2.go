//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type FivenetAccountsOauth2 struct {
	AccountID    uint64     `json:"account_id"`
	CreatedAt    *time.Time `json:"created_at"`
	Provider     string     `json:"provider"`
	ExternalID   string     `json:"external_id"`
	Username     string     `json:"username"`
	Avatar       string     `json:"avatar"`
	AccessToken  *string    `json:"access_token"`
	RefreshToken *string    `json:"refresh_token"`
	TokenType    *string    `json:"token_type"`
	Scope        *string    `json:"scope"`
	ExpiresIn    *int32     `json:"expires_in"`
	ObtainedAt   *time.Time `json:"obtained_at"`
}
