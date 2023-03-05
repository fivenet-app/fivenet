package session

import "encoding/gob"

const (
	UserSession    = "arpanet_user"
	UserIDKey      = "LoggedInUserID"
	CitizenInfoKey = "LoggedInCitizenInfo"
)

var Names = []string{
	UserSession,
}

type CitizenInfo struct {
	ID int `json:"id"`

	Identifier   string `json:"identifier"`
	ActiveCharID uint64 `json:"act_char_id"`

	Job      string `json:"job"`
	JobGrade int    `json:"job_grade"`
}

func init() {
	gob.Register(&CitizenInfo{})
}
