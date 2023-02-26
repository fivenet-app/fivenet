package model

import "time"

const TableNameAccounts = "arpanet_accounts"

type Accounts struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Enabled bool

	Password          string `gorm:"-" json:"-"`
	PasswordEncrypted string `json:"-"`
	Salt              string `json:"-"`

	ExpiredAt int64
}

func (*Accounts) TableName() string {
	return TableNameAccounts
}
