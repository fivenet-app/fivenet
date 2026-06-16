package model

import "time"

type FivenetDocumentsVisibilityCreator struct {
	TargetID   int64     `sql:"primary_key" json:"target_id"`
	CreatorID  int32     `                  json:"creator_id"`
	CreatorJob string    `                  json:"creator_job"`
	CreatedAt  time.Time `                  json:"created_at"`
}
