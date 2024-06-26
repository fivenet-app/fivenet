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

var FivenetMsgsThreadsUserAccess = newFivenetMsgsThreadsUserAccessTable("", "fivenet_msgs_threads_user_access", "")

type fivenetMsgsThreadsUserAccessTable struct {
	mysql.Table

	// Columns
	ID        mysql.ColumnInteger
	CreatedAt mysql.ColumnTimestamp
	ThreadID  mysql.ColumnInteger
	UserID    mysql.ColumnInteger
	Access    mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetMsgsThreadsUserAccessTable struct {
	fivenetMsgsThreadsUserAccessTable

	NEW fivenetMsgsThreadsUserAccessTable
}

// AS creates new FivenetMsgsThreadsUserAccessTable with assigned alias
func (a FivenetMsgsThreadsUserAccessTable) AS(alias string) *FivenetMsgsThreadsUserAccessTable {
	return newFivenetMsgsThreadsUserAccessTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetMsgsThreadsUserAccessTable with assigned schema name
func (a FivenetMsgsThreadsUserAccessTable) FromSchema(schemaName string) *FivenetMsgsThreadsUserAccessTable {
	return newFivenetMsgsThreadsUserAccessTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetMsgsThreadsUserAccessTable with assigned table prefix
func (a FivenetMsgsThreadsUserAccessTable) WithPrefix(prefix string) *FivenetMsgsThreadsUserAccessTable {
	return newFivenetMsgsThreadsUserAccessTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetMsgsThreadsUserAccessTable with assigned table suffix
func (a FivenetMsgsThreadsUserAccessTable) WithSuffix(suffix string) *FivenetMsgsThreadsUserAccessTable {
	return newFivenetMsgsThreadsUserAccessTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetMsgsThreadsUserAccessTable(schemaName, tableName, alias string) *FivenetMsgsThreadsUserAccessTable {
	return &FivenetMsgsThreadsUserAccessTable{
		fivenetMsgsThreadsUserAccessTable: newFivenetMsgsThreadsUserAccessTableImpl(schemaName, tableName, alias),
		NEW:                               newFivenetMsgsThreadsUserAccessTableImpl("", "new", ""),
	}
}

func newFivenetMsgsThreadsUserAccessTableImpl(schemaName, tableName, alias string) fivenetMsgsThreadsUserAccessTable {
	var (
		IDColumn        = mysql.IntegerColumn("id")
		CreatedAtColumn = mysql.TimestampColumn("created_at")
		ThreadIDColumn  = mysql.IntegerColumn("thread_id")
		UserIDColumn    = mysql.IntegerColumn("user_id")
		AccessColumn    = mysql.IntegerColumn("access")
		allColumns      = mysql.ColumnList{IDColumn, CreatedAtColumn, ThreadIDColumn, UserIDColumn, AccessColumn}
		mutableColumns  = mysql.ColumnList{CreatedAtColumn, ThreadIDColumn, UserIDColumn, AccessColumn}
	)

	return fivenetMsgsThreadsUserAccessTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		CreatedAt: CreatedAtColumn,
		ThreadID:  ThreadIDColumn,
		UserID:    UserIDColumn,
		Access:    AccessColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
