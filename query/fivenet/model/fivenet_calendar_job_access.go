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

type FivenetCalendarJobAccess struct {
	ID           uint64     `sql:"primary_key" json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	CalendarID   uint64     `json:"calendar_id"`
	Job          string     `json:"job"`
	MinimumGrade int32      `json:"minimum_grade"`
	Access       int16      `json:"access"`
}
