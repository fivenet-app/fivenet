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

type FivenetCentrumMarkers struct {
	ID          uint64     `sql:"primary_key" json:"id"`
	CreatedAt   *time.Time `json:"created_at"`
	Job         *string    `json:"job"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	X           *float64   `json:"x"`
	Y           *float64   `json:"y"`
	Color       *string    `json:"color"`
	Icon        *string    `json:"icon"`
	MarkerType  int16      `json:"marker_type"`
	MarkerData  *[]byte    `json:"marker_data"`
	CreatorID   *int32     `json:"creator_id"`
}
