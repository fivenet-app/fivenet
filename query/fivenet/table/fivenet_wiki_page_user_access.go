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

var FivenetWikiPageUserAccess = newFivenetWikiPageUserAccessTable("", "fivenet_wiki_page_user_access", "")

type fivenetWikiPageUserAccessTable struct {
	mysql.Table

	// Columns
	ID        mysql.ColumnInteger
	CreatedAt mysql.ColumnTimestamp
	PageID    mysql.ColumnInteger
	UserID    mysql.ColumnInteger
	Access    mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetWikiPageUserAccessTable struct {
	fivenetWikiPageUserAccessTable

	NEW fivenetWikiPageUserAccessTable
}

// AS creates new FivenetWikiPageUserAccessTable with assigned alias
func (a FivenetWikiPageUserAccessTable) AS(alias string) *FivenetWikiPageUserAccessTable {
	return newFivenetWikiPageUserAccessTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetWikiPageUserAccessTable with assigned schema name
func (a FivenetWikiPageUserAccessTable) FromSchema(schemaName string) *FivenetWikiPageUserAccessTable {
	return newFivenetWikiPageUserAccessTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetWikiPageUserAccessTable with assigned table prefix
func (a FivenetWikiPageUserAccessTable) WithPrefix(prefix string) *FivenetWikiPageUserAccessTable {
	return newFivenetWikiPageUserAccessTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetWikiPageUserAccessTable with assigned table suffix
func (a FivenetWikiPageUserAccessTable) WithSuffix(suffix string) *FivenetWikiPageUserAccessTable {
	return newFivenetWikiPageUserAccessTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetWikiPageUserAccessTable(schemaName, tableName, alias string) *FivenetWikiPageUserAccessTable {
	return &FivenetWikiPageUserAccessTable{
		fivenetWikiPageUserAccessTable: newFivenetWikiPageUserAccessTableImpl(schemaName, tableName, alias),
		NEW:                            newFivenetWikiPageUserAccessTableImpl("", "new", ""),
	}
}

func newFivenetWikiPageUserAccessTableImpl(schemaName, tableName, alias string) fivenetWikiPageUserAccessTable {
	var (
		IDColumn        = mysql.IntegerColumn("id")
		CreatedAtColumn = mysql.TimestampColumn("created_at")
		PageIDColumn    = mysql.IntegerColumn("page_id")
		UserIDColumn    = mysql.IntegerColumn("user_id")
		AccessColumn    = mysql.IntegerColumn("access")
		allColumns      = mysql.ColumnList{IDColumn, CreatedAtColumn, PageIDColumn, UserIDColumn, AccessColumn}
		mutableColumns  = mysql.ColumnList{CreatedAtColumn, PageIDColumn, UserIDColumn, AccessColumn}
	)

	return fivenetWikiPageUserAccessTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		CreatedAt: CreatedAtColumn,
		PageID:    PageIDColumn,
		UserID:    UserIDColumn,
		Access:    AccessColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}