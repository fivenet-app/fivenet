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

var FivenetCentrumDispatchers = newFivenetCentrumDispatchersTable("", "fivenet_centrum_dispatchers", "")

type fivenetCentrumDispatchersTable struct {
	mysql.Table

	// Columns
	Job    mysql.ColumnString
	UserID mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetCentrumDispatchersTable struct {
	fivenetCentrumDispatchersTable

	NEW fivenetCentrumDispatchersTable
}

// AS creates new FivenetCentrumDispatchersTable with assigned alias
func (a FivenetCentrumDispatchersTable) AS(alias string) *FivenetCentrumDispatchersTable {
	return newFivenetCentrumDispatchersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetCentrumDispatchersTable with assigned schema name
func (a FivenetCentrumDispatchersTable) FromSchema(schemaName string) *FivenetCentrumDispatchersTable {
	return newFivenetCentrumDispatchersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetCentrumDispatchersTable with assigned table prefix
func (a FivenetCentrumDispatchersTable) WithPrefix(prefix string) *FivenetCentrumDispatchersTable {
	return newFivenetCentrumDispatchersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetCentrumDispatchersTable with assigned table suffix
func (a FivenetCentrumDispatchersTable) WithSuffix(suffix string) *FivenetCentrumDispatchersTable {
	return newFivenetCentrumDispatchersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetCentrumDispatchersTable(schemaName, tableName, alias string) *FivenetCentrumDispatchersTable {
	return &FivenetCentrumDispatchersTable{
		fivenetCentrumDispatchersTable: newFivenetCentrumDispatchersTableImpl(schemaName, tableName, alias),
		NEW:                            newFivenetCentrumDispatchersTableImpl("", "new", ""),
	}
}

func newFivenetCentrumDispatchersTableImpl(schemaName, tableName, alias string) fivenetCentrumDispatchersTable {
	var (
		JobColumn      = mysql.StringColumn("job")
		UserIDColumn   = mysql.IntegerColumn("user_id")
		allColumns     = mysql.ColumnList{JobColumn, UserIDColumn}
		mutableColumns = mysql.ColumnList{JobColumn}
		defaultColumns = mysql.ColumnList{}
	)

	return fivenetCentrumDispatchersTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Job:    JobColumn,
		UserID: UserIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
