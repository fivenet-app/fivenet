//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type GksphoneSettings struct {
	ID          int32   `sql:"primary_key" json:"id"`
	Identifier  *string `json:"identifier"`
	Crypto      *string `json:"crypto"`
	PhoneNumber *string `json:"phone_number"`
	AvatarURL   *string `json:"avatar_url"`
}
