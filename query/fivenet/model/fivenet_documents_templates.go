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

type FivenetDocumentsTemplates struct {
	ID           uint64     `sql:"primary_key" json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	Weight       *uint32    `json:"weight"`
	CategoryID   *uint64    `json:"category_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Color        *string    `json:"color"`
	Icon         *string    `json:"icon"`
	ContentTitle string     `json:"content_title"`
	Content      string     `json:"content"`
	State        string     `json:"state"`
	Access       *string    `json:"access"`
	Schema       *string    `json:"schema"`
	Workflow     *string    `json:"workflow"`
	CreatorJob   string     `json:"creator_job"`
}
