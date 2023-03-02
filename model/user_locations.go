package model

import "time"

const TableNameUserLocation = "arpanet_user_locations"

type UserLocation struct {
	Identifier string `gorm:"primaryKey"`
	Job        string `gorm:"index"`

	X float32
	Y float32

	Hidden bool

	UpdatedAt time.Time
}

func (*UserLocation) TableName() string {
	return TableNameUserLocation
}
