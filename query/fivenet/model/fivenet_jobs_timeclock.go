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

type FivenetJobsTimeclock struct {
	Job       string     `json:"job"`
	UserID    int32      `json:"user_id"`
	Date      time.Time  `json:"date"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	SpentTime *float64   `json:"spent_time"`
}
