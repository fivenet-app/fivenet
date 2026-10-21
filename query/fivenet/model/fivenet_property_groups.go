package model

import "time"

type FivenetPropertyGroups struct {
	Scope     string     `sql:"primary_key" json:"scope"`
	GroupKey  string     `sql:"primary_key" json:"group_key"`
	Label     string     `                  json:"label"`
	SortOrder int32      `                  json:"sort_order"`
	DeletedAt *time.Time `                  json:"deleted_at"`
	CreatedAt time.Time  `                  json:"created_at"`
	UpdatedAt time.Time  `                  json:"updated_at"`
}
