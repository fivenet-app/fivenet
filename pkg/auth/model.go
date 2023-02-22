package auth

import "encoding/gob"

const (
	SessionName = "user"
	keyName     = "user"
)

type SessionInfo struct {
	ID         int    `json:"ID"`
	Identifier string `json:"identifier"`

	Job      string `json:"job"`
	JobGrade int    `json:"jobGrade"`

	Admin bool `json:"admin"`
}

func init() {
	gob.Register(&SessionInfo{})
}
