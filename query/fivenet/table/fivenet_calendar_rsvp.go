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

var FivenetCalendarRsvp = newFivenetCalendarRsvpTable("", "fivenet_calendar_rsvp", "")

type fivenetCalendarRsvpTable struct {
	mysql.Table

	// Columns
	EntryID   mysql.ColumnInteger
	CreatedAt mysql.ColumnTimestamp
	UserID    mysql.ColumnInteger
	Response  mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetCalendarRsvpTable struct {
	fivenetCalendarRsvpTable

	NEW fivenetCalendarRsvpTable
}

// AS creates new FivenetCalendarRsvpTable with assigned alias
func (a FivenetCalendarRsvpTable) AS(alias string) *FivenetCalendarRsvpTable {
	return newFivenetCalendarRsvpTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetCalendarRsvpTable with assigned schema name
func (a FivenetCalendarRsvpTable) FromSchema(schemaName string) *FivenetCalendarRsvpTable {
	return newFivenetCalendarRsvpTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetCalendarRsvpTable with assigned table prefix
func (a FivenetCalendarRsvpTable) WithPrefix(prefix string) *FivenetCalendarRsvpTable {
	return newFivenetCalendarRsvpTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetCalendarRsvpTable with assigned table suffix
func (a FivenetCalendarRsvpTable) WithSuffix(suffix string) *FivenetCalendarRsvpTable {
	return newFivenetCalendarRsvpTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetCalendarRsvpTable(schemaName, tableName, alias string) *FivenetCalendarRsvpTable {
	return &FivenetCalendarRsvpTable{
		fivenetCalendarRsvpTable: newFivenetCalendarRsvpTableImpl(schemaName, tableName, alias),
		NEW:                      newFivenetCalendarRsvpTableImpl("", "new", ""),
	}
}

func newFivenetCalendarRsvpTableImpl(schemaName, tableName, alias string) fivenetCalendarRsvpTable {
	var (
		EntryIDColumn   = mysql.IntegerColumn("entry_id")
		CreatedAtColumn = mysql.TimestampColumn("created_at")
		UserIDColumn    = mysql.IntegerColumn("user_id")
		ResponseColumn  = mysql.IntegerColumn("response")
		allColumns      = mysql.ColumnList{EntryIDColumn, CreatedAtColumn, UserIDColumn, ResponseColumn}
		mutableColumns  = mysql.ColumnList{CreatedAtColumn, ResponseColumn}
		defaultColumns  = mysql.ColumnList{CreatedAtColumn, ResponseColumn}
	)

	return fivenetCalendarRsvpTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		EntryID:   EntryIDColumn,
		CreatedAt: CreatedAtColumn,
		UserID:    UserIDColumn,
		Response:  ResponseColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
