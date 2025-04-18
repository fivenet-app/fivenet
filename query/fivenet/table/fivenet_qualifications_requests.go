//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/mysql"
)

var FivenetQualificationsRequests = newFivenetQualificationsRequestsTable("", "fivenet_qualifications_requests", "")

type fivenetQualificationsRequestsTable struct {
	mysql.Table

	// Columns
	CreatedAt       mysql.ColumnTimestamp
	DeletedAt       mysql.ColumnTimestamp
	QualificationID mysql.ColumnInteger
	UserID          mysql.ColumnInteger
	UserComment     mysql.ColumnString
	Status          mysql.ColumnInteger
	ApprovedAt      mysql.ColumnTimestamp
	ApproverComment mysql.ColumnString
	ApproverID      mysql.ColumnInteger
	ApproverJob     mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetQualificationsRequestsTable struct {
	fivenetQualificationsRequestsTable

	NEW fivenetQualificationsRequestsTable
}

// AS creates new FivenetQualificationsRequestsTable with assigned alias
func (a FivenetQualificationsRequestsTable) AS(alias string) *FivenetQualificationsRequestsTable {
	return newFivenetQualificationsRequestsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetQualificationsRequestsTable with assigned schema name
func (a FivenetQualificationsRequestsTable) FromSchema(schemaName string) *FivenetQualificationsRequestsTable {
	return newFivenetQualificationsRequestsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetQualificationsRequestsTable with assigned table prefix
func (a FivenetQualificationsRequestsTable) WithPrefix(prefix string) *FivenetQualificationsRequestsTable {
	return newFivenetQualificationsRequestsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetQualificationsRequestsTable with assigned table suffix
func (a FivenetQualificationsRequestsTable) WithSuffix(suffix string) *FivenetQualificationsRequestsTable {
	return newFivenetQualificationsRequestsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetQualificationsRequestsTable(schemaName, tableName, alias string) *FivenetQualificationsRequestsTable {
	return &FivenetQualificationsRequestsTable{
		fivenetQualificationsRequestsTable: newFivenetQualificationsRequestsTableImpl(schemaName, tableName, alias),
		NEW:                                newFivenetQualificationsRequestsTableImpl("", "new", ""),
	}
}

func newFivenetQualificationsRequestsTableImpl(schemaName, tableName, alias string) fivenetQualificationsRequestsTable {
	var (
		CreatedAtColumn       = mysql.TimestampColumn("created_at")
		DeletedAtColumn       = mysql.TimestampColumn("deleted_at")
		QualificationIDColumn = mysql.IntegerColumn("qualification_id")
		UserIDColumn          = mysql.IntegerColumn("user_id")
		UserCommentColumn     = mysql.StringColumn("user_comment")
		StatusColumn          = mysql.IntegerColumn("status")
		ApprovedAtColumn      = mysql.TimestampColumn("approved_at")
		ApproverCommentColumn = mysql.StringColumn("approver_comment")
		ApproverIDColumn      = mysql.IntegerColumn("approver_id")
		ApproverJobColumn     = mysql.StringColumn("approver_job")
		allColumns            = mysql.ColumnList{CreatedAtColumn, DeletedAtColumn, QualificationIDColumn, UserIDColumn, UserCommentColumn, StatusColumn, ApprovedAtColumn, ApproverCommentColumn, ApproverIDColumn, ApproverJobColumn}
		mutableColumns        = mysql.ColumnList{CreatedAtColumn, DeletedAtColumn, QualificationIDColumn, UserIDColumn, UserCommentColumn, StatusColumn, ApprovedAtColumn, ApproverCommentColumn, ApproverIDColumn, ApproverJobColumn}
		defaultColumns        = mysql.ColumnList{CreatedAtColumn, StatusColumn}
	)

	return fivenetQualificationsRequestsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		CreatedAt:       CreatedAtColumn,
		DeletedAt:       DeletedAtColumn,
		QualificationID: QualificationIDColumn,
		UserID:          UserIDColumn,
		UserComment:     UserCommentColumn,
		Status:          StatusColumn,
		ApprovedAt:      ApprovedAtColumn,
		ApproverComment: ApproverCommentColumn,
		ApproverID:      ApproverIDColumn,
		ApproverJob:     ApproverJobColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
