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

var FivenetQualificationsResults = newFivenetQualificationsResultsTable("", "fivenet_qualifications_results", "")

type fivenetQualificationsResultsTable struct {
	mysql.Table

	// Columns
	ID              mysql.ColumnInteger
	CreatedAt       mysql.ColumnTimestamp
	DeletedAt       mysql.ColumnTimestamp
	QualificationID mysql.ColumnInteger
	UserID          mysql.ColumnInteger
	Status          mysql.ColumnInteger
	Score           mysql.ColumnFloat
	Summary         mysql.ColumnString
	CreatorID       mysql.ColumnInteger
	CreatorJob      mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetQualificationsResultsTable struct {
	fivenetQualificationsResultsTable

	NEW fivenetQualificationsResultsTable
}

// AS creates new FivenetQualificationsResultsTable with assigned alias
func (a FivenetQualificationsResultsTable) AS(alias string) *FivenetQualificationsResultsTable {
	return newFivenetQualificationsResultsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetQualificationsResultsTable with assigned schema name
func (a FivenetQualificationsResultsTable) FromSchema(schemaName string) *FivenetQualificationsResultsTable {
	return newFivenetQualificationsResultsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetQualificationsResultsTable with assigned table prefix
func (a FivenetQualificationsResultsTable) WithPrefix(prefix string) *FivenetQualificationsResultsTable {
	return newFivenetQualificationsResultsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetQualificationsResultsTable with assigned table suffix
func (a FivenetQualificationsResultsTable) WithSuffix(suffix string) *FivenetQualificationsResultsTable {
	return newFivenetQualificationsResultsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetQualificationsResultsTable(schemaName, tableName, alias string) *FivenetQualificationsResultsTable {
	return &FivenetQualificationsResultsTable{
		fivenetQualificationsResultsTable: newFivenetQualificationsResultsTableImpl(schemaName, tableName, alias),
		NEW:                               newFivenetQualificationsResultsTableImpl("", "new", ""),
	}
}

func newFivenetQualificationsResultsTableImpl(schemaName, tableName, alias string) fivenetQualificationsResultsTable {
	var (
		IDColumn              = mysql.IntegerColumn("id")
		CreatedAtColumn       = mysql.TimestampColumn("created_at")
		DeletedAtColumn       = mysql.TimestampColumn("deleted_at")
		QualificationIDColumn = mysql.IntegerColumn("qualification_id")
		UserIDColumn          = mysql.IntegerColumn("user_id")
		StatusColumn          = mysql.IntegerColumn("status")
		ScoreColumn           = mysql.FloatColumn("score")
		SummaryColumn         = mysql.StringColumn("summary")
		CreatorIDColumn       = mysql.IntegerColumn("creator_id")
		CreatorJobColumn      = mysql.StringColumn("creator_job")
		allColumns            = mysql.ColumnList{IDColumn, CreatedAtColumn, DeletedAtColumn, QualificationIDColumn, UserIDColumn, StatusColumn, ScoreColumn, SummaryColumn, CreatorIDColumn, CreatorJobColumn}
		mutableColumns        = mysql.ColumnList{CreatedAtColumn, DeletedAtColumn, QualificationIDColumn, UserIDColumn, StatusColumn, ScoreColumn, SummaryColumn, CreatorIDColumn, CreatorJobColumn}
	)

	return fivenetQualificationsResultsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:              IDColumn,
		CreatedAt:       CreatedAtColumn,
		DeletedAt:       DeletedAtColumn,
		QualificationID: QualificationIDColumn,
		UserID:          UserIDColumn,
		Status:          StatusColumn,
		Score:           ScoreColumn,
		Summary:         SummaryColumn,
		CreatorID:       CreatorIDColumn,
		CreatorJob:      CreatorJobColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
