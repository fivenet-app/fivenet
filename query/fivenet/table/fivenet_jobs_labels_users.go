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

var FivenetJobsLabelsUsers = newFivenetJobsLabelsUsersTable("", "fivenet_jobs_labels_users", "")

type fivenetJobsLabelsUsersTable struct {
	mysql.Table

	// Columns
	UserID  mysql.ColumnInteger
	Job     mysql.ColumnString
	LabelID mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetJobsLabelsUsersTable struct {
	fivenetJobsLabelsUsersTable

	NEW fivenetJobsLabelsUsersTable
}

// AS creates new FivenetJobsLabelsUsersTable with assigned alias
func (a FivenetJobsLabelsUsersTable) AS(alias string) *FivenetJobsLabelsUsersTable {
	return newFivenetJobsLabelsUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetJobsLabelsUsersTable with assigned schema name
func (a FivenetJobsLabelsUsersTable) FromSchema(schemaName string) *FivenetJobsLabelsUsersTable {
	return newFivenetJobsLabelsUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetJobsLabelsUsersTable with assigned table prefix
func (a FivenetJobsLabelsUsersTable) WithPrefix(prefix string) *FivenetJobsLabelsUsersTable {
	return newFivenetJobsLabelsUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetJobsLabelsUsersTable with assigned table suffix
func (a FivenetJobsLabelsUsersTable) WithSuffix(suffix string) *FivenetJobsLabelsUsersTable {
	return newFivenetJobsLabelsUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetJobsLabelsUsersTable(schemaName, tableName, alias string) *FivenetJobsLabelsUsersTable {
	return &FivenetJobsLabelsUsersTable{
		fivenetJobsLabelsUsersTable: newFivenetJobsLabelsUsersTableImpl(schemaName, tableName, alias),
		NEW:                         newFivenetJobsLabelsUsersTableImpl("", "new", ""),
	}
}

func newFivenetJobsLabelsUsersTableImpl(schemaName, tableName, alias string) fivenetJobsLabelsUsersTable {
	var (
		UserIDColumn   = mysql.IntegerColumn("user_id")
		JobColumn      = mysql.StringColumn("job")
		LabelIDColumn  = mysql.IntegerColumn("label_id")
		allColumns     = mysql.ColumnList{UserIDColumn, JobColumn, LabelIDColumn}
		mutableColumns = mysql.ColumnList{UserIDColumn, JobColumn, LabelIDColumn}
		defaultColumns = mysql.ColumnList{}
	)

	return fivenetJobsLabelsUsersTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:  UserIDColumn,
		Job:     JobColumn,
		LabelID: LabelIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
