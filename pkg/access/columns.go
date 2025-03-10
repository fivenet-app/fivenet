package access

import (
	jet "github.com/go-jet/jet/v2/mysql"
)

type BaseAccessColumns struct {
	ID       jet.ColumnInteger
	TargetID jet.ColumnInteger
	Access   jet.ColumnInteger
}

type JobAccessColumns struct {
	BaseAccessColumns

	Job          jet.ColumnString
	MinimumGrade jet.ColumnInteger
}

type UserAccessColumns struct {
	BaseAccessColumns

	UserId jet.ColumnInteger
}

type QualificationAccessColumns struct {
	BaseAccessColumns

	QualificationId jet.ColumnInteger
}

type TargetTableColumns struct {
	ID        jet.ColumnInteger
	DeletedAt jet.ColumnTimestamp

	CreatorID  jet.ColumnInteger
	CreatorJob jet.ColumnString
}
