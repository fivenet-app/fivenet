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

type FivenetLawbooksLaws struct {
	ID            uint64     `sql:"primary_key" json:"id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	LawbookID     uint64     `json:"lawbook_id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	Fine          *uint64    `json:"fine"`
	DetentionTime *uint64    `json:"detention_time"`
	StvoPoints    *uint64    `json:"stvo_points"`
}