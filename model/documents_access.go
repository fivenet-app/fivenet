package model

import "time"

const AnyJobGradeHasAccess = -1

type AccessRole string

const (
	BlockedAccessRole = "blocked"
	ViewAccessRole    = "view"
	EditAccessRole    = "edit"
	LeaderAccessRole  = "leader"
	AdminAccessRole   = "admin"
)

const TableNameDocumentJobAccess = "arpanet_documents_job_access"

type DocumentJobAccess struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	DocumentID uint `gorm:"index"`

	Name         string     `json:"name"`
	MinimumGrade int        `json:"grade"`
	Access       AccessRole `gorm:"type:varchar(12)"`
}

// TableName DocumentJobAccess's table name
func (*DocumentJobAccess) TableName() string {
	return TableNameDocumentJobAccess
}

const TableNameDocumentUserAccess = "arpanet_documents_user_access"

type DocumentUserAccess struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	DocumentID uint       `gorm:"index"`
	Identifier string     `gorm:"index;type:varchar(64)"`
	Access     AccessRole `gorm:"type:varchar(12)"`
}

// TableName DocumentUserAccess's table name
func (*DocumentUserAccess) TableName() string {
	return TableNameDocumentUserAccess
}
