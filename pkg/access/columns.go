package access

import (
	"github.com/go-jet/jet/v2/mysql"
)

// BaseAccessColumns defines the common columns for access control tables.
type BaseAccessColumns struct {
	// ID is the primary key column.
	ID mysql.ColumnInteger
	// TargetID is the column referencing the target entity.
	TargetID mysql.ColumnInteger
	// Access is the column representing access level or permissions.
	Access mysql.ColumnInteger
}

// JobAccessColumns defines columns for job-based access control, embedding BaseAccessColumns.
type JobAccessColumns struct {
	BaseAccessColumns

	// Job is the column for the job name.
	Job mysql.ColumnString
	// MinimumGrade is the column for the minimum grade required for access.
	MinimumGrade mysql.ColumnInteger
}

// UserAccessColumns defines columns for user-based access control, embedding BaseAccessColumns.
type UserAccessColumns struct {
	BaseAccessColumns

	// UserID is the column for the user id.
	UserID mysql.ColumnInteger
}

// QualificationAccessColumns defines columns for qualification-based access control, embedding BaseAccessColumns.
type QualificationAccessColumns struct {
	BaseAccessColumns

	// QualificationID is the column for the qualification id.
	QualificationID mysql.ColumnInteger
}

// TargetTableColumns defines common columns for target tables in access control.
type TargetTableColumns struct {
	// ID is the primary key column.
	ID mysql.ColumnInteger
	// DeletedAt is the column for soft deletion timestamps.
	DeletedAt mysql.ColumnTimestamp

	// CreatorID is the column for the creator's user ID.
	CreatorID mysql.ColumnInteger
	// CreatorJob is the column for the creator's job or role.
	CreatorJob mysql.ColumnString
}
