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

var FivenetMailerSettingsBlocked = newFivenetMailerSettingsBlockedTable("", "fivenet_mailer_settings_blocked", "")

type fivenetMailerSettingsBlockedTable struct {
	mysql.Table

	// Columns
	EmailID     mysql.ColumnInteger
	TargetEmail mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetMailerSettingsBlockedTable struct {
	fivenetMailerSettingsBlockedTable

	NEW fivenetMailerSettingsBlockedTable
}

// AS creates new FivenetMailerSettingsBlockedTable with assigned alias
func (a FivenetMailerSettingsBlockedTable) AS(alias string) *FivenetMailerSettingsBlockedTable {
	return newFivenetMailerSettingsBlockedTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetMailerSettingsBlockedTable with assigned schema name
func (a FivenetMailerSettingsBlockedTable) FromSchema(schemaName string) *FivenetMailerSettingsBlockedTable {
	return newFivenetMailerSettingsBlockedTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetMailerSettingsBlockedTable with assigned table prefix
func (a FivenetMailerSettingsBlockedTable) WithPrefix(prefix string) *FivenetMailerSettingsBlockedTable {
	return newFivenetMailerSettingsBlockedTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetMailerSettingsBlockedTable with assigned table suffix
func (a FivenetMailerSettingsBlockedTable) WithSuffix(suffix string) *FivenetMailerSettingsBlockedTable {
	return newFivenetMailerSettingsBlockedTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetMailerSettingsBlockedTable(schemaName, tableName, alias string) *FivenetMailerSettingsBlockedTable {
	return &FivenetMailerSettingsBlockedTable{
		fivenetMailerSettingsBlockedTable: newFivenetMailerSettingsBlockedTableImpl(schemaName, tableName, alias),
		NEW:                               newFivenetMailerSettingsBlockedTableImpl("", "new", ""),
	}
}

func newFivenetMailerSettingsBlockedTableImpl(schemaName, tableName, alias string) fivenetMailerSettingsBlockedTable {
	var (
		EmailIDColumn     = mysql.IntegerColumn("email_id")
		TargetEmailColumn = mysql.StringColumn("target_email")
		allColumns        = mysql.ColumnList{EmailIDColumn, TargetEmailColumn}
		mutableColumns    = mysql.ColumnList{}
		defaultColumns    = mysql.ColumnList{}
	)

	return fivenetMailerSettingsBlockedTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		EmailID:     EmailIDColumn,
		TargetEmail: TargetEmailColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
