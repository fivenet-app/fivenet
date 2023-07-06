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

var FivenetCentrumUnitsUsers = newFivenetCentrumUnitsUsersTable("", "fivenet_centrum_units_users", "")

type fivenetCentrumUnitsUsersTable struct {
	mysql.Table

	// Columns
	UnitID     mysql.ColumnInteger
	UserID     mysql.ColumnInteger
	Identifier mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetCentrumUnitsUsersTable struct {
	fivenetCentrumUnitsUsersTable

	NEW fivenetCentrumUnitsUsersTable
}

// AS creates new FivenetCentrumUnitsUsersTable with assigned alias
func (a FivenetCentrumUnitsUsersTable) AS(alias string) *FivenetCentrumUnitsUsersTable {
	return newFivenetCentrumUnitsUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetCentrumUnitsUsersTable with assigned schema name
func (a FivenetCentrumUnitsUsersTable) FromSchema(schemaName string) *FivenetCentrumUnitsUsersTable {
	return newFivenetCentrumUnitsUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetCentrumUnitsUsersTable with assigned table prefix
func (a FivenetCentrumUnitsUsersTable) WithPrefix(prefix string) *FivenetCentrumUnitsUsersTable {
	return newFivenetCentrumUnitsUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetCentrumUnitsUsersTable with assigned table suffix
func (a FivenetCentrumUnitsUsersTable) WithSuffix(suffix string) *FivenetCentrumUnitsUsersTable {
	return newFivenetCentrumUnitsUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetCentrumUnitsUsersTable(schemaName, tableName, alias string) *FivenetCentrumUnitsUsersTable {
	return &FivenetCentrumUnitsUsersTable{
		fivenetCentrumUnitsUsersTable: newFivenetCentrumUnitsUsersTableImpl(schemaName, tableName, alias),
		NEW:                           newFivenetCentrumUnitsUsersTableImpl("", "new", ""),
	}
}

func newFivenetCentrumUnitsUsersTableImpl(schemaName, tableName, alias string) fivenetCentrumUnitsUsersTable {
	var (
		UnitIDColumn     = mysql.IntegerColumn("unit_id")
		UserIDColumn     = mysql.IntegerColumn("user_id")
		IdentifierColumn = mysql.StringColumn("identifier")
		allColumns       = mysql.ColumnList{UnitIDColumn, UserIDColumn, IdentifierColumn}
		mutableColumns   = mysql.ColumnList{IdentifierColumn}
	)

	return fivenetCentrumUnitsUsersTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UnitID:     UnitIDColumn,
		UserID:     UserIDColumn,
		Identifier: IdentifierColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
