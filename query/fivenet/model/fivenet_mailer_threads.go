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

type FivenetMailerThreads struct {
	ID             uint64     `sql:"primary_key" json:"id"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	Title          string     `json:"title"`
	CreatorEmailID uint64     `json:"creator_email_id"`
	CreatorID      *int32     `json:"creator_id"`
	CreatorEmail   string     `json:"creator_email"`
}
