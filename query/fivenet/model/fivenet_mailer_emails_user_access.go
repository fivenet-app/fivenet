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

type FivenetMailerEmailsUserAccess struct {
	ID        uint64     `sql:"primary_key" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	EmailID   uint64     `json:"email_id"`
	UserID    int32      `json:"user_id"`
	Access    int16      `json:"access"`
}
