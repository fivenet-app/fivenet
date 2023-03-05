package model

import "time"

type AcitvityType string

const (
	ChangedActivityType   AcitvityType = "changed"
	CreatedActivityType   AcitvityType = "created"
	MentionedActivityType AcitvityType = "mentioned"
)

const TableNameUserActivity = "arpanet_user_activity"

// TODO Create a second table that is used to map the user IDs to the UserActivity entries?
type UserActivity struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	TargetUserID int32 `gorm:"index;type:int(11);not null"`
	CauseUserID  int32 `gorm:"index;type:int(11);not null"`

	Type AcitvityType `gorm:"column:type" json:"type"`

	Key      string `gorm:"type:varchar(64)" json:"key"`
	OldValue string `gorm:"type:varchar(256)" json:"old_value"`
	NewValue string `gorm:"type:varchar(256)" json:"new_value"`

	Reason string `json:"reason"`
}

// TableName UserActivity's table name
func (*UserActivity) TableName() string {
	return TableNameUserActivity
}
