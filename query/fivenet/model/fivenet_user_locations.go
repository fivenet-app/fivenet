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

type FivenetUserLocations struct {
	Identifier string     `sql:"primary_key" json:"identifier"`
	Job        string     `json:"job"`
	X          *float64   `json:"x"`
	Y          *float64   `json:"y"`
	Hidden     *bool      `json:"hidden"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
