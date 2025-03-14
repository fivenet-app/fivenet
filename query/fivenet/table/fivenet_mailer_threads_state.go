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

var FivenetMailerThreadsState = newFivenetMailerThreadsStateTable("", "fivenet_mailer_threads_state", "")

type fivenetMailerThreadsStateTable struct {
	mysql.Table

	// Columns
	ThreadID  mysql.ColumnInteger
	EmailID   mysql.ColumnInteger
	LastRead  mysql.ColumnTimestamp
	Unread    mysql.ColumnBool
	Important mysql.ColumnBool
	Favorite  mysql.ColumnBool
	Muted     mysql.ColumnBool
	Archived  mysql.ColumnBool

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetMailerThreadsStateTable struct {
	fivenetMailerThreadsStateTable

	NEW fivenetMailerThreadsStateTable
}

// AS creates new FivenetMailerThreadsStateTable with assigned alias
func (a FivenetMailerThreadsStateTable) AS(alias string) *FivenetMailerThreadsStateTable {
	return newFivenetMailerThreadsStateTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetMailerThreadsStateTable with assigned schema name
func (a FivenetMailerThreadsStateTable) FromSchema(schemaName string) *FivenetMailerThreadsStateTable {
	return newFivenetMailerThreadsStateTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetMailerThreadsStateTable with assigned table prefix
func (a FivenetMailerThreadsStateTable) WithPrefix(prefix string) *FivenetMailerThreadsStateTable {
	return newFivenetMailerThreadsStateTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetMailerThreadsStateTable with assigned table suffix
func (a FivenetMailerThreadsStateTable) WithSuffix(suffix string) *FivenetMailerThreadsStateTable {
	return newFivenetMailerThreadsStateTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetMailerThreadsStateTable(schemaName, tableName, alias string) *FivenetMailerThreadsStateTable {
	return &FivenetMailerThreadsStateTable{
		fivenetMailerThreadsStateTable: newFivenetMailerThreadsStateTableImpl(schemaName, tableName, alias),
		NEW:                            newFivenetMailerThreadsStateTableImpl("", "new", ""),
	}
}

func newFivenetMailerThreadsStateTableImpl(schemaName, tableName, alias string) fivenetMailerThreadsStateTable {
	var (
		ThreadIDColumn  = mysql.IntegerColumn("thread_id")
		EmailIDColumn   = mysql.IntegerColumn("email_id")
		LastReadColumn  = mysql.TimestampColumn("last_read")
		UnreadColumn    = mysql.BoolColumn("unread")
		ImportantColumn = mysql.BoolColumn("important")
		FavoriteColumn  = mysql.BoolColumn("favorite")
		MutedColumn     = mysql.BoolColumn("muted")
		ArchivedColumn  = mysql.BoolColumn("archived")
		allColumns      = mysql.ColumnList{ThreadIDColumn, EmailIDColumn, LastReadColumn, UnreadColumn, ImportantColumn, FavoriteColumn, MutedColumn, ArchivedColumn}
		mutableColumns  = mysql.ColumnList{LastReadColumn, UnreadColumn, ImportantColumn, FavoriteColumn, MutedColumn, ArchivedColumn}
		defaultColumns  = mysql.ColumnList{UnreadColumn, ImportantColumn, FavoriteColumn, MutedColumn, ArchivedColumn}
	)

	return fivenetMailerThreadsStateTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ThreadID:  ThreadIDColumn,
		EmailID:   EmailIDColumn,
		LastRead:  LastReadColumn,
		Unread:    UnreadColumn,
		Important: ImportantColumn,
		Favorite:  FavoriteColumn,
		Muted:     MutedColumn,
		Archived:  ArchivedColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
