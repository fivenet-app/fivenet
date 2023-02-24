package model

import "gorm.io/gorm"

const TableNameDocument = "arpanet_documents"

type Document struct {
	gorm.Model

	Title       string      `json:"title"`
	Content     string      `json:"content"`
	ContentType ContentType `gorm:"column:content_type;type:varchar(24)" json:"content_type"`
	Creator     string      `gorm:"column:creator;index:arpanet_documents_FK,priority:1" json:"creator"`
	Public      bool        `gorm:"column:public" json:"public"`
	Jobs        []DocumentJobAccess
	Users       []DocumentUserAccess
}

// TableName Document's table name
func (*Document) TableName() string {
	return TableNameDocument
}
