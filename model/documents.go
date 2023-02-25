package model

import (
	"time"

	"gorm.io/gorm"
)

type DocumentType string

const (
	PlainDocument DocumentType = "plain"
	FormDocument  DocumentType = "form"
)

const TableNameDocument = "arpanet_documents"

type Document struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Type       DocumentType `gorm:"column:content_type;type:varchar(24)" json:"content_type"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	Creator    string       `gorm:"column:creator;index:arpanet_documents_FK,priority:1" json:"creator"`
	CreatorJob string       `gorm:"column:creator_job;index:arpanet_documents_creator_job_FK,priority:1" json:"creator_job"`
	Public     bool         `gorm:"column:public;default:0" json:"public"`
	ResponseID uint         `gorm:"index" json:"response_id"`
	Responses  []*Document  `gorm:"foreignKey:ResponseID" json:"response"`

	Jobs  []DocumentJobAccess  `json:"jobs"`
	Users []DocumentUserAccess `json:"users"`
}

// TableName Document's table name
func (*Document) TableName() string {
	return TableNameDocument
}
