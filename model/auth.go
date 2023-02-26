package model

import "time"

const TableNameOAuth2Token = "arpanet_oauth2_token"

type OAuth2Token struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time

	ExpiredAt int64

	Code    string `gorm:"type:varchar(512);index"`
	Access  string `gorm:"type:varchar(512)"`
	Refresh string `gorm:"type:varchar(512)"`
	Data    string `gorm:"type:text"`
}

// TableName OAuth2Token's table name
func (*OAuth2Token) TableName() string {
	return TableNameOAuth2Token
}
