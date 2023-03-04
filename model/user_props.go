package model

const TableNameUserProps = "arpanet_user_props"

type UserProps struct {
	Identifier string `gorm:"primaryKey;type:varchar(64)"`

	Wanted *bool `gorm:"index"`
}

// TableName UserProps' table name
func (*UserProps) TableName() string {
	return TableNameUserProps
}
