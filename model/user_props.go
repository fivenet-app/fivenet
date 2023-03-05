package model

const TableNameUserProps = "arpanet_user_props"

type UserProps struct {
	UserID int32 `gorm:"index;type:int(11)"`

	Wanted *bool `gorm:"index"`
}

// TableName UserProps' table name
func (*UserProps) TableName() string {
	return TableNameUserProps
}
