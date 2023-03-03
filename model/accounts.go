package model

import (
	"time"

	"github.com/Permify/go-role/models"
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

const TableNameAccountUser = "arpanet_accounts_user"

type AccountUser struct {
	ID         uint   `gorm:"primarykey"`
	Identifier string `gorm:"index"`

	// permify
	Roles       []models.Role       `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permissions []models.Permission `gorm:"many2many:user_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*AccountUser) TableName() string {
	return TableNameAccountUser
}
