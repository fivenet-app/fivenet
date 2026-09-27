package model

import "time"

type FivenetJobAssets struct {
	Job             string    `sql:"primary_key" json:"job"`
	FileID          int64     `sql:"primary_key" json:"file_id"`
	CreatedByUserID int32     `                  json:"created_by_user_id"`
	DisplayName     string    `                  json:"display_name"`
	CreatedAt       time.Time `                  json:"created_at"`
}
