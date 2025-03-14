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

var FivenetJobCitizenLabels = newFivenetJobCitizenLabelsTable("", "fivenet_job_citizen_labels", "")

type fivenetJobCitizenLabelsTable struct {
	mysql.Table

	// Columns
	ID    mysql.ColumnInteger
	Job   mysql.ColumnString
	Name  mysql.ColumnString
	Color mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetJobCitizenLabelsTable struct {
	fivenetJobCitizenLabelsTable

	NEW fivenetJobCitizenLabelsTable
}

// AS creates new FivenetJobCitizenLabelsTable with assigned alias
func (a FivenetJobCitizenLabelsTable) AS(alias string) *FivenetJobCitizenLabelsTable {
	return newFivenetJobCitizenLabelsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetJobCitizenLabelsTable with assigned schema name
func (a FivenetJobCitizenLabelsTable) FromSchema(schemaName string) *FivenetJobCitizenLabelsTable {
	return newFivenetJobCitizenLabelsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetJobCitizenLabelsTable with assigned table prefix
func (a FivenetJobCitizenLabelsTable) WithPrefix(prefix string) *FivenetJobCitizenLabelsTable {
	return newFivenetJobCitizenLabelsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetJobCitizenLabelsTable with assigned table suffix
func (a FivenetJobCitizenLabelsTable) WithSuffix(suffix string) *FivenetJobCitizenLabelsTable {
	return newFivenetJobCitizenLabelsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetJobCitizenLabelsTable(schemaName, tableName, alias string) *FivenetJobCitizenLabelsTable {
	return &FivenetJobCitizenLabelsTable{
		fivenetJobCitizenLabelsTable: newFivenetJobCitizenLabelsTableImpl(schemaName, tableName, alias),
		NEW:                          newFivenetJobCitizenLabelsTableImpl("", "new", ""),
	}
}

func newFivenetJobCitizenLabelsTableImpl(schemaName, tableName, alias string) fivenetJobCitizenLabelsTable {
	var (
		IDColumn       = mysql.IntegerColumn("id")
		JobColumn      = mysql.StringColumn("job")
		NameColumn     = mysql.StringColumn("name")
		ColorColumn    = mysql.StringColumn("color")
		allColumns     = mysql.ColumnList{IDColumn, JobColumn, NameColumn, ColorColumn}
		mutableColumns = mysql.ColumnList{JobColumn, NameColumn, ColorColumn}
		defaultColumns = mysql.ColumnList{ColorColumn}
	)

	return fivenetJobCitizenLabelsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:    IDColumn,
		Job:   JobColumn,
		Name:  NameColumn,
		Color: ColorColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
