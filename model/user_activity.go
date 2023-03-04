package model

type AcitvityType string

const (
	ChangedType AcitvityType = "changed"
)

const TableNameUserActivity = "arpanet_documents_mentions"

type UserActivity struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Identifier string `gorm:"index;type:varchar(64)" json:"identifier"`

	AcitvityType AcitvityType `gorm:"column:type" json:"type"`

	OldValue        string `json:"old_value"`
	NewValue        string `json:"new_value"`
	Reason          string `json:"reason"`
	CauseIdentifier string `gorm:"index;type:varchar(64)" json:"cause"`
}

// TableName UserActivity's table name
func (*UserActivity) TableName() string {
	return TableNameUserActivity
}
