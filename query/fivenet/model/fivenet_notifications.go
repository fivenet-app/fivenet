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

type FivenetNotifications struct {
	ID        uint64     `sql:"primary_key" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
	UserID    *int32     `json:"user_id"`
	Job       *string    `json:"job"`
	Title     string     `json:"title"`
	Type      string     `json:"type"`
	Content   *string    `json:"content"`
	Category  int16      `json:"category"`
	Data      *string    `json:"data"`
	Starred   *bool      `json:"starred"`
}
