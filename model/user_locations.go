package model

import "time"

const TableNameUserLocation = "arpanet_user_locations"

type UserLocation struct {
	Identifier string `gorm:"primaryKey;type:varchar(64)"`
	Job        string `gorm:"index;type:varchar(20)"`

	X float32
	Y float32

	Hidden bool

	UpdatedAt time.Time
}

func (*UserLocation) TableName() string {
	return TableNameUserLocation
}
