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

type ArpanetDocumentsUserAccess struct {
	ID         uint64     `sql:"primary_key" json:"id"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DocumentID *uint64    `json:"document_id"`
	UserID     int32      `json:"user_id"`
	Access     int16      `json:"access"`
	CreatorID  int32      `json:"creator_id"`
}
