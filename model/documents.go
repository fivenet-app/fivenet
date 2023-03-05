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
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Type       DocumentType `gorm:"column:content_type;type:varchar(24)" json:"content_type"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	CreatorID  int32        `gorm:"column:creator;type:int(11);index" json:"creator"`
	CreatorJob string       `gorm:"column:creator_job;type:varchar(20);index" json:"creator_job"`
	Public     bool         `gorm:"column:public;default:0" json:"public"`

	ResponseID uint        `gorm:"index" json:"response_id"`
	Responses  []*Document `gorm:"foreignKey:ResponseID" json:"response"`

	Mentions []DocumentMentions `json:"mentions"`

	JobAccess  []DocumentJobAccess  `json:"jobAccess"`
	UserAccess []DocumentUserAccess `json:"userAccess"`
}

// TableName Document's table name
func (*Document) TableName() string {
	return TableNameDocument
}

const TableNameDocumentMentions = "arpanet_documents_mentions"

type DocumentMentions struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	DocumentID uint   `gorm:"index" json:"document_id"`
	Identifier string `gorm:"index;type:varchar(64)" json:"identifier"`
}

// TableName DocumentMentions's table name
func (*DocumentMentions) TableName() string {
	return TableNameDocumentMentions
}
