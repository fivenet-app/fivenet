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

type FivenetDocumentsCategories struct {
	ID          uint64     `sql:"primary_key" json:"id"`
	CreatedAt   *time.Time `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Name        string     `json:"name"`
	SortKey     *string    `json:"sort_key"`
	Description *string    `json:"description"`
	Job         string     `json:"job"`
	Color       *string    `json:"color"`
	Icon        *string    `json:"icon"`
}
