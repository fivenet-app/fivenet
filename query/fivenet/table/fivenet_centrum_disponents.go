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

var FivenetCentrumDisponents = newFivenetCentrumDisponentsTable("", "fivenet_centrum_disponents", "")

type fivenetCentrumDisponentsTable struct {
	mysql.Table

	// Columns
	Job    mysql.ColumnString
	UserID mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetCentrumDisponentsTable struct {
	fivenetCentrumDisponentsTable

	NEW fivenetCentrumDisponentsTable
}

// AS creates new FivenetCentrumDisponentsTable with assigned alias
func (a FivenetCentrumDisponentsTable) AS(alias string) *FivenetCentrumDisponentsTable {
	return newFivenetCentrumDisponentsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetCentrumDisponentsTable with assigned schema name
func (a FivenetCentrumDisponentsTable) FromSchema(schemaName string) *FivenetCentrumDisponentsTable {
	return newFivenetCentrumDisponentsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetCentrumDisponentsTable with assigned table prefix
func (a FivenetCentrumDisponentsTable) WithPrefix(prefix string) *FivenetCentrumDisponentsTable {
	return newFivenetCentrumDisponentsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetCentrumDisponentsTable with assigned table suffix
func (a FivenetCentrumDisponentsTable) WithSuffix(suffix string) *FivenetCentrumDisponentsTable {
	return newFivenetCentrumDisponentsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetCentrumDisponentsTable(schemaName, tableName, alias string) *FivenetCentrumDisponentsTable {
	return &FivenetCentrumDisponentsTable{
		fivenetCentrumDisponentsTable: newFivenetCentrumDisponentsTableImpl(schemaName, tableName, alias),
		NEW:                           newFivenetCentrumDisponentsTableImpl("", "new", ""),
	}
}

func newFivenetCentrumDisponentsTableImpl(schemaName, tableName, alias string) fivenetCentrumDisponentsTable {
	var (
		JobColumn      = mysql.StringColumn("job")
		UserIDColumn   = mysql.IntegerColumn("user_id")
		allColumns     = mysql.ColumnList{JobColumn, UserIDColumn}
		mutableColumns = mysql.ColumnList{JobColumn}
		defaultColumns = mysql.ColumnList{}
	)

	return fivenetCentrumDisponentsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Job:    JobColumn,
		UserID: UserIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
