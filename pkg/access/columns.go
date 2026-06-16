package access

import "github.com/go-jet/jet/v2/mysql"

type BaseAccessColumns struct {
	ID       mysql.ColumnInteger
	TargetID mysql.ColumnInteger
	Access   mysql.ColumnInteger
}

type VisibilityColumns struct {
	BaseAccessColumns

	RuleKind   mysql.ColumnInteger
	SubjectID  mysql.ColumnInteger
	CreatorID  mysql.ColumnInteger
	CreatorJob mysql.ColumnString
	Effect     mysql.ColumnBool
}
