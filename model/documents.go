package model

import "gorm.io/gorm"

type ContentType string

const (
	PlaintextContentType ContentType = "plaintext"
	MarkdownContentType  ContentType = "markdown"
)

const TableNameDocument = "arpanet_documents"

type Document struct {
	gorm.Model

	Title       string
	Content     string
	ContentType ContentType
	//Creator     Citizen
}

// TableName Document's table name
func (*Document) TableName() string {
	return TableNameDocument
}
