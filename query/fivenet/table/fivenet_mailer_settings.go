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

var FivenetMailerSettings = newFivenetMailerSettingsTable("", "fivenet_mailer_settings", "")

type fivenetMailerSettingsTable struct {
	mysql.Table

	// Columns
	EmailID   mysql.ColumnInteger
	Signature mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetMailerSettingsTable struct {
	fivenetMailerSettingsTable

	NEW fivenetMailerSettingsTable
}

// AS creates new FivenetMailerSettingsTable with assigned alias
func (a FivenetMailerSettingsTable) AS(alias string) *FivenetMailerSettingsTable {
	return newFivenetMailerSettingsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetMailerSettingsTable with assigned schema name
func (a FivenetMailerSettingsTable) FromSchema(schemaName string) *FivenetMailerSettingsTable {
	return newFivenetMailerSettingsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetMailerSettingsTable with assigned table prefix
func (a FivenetMailerSettingsTable) WithPrefix(prefix string) *FivenetMailerSettingsTable {
	return newFivenetMailerSettingsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetMailerSettingsTable with assigned table suffix
func (a FivenetMailerSettingsTable) WithSuffix(suffix string) *FivenetMailerSettingsTable {
	return newFivenetMailerSettingsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetMailerSettingsTable(schemaName, tableName, alias string) *FivenetMailerSettingsTable {
	return &FivenetMailerSettingsTable{
		fivenetMailerSettingsTable: newFivenetMailerSettingsTableImpl(schemaName, tableName, alias),
		NEW:                        newFivenetMailerSettingsTableImpl("", "new", ""),
	}
}

func newFivenetMailerSettingsTableImpl(schemaName, tableName, alias string) fivenetMailerSettingsTable {
	var (
		EmailIDColumn   = mysql.IntegerColumn("email_id")
		SignatureColumn = mysql.StringColumn("signature")
		allColumns      = mysql.ColumnList{EmailIDColumn, SignatureColumn}
		mutableColumns  = mysql.ColumnList{SignatureColumn}
		defaultColumns  = mysql.ColumnList{}
	)

	return fivenetMailerSettingsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		EmailID:   EmailIDColumn,
		Signature: SignatureColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
