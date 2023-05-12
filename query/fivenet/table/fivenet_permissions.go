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

var FivenetPermissions = newFivenetPermissionsTable("", "fivenet_permissions", "")

type fivenetPermissionsTable struct {
	mysql.Table

	// Columns
	ID        mysql.ColumnInteger
	CreatedAt mysql.ColumnTimestamp
	Category  mysql.ColumnString
	Name      mysql.ColumnString
	GuardName mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetPermissionsTable struct {
	fivenetPermissionsTable

	NEW fivenetPermissionsTable
}

// AS creates new FivenetPermissionsTable with assigned alias
func (a FivenetPermissionsTable) AS(alias string) *FivenetPermissionsTable {
	return newFivenetPermissionsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetPermissionsTable with assigned schema name
func (a FivenetPermissionsTable) FromSchema(schemaName string) *FivenetPermissionsTable {
	return newFivenetPermissionsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetPermissionsTable with assigned table prefix
func (a FivenetPermissionsTable) WithPrefix(prefix string) *FivenetPermissionsTable {
	return newFivenetPermissionsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetPermissionsTable with assigned table suffix
func (a FivenetPermissionsTable) WithSuffix(suffix string) *FivenetPermissionsTable {
	return newFivenetPermissionsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetPermissionsTable(schemaName, tableName, alias string) *FivenetPermissionsTable {
	return &FivenetPermissionsTable{
		fivenetPermissionsTable: newFivenetPermissionsTableImpl(schemaName, tableName, alias),
		NEW:                     newFivenetPermissionsTableImpl("", "new", ""),
	}
}

func newFivenetPermissionsTableImpl(schemaName, tableName, alias string) fivenetPermissionsTable {
	var (
		IDColumn        = mysql.IntegerColumn("id")
		CreatedAtColumn = mysql.TimestampColumn("created_at")
		CategoryColumn  = mysql.StringColumn("category")
		NameColumn      = mysql.StringColumn("name")
		GuardNameColumn = mysql.StringColumn("guard_name")
		allColumns      = mysql.ColumnList{IDColumn, CreatedAtColumn, CategoryColumn, NameColumn, GuardNameColumn}
		mutableColumns  = mysql.ColumnList{CreatedAtColumn, CategoryColumn, NameColumn, GuardNameColumn}
	)

	return fivenetPermissionsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		CreatedAt: CreatedAtColumn,
		Category:  CategoryColumn,
		Name:      NameColumn,
		GuardName: GuardNameColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
