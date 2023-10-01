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

type FivenetCentrumDispatches struct {
	ID          uint64     `sql:"primary_key" json:"id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Job         string     `json:"job"`
	Message     string     `json:"message"`
	Description *string    `json:"description"`
	Attributes  *string    `json:"attributes"`
	X           *float64   `json:"x"`
	Y           *float64   `json:"y"`
	Postal      *string    `json:"postal"`
	Anon        bool       `json:"anon"`
	CreatorID   int32      `json:"creator_id"`
}
