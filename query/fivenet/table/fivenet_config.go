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

var FivenetConfig = newFivenetConfigTable("", "fivenet_config", "")

type fivenetConfigTable struct {
	mysql.Table

	// Columns
	Key          mysql.ColumnInteger
	CreatedAt    mysql.ColumnTimestamp
	UpdatedAt    mysql.ColumnTimestamp
	AppConfig    mysql.ColumnString
	PluginConfig mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetConfigTable struct {
	fivenetConfigTable

	NEW fivenetConfigTable
}

// AS creates new FivenetConfigTable with assigned alias
func (a FivenetConfigTable) AS(alias string) *FivenetConfigTable {
	return newFivenetConfigTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetConfigTable with assigned schema name
func (a FivenetConfigTable) FromSchema(schemaName string) *FivenetConfigTable {
	return newFivenetConfigTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetConfigTable with assigned table prefix
func (a FivenetConfigTable) WithPrefix(prefix string) *FivenetConfigTable {
	return newFivenetConfigTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetConfigTable with assigned table suffix
func (a FivenetConfigTable) WithSuffix(suffix string) *FivenetConfigTable {
	return newFivenetConfigTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetConfigTable(schemaName, tableName, alias string) *FivenetConfigTable {
	return &FivenetConfigTable{
		fivenetConfigTable: newFivenetConfigTableImpl(schemaName, tableName, alias),
		NEW:                newFivenetConfigTableImpl("", "new", ""),
	}
}

func newFivenetConfigTableImpl(schemaName, tableName, alias string) fivenetConfigTable {
	var (
		KeyColumn          = mysql.IntegerColumn("key")
		CreatedAtColumn    = mysql.TimestampColumn("created_at")
		UpdatedAtColumn    = mysql.TimestampColumn("updated_at")
		AppConfigColumn    = mysql.StringColumn("app_config")
		PluginConfigColumn = mysql.StringColumn("plugin_config")
		allColumns         = mysql.ColumnList{KeyColumn, CreatedAtColumn, UpdatedAtColumn, AppConfigColumn, PluginConfigColumn}
		mutableColumns     = mysql.ColumnList{CreatedAtColumn, UpdatedAtColumn, AppConfigColumn, PluginConfigColumn}
		defaultColumns     = mysql.ColumnList{CreatedAtColumn}
	)

	return fivenetConfigTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Key:          KeyColumn,
		CreatedAt:    CreatedAtColumn,
		UpdatedAt:    UpdatedAtColumn,
		AppConfig:    AppConfigColumn,
		PluginConfig: PluginConfigColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
