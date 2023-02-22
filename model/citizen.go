package model

import "time"

type Citizen struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Identifier  string `gorm:"type:varchar(64)"`
	FirstName   string `gorm:"column:firstname;type:varchar(50)"`
	LastName    string `gorm:"column:lastname;type:varchar(50)"`
	DateOfBirth string `gorm:"column:dateofbirth;type:varchar(25)"`
	Group       string `gorm:"type:varchar(50)"`
	Job         Job
	JobGrade    JobGrade
	Sex         Sex `gorm:"type:varchar(1)"`
	Height      int
	PhoneNumber string   `gorm:"column:phone_number;type:varchar(10)"`
	Visum       int      `gorm:"column:visum"`
	Playtime    int      `gorm:"column:playtime"`
	Accounts    Accounts `gorm:"serializer:json"`
}

func (Citizen) TableName() string {
	return "users"
}

type Sex string

const (
	MaleSex   Sex = "m"
	FemaleSex Sex = "f"
)

type Accounts struct {
	BlackMoney int `json:"black_money"`
	Bank       int `json:"bank"`
	Cash       int `json:"money"`
}
