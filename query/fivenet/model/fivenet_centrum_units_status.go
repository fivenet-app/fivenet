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

type FivenetCentrumUnitsStatus struct {
	ID        uint64     `sql:"primary_key" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UnitID    uint64     `json:"unit_id"`
	Status    int16      `json:"status"`
	Reason    *string    `json:"reason"`
	Code      *string    `json:"code"`
	UserID    *int32     `json:"user_id"`
	X         *float64   `json:"x"`
	Y         *float64   `json:"y"`
}
