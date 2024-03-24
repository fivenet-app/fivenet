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

var FivenetDocumentsComments = newFivenetDocumentsCommentsTable("", "fivenet_documents_comments", "")

type fivenetDocumentsCommentsTable struct {
	mysql.Table

	// Columns
	ID         mysql.ColumnInteger
	CreatedAt  mysql.ColumnTimestamp
	UpdatedAt  mysql.ColumnTimestamp
	DeletedAt  mysql.ColumnTimestamp
	DocumentID mysql.ColumnInteger
	Comment    mysql.ColumnString
	CreatorID  mysql.ColumnInteger
	CreatorJob mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetDocumentsCommentsTable struct {
	fivenetDocumentsCommentsTable

	NEW fivenetDocumentsCommentsTable
}

// AS creates new FivenetDocumentsCommentsTable with assigned alias
func (a FivenetDocumentsCommentsTable) AS(alias string) *FivenetDocumentsCommentsTable {
	return newFivenetDocumentsCommentsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetDocumentsCommentsTable with assigned schema name
func (a FivenetDocumentsCommentsTable) FromSchema(schemaName string) *FivenetDocumentsCommentsTable {
	return newFivenetDocumentsCommentsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetDocumentsCommentsTable with assigned table prefix
func (a FivenetDocumentsCommentsTable) WithPrefix(prefix string) *FivenetDocumentsCommentsTable {
	return newFivenetDocumentsCommentsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetDocumentsCommentsTable with assigned table suffix
func (a FivenetDocumentsCommentsTable) WithSuffix(suffix string) *FivenetDocumentsCommentsTable {
	return newFivenetDocumentsCommentsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetDocumentsCommentsTable(schemaName, tableName, alias string) *FivenetDocumentsCommentsTable {
	return &FivenetDocumentsCommentsTable{
		fivenetDocumentsCommentsTable: newFivenetDocumentsCommentsTableImpl(schemaName, tableName, alias),
		NEW:                           newFivenetDocumentsCommentsTableImpl("", "new", ""),
	}
}

func newFivenetDocumentsCommentsTableImpl(schemaName, tableName, alias string) fivenetDocumentsCommentsTable {
	var (
		IDColumn         = mysql.IntegerColumn("id")
		CreatedAtColumn  = mysql.TimestampColumn("created_at")
		UpdatedAtColumn  = mysql.TimestampColumn("updated_at")
		DeletedAtColumn  = mysql.TimestampColumn("deleted_at")
		DocumentIDColumn = mysql.IntegerColumn("document_id")
		CommentColumn    = mysql.StringColumn("comment")
		CreatorIDColumn  = mysql.IntegerColumn("creator_id")
		CreatorJobColumn = mysql.StringColumn("creator_job")
		allColumns       = mysql.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, DocumentIDColumn, CommentColumn, CreatorIDColumn, CreatorJobColumn}
		mutableColumns   = mysql.ColumnList{CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, DocumentIDColumn, CommentColumn, CreatorIDColumn, CreatorJobColumn}
	)

	return fivenetDocumentsCommentsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:         IDColumn,
		CreatedAt:  CreatedAtColumn,
		UpdatedAt:  UpdatedAtColumn,
		DeletedAt:  DeletedAtColumn,
		DocumentID: DocumentIDColumn,
		Comment:    CommentColumn,
		CreatorID:  CreatorIDColumn,
		CreatorJob: CreatorJobColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
