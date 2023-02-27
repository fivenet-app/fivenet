package session

import "encoding/gob"

const (
	UserSession = "arpanet_user"
	UserIDKey   = "LoggedInUserID"
	UserInfoKey = "LoggedInUserInfo"
)

var Names = []string{
	UserSession,
}

type UserInfo struct {
	ID int `json:"id"`

	Identifier string `json:"identifier"`
	CharIndex  int    `json:"char_index"`

	Job      string `json:"job"`
	JobGrade int    `json:"job_grade"`
}

func init() {
	gob.Register(&UserInfo{})
}
