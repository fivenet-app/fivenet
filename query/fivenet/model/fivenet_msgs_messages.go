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

type FivenetMsgsMessages struct {
	ID        uint64     `sql:"primary_key" json:"id"`
	ThreadID  uint64     `json:"thread_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Message   string     `json:"message"`
	Data      *string    `json:"data"`
	CreatorID int32      `json:"creator_id"`
}
