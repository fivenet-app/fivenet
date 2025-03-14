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

var FivenetCalendarEntries = newFivenetCalendarEntriesTable("", "fivenet_calendar_entries", "")

type fivenetCalendarEntriesTable struct {
	mysql.Table

	// Columns
	ID         mysql.ColumnInteger
	CreatedAt  mysql.ColumnTimestamp
	UpdatedAt  mysql.ColumnTimestamp
	DeletedAt  mysql.ColumnTimestamp
	CalendarID mysql.ColumnInteger
	Job        mysql.ColumnString
	StartTime  mysql.ColumnTimestamp
	EndTime    mysql.ColumnTimestamp
	Title      mysql.ColumnString
	Content    mysql.ColumnString
	Closed     mysql.ColumnBool
	RsvpOpen   mysql.ColumnBool
	CreatorID  mysql.ColumnInteger
	CreatorJob mysql.ColumnString
	Recurring  mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetCalendarEntriesTable struct {
	fivenetCalendarEntriesTable

	NEW fivenetCalendarEntriesTable
}

// AS creates new FivenetCalendarEntriesTable with assigned alias
func (a FivenetCalendarEntriesTable) AS(alias string) *FivenetCalendarEntriesTable {
	return newFivenetCalendarEntriesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetCalendarEntriesTable with assigned schema name
func (a FivenetCalendarEntriesTable) FromSchema(schemaName string) *FivenetCalendarEntriesTable {
	return newFivenetCalendarEntriesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetCalendarEntriesTable with assigned table prefix
func (a FivenetCalendarEntriesTable) WithPrefix(prefix string) *FivenetCalendarEntriesTable {
	return newFivenetCalendarEntriesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetCalendarEntriesTable with assigned table suffix
func (a FivenetCalendarEntriesTable) WithSuffix(suffix string) *FivenetCalendarEntriesTable {
	return newFivenetCalendarEntriesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetCalendarEntriesTable(schemaName, tableName, alias string) *FivenetCalendarEntriesTable {
	return &FivenetCalendarEntriesTable{
		fivenetCalendarEntriesTable: newFivenetCalendarEntriesTableImpl(schemaName, tableName, alias),
		NEW:                         newFivenetCalendarEntriesTableImpl("", "new", ""),
	}
}

func newFivenetCalendarEntriesTableImpl(schemaName, tableName, alias string) fivenetCalendarEntriesTable {
	var (
		IDColumn         = mysql.IntegerColumn("id")
		CreatedAtColumn  = mysql.TimestampColumn("created_at")
		UpdatedAtColumn  = mysql.TimestampColumn("updated_at")
		DeletedAtColumn  = mysql.TimestampColumn("deleted_at")
		CalendarIDColumn = mysql.IntegerColumn("calendar_id")
		JobColumn        = mysql.StringColumn("job")
		StartTimeColumn  = mysql.TimestampColumn("start_time")
		EndTimeColumn    = mysql.TimestampColumn("end_time")
		TitleColumn      = mysql.StringColumn("title")
		ContentColumn    = mysql.StringColumn("content")
		ClosedColumn     = mysql.BoolColumn("closed")
		RsvpOpenColumn   = mysql.BoolColumn("rsvp_open")
		CreatorIDColumn  = mysql.IntegerColumn("creator_id")
		CreatorJobColumn = mysql.StringColumn("creator_job")
		RecurringColumn  = mysql.StringColumn("recurring")
		allColumns       = mysql.ColumnList{IDColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, CalendarIDColumn, JobColumn, StartTimeColumn, EndTimeColumn, TitleColumn, ContentColumn, ClosedColumn, RsvpOpenColumn, CreatorIDColumn, CreatorJobColumn, RecurringColumn}
		mutableColumns   = mysql.ColumnList{CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, CalendarIDColumn, JobColumn, StartTimeColumn, EndTimeColumn, TitleColumn, ContentColumn, ClosedColumn, RsvpOpenColumn, CreatorIDColumn, CreatorJobColumn, RecurringColumn}
		defaultColumns   = mysql.ColumnList{CreatedAtColumn, ClosedColumn, RsvpOpenColumn}
	)

	return fivenetCalendarEntriesTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:         IDColumn,
		CreatedAt:  CreatedAtColumn,
		UpdatedAt:  UpdatedAtColumn,
		DeletedAt:  DeletedAtColumn,
		CalendarID: CalendarIDColumn,
		Job:        JobColumn,
		StartTime:  StartTimeColumn,
		EndTime:    EndTimeColumn,
		Title:      TitleColumn,
		Content:    ContentColumn,
		Closed:     ClosedColumn,
		RsvpOpen:   RsvpOpenColumn,
		CreatorID:  CreatorIDColumn,
		CreatorJob: CreatorJobColumn,
		Recurring:  RecurringColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
