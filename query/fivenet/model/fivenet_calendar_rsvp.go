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

type FivenetCalendarRsvp struct {
	EntryID   uint64     `sql:"primary_key" json:"entry_id"`
	CreatedAt *time.Time `json:"created_at"`
	UserID    int32      `sql:"primary_key" json:"user_id"`
	Response  int16      `json:"response"`
}
