package access

import (
	jet "github.com/go-jet/jet/v2/mysql"
)

// BaseAccessColumns defines the common columns for access control tables.
type BaseAccessColumns struct {
	// ID is the primary key column.
	ID jet.ColumnInteger
	// TargetID is the column referencing the target entity.
	TargetID jet.ColumnInteger
	// Access is the column representing access level or permissions.
	Access jet.ColumnInteger
}

// JobAccessColumns defines columns for job-based access control, embedding BaseAccessColumns.
type JobAccessColumns struct {
	BaseAccessColumns

	// Job is the column for the job name or identifier.
	Job jet.ColumnString
	// MinimumGrade is the column for the minimum grade required for access.
	MinimumGrade jet.ColumnInteger
}

// UserAccessColumns defines columns for user-based access control, embedding BaseAccessColumns.
type UserAccessColumns struct {
	BaseAccessColumns

	// UserId is the column for the user identifier.
	UserId jet.ColumnInteger
}

// QualificationAccessColumns defines columns for qualification-based access control, embedding BaseAccessColumns.
type QualificationAccessColumns struct {
	BaseAccessColumns

	// QualificationId is the column for the qualification identifier.
	QualificationId jet.ColumnInteger
}

// TargetTableColumns defines common columns for target tables in access control.
type TargetTableColumns struct {
	// ID is the primary key column.
	ID jet.ColumnInteger
	// DeletedAt is the column for soft deletion timestamps.
	DeletedAt jet.ColumnTimestamp

	// CreatorID is the column for the creator's user ID.
	CreatorID jet.ColumnInteger
	// CreatorJob is the column for the creator's job or role.
	CreatorJob jet.ColumnString
}
