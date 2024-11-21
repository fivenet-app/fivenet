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

var FivenetMailerThreads = newFivenetMailerThreadsTable("", "fivenet_mailer_threads", "")

type fivenetMailerThreadsTable struct {
	mysql.Table

	// Columns
	ID             mysql.ColumnInteger
	CreatedAt      mysql.ColumnTimestamp
	UpdatedAt      mysql.ColumnTimestamp
	DeletedAt      mysql.ColumnTimestamp
	Title          mysql.ColumnString
	CreatorEmailID mysql.ColumnInteger
	CreatorID      mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetMailerThreadsTable struct {
	fivenetMailerThreadsTable

	NEW fivenetMailerThreadsTable
}

// AS creates new FivenetMailerThreadsTable with assigned alias
func (a FivenetMailerThreadsTable) AS(alias string) *FivenetMailerThreadsTable {
	return newFivenetMailerThreadsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetMailerThreadsTable with assigned schema name
func (a FivenetMailerThreadsTable) FromSchema(schemaName string) *FivenetMailerThreadsTable {
	return newFivenetMailerThreadsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetMailerThreadsTable with assigned table prefix
func (a FivenetMailerThreadsTable) WithPrefix(prefix string) *FivenetMailerThreadsTable {
	return newFivenetMailerThreadsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetMailerThreadsTable with assigned table suffix
func (a FivenetMailerThreadsTable) WithSuffix(suffix string) *FivenetMailerThreadsTable {
	return newFivenetMailerThreadsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetMailerThreadsTable(schemaName, tableName, alias string) *FivenetMailerThreadsTable {
	return &FivenetMailerThreadsTable{
		fivenetMailerThreadsTable: newFivenetMailerThreadsTableImpl(schemaName, tableName, alias),
		NEW:                       newFivenetMailerThreadsTableImpl("", "new", ""),
	}
}

func newFivenetMailerThreadsTableImpl(schemaName, tableName, alias string) fivenetMailerThreadsTable {
	var (
		IDColumn             = mysql.IntegerColumn("id")
		CreatedAtColumn      = mysql.TimestampColumn("created_at")
		UpdatedAtColumn      = mysql.TimestampColumn("updated_at")
		DeletedAtColumn      = mysql.TimestampColumn("deleted_at")
		TitleColumn          = mysql.StringColumn("title")
		CreatorEmailIDColumn = mysql.IntegerColumn("creator_email_id")
		CreatorIDColumn      = mysql.IntegerColumn("creator_id")
		allColumns           = mysql.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, TitleColumn, CreatorEmailIDColumn, CreatorIDColumn}
		mutableColumns       = mysql.ColumnList{CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, TitleColumn, CreatorEmailIDColumn, CreatorIDColumn}
	)

	return fivenetMailerThreadsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		CreatedAt:      CreatedAtColumn,
		UpdatedAt:      UpdatedAtColumn,
		DeletedAt:      DeletedAtColumn,
		Title:          TitleColumn,
		CreatorEmailID: CreatorEmailIDColumn,
		CreatorID:      CreatorIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}