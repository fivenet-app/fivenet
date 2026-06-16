package model

import "time"

type FivenetWikiPagesVisibilityCreator struct {
	TargetID   int64     `sql:"primary_key" json:"target_id"`
	CreatorID  int32     `                  json:"creator_id"`
	CreatorJob string    `                  json:"creator_job"`
	CreatedAt  time.Time `                  json:"created_at"`
}
