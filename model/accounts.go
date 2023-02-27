package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const TableNameAccounts = "arpanet_accounts"

type Account struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Enabled bool

	Username       string `gorm:"index;type:varchar(24)"`
	Password       string `gorm:"-" json:"-"`
	HashedPassword []byte `gorm:"column:password;type:varchar(64)" json:"-"`
	License        string `gorm:"index;type:varchar(64)"`
	ActiveChar     int
}

func (*Account) TableName() string {
	return TableNameAccounts
}

func (a *Account) SetPassword(password string) (string, error) {
	var err error
	a.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(a.HashedPassword), err
}

func (a *Account) CheckPassword(input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.HashedPassword), []byte(input))
	return err == nil
}
