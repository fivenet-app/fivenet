package model

import "time"

const JobGradeEveryoneAccess = -1

const TableNameDocumentJobAccess = "arpanet_documents_job_access"

type DocumentJobAccess struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	DocumentID uint

	Name         string `json:"name"`
	MinimumGrade int    `json:"grade"`
}

// TableName DocumentJobAccess's table name
func (*DocumentJobAccess) TableName() string {
	return TableNameDocumentJobAccess
}

const TableNameDocumentUserAccess = "arpanet_documents_user_access"

type DocumentUserAccess struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	DocumentID uint
}

// TableName DocumentUserAccess's table name
func (*DocumentUserAccess) TableName() string {
	return TableNameDocumentUserAccess
}
