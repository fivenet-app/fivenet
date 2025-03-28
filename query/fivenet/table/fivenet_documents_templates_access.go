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

var FivenetDocumentsTemplatesAccess = newFivenetDocumentsTemplatesAccessTable("", "fivenet_documents_templates_access", "")

type fivenetDocumentsTemplatesAccessTable struct {
	mysql.Table

	// Columns
	ID           mysql.ColumnInteger
	TargetID     mysql.ColumnInteger
	Job          mysql.ColumnString
	MinimumGrade mysql.ColumnInteger
	Access       mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetDocumentsTemplatesAccessTable struct {
	fivenetDocumentsTemplatesAccessTable

	NEW fivenetDocumentsTemplatesAccessTable
}

// AS creates new FivenetDocumentsTemplatesAccessTable with assigned alias
func (a FivenetDocumentsTemplatesAccessTable) AS(alias string) *FivenetDocumentsTemplatesAccessTable {
	return newFivenetDocumentsTemplatesAccessTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetDocumentsTemplatesAccessTable with assigned schema name
func (a FivenetDocumentsTemplatesAccessTable) FromSchema(schemaName string) *FivenetDocumentsTemplatesAccessTable {
	return newFivenetDocumentsTemplatesAccessTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetDocumentsTemplatesAccessTable with assigned table prefix
func (a FivenetDocumentsTemplatesAccessTable) WithPrefix(prefix string) *FivenetDocumentsTemplatesAccessTable {
	return newFivenetDocumentsTemplatesAccessTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetDocumentsTemplatesAccessTable with assigned table suffix
func (a FivenetDocumentsTemplatesAccessTable) WithSuffix(suffix string) *FivenetDocumentsTemplatesAccessTable {
	return newFivenetDocumentsTemplatesAccessTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetDocumentsTemplatesAccessTable(schemaName, tableName, alias string) *FivenetDocumentsTemplatesAccessTable {
	return &FivenetDocumentsTemplatesAccessTable{
		fivenetDocumentsTemplatesAccessTable: newFivenetDocumentsTemplatesAccessTableImpl(schemaName, tableName, alias),
		NEW:                                  newFivenetDocumentsTemplatesAccessTableImpl("", "new", ""),
	}
}

func newFivenetDocumentsTemplatesAccessTableImpl(schemaName, tableName, alias string) fivenetDocumentsTemplatesAccessTable {
	var (
		IDColumn           = mysql.IntegerColumn("id")
		TargetIDColumn     = mysql.IntegerColumn("target_id")
		JobColumn          = mysql.StringColumn("job")
		MinimumGradeColumn = mysql.IntegerColumn("minimum_grade")
		AccessColumn       = mysql.IntegerColumn("access")
		allColumns         = mysql.ColumnList{IDColumn, TargetIDColumn, JobColumn, MinimumGradeColumn, AccessColumn}
		mutableColumns     = mysql.ColumnList{TargetIDColumn, JobColumn, MinimumGradeColumn, AccessColumn}
		defaultColumns     = mysql.ColumnList{}
	)

	return fivenetDocumentsTemplatesAccessTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:           IDColumn,
		TargetID:     TargetIDColumn,
		Job:          JobColumn,
		MinimumGrade: MinimumGradeColumn,
		Access:       AccessColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
