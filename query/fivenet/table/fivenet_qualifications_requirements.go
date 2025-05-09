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

var FivenetQualificationsRequirements = newFivenetQualificationsRequirementsTable("", "fivenet_qualifications_requirements", "")

type fivenetQualificationsRequirementsTable struct {
	mysql.Table

	// Columns
	ID                    mysql.ColumnInteger
	CreatedAt             mysql.ColumnTimestamp
	QualificationID       mysql.ColumnInteger
	TargetQualificationID mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetQualificationsRequirementsTable struct {
	fivenetQualificationsRequirementsTable

	NEW fivenetQualificationsRequirementsTable
}

// AS creates new FivenetQualificationsRequirementsTable with assigned alias
func (a FivenetQualificationsRequirementsTable) AS(alias string) *FivenetQualificationsRequirementsTable {
	return newFivenetQualificationsRequirementsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetQualificationsRequirementsTable with assigned schema name
func (a FivenetQualificationsRequirementsTable) FromSchema(schemaName string) *FivenetQualificationsRequirementsTable {
	return newFivenetQualificationsRequirementsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetQualificationsRequirementsTable with assigned table prefix
func (a FivenetQualificationsRequirementsTable) WithPrefix(prefix string) *FivenetQualificationsRequirementsTable {
	return newFivenetQualificationsRequirementsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetQualificationsRequirementsTable with assigned table suffix
func (a FivenetQualificationsRequirementsTable) WithSuffix(suffix string) *FivenetQualificationsRequirementsTable {
	return newFivenetQualificationsRequirementsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetQualificationsRequirementsTable(schemaName, tableName, alias string) *FivenetQualificationsRequirementsTable {
	return &FivenetQualificationsRequirementsTable{
		fivenetQualificationsRequirementsTable: newFivenetQualificationsRequirementsTableImpl(schemaName, tableName, alias),
		NEW:                                    newFivenetQualificationsRequirementsTableImpl("", "new", ""),
	}
}

func newFivenetQualificationsRequirementsTableImpl(schemaName, tableName, alias string) fivenetQualificationsRequirementsTable {
	var (
		IDColumn                    = mysql.IntegerColumn("id")
		CreatedAtColumn             = mysql.TimestampColumn("created_at")
		QualificationIDColumn       = mysql.IntegerColumn("qualification_id")
		TargetQualificationIDColumn = mysql.IntegerColumn("target_qualification_id")
		allColumns                  = mysql.ColumnList{IDColumn, CreatedAtColumn, QualificationIDColumn, TargetQualificationIDColumn}
		mutableColumns              = mysql.ColumnList{CreatedAtColumn, QualificationIDColumn, TargetQualificationIDColumn}
		defaultColumns              = mysql.ColumnList{CreatedAtColumn}
	)

	return fivenetQualificationsRequirementsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                    IDColumn,
		CreatedAt:             CreatedAtColumn,
		QualificationID:       QualificationIDColumn,
		TargetQualificationID: TargetQualificationIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
