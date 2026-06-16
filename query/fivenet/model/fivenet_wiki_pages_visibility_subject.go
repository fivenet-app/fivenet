package model

import "time"

type FivenetWikiPagesVisibilitySubject struct {
	TargetID  int64     `sql:"primary_key" json:"target_id"`
	SubjectID int64     `                  json:"subject_id"`
	Access    int32     `                  json:"access"`
	Effect    bool      `                  json:"effect"`
	CreatedAt time.Time `                  json:"created_at"`
}
