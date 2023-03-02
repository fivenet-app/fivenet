package model

const TableNameUserProps = "arpanet_users_props"

type UserProps struct {
	Identifier string `gorm:"primaryKey"`

	Wanted bool `gorm:"index"`
}

// TableName UserProps' table name
func (*UserProps) TableName() string {
	return TableNameUserProps
}
