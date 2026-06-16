package model

import "time"

type FivenetWikiPagesVisibilityPublic struct {
	TargetID  int64     `sql:"primary_key" json:"target_id"`
	CreatedAt time.Time `                  json:"created_at"`
}
