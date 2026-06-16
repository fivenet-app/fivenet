package model

import "time"

type FivenetDocumentsVisibilityPublic struct {
	TargetID  int64     `sql:"primary_key" json:"target_id"`
	CreatedAt time.Time `                  json:"created_at"`
}
