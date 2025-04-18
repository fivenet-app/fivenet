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

var FivenetDocumentsReferences = newFivenetDocumentsReferencesTable("", "fivenet_documents_references", "")

type fivenetDocumentsReferencesTable struct {
	mysql.Table

	// Columns
	ID               mysql.ColumnInteger
	CreatedAt        mysql.ColumnTimestamp
	DeletedAt        mysql.ColumnTimestamp
	SourceDocumentID mysql.ColumnInteger
	Reference        mysql.ColumnInteger
	TargetDocumentID mysql.ColumnInteger
	CreatorID        mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetDocumentsReferencesTable struct {
	fivenetDocumentsReferencesTable

	NEW fivenetDocumentsReferencesTable
}

// AS creates new FivenetDocumentsReferencesTable with assigned alias
func (a FivenetDocumentsReferencesTable) AS(alias string) *FivenetDocumentsReferencesTable {
	return newFivenetDocumentsReferencesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetDocumentsReferencesTable with assigned schema name
func (a FivenetDocumentsReferencesTable) FromSchema(schemaName string) *FivenetDocumentsReferencesTable {
	return newFivenetDocumentsReferencesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetDocumentsReferencesTable with assigned table prefix
func (a FivenetDocumentsReferencesTable) WithPrefix(prefix string) *FivenetDocumentsReferencesTable {
	return newFivenetDocumentsReferencesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetDocumentsReferencesTable with assigned table suffix
func (a FivenetDocumentsReferencesTable) WithSuffix(suffix string) *FivenetDocumentsReferencesTable {
	return newFivenetDocumentsReferencesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetDocumentsReferencesTable(schemaName, tableName, alias string) *FivenetDocumentsReferencesTable {
	return &FivenetDocumentsReferencesTable{
		fivenetDocumentsReferencesTable: newFivenetDocumentsReferencesTableImpl(schemaName, tableName, alias),
		NEW:                             newFivenetDocumentsReferencesTableImpl("", "new", ""),
	}
}

func newFivenetDocumentsReferencesTableImpl(schemaName, tableName, alias string) fivenetDocumentsReferencesTable {
	var (
		IDColumn               = mysql.IntegerColumn("id")
		CreatedAtColumn        = mysql.TimestampColumn("created_at")
		DeletedAtColumn        = mysql.TimestampColumn("deleted_at")
		SourceDocumentIDColumn = mysql.IntegerColumn("source_document_id")
		ReferenceColumn        = mysql.IntegerColumn("reference")
		TargetDocumentIDColumn = mysql.IntegerColumn("target_document_id")
		CreatorIDColumn        = mysql.IntegerColumn("creator_id")
		allColumns             = mysql.ColumnList{IDColumn, CreatedAtColumn, DeletedAtColumn, SourceDocumentIDColumn, ReferenceColumn, TargetDocumentIDColumn, CreatorIDColumn}
		mutableColumns         = mysql.ColumnList{CreatedAtColumn, DeletedAtColumn, SourceDocumentIDColumn, ReferenceColumn, TargetDocumentIDColumn, CreatorIDColumn}
		defaultColumns         = mysql.ColumnList{CreatedAtColumn}
	)

	return fivenetDocumentsReferencesTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:               IDColumn,
		CreatedAt:        CreatedAtColumn,
		DeletedAt:        DeletedAtColumn,
		SourceDocumentID: SourceDocumentIDColumn,
		Reference:        ReferenceColumn,
		TargetDocumentID: TargetDocumentIDColumn,
		CreatorID:        CreatorIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
