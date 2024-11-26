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

var FivenetCentrumUnitsJobAccess = newFivenetCentrumUnitsJobAccessTable("", "fivenet_centrum_units_job_access", "")

type fivenetCentrumUnitsJobAccessTable struct {
	mysql.Table

	// Columns
	ID           mysql.ColumnInteger
	CreatedAt    mysql.ColumnTimestamp
	UnitID       mysql.ColumnInteger
	Job          mysql.ColumnString
	MinimumGrade mysql.ColumnInteger
	Access       mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetCentrumUnitsJobAccessTable struct {
	fivenetCentrumUnitsJobAccessTable

	NEW fivenetCentrumUnitsJobAccessTable
}

// AS creates new FivenetCentrumUnitsJobAccessTable with assigned alias
func (a FivenetCentrumUnitsJobAccessTable) AS(alias string) *FivenetCentrumUnitsJobAccessTable {
	return newFivenetCentrumUnitsJobAccessTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetCentrumUnitsJobAccessTable with assigned schema name
func (a FivenetCentrumUnitsJobAccessTable) FromSchema(schemaName string) *FivenetCentrumUnitsJobAccessTable {
	return newFivenetCentrumUnitsJobAccessTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetCentrumUnitsJobAccessTable with assigned table prefix
func (a FivenetCentrumUnitsJobAccessTable) WithPrefix(prefix string) *FivenetCentrumUnitsJobAccessTable {
	return newFivenetCentrumUnitsJobAccessTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetCentrumUnitsJobAccessTable with assigned table suffix
func (a FivenetCentrumUnitsJobAccessTable) WithSuffix(suffix string) *FivenetCentrumUnitsJobAccessTable {
	return newFivenetCentrumUnitsJobAccessTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetCentrumUnitsJobAccessTable(schemaName, tableName, alias string) *FivenetCentrumUnitsJobAccessTable {
	return &FivenetCentrumUnitsJobAccessTable{
		fivenetCentrumUnitsJobAccessTable: newFivenetCentrumUnitsJobAccessTableImpl(schemaName, tableName, alias),
		NEW:                               newFivenetCentrumUnitsJobAccessTableImpl("", "new", ""),
	}
}

func newFivenetCentrumUnitsJobAccessTableImpl(schemaName, tableName, alias string) fivenetCentrumUnitsJobAccessTable {
	var (
		IDColumn           = mysql.IntegerColumn("id")
		CreatedAtColumn    = mysql.TimestampColumn("created_at")
		UnitIDColumn       = mysql.IntegerColumn("unit_id")
		JobColumn          = mysql.StringColumn("job")
		MinimumGradeColumn = mysql.IntegerColumn("minimum_grade")
		AccessColumn       = mysql.IntegerColumn("access")
		allColumns         = mysql.ColumnList{IDColumn, CreatedAtColumn, UnitIDColumn, JobColumn, MinimumGradeColumn, AccessColumn}
		mutableColumns     = mysql.ColumnList{CreatedAtColumn, UnitIDColumn, JobColumn, MinimumGradeColumn, AccessColumn}
	)

	return fivenetCentrumUnitsJobAccessTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:           IDColumn,
		CreatedAt:    CreatedAtColumn,
		UnitID:       UnitIDColumn,
		Job:          JobColumn,
		MinimumGrade: MinimumGradeColumn,
		Access:       AccessColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}