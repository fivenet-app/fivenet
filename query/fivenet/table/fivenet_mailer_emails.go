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

var FivenetMailerEmails = newFivenetMailerEmailsTable("", "fivenet_mailer_emails", "")

type fivenetMailerEmailsTable struct {
	mysql.Table

	// Columns
	ID           mysql.ColumnInteger
	CreatedAt    mysql.ColumnTimestamp
	UpdatedAt    mysql.ColumnTimestamp
	DeletedAt    mysql.ColumnTimestamp
	Deactivated  mysql.ColumnBool
	Job          mysql.ColumnString
	UserID       mysql.ColumnInteger
	Email        mysql.ColumnString
	EmailChanged mysql.ColumnTimestamp
	Label        mysql.ColumnString
	Internal     mysql.ColumnBool
	CreatorID    mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetMailerEmailsTable struct {
	fivenetMailerEmailsTable

	NEW fivenetMailerEmailsTable
}

// AS creates new FivenetMailerEmailsTable with assigned alias
func (a FivenetMailerEmailsTable) AS(alias string) *FivenetMailerEmailsTable {
	return newFivenetMailerEmailsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetMailerEmailsTable with assigned schema name
func (a FivenetMailerEmailsTable) FromSchema(schemaName string) *FivenetMailerEmailsTable {
	return newFivenetMailerEmailsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetMailerEmailsTable with assigned table prefix
func (a FivenetMailerEmailsTable) WithPrefix(prefix string) *FivenetMailerEmailsTable {
	return newFivenetMailerEmailsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetMailerEmailsTable with assigned table suffix
func (a FivenetMailerEmailsTable) WithSuffix(suffix string) *FivenetMailerEmailsTable {
	return newFivenetMailerEmailsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetMailerEmailsTable(schemaName, tableName, alias string) *FivenetMailerEmailsTable {
	return &FivenetMailerEmailsTable{
		fivenetMailerEmailsTable: newFivenetMailerEmailsTableImpl(schemaName, tableName, alias),
		NEW:                      newFivenetMailerEmailsTableImpl("", "new", ""),
	}
}

func newFivenetMailerEmailsTableImpl(schemaName, tableName, alias string) fivenetMailerEmailsTable {
	var (
		IDColumn           = mysql.IntegerColumn("id")
		CreatedAtColumn    = mysql.TimestampColumn("created_at")
		UpdatedAtColumn    = mysql.TimestampColumn("updated_at")
		DeletedAtColumn    = mysql.TimestampColumn("deleted_at")
		DeactivatedColumn  = mysql.BoolColumn("deactivated")
		JobColumn          = mysql.StringColumn("job")
		UserIDColumn       = mysql.IntegerColumn("user_id")
		EmailColumn        = mysql.StringColumn("email")
		EmailChangedColumn = mysql.TimestampColumn("email_changed")
		LabelColumn        = mysql.StringColumn("label")
		InternalColumn     = mysql.BoolColumn("internal")
		CreatorIDColumn    = mysql.IntegerColumn("creator_id")
		allColumns         = mysql.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, DeactivatedColumn, JobColumn, UserIDColumn, EmailColumn, EmailChangedColumn, LabelColumn, InternalColumn, CreatorIDColumn}
		mutableColumns     = mysql.ColumnList{CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, DeactivatedColumn, JobColumn, UserIDColumn, EmailColumn, EmailChangedColumn, LabelColumn, InternalColumn, CreatorIDColumn}
		defaultColumns     = mysql.ColumnList{CreatedAtColumn, DeactivatedColumn, InternalColumn}
	)

	return fivenetMailerEmailsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:           IDColumn,
		CreatedAt:    CreatedAtColumn,
		UpdatedAt:    UpdatedAtColumn,
		DeletedAt:    DeletedAtColumn,
		Deactivated:  DeactivatedColumn,
		Job:          JobColumn,
		UserID:       UserIDColumn,
		Email:        EmailColumn,
		EmailChanged: EmailChangedColumn,
		Label:        LabelColumn,
		Internal:     InternalColumn,
		CreatorID:    CreatorIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
