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
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	DocumentID uint `gorm:"index" json:"document_id"`

	Name         string     `json:"name"`
	MinimumGrade int        `json:"minimum_grade"`
	Access       AccessRole `gorm:"type:varchar(12)" json:"access"`
}

// TableName DocumentJobAccess's table name
func (*DocumentJobAccess) TableName() string {
	return TableNameDocumentJobAccess
}

const TableNameDocumentUserAccess = "arpanet_documents_user_access"

type DocumentUserAccess struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	DocumentID uint       `gorm:"index" json:"document_id"`
	UserID     int32      `gorm:"index;type:int(11)" json:"identifier"`
	Access     AccessRole `gorm:"type:varchar(12)" json:"access"`
}

// TableName DocumentUserAccess's table name
func (*DocumentUserAccess) TableName() string {
	return TableNameDocumentUserAccess
}
