package auth

import "encoding/gob"

const (
	SessionName = "user"
	keyName     = "user"
)

type SessionInfo struct {
	ID         int    `json:"ID"`
	Identifier string `json:"identifier"`
	CharIndex  int    `json:"charIndex"`

	Job      string `json:"job"`
	JobGrade int    `json:"jobGrade"`
}

func init() {
	gob.Register(&SessionInfo{})
}
